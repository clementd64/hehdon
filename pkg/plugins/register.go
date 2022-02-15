package plugins

import (
	"errors"
	"reflect"
)

type PluginWrapper struct {
	Name   string
	Type   reflect.Type
	Config reflect.Type
}

var plugins = make(map[string]*PluginWrapper)

func Register(name string, plugin interface{}) {
	if _, found := plugins[name]; found {
		panic("plugin " + name + " already registered")
	}

	t := reflect.TypeOf(plugin)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	wrapper := &PluginWrapper{
		Name: name,
		Type: t,
	}

	for i := 0; i < t.NumField(); i++ {
		if tag := t.Field(i).Tag.Get("hehdon"); tag == "config" {
			wrapper.Config = t.Field(i).Type
			break
		}
	}
	if wrapper.Config == nil {
		panic("plugin " + name + " don't have any configuration field")
	}

	plugins[name] = wrapper
}

func Get(name string) (*PluginWrapper, error) {
	plugin, ok := plugins[name]
	if !ok {
		return nil, errors.New("plugin " + name + " not found")
	}
	return plugin, nil
}

func (w *PluginWrapper) New(config interface{}, wrapper interface{}) error {
	p := reflect.New(w.Type)

	for i := 0; i < p.Elem().NumField(); i++ {
		if tag := p.Elem().Type().Field(i).Tag.Get("hehdon"); tag == "config" {
			p.Elem().Field(i).Set(reflect.ValueOf(config))
			break
		}
	}

	if err := WrapPlugin(p.Interface(), wrapper); err != nil {
		return errors.New("plugin " + w.Name + ": " + err.Error())
	}
	return nil
}
