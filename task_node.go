package main

import "github.com/kitproj/kit/internal/types"

type taskNode struct {
	task   types.Task
	status string
}

func (n taskNode) busy() bool {
	return n.status == "running" || n.status == "starting" || n.status == "waiting"
}

func (n taskNode) blocked() bool {
	if n.task.IsService() {
		return n.status != "running"
	} else {
		return n.status != "succeeded" && n.status != "skipped"
	}
}
