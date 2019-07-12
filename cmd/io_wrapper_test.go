package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	b := bytes.NewBuffer(nil)
	w := &ioWrapper{prefix: "PREFIX", inner: b}

	w.Write([]byte("hello\nbye\n"))

	require.Equal(t, "PREFIXhello\nPREFIXbye\n", b.String())
}

func TestDouble(t *testing.T) {
	b := bytes.NewBuffer(nil)
	w := &ioWrapper{prefix: "PRE", inner: b}
	w2 := &ioWrapper{prefix: "FIX", inner: w}

	w2.Write([]byte("hello\nbye\n"))

	require.Equal(t, "PREFIXhello\nPREFIXbye\n", b.String())
}
