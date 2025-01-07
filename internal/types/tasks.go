package types

import "encoding/json"

type Tasks map[string]Task

// the legacy format for tasks was an array of named tasks
func (t *Tasks) UnmarshalJSON(data []byte) error {
	*t = map[string]Task{}
	if data[0] == '[' {
		var x []NamedTask
		if err := json.Unmarshal(data, &x); err != nil {
			return err
		}
		for _, task := range x {
			(*t)[task.Name] = task.Task
		}
		return nil
	}
	var x = map[string]Task{}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	for name, task := range x {
		(*t)[name] = task
	}
	return nil
}
