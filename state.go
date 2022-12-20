package main

import (
	"io"
	"strings"
)

type state struct {
	phase phase
	log   logEntry
	msg   string
	kill  func() error
}

type WriteFunc func(p []byte) (n int, err error)

func (w WriteFunc) Write(p []byte) (n int, err error) {
	return w(p)
}

func (s *state) Stdout() io.Writer {
	return WriteFunc(func(p []byte) (n int, err error) {
		s.log = logEntry{"info", last(p)}
		return len(p), nil
	})
}

func last(p []byte) string {
	parts := strings.Split(strings.TrimSpace(string(p)), "\n")
	return parts[len(parts)-1]
}

func (s *state) Stderr() io.Writer {
	return WriteFunc(func(p []byte) (n int, err error) {
		s.log = logEntry{"error", last(p)}
		return len(p), nil
	})
}

var states = map[string]*state{}
