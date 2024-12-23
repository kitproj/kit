package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnviron(t *testing.T) {

	environ, err := Environ(PodSpec{
		Envfile: Envfile{"testdata/spec.env"},
		Env: EnvVars{
			{
				Name:  "BAR",
				Value: "2",
			},
			{
				Name:      "GAK",
				ValueFrom: &EnvVarSource{File: "testdata/six"},
			},
		},
	}, Task{
		Envfile: Envfile{"testdata/task.env"},
		Env: EnvVars{
			{
				Name:  "QUX",
				Value: "4",
			},
			{
				Name:  "FUZ",
				Value: "5",
			},
		},
	})

	assert.NoError(t, err)

	assert.ElementsMatch(t, []string{"FOO=1", "BAR=2", "BAZ=3", "QUX=4", "FUZ=5", "GAK=6"}, environ)

}
