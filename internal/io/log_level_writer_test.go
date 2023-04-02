package io

import (
	"bytes"
	"testing"

	"github.com/kitproj/kit/internal/types"
	"github.com/stretchr/testify/assert"
)

func Test_logLevelWriter_Write(t *testing.T) {
	writer := &bytes.Buffer{}
	w := NewLogLevelWriter(types.LogLevelInfo, writer)
	_, err := w.Write([]byte("DEBUG: hello world\n"))
	assert.NoError(t, err)
	_, err = w.Write([]byte("INFO: hello world\n"))
	assert.NoError(t, err)
	assert.Equal(t, "INFO: hello world\n", writer.String())
}
