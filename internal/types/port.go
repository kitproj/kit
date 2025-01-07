package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// A port to expose.
type Port struct {
	// The container port to expose
	ContainerPort uint16 `json:"containerPort,omitempty"`
	// The host port to route to the container port
	HostPort uint16 `json:"hostPort,omitempty"`
}

func (p *Port) UnmarshalJSON(data []byte) error {
	if data[0] == '{' {
		var x struct {
			ContainerPort uint16 `json:"containerPort"`
			HostPort      uint16 `json:"hostPort"`
		}
		if err := json.Unmarshal(data, &x); err != nil {
			return err
		}
		p.ContainerPort = x.ContainerPort
		p.HostPort = x.HostPort
		return nil
	}
	var x string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	return p.Unstring(x)
}

func (p Port) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Port) Unstring(s string) error {
	parts := strings.Split(s, ":")
	containerPort, err := strconv.ParseUint(parts[0], 10, 16)
	p.ContainerPort = uint16(containerPort)
	switch len(parts) {
	case 1:
		p.HostPort = p.ContainerPort
		return err
	case 2:
		hostPort, err := strconv.ParseUint(parts[1], 10, 16)
		p.HostPort = uint16(hostPort)
		return err
	default:
		return fmt.Errorf("invalid port string %q", s)
	}
}

func (p Port) String() string {
	if p.GetHostPort() == p.ContainerPort {
		return fmt.Sprint(p.ContainerPort)
	}
	return fmt.Sprintf("%d:%d", p.ContainerPort, p.GetHostPort())
}

func (p Port) GetHostPort() uint16 {
	if p.HostPort == 0 {
		return p.ContainerPort
	}
	return p.HostPort
}
