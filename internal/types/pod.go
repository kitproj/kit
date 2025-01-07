package types

import (
	"encoding/json"
	"time"
)

type Pod PodSpec

// when unmarshalling legacy format, we need to convert it to the new format
func (p *Pod) UnmarshalJSON(data []byte) error {
	// legacy format has a field named "spec"
	var x podV1
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if len(x.Spec.Tasks) > 0 {
		*p = Pod(x.Spec)
		return nil
	}
	// otherwise, normal unmarshall
	return json.Unmarshal(data, (*PodSpec)(p))
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
	// Environment variables to set in the container or on the host
	Env EnvVars `json:"env,omitempty"`
	// Environment file (e.g. .env) to use
	Envfile Envfile `json:"envfile,omitempty"`
}

func (s *PodSpec) GetTerminationGracePeriod() time.Duration {
	if s.TerminationGracePeriodSeconds != nil {
		return time.Duration(*s.TerminationGracePeriodSeconds) * time.Second
	}
	return 3 * time.Second
}

// Retuns the environment variables for the spec.
func (s *PodSpec) Environ() ([]string, error) {
	environ, err := s.Envfile.Environ("")
	if err != nil {
		return nil, err
	}
	e, err := s.Env.Environ()
	return append(environ, e...), err
}
