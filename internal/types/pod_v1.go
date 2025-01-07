package types

// Deprecated: use PodSpec instead.
type podV1 struct {
	// The specification of tasks to run.
	Spec PodSpec `json:"spec"`
	// APIVersion must be `kit/v1`.
	// Deprecated: ignored.
	ApiVersion string `json:"apiVersion,omitempty"`
	// Kind must be `Tasks`.
	// Deprecated: ignored.
	Kind string `json:"kind,omitempty"`
	// Metadata is the metadata for the pod.
	// Deprecated: ignored.
	Metadata Metadata `json:"metadata,omitempty"`
}
