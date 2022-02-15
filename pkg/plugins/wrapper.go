package plugins

import (
	"errors"
	"fmt"
	"reflect"
)

func WrapPlugin(plugin interface{}, wrapper interface{}) error {
	p := reflect.ValueOf(plugin)

	if p.Kind() != reflect.Ptr || p.Elem().Kind() != reflect.Struct {
		return errors.New("plugin must be a pointer to a struct")
	}

	w := reflect.ValueOf(wrapper).Elem()
	for i := 0; i < w.Type().NumField(); i++ {
		if flag, ok := w.Type().Field(i).Tag.Lookup("hehdon"); ok {
			method, ok := method(p, w.Type().Field(i).Name)
			if !ok && flag == "required" {
				return errors.New("method " + w.Type().Field(i).Name + " is required")
			}
			if !ok {
				method = defaultMethod(w.Type().Field(i).Type)
			}
			if !method.Type().AssignableTo(w.Field(i).Type()) {
				return fmt.Errorf("invalid method %s (need %s, found %s)", w.Type().Field(i).Name, w.Field(i).Type(), method.Type())
			}
			w.Field(i).Set(method)
		}
	}

	return nil
}

func method(plugin reflect.Value, name string) (reflect.Value, bool) {
	method := plugin.MethodByName(name)
	if method.Kind() == reflect.Func {
		return method, true
	}

	method = plugin.Elem().FieldByName(name)
	if method.Kind() == reflect.Func {
		return method, true
	}

	return reflect.Value{}, false
}

func defaultMethod(t reflect.Type) reflect.Value {
	return reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
		out := []reflect.Value{}
		for i := 0; i < t.NumOut(); i++ {
			out = append(out, reflect.New(t.Out(i)).Elem())
		}
		return out
	})
}
