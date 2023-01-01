package types

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/util/intstr"
)

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

type EnvVars []EnvVar

func (v EnvVars) Environ() []string {
	var environ []string
	for _, env := range v {
		environ = append(environ, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}
	return environ
}

type Build struct {
	Command    []string `json:"command,omitempty"`
	Args       []string `json:"args,omitempty"`
	WorkingDir string   `json:"workingDir,omitempty"`
	Env        EnvVars  `json:"env,omitempty"`
	Watch      []string `json:"watch,omitempty"`
	Mutex      string   `json:"mutex,omitempty"`
}

func (b *Build) DeepCopy() *Build {
	if b == nil {
		return nil
	}
	return &Build{
		Command:    b.Command,
		Args:       b.Args,
		WorkingDir: b.WorkingDir,
		Env:        b.Env,
		Watch:      b.Watch,
		Mutex:      b.Mutex,
	}
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
	Env             EnvVars         `json:"env,omitempty"`
	Ports           []ContainerPort `json:"ports,omitempty"`
	VolumeMounts    []VolumeMount   `json:"volumeMounts,omitempty"`
	TTY             bool            `json:"tty,omitempty"`
	Build           *Build          `json:"build,omitempty"`
}

type Pod struct {
	Spec       Spec      `json:"spec"`
	ApiVersion string    `json:"apiVersion,omitempty"`
	Kind       string    `json:"kind,omitempty"`
	Metadata   *Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	Name string `json:"name"`
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

func (p *Probe) DeepCopy() *Probe {
	if p == nil {
		return nil
	}
	return &Probe{
		InitialDelaySeconds: p.InitialDelaySeconds,
		PeriodSeconds:       p.PeriodSeconds,
		TCPSocket:           p.TCPSocket.DeepCopy(),
		HTTPGet:             p.HTTPGet.DeepCopy(),
		SuccessThreshold:    p.SuccessThreshold,
		FailureThreshold:    p.FailureThreshold,
	}
}

type TCPSocketAction struct {
	Port intstr.IntOrString `json:"port"`
}

func (a *TCPSocketAction) DeepCopy() *TCPSocketAction {
	if a == nil {
		return nil
	}
	return &TCPSocketAction{Port: a.Port}
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

func (a *HTTPGetAction) DeepCopy() *HTTPGetAction {
	if a == nil {
		return nil
	}
	return &HTTPGetAction{
		Scheme: a.Scheme,
		Port:   a.Port,
		Path:   a.Path,
	}

}

type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

type HostPath struct {
	Path string `json:"path"`
}

type Volume struct {
	Name     string   `json:"name"`
	HostPath HostPath `json:"hostPath"`
}

type Spec struct {
	TerminationGracePeriodSeconds *int32      `json:"terminationGracePeriodSeconds,omitempty"`
	InitContainers                []Container `json:"initContainers,omitempty"`
	Containers                    []Container `json:"containers,omitempty"`
	Volumes                       []Volume    `json:"volumes,omitempty"`
}

func (s Spec) GetTerminationGracePeriod() time.Duration {
	if s.TerminationGracePeriodSeconds != nil {
		return time.Duration(*s.TerminationGracePeriodSeconds) * time.Second
	}
	return 30 * time.Second
}

type Status corev1.PodStatus

func (s *Status) GetContainerStatus(name string) *corev1.ContainerStatus {
	for i, x := range s.InitContainerStatuses {
		if x.Name == name {
			return &s.InitContainerStatuses[i]
		}
	}
	for i, x := range s.ContainerStatuses {
		if x.Name == name {
			return &s.ContainerStatuses[i]
		}
	}
	return nil
}
