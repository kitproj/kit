package io

import (
	"bytes"
	"io"
	"sync"
)

// prefixWriter is an io.Writer that prefixes each line with a string
type prefixWriter struct {
	mu     *sync.Mutex
	prefix string
	writer io.Writer
	buffer *bytes.Buffer
}

func NewPrefixWriter(prefix string, buffer io.Writer) io.Writer {
	return &prefixWriter{
		mu:     &sync.Mutex{},
		prefix: prefix,
		writer: buffer,
		buffer: bytes.NewBufferString(prefix),
	}
}

func (w *prefixWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, c := range p {
		if _, err := w.buffer.Write([]byte{c}); err != nil {
			return 0, err
		}
		if c == '\n' {
			if _, err := w.writer.Write(w.buffer.Bytes()); err != nil {
				return 0, err
			}
			w.buffer.Reset()
			if _, err := w.buffer.Write([]byte(w.prefix)); err != nil {
				return 0, err
			}
		}
	}
	// we must return the original length of p, not the length of the output, because we'll get broken pipes if we don't (SIGPIPE)
	return len(p), nil
}
