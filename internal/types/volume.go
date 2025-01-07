package types

type Volume struct {
	// Volume's name.
	Name string `json:"name"`
	// HostPath represents a pre-existing file or directory on the host machine that is directly exposed to the container.
	HostPath HostPath `json:"hostPath"`
}
