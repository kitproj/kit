package types

// VolumeMount describes a mounting of a Volume within a container.
type VolumeMount struct {
	// This must match the name of a volume.
	Name string `json:"name"`
	// Path within the container at which the volume should be mounted.
	MountPath string `json:"mountPath"`
}
