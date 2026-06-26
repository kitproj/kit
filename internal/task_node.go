package internal

import (
	"encoding/json"
	"sync"

	"github.com/kitproj/kit/internal/types"
)

type TaskNode struct {
	Name string     `json:"name"`
	Task types.Task `json:"task"`
	// logFile is the log file path
	logFile string
	// Phase represents the current execution state of the task:
	// - "pending": Task is registered but waiting for dependencies to complete
	// - "waiting": Task's dependencies are satisfied and is ready to execute
	// - "starting": Task is initializing but not yet fully running
	// - "running": Task is actively executing
	// - "stalled": Task is running but not making progress
	// - "succeeded": Task completed successfully
	// - "failed": Task completed with errors
	// - "cancelled": Task was manually stopped
	// - "skipped": Task was intentionally not executed
	Phase string `json:"phase"`
	// the message for the task phase, e.g. "exit code 1'
	Message string `json:"message,omitempty"`
	// metrics for resource usage tracking
	Metrics *types.Metrics `json:"metrics,omitempty"`
	// cancel function
	cancel func()
	// mu serializes execution so two instances of a task don't run at once
	mu *sync.Mutex
	// status guards the Phase, Message, Metrics and cancel fields, which are
	// written by task/timer/metrics goroutines and read by the main loop and the server
	status sync.RWMutex
}

func (n *TaskNode) setStatus(phase, message string) {
	n.status.Lock()
	n.Phase = phase
	n.Message = message
	n.status.Unlock()
}

func (n *TaskNode) setMetrics(m *types.Metrics) {
	n.status.Lock()
	n.Metrics = m
	n.status.Unlock()
}

func (n *TaskNode) getPhase() string {
	n.status.RLock()
	defer n.status.RUnlock()
	return n.Phase
}

func (n *TaskNode) getMessage() string {
	n.status.RLock()
	defer n.status.RUnlock()
	return n.Message
}

func (n *TaskNode) setCancel(cancel func()) {
	n.status.Lock()
	n.cancel = cancel
	n.status.Unlock()
}

func (n *TaskNode) doCancel() {
	n.status.RLock()
	cancel := n.cancel
	n.status.RUnlock()
	cancel()
}

// MarshalJSON takes the status lock so the phase/message/metrics fields are
// read consistently while task goroutines may be writing them.
func (n *TaskNode) MarshalJSON() ([]byte, error) {
	type alias TaskNode
	n.status.RLock()
	defer n.status.RUnlock()
	return json.Marshal((*alias)(n))
}

func (n *TaskNode) blocked() bool {
	switch n.getPhase() {
	case "running", "stalled":
		return n.Task.GetType() == types.TaskTypeJob
	case "succeeded", "skipped":
		return false
	default:
		return true
	}
}
