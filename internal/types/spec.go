package types

import "time"

// Task is a unit of work that should be run.
type Spec struct {
	// Port is the port that should be used for the user interface.
	Port *uint16 `json:"port,omitempty"`
	// TerminationGracePeriodSeconds is the grace period for terminating the workflow.
	TerminationGracePeriodSeconds *int32 `json:"terminationGracePeriodSeconds,omitempty"`
	// Tasks is a list of tasks that should be run.
	Tasks Tasks `json:"tasks,omitempty"`
	// Volumes is a list of volumes that can be mounted by containers belonging to the workflow.
	Volumes []Volume `json:"volumes,omitempty"`
	// Semaphores is a list of semaphores that can be acquired by tasks.
	Semaphores map[string]int `json:"semaphores,omitempty"`
	// Environment variables to set in the container or on the host
	Env EnvVars `json:"env,omitempty"`
	// Environment file (e.g. .env) to use
	Envfile Envfile `json:"envfile,omitempty"`
	// EnableMetrics enables resource usage metrics collection for all tasks
	EnableMetrics *bool `json:"enableMetrics,omitempty"`
}

func (s *Spec) GetTerminationGracePeriod() time.Duration {
	if s.TerminationGracePeriodSeconds != nil {
		return time.Duration(*s.TerminationGracePeriodSeconds) * time.Second
	}
	return 3 * time.Second
}

func (s *Spec) GetEnableMetrics() bool {
	if s.EnableMetrics != nil {
		return *s.EnableMetrics
	}
	return true // Default to enabled
}

// Retuns the environment variables for the spec.
func (s *Spec) Environ() ([]string, error) {
	environ, err := s.Envfile.Environ("")
	if err != nil {
		return nil, err
	}
	e, err := s.Env.Environ()
	return append(environ, e...), err
}
