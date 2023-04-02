package io

import (
	"bytes"
	"io"
	"sync"

	"github.com/kitproj/kit/internal/types"
)

// logColorizer is a writer that colors log levels.
type logColorizer struct {
	mu     *sync.Mutex
	buffer *bytes.Buffer
	writer io.Writer
}

func NewLogColorizer(writer io.Writer) io.Writer {
	return &logColorizer{
		mu:     &sync.Mutex{},
		buffer: &bytes.Buffer{},
		writer: writer,
	}
}

func (w *logColorizer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, c := range p {
		next := c < 'A' || c > 'Z'
		if next {
			s := w.buffer.String()
			level := types.LogLevel(s)
			if _, err := w.writer.Write([]byte(level.Color())); err != nil {
				return 0, err

			}
			w.buffer.Reset()
		}
		if _, err := w.buffer.Write([]byte{c}); err != nil {
			return 0, err
		}
		if next {
			if _, err := w.writer.Write(w.buffer.Bytes()); err != nil {
				return 0, err
			}
			w.buffer.Reset()
		}
	}
	// we must return the original length of p, not the length of the output, because we'll get broken pipes if we don't (SIGPIPE)
	return len(p), nil
}
