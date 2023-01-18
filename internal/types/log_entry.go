package types

import (
	"fmt"
	"io"
	"strings"
)

type LogEntry struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

func (s *LogEntry) Stdout() io.Writer {
	return writeFunc(func(p []byte) (n int, err error) {
		*s = LogEntry{"info", last(p)}
		return len(p), nil
	})
}

func last(p []byte) string {
	parts := strings.Split(strings.TrimSpace(string(p)), "\n")
	return parts[len(parts)-1]
}

func (s *LogEntry) Stderr() io.Writer {
	return writeFunc(func(p []byte) (n int, err error) {
		*s = LogEntry{"error", last(p)}
		return len(p), nil
	})
}

func (s *LogEntry) IsError() bool {
	return s.Level == "error"
}

func (s *LogEntry) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(s.Stdout(), format, args...)
}
