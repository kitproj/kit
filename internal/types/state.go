package types

import (
	"io"
	"strings"
)

type State struct {
	Phase Phase
	Log   LogEntry
}

type WriteFunc func(p []byte) (n int, err error)

func (w WriteFunc) Write(p []byte) (n int, err error) {
	return w(p)
}

func (s *State) Stdout() io.Writer {
	return WriteFunc(func(p []byte) (n int, err error) {
		s.Log = LogEntry{"info", last(p)}
		return len(p), nil
	})
}

func last(p []byte) string {
	parts := strings.Split(strings.TrimSpace(string(p)), "\n")
	return parts[len(parts)-1]
}

func (s *State) Stderr() io.Writer {
	return WriteFunc(func(p []byte) (n int, err error) {
		s.Log = LogEntry{"error", last(p)}
		return len(p), nil
	})
}
