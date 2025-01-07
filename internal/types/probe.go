package types

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// A probe to check if the task is alive, it will be restarted if not.
type Probe struct {
	// The action to perform.
	TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty"`
	// The action to perform.
	HTTPGet *HTTPGetAction `json:"httpGet,omitempty"`
	// Number of seconds after the process has started before the probe is initiated.
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty"`
	// How often (in seconds) to perform the probe.
	PeriodSeconds int32 `json:"periodSeconds,omitempty"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	SuccessThreshold int32 `json:"successThreshold,omitempty"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
}

func (p *Probe) UnmarshalJSON(data []byte) error {
	if data[0] == '{' {
		x := struct {
			TCPSocket           *TCPSocketAction `json:"tcpSocket,omitempty"`
			HTTPGet             *HTTPGetAction   `json:"httpGet,omitempty"`
			InitialDelaySeconds int32            `json:"initialDelaySeconds,omitempty"`
			PeriodSeconds       int32            `json:"periodSeconds,omitempty"`
			SuccessThreshold    int32            `json:"successThreshold,omitempty"`
			FailureThreshold    int32            `json:"failureThreshold,omitempty"`
		}{}
		if err := json.Unmarshal(data, &x); err != nil {
			return err
		}
		p.TCPSocket = x.TCPSocket
		p.HTTPGet = x.HTTPGet
		p.InitialDelaySeconds = x.InitialDelaySeconds
		p.PeriodSeconds = x.PeriodSeconds
		p.SuccessThreshold = x.SuccessThreshold
		p.FailureThreshold = x.FailureThreshold
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return p.Unstring(s)
}

func (p Probe) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p Probe) String() string {
	return p.URL().String()
}

func (p *Probe) Unstring(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	port := parsePort(u.Port())
	if u.Scheme == "tcp" {
		p.TCPSocket = &TCPSocketAction{Port: port}
	} else {
		p.HTTPGet = &HTTPGetAction{
			Scheme: u.Scheme,
			Port:   port,
			Path:   u.Path,
		}
	}

	q := u.Query()
	successThreshold, _ := strconv.ParseInt(q.Get("successThreshold"), 10, 32)
	p.SuccessThreshold = int32(successThreshold)
	failureThreshold, _ := strconv.ParseInt(q.Get("failureThreshold"), 10, 32)
	p.FailureThreshold = int32(failureThreshold)
	period, _ := time.ParseDuration(q.Get("period"))
	p.PeriodSeconds = int32(period.Seconds())
	initialDelay, _ := time.ParseDuration(q.Get("initialDelay"))
	p.InitialDelaySeconds = int32(initialDelay.Seconds())
	return err
}

func parsePort(s string) uint16 {
	port, _ := strconv.ParseUint(s, 10, 16)
	return uint16(port)
}

func (p Probe) URL() *url.URL {
	var u *url.URL
	if p.TCPSocket != nil {
		u = p.TCPSocket.URL()
	} else {
		u = p.HTTPGet.URL()
	}
	var x = url.Values{}
	x.Add("initialDelay", p.GetInitialDelay().String())
	x.Add("period", p.GetPeriod().String())
	x.Add("successThreshold", fmt.Sprint(p.GetSuccessThreshold()))
	x.Add("failureThreshold", fmt.Sprint(p.GetFailureThreshold()))
	u.RawQuery = x.Encode()
	return u
}

func (p Probe) GetInitialDelay() time.Duration {
	if p.InitialDelaySeconds == 0 {
		return p.GetPeriod()
	}
	return time.Duration(p.InitialDelaySeconds) * time.Second
}

func (p Probe) GetPeriod() time.Duration {
	if p.PeriodSeconds == 0 {
		return 5 * time.Second
	}
	return time.Duration(p.PeriodSeconds) * time.Second
}

func (p Probe) GetFailureThreshold() int {
	if p.FailureThreshold == 0 {
		return 20 // 1m
	}
	return int(p.FailureThreshold)
}

func (p Probe) GetSuccessThreshold() int {
	if p.SuccessThreshold == 0 {
		return 1 // 3s
	}
	return int(p.SuccessThreshold)
}
