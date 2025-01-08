package types

import (
	"fmt"
	"net/url"
)

// TCPSocketAction describes an action based on opening a socket
type TCPSocketAction struct {
	// Port number of the port to probe.
	Port uint16 `json:"port"`
}

func (a TCPSocketAction) URL() *url.URL {
	return &url.URL{Scheme: "tcp", Host: fmt.Sprintf("localhost:%v", a.Port)}
}
