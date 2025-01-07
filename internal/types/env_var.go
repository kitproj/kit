package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

// A environment variable.
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (v EnvVar) String() (string, error) {
	return fmt.Sprintf("%s=%s", v.Name, v.Value), nil
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
	s, err := v.String()
	if err != nil {
		return nil, err
	}
	return json.Marshal(s)
}
