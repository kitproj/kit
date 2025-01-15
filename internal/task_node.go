package internal

import (
	"sync"

	"github.com/kitproj/kit/internal/types"
)

type TaskNode struct {
	name string
	task types.Task
	// the phase of the task, e.g. "pending", "waiting", "running", "succeeded", "failed"
	phase string
	// the message for the task phase, e.g. "exit code 1'
	message string
	// cancel function
	cancel func()
	// a mutex
	mu *sync.Mutex
}

func (n TaskNode) blocked() bool {
	if n.task.GetType() == types.TaskTypeService {
		// skipped services are succeeded
		return n.phase != "running" && n.phase != "succeeded"
	} else {
		return n.phase != "succeeded"
	}
}
