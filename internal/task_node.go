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
	if n.task.GetType() == types.TaskTypeService {
		// skipped services are succeeded
		return n.Phase != "running" && n.Phase != "succeeded"
	} else {
		return n.Phase != "succeeded"
	}
}
