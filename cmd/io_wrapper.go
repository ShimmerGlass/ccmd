package cmd

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
)

func init() {
	color.NoColor = false
}

var wrapperColors = []color.Attribute{
	color.FgRed,
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
}

type consoleWriter struct {
	spliced         bool
	splicedPrefix   string
	maxPrefixLength int
	inner           io.Writer
	lock            sync.Mutex
}

func (w *consoleWriter) Write(prefix string, b []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	p := []byte(prefix)
	r := []byte{}

	if w.spliced && w.splicedPrefix != prefix {
		r = append(r, '\n')
	}
	if !w.spliced || w.splicedPrefix != prefix {
		r = append(r, p...)
		if len(p) <= w.maxPrefixLength {
			r = append(r, bytes.Repeat([]byte{' '}, w.maxPrefixLength-len(p))...)
		}
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

		r = append(r, b[offset:offset+nextL+1]...)

		if offset+nextL < len(b)-1 {
			r = append(r, p...)
		}

		offset = offset + nextL + 1
	}

	w.spliced = b[len(b)-1] != '\n'
	w.splicedPrefix = prefix

	_, err := w.inner.Write(r)
	return len(b), err
}

type writerAdapter struct {
	prefix string
	inner  *consoleWriter
}

func (w *writerAdapter) Write(b []byte) (int, error) {
	return w.inner.Write(w.prefix, b)
}

func getPrefix(args map[string]string) string {
	if len(args) == 0 {
		return ""
	}

	var prefix string

	if len(args) == 1 {
		for _, v := range args {
			prefix = v + ": "
			break
		}
	} else {
		parts := []string{}
		for k, v := range args {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}

		sort.Strings(parts)
		prefix = strings.Join(parts, " ") + ": "
	}

	h := fnv.New32()
	h.Write([]byte(prefix))
	i := h.Sum32()

	att := wrapperColors[int(i)%len(wrapperColors)]
	d := color.New(att)
	return d.Sprint(prefix)
}
