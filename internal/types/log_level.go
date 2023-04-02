package types

import (
	"strings"

	"github.com/fatih/color"
)

type LogLevel string

const (
	LogLevelOff   LogLevel = ""
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

func (l LogLevel) precedence() int {
	switch l {
	case LogLevelDebug:
		return 0
	case LogLevelInfo:
		return 1
	case LogLevelWarn:
		return 2
	case LogLevelError:
		return 3
	default:
		return 4
	}
}

func (l LogLevel) Less(o LogLevel) bool {
	return l.precedence() < o.precedence()
}

func (l LogLevel) Color() string {
	v := l.String()
	switch l {
	case LogLevelDebug:
		return color.New(color.Faint).Sprint(v)
	case LogLevelInfo:
		return color.BlueString(v)
	case LogLevelWarn:
		return color.YellowString(v)
	case LogLevelError:
		return color.RedString(v)
	default:
		return v
	}
}

func (l LogLevel) String() string {
	return string(l)
}

func LogLevelOf(s string) LogLevel {
	switch {
	case strings.Contains(s, "DEBUG"):
		return LogLevelDebug
	case strings.Contains(s, "WARN"):
		return LogLevelWarn
	case strings.Contains(s, "ERROR"):
		return LogLevelError
	default:
		return LogLevelInfo
	}
}
