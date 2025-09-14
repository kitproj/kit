package types

type Metrics struct {
	CPU float64 `json:"cpu"` // CPU usage in percentage
	Mem uint64  `json:"mem"` // Memory usage in bytes
}
