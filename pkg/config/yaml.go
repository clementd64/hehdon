package config

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"

	"sigs.k8s.io/yaml"
)

const yamlSeparator = "\n---"

// splitYAMLDocument is a bufio.SplitFunc for splitting YAML streams into individual documents.
func splitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	sep := len([]byte(yamlSeparator))
	if i := bytes.Index(data, []byte(yamlSeparator)); i >= 0 {
		// We have a potential document terminator
		i += sep
		after := data[i:]
		if len(after) == 0 {
			// we can't read any more characters
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}
			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + j + 1, data[0 : i-sep], nil
		}
		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func ParseReader(r io.Reader) ([]Config, error) {
	configs := []Config{}

	scanner := bufio.NewScanner(r)
	scanner.Split(splitYAMLDocument)

	for scanner.Scan() {
		config := Config{}
		if err := yaml.Unmarshal(scanner.Bytes(), &config); err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

func ParseFile(filename string) ([]Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return ParseReader(file)
}

func Parse(file string) ([]Config, error) {
	return ParseReader(strings.NewReader(file))
}
