package types

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// HTTPGetAction describes an action based on HTTP Locks requests.
type HTTPGetAction struct {
	// Scheme to use for connecting to the host. Defaults to HTTP.
	Scheme string `json:"scheme,omitempty"`
	// Number of the port
	Port uint16 `json:"port,omitempty"`
	// Path to access on the HTTP server.
	Path string `json:"path,omitempty"`
}

func (a HTTPGetAction) URL() *url.URL {
	return &url.URL{Scheme: a.GetProto(), Host: fmt.Sprintf("localhost:%v", a.Port), Path: a.Path}
}

func (a *HTTPGetAction) Unstring(s string) error {
	x, err := url.Parse(s)
	if err != nil {
		return err
	}
	a.Scheme = x.Scheme
	port, _ := strconv.ParseUint(x.Port(), 10, 16)
	a.Port = uint16(port)
	a.Path = x.Path
	return nil
}

func (a HTTPGetAction) GetProto() string {
	if a.Scheme == "" {
		return "http"
	}
	return strings.ToLower(a.Scheme)
}

func (a HTTPGetAction) GetURL() string {
	return fmt.Sprintf("%s://localhost:%v%s", a.GetProto(), a.GetPort(), a.Path)
}

func (a HTTPGetAction) GetPort() uint16 {
	if a.Port > 0 {
		return a.Port
	}
	if a.GetProto() == "https" {
		return 443
	}
	return 80
}
