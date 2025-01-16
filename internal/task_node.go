package internal

import (
	"sync"

	"github.com/kitproj/kit/internal/types"
)

type TaskNode struct {
	Name string `json:"name"`
	task types.Task
	// the phase of the task, e.g. "pending", "waiting", "running", "succeeded", "failed"
	Phase string `json:"phase"`
	// the message for the task phase, e.g. "exit code 1'
	Message string `json:"message"`
	// cancel function
	cancel func()
	// a mutex
	mu *sync.Mutex
}

func (n TaskNode) blocked() bool {
	switch n.Phase {
	case "running":
		return n.task.GetType() == types.TaskTypeJob
	case "succeeded":
		return false
	default:
		return true
	}
}
