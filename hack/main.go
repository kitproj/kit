package main

import (
	"os"

	"sigs.k8s.io/yaml"

	"github.com/alexec/kit/internal/types"
	"github.com/invopop/jsonschema"
)

func main() {
	r := new(jsonschema.Reflector)
	_ = r.AddGoComments("github.com/alexec/kit", "./")
	s := r.Reflect(types.Pod{})
	data, _ := yaml.Marshal(s)
	_ = os.WriteFile("pod.schema.yaml", data, 0o777)

}
