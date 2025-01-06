package main

import "github.com/kitproj/kit/internal/types"

type taskNode struct {
	name string
	task types.Task
	// the phase of the task, e.g. "waiting", "running", "succeeded", "failed"
	phase string
	// the message for the task phase, e.g. "exit code 1'
	message string
	// cancel function
	cancel func()
}

func (n taskNode) blocked() bool {
	if n.task.IsService() {
		return n.phase != "running"
	} else {
		return n.phase != "succeeded"
	}
}
