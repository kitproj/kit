package types

// Deprecated: use Spec instead.
type workflowV1 struct {
	// The specification of tasks to run.
	Spec Spec `json:"spec"`
	// APIVersion must be `kit/v1`.
	// Deprecated: ignored.
	ApiVersion string `json:"apiVersion,omitempty"`
	// Kind must be `Tasks`.
	// Deprecated: ignored.
	Kind string `json:"kind,omitempty"`
	// Metadata is the metadata for the workflow.
	// Deprecated: ignored.
	Metadata Metadata `json:"metadata,omitempty"`
}
