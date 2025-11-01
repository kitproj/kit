package types

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
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
// Variable references in values (using ${VAR} or $VAR syntax) are expanded using os.Expand.
// Variables can reference other variables in the same EnvVars map or from the OS environment.
// If any referenced variables are not found, an error is returned.
func (v EnvVars) Environ() ([]string, error) {
	var environ []string
	var missingVars []string

	// Create a mapping function that looks up variables in this map first, then OS environment
	// and tracks missing variables
	mapping := func(key string) string {
		// First check if the variable is defined in this EnvVars map
		if val, ok := v[key]; ok {
			return val
		}
		// Otherwise check OS environment
		val, exists := os.LookupEnv(key)
		if !exists {
			missingVars = append(missingVars, key)
		}
		return val
	}

	for k, val := range v {
		// Expand variable references in the value
		expandedValue := os.Expand(val, mapping)
		environ = append(environ, fmt.Sprintf("%s=%s", k, expandedValue))
	}

	// If any variables were missing, return an error
	if len(missingVars) > 0 {
		// Remove duplicates and sort for consistent error messages
		uniqueVars := make(map[string]bool)
		for _, v := range missingVars {
			uniqueVars[v] = true
		}
		var sortedVars []string
		for v := range uniqueVars {
			sortedVars = append(sortedVars, v)
		}
		sort.Strings(sortedVars)
		return nil, fmt.Errorf("undefined environment variables: %s", strings.Join(sortedVars, ", "))
	}

	return environ, nil
}
