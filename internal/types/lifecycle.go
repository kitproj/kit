package types

// LifecycleHook defines a command to run at a specific point in the task lifecycle.
type LifecycleHook struct {
	// The command to run.
	Command Strings `json:"command,omitempty"`
	// The shell script to run, instead of command.
	Sh string `json:"sh,omitempty"`
}

// GetCommand returns the command to run, handling both command and sh forms.
func (h *LifecycleHook) GetCommand() Strings {
	if h == nil {
		return nil
	}
	if len(h.Command) > 0 {
		return h.Command
	}
	if h.Sh != "" {
		return []string{"sh", "-c", h.Sh}
	}
	return nil
}

// Lifecycle describes actions that the system should take in response to lifecycle events.
type Lifecycle struct {
	// OnSuccess is the hook to run after the task succeeds.
	OnSuccess *LifecycleHook `json:"onSuccess,omitempty"`
	// OnFailure is the hook to run after the task fails.
	OnFailure *LifecycleHook `json:"onFailure,omitempty"`
}

// GetOnSuccessHook returns the OnSuccess hook, or nil if the Lifecycle is nil.
func (l *Lifecycle) GetOnSuccessHook() *LifecycleHook {
	if l == nil {
		return nil
	}
	return l.OnSuccess
}

// GetOnFailureHook returns the OnFailure hook, or nil if the Lifecycle is nil.
func (l *Lifecycle) GetOnFailureHook() *LifecycleHook {
	if l == nil {
		return nil
	}
	return l.OnFailure
}
