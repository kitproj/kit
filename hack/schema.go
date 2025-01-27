package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/kitproj/kit/internal/types"
)

func updateSchema() error {
	log.Println("updating schema")
	r := new(jsonschema.Reflector)
	_ = r.AddGoComments("github.com/kitproj/kit", "./")
	s := r.Reflect(types.Workflow{})
	for i, definition := range s.Definitions {
		definition.Title = i
		s.Definitions[i] = definition
		if definition.Properties != nil {
			for _, name := range definition.Properties.Keys() {
				value, _ := definition.Properties.Get(name)
				property := value.(*jsonschema.Schema)
				property.Title = name
			}
		}
	}
	data, _ := json.MarshalIndent(s, "", "  ")
	if err := os.WriteFile("schema/workflow.schema.json", data, 0o777); err != nil {
		return fmt.Errorf("failed to write schema/workflow.schema.json: %w", err)
	}
	return nil
}
