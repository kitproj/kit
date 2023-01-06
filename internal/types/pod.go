package types

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/intstr"
)

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Port struct {
	ContainerPort uint16 `json:"containerPort,omitempty"`
	HostPort      uint16 `json:"hostPort"`
}

func (p Port) GetHostPort() uint16 {
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

func (t *Task) HasMutex() bool {
	return t != nil && t.Mutex != ""
}

type Task struct {
	Name            string        `json:"name"`
	Image           string        `json:"image,omitempty"`
	ImagePullPolicy string        `json:"imagePullPolicy,omitempty"`
	LivenessProbe   *Probe        `json:"livenessProbe,omitempty"`
	ReadinessProbe  *Probe        `json:"readinessProbe,omitempty"`
	Command         []string      `json:"command,omitempty"`
	Args            []string      `json:"args,omitempty"`
	WorkingDir      string        `json:"workingDir,omitempty"`
	Env             EnvVars       `json:"env,omitempty"`
	Ports           []Port        `json:"ports,omitempty"`
	VolumeMounts    []VolumeMount `json:"volumeMounts,omitempty"`
	TTY             bool          `json:"tty,omitempty"`
	Watch           []string      `json:"watch,omitempty"`
	Mutex           string        `json:"mutex,omitempty"`
	Dependencies    []string      `json:"dependencies,omitempty"`
}

func (t *Task) IsBackground() bool {
	return t.ReadinessProbe != nil && t.LivenessProbe != nil
}

func (t Task) GetHostPorts() []uint16 {
	var ports []uint16
	for _, p := range t.Ports {
		ports = append(ports, p.GetHostPort())
	}
	return ports
}

type Pod struct {
	Spec       PodSpec  `json:"spec"`
	ApiVersion string   `json:"apiVersion,omitempty"`
	Kind       string   `json:"kind,omitempty"`
	Metadata   Metadata `json:"metadata"`
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
	return fmt.Sprintf("%s://localhost:%s%s", a.GetProto(), a.GetPort(), a.Path)
}

func (a HTTPGetAction) GetPort() string {
	if a.Port == nil {
		if a.GetProto() == "http" {
			return "80"
		} else {
			return "443"
		}
	}
	return a.Port.String()
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

type Tasks []Task

func (t Tasks) GetLeaves() Tasks {
	var out Tasks
	for _, t := range t {
		if len(t.Dependencies) == 0 {
			out = append(out, t)
		}
	}
	return out
}

func (t Tasks) GetDownstream(name string) Tasks {
	var out Tasks
	for _, downstream := range t {
		for _, upstream := range downstream.Dependencies {
			if upstream == name {
				out = append(out, downstream)
			}
		}
	}
	return out
}

func (t Tasks) NeededFor(names []string) Tasks {
	todo := make(chan string, len(t))
	for _, name := range names {
		todo <- name
	}
	found := map[string]bool{}
	for name := range todo {
		found[name] = true
		for _, d := range t.Get(name).Dependencies {
			todo <- d
		}
		if len(todo) == 0 {
			close(todo)
		}
	}
	var out Tasks
	for name := range found {
		out = append(out, t.Get(name))
	}
	return out
}

func (t Tasks) Get(name string) Task {
	for _, task := range t {
		if task.Name == name {
			return task
		}

	}
	panic(fmt.Errorf("not task named %q", name))
}

type PodSpec struct {
	TerminationGracePeriodSeconds *int32   `json:"terminationGracePeriodSeconds,omitempty"`
	Tasks                         Tasks    `json:"tasks,omitempty"`
	Volumes                       []Volume `json:"volumes,omitempty"`
}

func (s PodSpec) GetTerminationGracePeriod() time.Duration {
	if s.TerminationGracePeriodSeconds != nil {
		return time.Duration(*s.TerminationGracePeriodSeconds) * time.Second
	}
	return 30 * time.Second
}

type Status struct {
	TaskStatuses []*TaskStatus
}

type TaskStateWaiting struct {
	Reason string
}

type TaskStateRunning struct {
}

type TaskStateTerminated struct {
	Reason string
}

type TaskState struct {
	Waiting    *TaskStateWaiting
	Running    *TaskStateRunning
	Terminated *TaskStateTerminated
}

func (s TaskStatus) GetReason() string {
	if s.State.Waiting != nil {
		return s.State.Waiting.Reason
	} else if s.State.Running != nil {
		if s.Ready {
			return "ready"
		} else {
			return "running"
		}
	} else if s.State.Terminated != nil {
		return s.State.Terminated.Reason
	}
	return "unknown"
}

func (s *TaskStatus) IsSuccess() bool {
	return s.IsTerminated() && s.State.Terminated.Reason == "success"
}

func (s TaskStatus) Failed() bool {
	return s.IsTerminated() && !s.IsSuccess()
}

func (s *TaskStatus) IsTerminated() bool {
	return s != nil && s.State.Terminated != nil
}

func (s *TaskStatus) IsReady() bool {
	return s != nil && s.Ready
}

func (s TaskStatus) IsFulfilled() bool {
	return s.IsSuccess() || s.IsReady()
}

type TaskStatus struct {
	Name  string
	Ready bool
	State TaskState
}

func (s *Status) GetStatus(name string) *TaskStatus {
	for i, x := range s.TaskStatuses {
		if x.Name == name {
			return s.TaskStatuses[i]
		}
	}
	return nil
}
