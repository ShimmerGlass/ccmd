package cmd

import (
	"bytes"
	"io"
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
	lock            sync.Mutex
}

func (w *consoleWriter) Write(target io.Writer, prefix string, b []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	p := []byte(prefix + ": ")
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

	_, err := target.Write(r)
	return len(b), err
}

type writerAdapter struct {
	prefix string
	inner  *consoleWriter
	target io.Writer
}

func (w *writerAdapter) Write(b []byte) (int, error) {
	return w.inner.Write(w.target, w.prefix, b)
}
