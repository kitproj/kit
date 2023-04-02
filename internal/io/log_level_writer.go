package io

import (
	"bytes"
	"io"
	"sync"

	"github.com/kitproj/kit/internal/types"
)

// LogLevelWriter is a wrapper of io.Writer that only writes logs with minimuw level
type LogLevelWriter interface {
	io.Writer
	SetLogLevel(level types.LogLevel)
	GetLogLevel() types.LogLevel
}
type logLevelWriter struct {
	mu       *sync.Mutex
	LogLevel types.LogLevel
	buffer   *bytes.Buffer
	writer   io.Writer
}

func NewLogLevelWriter(logLevel types.LogLevel, writer io.Writer) LogLevelWriter {
	return &logLevelWriter{
		mu:       &sync.Mutex{},
		LogLevel: logLevel,
		buffer:   &bytes.Buffer{},
		writer:   writer,
	}
}
func (w *logLevelWriter) SetLogLevel(level types.LogLevel) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.LogLevel = level
}
func (w *logLevelWriter) GetLogLevel() types.LogLevel {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.LogLevel
}

func (w *logLevelWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, c := range p {
		if _, err := w.buffer.Write([]byte{c}); err != nil {
			return 0, err
		}
		if c == '\n' {
			s := w.buffer.String()
			level := types.LogLevelOf(s)
			if !level.Less(w.LogLevel) {
				if _, err := w.writer.Write(w.buffer.Bytes()); err != nil {
					return 0, err
				}
			}
			w.buffer.Reset()
		}
	}
	// we must return the original length of p, not the length of the output, because we'll get broken pipes if we don't (SIGPIPE)
	return len(p), nil
}
