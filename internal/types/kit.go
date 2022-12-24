package types

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/intstr"
)

type Metadata struct {
	Name string `json:"name"`
}

type Kit struct {
	Spec       Spec      `json:"spec"`
	ApiVersion string    `json:"apiVersion,omitempty"`
	Kind       string    `json:"kind,omitempty"`
	Metadata   *Metadata `json:"metadata,omitempty"`
}

type Spec struct {
	TerminationGracePeriodSeconds *int32      `json:"terminationGracePeriodSeconds,omitempty"`
	InitContainers                []Container `json:"initContainers,omitempty"`
	Containers                    []Container `json:"containers,omitempty"`
}

func (s Spec) GetTerminationGracePeriod() time.Duration {
	if s.TerminationGracePeriodSeconds != nil {
		return time.Duration(*s.TerminationGracePeriodSeconds) * time.Second
	}
	return 30 * time.Second
}

type TCPSocketAction struct {
	Port intstr.IntOrString `json:"port"`
}

type HTTPGetAction struct {
	Scheme string              `json:"scheme,omitempty"`
	Port   *intstr.IntOrString `json:"port,omitempty"`
	Path   string              `json:"path,omitempty"`
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

func (a HTTPGetAction) GetPort() int32 {
	if a.Port == nil {
		return 80
	}
	return a.Port.IntVal
}

type Probe struct {
	InitialDelaySeconds int32            `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int32            `json:"periodSeconds,omitempty"`
	TCPSocket           *TCPSocketAction `json:"tcpSocket,omitempty"`
	HTTPGet             *HTTPGetAction   `json:"httpGet,omitempty"`
	SuccessThreshold    int32            `json:"successThreshold,omitempty"`
	FailureThreshold    int32            `json:"failureThreshold,omitempty"`
}

func (p Probe) GetInitialDelay() time.Duration {
	return time.Duration(p.InitialDelaySeconds) * time.Second
}

func (p Probe) GetPeriod() time.Duration {
	if p.PeriodSeconds == 0 {
		return 10 * time.Second
	}
	return time.Duration(p.PeriodSeconds) * time.Second
}

func (p Probe) GetFailureThreshold() int {
	if p.FailureThreshold == 0 {
		return 3
	}
	return int(p.FailureThreshold)
}

func (p Probe) GetSuccessThreshold() int {
	if p.SuccessThreshold == 0 {
		return 1
	}
	return int(p.SuccessThreshold)
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ContainerPort struct {
	ContainerPort int `json:"containerPort,omitempty"`
	HostPort      int `json:"hostPort"`
}

func (p ContainerPort) GetHostPort() int {
	if p.HostPort == 0 {
		return p.ContainerPort
	}
	return p.HostPort
}

type Container struct {
	Name            string          `json:"name"`
	Image           string          `json:"image,omitempty"`
	ImagePullPolicy string          `json:"imagePullPolicy,omitempty"`
	LivenessProbe   *Probe          `json:"livenessProbe,omitempty"`
	ReadinessProbe  *Probe          `json:"readinessProbe,omitempty"`
	Command         []string        `json:"command,omitempty"`
	Args            []string        `json:"args,omitempty"`
	WorkingDir      string          `json:"workingDir,omitempty"`
	Env             []EnvVar        `json:"env,omitempty"`
	Ports           []ContainerPort `json:"ports,omitempty"`
}
