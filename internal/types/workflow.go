package types

import (
	"encoding/json"
)

type Workflow Spec

// when unmarshalling legacy format, we need to convert it to the new format
func (p *Workflow) UnmarshalJSON(data []byte) error {
	// legacy format has a field named "spec"
	var x workflowV1
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if len(x.Spec.Tasks) > 0 {
		*p = Workflow(x.Spec)
		return nil
	}
	// otherwise, normal unmarshall
	return json.Unmarshal(data, (*Spec)(p))
}
