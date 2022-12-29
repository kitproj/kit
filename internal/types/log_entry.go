package types

import (
	"github.com/fatih/color"
	"io"
	"strings"
)

type LogEntry struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

func (e *LogEntry) String() string {
	if e.Level == "error" {
		return color.YellowString(e.Msg)
	}
	return e.Msg
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
