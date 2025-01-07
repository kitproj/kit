package internal

import "github.com/kitproj/kit/internal/types"

type TaskNode struct {
	Name string
	Task types.Task
	// the phase of the task, e.g. "waiting", "running", "succeeded", "failed"
	Phase string
	// the message for the task phase, e.g. "exit code 1'
	message string
	// cancel function
	Cancel func()
}

func (n TaskNode) blocked() bool {
	if n.Task.IsService() {
		return n.Phase != "running"
	} else {
		return n.Phase != "succeeded"
	}
}
