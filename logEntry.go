package main

import "github.com/fatih/color"

type LogEntry struct {
	level string
	msg   string
}

func (e LogEntry) String() string {
	if e.level == "error" {
		return color.YellowString(e.msg)
	}
	return e.msg
}
