package cmd

import (
	"fmt"
	"reflect"
)

func getArgs(in interface{}) (map[string]string, error) {
	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("unexpected type %T, expected struct", in)
	}

	res := map[string]string{}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		name := field.Tag.Get("tmpl")
		if name == "" {
			continue
		}

		res[name] = fmt.Sprint(v.Field(i).Interface())
	}

	return res, nil
}
