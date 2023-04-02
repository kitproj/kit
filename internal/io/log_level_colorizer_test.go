package io

import (
	"bytes"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func Test_logColorizer_Write(t *testing.T) {
	writer := &bytes.Buffer{}
	w := NewLogColorizer(writer)
	_, err := w.Write([]byte("INFO: hello world\n"))
	assert.NoError(t, err)
	assert.Equal(t, color.BlueString("INFO")+": hello world\n", writer.String())
}
