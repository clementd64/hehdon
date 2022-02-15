package config

import (
	"github.com/clementd64/hehdon/pkg/plugins"
	"sigs.k8s.io/yaml"
)

type config struct {
	Metadata ConfigMetadata `json:"metadata"`
}

type ConfigMetadata struct {
	Name string `json:"name"`
}

type PluginSpec struct {
	Type   string                 `json:"type"`
	Spec   interface{}            `json:"spec"`
	Plugin *plugins.PluginWrapper `json:"-"`
}

type Config struct {
	config
	PluginSpec
}

// Unmarshal field separatly to avoid custom UnmarshalJSON problem with inline
func (c *Config) UnmarshalJSON(b []byte) error {
	if err := yaml.Unmarshal(b, &c.PluginSpec); err != nil {
		return err
	}
	return yaml.Unmarshal(b, &c.config)
}
