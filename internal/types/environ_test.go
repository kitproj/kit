package types

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnviron(t *testing.T) {

	os.Clearenv()
	err := os.Setenv("FUZ", "5")
	assert.NoError(t, err)

	environ, err := Environ(PodSpec{
		Envfile: "testdata/spec.env",
		Env: EnvVars{
			{
				Name:  "BAR",
				Value: "2",
			},
		},
	}, Task{
		Envfile: "testdata/task.env",
		Env: EnvVars{
			{
				Name:  "QUX",
				Value: "4",
			},
			{
				Name:  "FUZ",
				Value: "0",
			},
		},
	})

	assert.NoError(t, err)

	assert.ElementsMatch(t, []string{"FOO=1", "BAR=2", "BAZ=3", "QUX=4", "FUZ=5"}, environ)

}
