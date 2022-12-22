package types

import "github.com/fatih/color"

type LogEntry struct {
	Level string
	Msg   string
}

func (e LogEntry) String() string {
	if e.Level == "error" {
		return color.YellowString(e.Msg)
	}
	return e.Msg
}
