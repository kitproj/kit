package util

import (
	"strings"
	"sync"
)

type LastNLinesWriter struct {
	mu    sync.RWMutex
	size  int
	lines []string
}

func NewLastNLinesWriter(size int) *LastNLinesWriter {
	return &LastNLinesWriter{
		mu:    sync.RWMutex{},
		size:  size,
		lines: make([]string, 0, size),
	}
}

func (w *LastNLinesWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	// TODO: this will not work well for large p, as it essentially assume that we will get a full line, which is not always the case
	lines := strings.TrimSpace(string(p))
	for _, line := range strings.Split(lines, "\n") {
		w.lines = append(w.lines, line)
	}
	for len(w.lines) > w.size {
		w.lines = w.lines[1:]
	}
	return len(p), nil
}

func (w *LastNLinesWriter) Lines() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.lines
}
