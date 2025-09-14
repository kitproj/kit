package types

type Metrics struct {
	CPU uint64 `json:"cpu"` // CPU usage in millicores
	Mem uint64 `json:"mem"` // Memory usage in bytes
}
