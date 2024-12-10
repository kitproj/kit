package types

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEnviron(t *testing.T) {

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
		},
	})

	assert.NoError(t, err)

	assert.True(t, len(environ) >= 4, "has at least 4 elements (including OS values))")

	// remove OS values
	environ = environ[:len(environ)-len(os.Environ())]

	assert.Equal(t, []string{"FOO=1", "BAR=2", "BAZ=3", "QUX=4"}, environ)

}
