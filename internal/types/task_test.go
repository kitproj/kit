package types

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTask_AllTargetsExist(t *testing.T) {
	// touch testdata/younger
	err := os.Chtimes("testdata/younger", time.Now(), time.Now())
	assert.NoError(t, err)

	tests := []struct {
		name    string
		sources Strings
		targets Strings
		exist   bool
	}{
		{name: "No source, no target", sources: nil, targets: nil, exist: false},
		{name: "Source, no target", sources: Strings{"testdata"}, targets: nil, exist: false},
		{name: "Target, no source", sources: nil, targets: Strings{"testdata"}, exist: true},
		{name: "Missing source", sources: Strings{"missing"}, targets: Strings{"testdata"}, exist: true},
		{name: "Missing targets", sources: Strings{"testdata"}, targets: Strings{"missing"}, exist: false},
		{name: "Target younger than source", sources: Strings{"testdata/older"}, targets: Strings{"testdata/younger"}, exist: true},
		{name: "Target older than source", sources: Strings{"testdata/younger"}, targets: Strings{"testdata/older"}, exist: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := &Task{Watch: test.sources, Targets: test.targets}
			assert.Equal(t, test.exist, task.Skip())
		})
	}
}

func TestTask_Validate(t *testing.T) {
	t.Run("NoExecField", func(t *testing.T) {
		task := &Task{}
		assert.NoError(t, task.Validate())
	})
	t.Run("CommandOnly", func(t *testing.T) {
		task := &Task{Command: Strings{"echo", "hello"}}
		assert.NoError(t, task.Validate())
	})
	t.Run("ShOnly", func(t *testing.T) {
		task := &Task{Sh: "echo hello"}
		assert.NoError(t, task.Validate())
	})
	t.Run("ImageOnly", func(t *testing.T) {
		task := &Task{Image: "nginx"}
		assert.NoError(t, task.Validate())
	})
	t.Run("ManifestsOnly", func(t *testing.T) {
		task := &Task{Manifests: Strings{"deploy.yaml"}}
		assert.NoError(t, task.Validate())
	})
	t.Run("CommandAndSh", func(t *testing.T) {
		task := &Task{Command: Strings{"echo"}, Sh: "echo hello"}
		assert.EqualError(t, task.Validate(), "only one of [command sh] is allowed")
	})
	t.Run("ShAndImage", func(t *testing.T) {
		task := &Task{Sh: "echo hello", Image: "nginx"}
		assert.NoError(t, task.Validate())
	})
	t.Run("CommandAndManifests", func(t *testing.T) {
		task := &Task{Command: Strings{"echo"}, Manifests: Strings{"deploy.yaml"}}
		assert.EqualError(t, task.Validate(), "only one of [command manifests] is allowed")
	})
	t.Run("CommandAndImage", func(t *testing.T) {
		task := &Task{Command: Strings{"echo"}, Image: "nginx"}
		assert.EqualError(t, task.Validate(), "only one of [command image] is allowed")
	})
	t.Run("ShAndManifests", func(t *testing.T) {
		task := &Task{Sh: "echo hello", Manifests: Strings{"deploy.yaml"}}
		assert.EqualError(t, task.Validate(), "only one of [sh manifests] is allowed")
	})
	t.Run("ImageAndManifests", func(t *testing.T) {
		task := &Task{Image: "nginx", Manifests: Strings{"deploy.yaml"}}
		assert.EqualError(t, task.Validate(), "only one of [image manifests] is allowed")
	})
	t.Run("ThreeFields", func(t *testing.T) {
		task := &Task{Command: Strings{"echo"}, Sh: "echo hello", Image: "nginx"}
		assert.EqualError(t, task.Validate(), "only one of [command sh image] is allowed")
	})
}

func TestTask_GetType(t *testing.T) {
	t.Run("Defined", func(t *testing.T) {
		task := &Task{Type: TaskTypeService}
		assert.Equal(t, TaskTypeService, task.GetType())
	})
	t.Run("Job", func(t *testing.T) {
		task := &Task{}
		assert.Equal(t, TaskTypeJob, task.GetType())
	})
	t.Run("Ports", func(t *testing.T) {
		task := &Task{Ports: []Port{{}}}
		assert.Equal(t, TaskTypeService, task.GetType())
	})
	t.Run("LivenessProbe", func(t *testing.T) {
		task := &Task{LivenessProbe: &Probe{}}
		assert.Equal(t, TaskTypeService, task.GetType())
	})
	t.Run("ReadinessProbe", func(t *testing.T) {
		task := &Task{ReadinessProbe: &Probe{}}
		assert.Equal(t, TaskTypeService, task.GetType())
	})
}
