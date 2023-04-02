package io

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_prefixWriter_Write(t *testing.T) {
	buffer := &bytes.Buffer{}
	writer := NewPrefixWriter("foo: ", buffer)
	_, err := writer.Write([]byte("hello world\n"))
	assert.NoError(t, err)
	assert.Equal(t, "foo: hello world\n", buffer.String())
}
