package types

import (
	"encoding/json"
	"fmt"
)

// A list of environment variables.
type EnvVars map[string]string

// the legacy format for env vars was an array of named env vars
func (v *EnvVars) UnmarshalJSON(data []byte) error {
	*v = EnvVars{}
	if data[0] == '[' {
		var x []EnvVar
		if err := json.Unmarshal(data, &x); err != nil {
			return err
		}
		for _, env := range x {
			(*v)[env.Name] = env.Value
		}
		return nil
	}
	var x = map[string]string{}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	for name, value := range x {
		(*v)[name] = value
	}
	return nil
}

// Environ returns a list of environment variables. If an environment variable is defined in both the task and the host, the host value is used.
func (v EnvVars) Environ() ([]string, error) {
	var environ []string
	for k, v := range v {
		environ = append(environ, fmt.Sprintf("%s=%s", k, v))
	}
	return environ, nil
}
