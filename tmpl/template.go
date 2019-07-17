package tmpl

import (
	"fmt"
	"strings"
)

type textElem string
type varElem string

type Template struct {
	Vars  []string
	elems []interface{}
}

func Parse(in string) (*Template, error) {
	t := &Template{}

	for {
		nextO := strings.IndexByte(in, '{')
		nextC := strings.IndexByte(in, '}')

		// end escape
		if nextC != -1 && (nextC < nextO || nextO == -1) {
			if nextC == len(in)-1 {
				return nil, fmt.Errorf("unterminated expression at %d", nextO)
			}
			if in[nextC+1] == '}' {
				t.elems = append(t.elems, textElem(in[:nextC+1]))
				in = in[nextC+2:]
				continue
			}
			return nil, fmt.Errorf("unexpected '}' at %d", nextC)
		}

		// end
		if nextO == -1 {
			t.elems = append(t.elems, textElem(in))
			break
		}

		if nextO == len(in)-1 {
			return nil, fmt.Errorf("unterminated expression at %d", nextO)
		}

		// start escape
		if in[nextO+1] == '{' {
			t.elems = append(t.elems, textElem(in[:nextO+1]))
			in = in[nextO+2:]
			continue
		}

		// var
		if nextC == -1 {
			return nil, fmt.Errorf("unterminated expression at %d", nextO)
		}

		if nextO > 0 {
			t.elems = append(t.elems, textElem(in[:nextO]))
		}
		v := in[nextO+1 : nextC]
		t.Vars = append(t.Vars, v)
		t.elems = append(t.elems, varElem(v))
		in = in[nextC+1:]
	}

	return t, nil
}

func (t *Template) Exec(args map[string]string) string {
	res := ""

	for _, e := range t.elems {
		switch v := e.(type) {
		case textElem:
			res += string(v)
		case varElem:
			val, ok := args[string(v)]
			if !ok {
				res += "{" + string(v) + "}"
			} else {
				res += fmt.Sprint(val)
			}
		default:
			panic(fmt.Sprintf("bad elem %T", e))
		}
	}

	return res
}

func Exec(tmpl string, args map[string]string) (string, error) {
	t, err := Parse(tmpl)
	if err != nil {
		return "", err
	}

	return t.Exec(args), nil
}
