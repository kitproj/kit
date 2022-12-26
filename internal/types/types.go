package types

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
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

type ContainerStatus struct {
	Name  string   `json:"name,omitempty"`
	Phase Phase    `json:"phase,omitempty"`
	Log   LogEntry `json:"log"`
}

type WriteFunc func(p []byte) (n int, err error)

func (w WriteFunc) Write(p []byte) (n int, err error) {
	return w(p)
}

func (s *ContainerStatus) Stdout() io.Writer {
	return WriteFunc(func(p []byte) (n int, err error) {
		s.Log = LogEntry{"info", last(p)}
		return len(p), nil
	})
}

func last(p []byte) string {
	parts := strings.Split(strings.TrimSpace(string(p)), "\n")
	return parts[len(parts)-1]
}

func (s *ContainerStatus) Stderr() io.Writer {
	return WriteFunc(func(p []byte) (n int, err error) {
		s.Log = LogEntry{"error", last(p)}
		return len(p), nil
	})
}

type Kit struct {
	Spec       Spec      `json:"spec"`
	ApiVersion string    `json:"apiVersion,omitempty"`
	Kind       string    `json:"kind,omitempty"`
	Metadata   *Metadata `json:"metadata,omitempty"`
	Status     *Status   `json:"status,omitempty"`
}

func (k Kit) GetContainers() []Container {
	if len(k.Status.ContainerStatuses) > 0 {
		return k.Spec.Containers
	}
	return k.Spec.InitContainers
}

type ContainerStatuses []*ContainerStatus

func (s ContainerStatuses) Get(name string) *ContainerStatus {
	for _, status := range s {
		if status.Name == name {
			return status
		}
	}
	panic(name)
}

func (k Kit) GetContainerStatuses() ContainerStatuses {
	if len(k.Status.ContainerStatuses) > 0 {
		return k.Status.ContainerStatuses
	}
	return k.Status.InitContainerStatuses
}

type LogEntry struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

func (e LogEntry) String() string {
	if e.Level == "error" {
		return color.YellowString(e.Msg)
	}
	return e.Msg
}

type Metadata struct {
	Name string `json:"name"`
}

type Phase string

const (
	CreatingPhase Phase = "creating"
	ExcludedPhase Phase = "excluded"
	BuildingPhase Phase = "building"
	RunningPhase  Phase = "running"
	LivePhase     Phase = "live"
	DeadPhase     Phase = "dead"
	ReadyPhase    Phase = "ready"
	UnreadyPhase  Phase = "unready"
	ExitedPhase   Phase = "exited"
	ErrorPhase    Phase = "error"
)

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

type Status struct {
	InitContainerStatuses []*ContainerStatus `json:"initContainerStatuses,omitempty"`
	ContainerStatuses     []*ContainerStatus `json:"containerStatuses,omitempty"`
}
