package types

// Deprecated: only used for legacy unmarshalling.
type NamedTask struct {
	// The name of the task, must be unique
	Name string `json:"name"`
	Task
}
