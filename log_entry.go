package main

import "github.com/fatih/color"

type logEntry struct {
	level string
	msg   string
}

func (e logEntry) String() string {
	msg := e.msg
	if len(msg) > 100 {
		msg = msg[0:99]
	}
	if e.level == "error" {
		return color.YellowString(msg)
	}
	return msg
}
