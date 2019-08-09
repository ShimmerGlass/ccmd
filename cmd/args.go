package cmd

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"sort"
	"strings"

	"github.com/fatih/color"
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

func getVarsDesc(args map[string]string, short bool) string {
	if len(args) == 0 {
		return ""
	}

	var prefix string

	if len(args) == 1 && short {
		for _, v := range args {
			prefix = v
			break
		}
	} else {
		parts := []string{}
		for k, v := range args {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}

		sort.Strings(parts)
		prefix = strings.Join(parts, " ")
	}

	h := fnv.New32()
	h.Write([]byte(prefix))
	i := h.Sum32()

	att := wrapperColors[int(i)%len(wrapperColors)]
	d := color.New(att)
	return d.Sprint(prefix)
}
