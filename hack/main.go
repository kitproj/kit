package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/alexec/kit/internal/types"
	"github.com/invopop/jsonschema"
)

func main() {
	r := new(jsonschema.Reflector)
	_ = r.AddGoComments("github.com/alexec/kit", "./")
	s := r.Reflect(types.Pod{})
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
	if err := os.WriteFile("schema/pod.schema.json", data, 0o777); err != nil {
		log.Fatal(err)
	}
}
