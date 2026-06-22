package types

// Lifecycle describes actions that the system should take in response to lifecycle events.
type Lifecycle struct {
	// OnSuccess is the name of the task to run after the task/graph succeeds.
	OnSuccess string `json:"onSuccess,omitempty"`
	// OnFailure is the name of the task to run after the task/graph fails.
	OnFailure string `json:"onFailure,omitempty"`
}

// GetOnSuccess returns the OnSuccess task name, or empty string if the Lifecycle is nil.
func (l *Lifecycle) GetOnSuccess() string {
	if l == nil {
		return ""
	}
	return l.OnSuccess
}

// GetOnFailure returns the OnFailure task name, or empty string if the Lifecycle is nil.
func (l *Lifecycle) GetOnFailure() string {
	if l == nil {
		return ""
	}
	return l.OnFailure
}
