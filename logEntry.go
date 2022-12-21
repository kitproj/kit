package main

import "github.com/fatih/color"

type logEntry struct {
	level string
	msg   string
}

func (e logEntry) String() string {
	if e.level == "error" {
		return color.YellowString(e.msg)
	}
	return e.msg
}
