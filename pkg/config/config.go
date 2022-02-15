package config

import (
	"encoding/json"
	"reflect"

	"github.com/clementd64/hehdon/pkg/plugins"
)

func (s *PluginSpec) UnmarshalJSON(b []byte) error {
	var unstructuredType struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(b, &unstructuredType); err != nil {
		return err
	}
	s.Type = unstructuredType.Type

	plugin, err := plugins.Get(s.Type)
	if err != nil {
		return err
	}

	spec := reflect.New(reflect.StructOf([]reflect.StructField{
		{
			Name: "Spec",
			Type: plugin.Config,
			Tag:  `json:"spec"`,
		},
	}))

	if err := json.Unmarshal(b, spec.Interface()); err != nil {
		return err
	}

	s.Spec = spec.Elem().FieldByName("Spec").Interface()
	s.Plugin = plugin

	return nil
}

func (s *PluginSpec) New(wrapper interface{}) error {
	return s.Plugin.New(s.Spec, wrapper)
}
