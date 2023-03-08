package types

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// A environment variable.
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (v EnvVar) String() string {
	return fmt.Sprintf("%s=%s", v.Name, v.Value)
}

func (v *EnvVar) Unstring(s string) error {
	parts := strings.Split(s, "=")
	switch len(parts) {
	case 2:
		v.Name = parts[0]
		v.Value = parts[1]
		return nil
	default:
		return fmt.Errorf("invalid EnvVar string %q", s)
	}
}

func (v *EnvVar) UnmarshalJSON(data []byte) error {
	if data[0] == '{' {
		var x struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}
		if err := json.Unmarshal(data, &x); err != nil {
			return err
		}
		v.Name = x.Name
		v.Value = x.Value
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return v.Unstring(s)
}

func (v EnvVar) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// A list of ports to expose.
type Ports []Port

func (p *Ports) UnmarshalJSON(data []byte) error {
	if data[0] == '[' {
		var x []Port
		if err := json.Unmarshal(data, &x); err != nil {
			return err
		}
		for _, port := range x {
			*p = append(*p, port)
		}
		return nil
	}
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*p = append(*p, Port{ContainerPort: uint16(i)})
		return nil
	}
	var x = Strings{}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	for _, port := range x {
		y := Port{}
		if err := y.Unstring(port); err != nil {
			return err
		}
		*p = append(*p, y)
	}

	return nil
}

func (p Ports) MarshalJSON() ([]byte, error) {
	var x Strings
	for _, port := range p {
		x = append(x, port.String())
	}
	return json.Marshal(x)
}

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
	switch len(parts) {
	case 1:
		containerPort, err := strconv.ParseUint(s, 10, 16)
		p.ContainerPort = uint16(containerPort)
		return err
	case 2:
		containerPort, err := strconv.ParseUint(parts[0], 10, 16)
		p.ContainerPort = uint16(containerPort)
		hostPort, err := strconv.ParseUint(parts[1], 10, 16)
		p.HostPort = uint16(hostPort)
		return err
	default:
		return fmt.Errorf("invalid port string %q", s)
	}
}

func (p Port) String() string {
	if p.HostPort == 0 {
		return fmt.Sprint(p.ContainerPort)
	}
	if p.ContainerPort == 0 {
		return fmt.Sprint(p.HostPort)
	}
	return fmt.Sprintf("%d:%d", p.ContainerPort, p.HostPort)
}

func (p Port) GetHostPort() uint16 {
	if p.HostPort == 0 {
		return p.ContainerPort
	}
	return p.HostPort
}

// A list of environment variables.
type EnvVars []EnvVar

// Environ returns a list of environment variables. If an environment variable is defined in both the task and the host, the host value is used.
func (v EnvVars) Environ() []string {
	osEnviron := make(map[string]string)
	for _, env := range v {
		osEnviron[env.Name] = env.Value
	}
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		n, v := parts[0], parts[1]
		if osEnviron[n] != "" {
			osEnviron[n] = v
		}
	}
	var environ []string
	for k, v := range osEnviron {
		environ = append(environ, fmt.Sprintf("%s=%s", k, v))
	}
	return environ
}

func (t *Task) HasMutex() bool {
	return t != nil && t.Mutex != ""
}

// A task is a container or a command to run.
type Task struct {
	// The name of the task, must be unique
	Name string `json:"name"`
	// Either the container image to run, or a directory containing a Dockerfile. If omitted, the process runs on the host.
	Image string `json:"image,omitempty"`
	// Pull policy, e.g. Always, Never, IfNotPresent
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
	// A probe to check if the task is alive, it will be restarted if not. If omitted, the task is assumed to be alive.
	LivenessProbe *Probe `json:"livenessProbe,omitempty"`
	// A probe to check if the task is ready to serve requests. If omitted, the task is assumed to be ready if when the first port is open.
	ReadinessProbe *Probe `json:"readinessProbe,omitempty"`
	// The command to run in the container or on the host. If both the image and the command are omitted, this is a noop.
	Command Strings `json:"command,omitempty"`
	// The arguments to pass to the command
	Args Strings `json:"args,omitempty"`
	// The working directory in the container or on the host
	WorkingDir string `json:"workingDir,omitempty"`
	// The user to run the task as.
	User string `json:"user,omitempty"`
	// Environment variables to set in the container or on the host
	Env EnvVars `json:"env,omitempty"`
	// The ports to expose
	Ports Ports `json:"ports,omitempty"`
	// Volumes to mount in the container
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`
	// Use a pseudo-TTY
	TTY bool `json:"tty,omitempty"`
	// A list of files to watch for changes, and restart the task if they change
	Watch Strings `json:"watch,omitempty"`
	// A mutex to prevent multiple tasks with the same mutex from running at the same time
	Mutex string `json:"mutex,omitempty"`
	// A semaphore to limit the number of tasks with the same semaphore that can run at the same time
	Semaphore string `json:"semaphore,omitempty"`
	// A list of tasks to run before this task
	Dependencies Strings `json:"dependencies,omitempty"`
	// The restart policy, e.g. Always, Never, OnFailure
	RestartPolicy string `json:"restartPolicy,omitempty"`
}

func (t *Task) IsBackground() bool {
	return t != nil && t.GetReadinessProbe() != nil && t.GetLivenessProbe() != nil
}

func (t Task) GetHostPorts() []uint16 {
	var ports []uint16
	for _, p := range t.Ports {
		ports = append(ports, p.GetHostPort())
	}
	return ports
}

func (t *Task) GetReadinessProbe() *Probe {
	if t == nil {
		return nil
	}
	if t.ReadinessProbe != nil {
		return t.ReadinessProbe
	}
	if len(t.Ports) > 0 {
		return &Probe{TCPSocket: &TCPSocketAction{Port: t.Ports[0].GetHostPort()}}
	}
	return nil
}

func (t *Task) GetLivenessProbe() *Probe {
	if t == nil {
		return nil
	}
	if t.LivenessProbe != nil {
		return t.LivenessProbe
	}
	return nil

}

func (t *Task) GetRestartPolicy() string {
	if t.RestartPolicy != "" {
		return t.RestartPolicy
	}
	return "OnFailure"
}

type Pod struct {
	// The specification of tasks to run.
	Spec PodSpec `json:"spec"`
	// APIVersion must be `kit/v1`.
	ApiVersion string `json:"apiVersion,omitempty"`
	// Kind must be `Tasks`.
	Kind string `json:"kind,omitempty"`
	// Metadata is the metadata for the pod.
	Metadata Metadata `json:"metadata"`
}

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
		return 3 * time.Second
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

// TCPSocketAction describes an action based on opening a socket
type TCPSocketAction struct {
	// Port number of the port to probe.
	Port uint16 `json:"port"`
}

func (a TCPSocketAction) URL() *url.URL {
	return &url.URL{Scheme: "tcp", Host: fmt.Sprintf(":%v", a.Port)}
}

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
	return &url.URL{Scheme: a.GetProto(), Host: fmt.Sprintf(":%v", a.Port), Path: a.Path}
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

// VolumeMount describes a mounting of a Volume within a container.
type VolumeMount struct {
	// This must match the name of a volume.
	Name string `json:"name"`
	// Path within the container at which the volume should be mounted.
	MountPath string `json:"mountPath"`
}

type HostPath struct {
	// Path of the directory on the host.
	Path string `json:"path"`
}

type Volume struct {
	// Volume's name.
	Name string `json:"name"`
	// HostPath represents a pre-existing file or directory on the host machine that is directly exposed to the container.
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
	var todo []string
	for _, name := range names {
		todo = append(todo, name)
	}
	done := map[string]bool{}
	for len(todo) > 0 {
		name := todo[0]
		todo = todo[1:]
		done[name] = true
		for _, d := range t.Get(name).Dependencies {
			if !done[d] {
				todo = append(todo, d)
			}
		}
	}
	var out Tasks
	for name := range done {
		out = append(out, t.Get(name))
	}
	return out
}
func (t Tasks) Has(name string) bool {
	for _, task := range t {
		if task.Name == name {
			return true
		}
	}
	return false
}
func (t Tasks) Get(name string) Task {
	for _, task := range t {
		if task.Name == name {
			return task
		}

	}
	panic(fmt.Errorf("no task named %q", name))
}

// Task is a unit of work that should be run.
type PodSpec struct {
	// TerminationGracePeriodSeconds is the grace period for terminating the pod.
	TerminationGracePeriodSeconds *int32 `json:"terminationGracePeriodSeconds,omitempty"`
	// Tasks is a list of tasks that should be run.
	Tasks Tasks `json:"tasks,omitempty"`
	// Volumes is a list of volumes that can be mounted by containers belonging to the pod.
	Volumes []Volume `json:"volumes,omitempty"`
	// Semaphores is a list of semaphores that can be acquired by tasks.
	Semaphores map[string]int `json:"semaphores,omitempty"`
}

func (s PodSpec) GetTerminationGracePeriod() time.Duration {
	if s.TerminationGracePeriodSeconds != nil {
		return time.Duration(*s.TerminationGracePeriodSeconds) * time.Second
	}
	return 3 * time.Second
}
