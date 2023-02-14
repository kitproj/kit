package types

type Metadata struct {
	// Name is the name of the resource.
	Name string `json:"name"`
	// Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata.
	Annotations map[string]string `json:"annotations,omitempty"`
}
