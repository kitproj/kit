package types

import (
	"fmt"
	"hash/adler32"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (t *Task) HasMutex() bool {
	return t != nil && t.Mutex != ""
}

// A task is a container or a command to run.
type Task struct {
	// Type is the type of the task: "service" or "job". If omitted, if there are ports, it's a service, otherwise it's a job.
	// This is only needed when you have service that does not listen on ports.
	// Services are running in the background.
	Type TaskType `json:"type,omitempty"`
	// Where to log the output of the task. E.g. if the task is verbose. Defaults to /dev/stdout. Maybe a file, or /dev/null.
	Log string `json:"log,omitempty"`
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
	// The shell script to run, instead of the command
	Sh string `json:"sh,omitempty"`
	// A directories or files of Kubernetes manifests to apply. Once running the task will wait for the resources to be ready.
	Manifests Strings `json:"manifests,omitempty"`
	// The namespace to run the Kubernetes resource in. Defaults to the namespace of the current Kubernetes context.
	Namespace string `json:"namespace,omitempty"`
	// The working directory in the container or on the host
	WorkingDir string `json:"workingDir,omitempty"`
	// The user to run the task as.
	User string `json:"user,omitempty"`
	// Environment variables to set in the container or on the host
	Env EnvVars `json:"env,omitempty"`
	// Environment file (e.g. .env) to use
	Envfile Envfile `json:"envfile,omitempty"`
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
	// A list of files this task will create. If these exist, and they're newer than the watched files, the task is skipped.
	Targets Strings `json:"targets,omitempty"`
	// The restart policy, e.g. Always, Never, OnFailure. Defaults depends on the type of task.
	RestartPolicy string `json:"restartPolicy,omitempty"`
	// The timeout for the task to be considered stalled. If omitted, the task will be considered stalled after 30 seconds of no activity.
	StalledTimeout *metav1.Duration `json:"stalledTimeout,omitempty"`
	// The group this task belongs to. Tasks in the same group will be visually grouped together in the UI.
	Group string `json:"group,omitempty"`
	// Whether this is the default task to run if no task is specified.
	Default bool `json:"default,omitempty"`
}

func (t *Task) GetHostPorts() []uint16 {
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
	if t.GetType() == TaskTypeService {
		return "Always"
	}
	return "Never"
}

func (t *Task) String() string {
	if t.Image != "" {
		return t.Image
	}
	if len(t.GetCommand()) > 0 {
		return t.GetCommand().String()
	}
	if t.Args != nil {
		return t.Args.String()
	}
	return "noop"
}

func (t *Task) Environ() ([]string, error) {
	environ, err := t.Envfile.Environ(t.WorkingDir)
	if err != nil {
		return nil, err
	}
	s, err := t.Env.Environ()
	return append(environ, s...), err
}

// CreateBootstrapFile creates a temporary bootstrap file for shell scripts
// and returns its path. The file is named using the pattern: <name>-<hash>.sh
// where hash is an 8-character hex string.
func (t *Task) CreateBootstrapFile(name string) (string, error) {
	if t.Sh == "" {
		return "", fmt.Errorf("no shell script to create bootstrap file for")
	}
	
	// Calculate hash of the script content
	hash := adler32.Checksum([]byte(t.Sh))
	
	// Create filename with format: name-XXXXXXXX.sh (8 hex chars)
	// Using %08x to ensure consistent 8-character hex output with zero-padding
	filename := fmt.Sprintf("%s-%08x.sh", name, hash)
	
	// Create the file in the working directory (or current directory if not set)
	workDir := t.WorkingDir
	if workDir == "" {
		workDir = "."
	}
	
	filepath := filepath.Join(workDir, filename)
	
	// Write the script content to the file
	err := os.WriteFile(filepath, []byte(t.Sh), 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create bootstrap file: %w", err)
	}
	
	return filepath, nil
}

func (t *Task) GetCommand() Strings {
	if len(t.Command) > 0 {
		return t.Command
	}
	if t.Sh != "" {
		return []string{"sh", "-c", t.Sh}
	}
	return nil
}

// Skip Determines if all the targets exist. And if they're all newer that the newest source file.
func (t *Task) Skip() bool {
	// if there are no targets, we must run the task
	if len(t.Targets) == 0 {
		return false
	}

	youngestSource := time.Time{}
	for _, source := range t.Watch {
		stat, err := os.Stat(filepath.Join(t.WorkingDir, source))
		if err != nil {
			continue
		}
		if stat.ModTime().After(youngestSource) {
			youngestSource = stat.ModTime()
		}
	}

	oldestTarget := time.Now()
	for _, target := range t.Targets {
		stat, err := os.Stat(filepath.Join(t.WorkingDir, target))
		// if the target does not exist, we must run the task
		if err != nil {
			return false
		}
		if stat.ModTime().Before(oldestTarget) {
			oldestTarget = stat.ModTime()
		}
	}

	return oldestTarget.After(youngestSource)
}

func (t *Task) GetType() TaskType {
	if t.Type != "" {
		return t.Type
	}
	if len(t.Ports) > 0 || t.LivenessProbe != nil || t.ReadinessProbe != nil {
		return TaskTypeService
	}
	return TaskTypeJob

}

func (t *Task) GetStalledTimeout() time.Duration {
	if t.StalledTimeout != nil {
		return t.StalledTimeout.Duration
	}
	return 30 * time.Second
}
