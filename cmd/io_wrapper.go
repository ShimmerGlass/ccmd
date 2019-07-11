package cmd

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"sort"
	"strings"

	"github.com/fatih/color"
)

var wrapperColors = []color.Attribute{
	color.FgRed,
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
}

type ioWrapper struct {
	spliced bool
	prefix  string
	inner   io.Writer
}

func (w *ioWrapper) Write(b []byte) (int, error) {
	p := []byte(w.prefix)
	r := []byte{}

	if !w.spliced {
		r = append(r, p...)
	}

	offset := 0
	for {
		if offset >= len(b)-1 {
			break
		}
		nextL := bytes.Index(b[offset:], []byte("\n"))
		if nextL == -1 {
			r = append(r, b[offset:]...)
			w.spliced = true
			break
		}

		r = append(r, p...)
		r = append(r, b[offset:offset+nextL+1]...)
		offset = nextL + 1
	}

	w.inner.Write(r)

	return len(b), nil
}

func newIOWrapper(args map[string]string, used []string, lvl int, target io.Writer) io.Writer {
	sort.Strings(used)
	parts := []string{}
	for _, u := range used {
		v, ok := args[u]
		if !ok {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", u, v))
	}

	if len(parts) == 0 {
		return target
	}

	prefix := strings.Join(parts, ",") + ": "

	h := fnv.New32()
	h.Write([]byte(prefix))
	i := h.Sum32()

	att := wrapperColors[int(i)%len(wrapperColors)]
	d := color.New(att)
	prefix = d.Sprint(prefix)

	return &ioWrapper{
		prefix: prefix,
		inner:  target,
	}
}
