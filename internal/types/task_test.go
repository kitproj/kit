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

func TestTask_IsService(t *testing.T) {
	t.Run("Job", func(t *testing.T) {
		task := &Task{}
		assert.False(t, task.IsService())
	})
	t.Run("Ports", func(t *testing.T) {
		task := &Task{Ports: []Port{{}}}
		assert.True(t, task.IsService())
	})
	t.Run("LivenessProbe", func(t *testing.T) {
		task := &Task{LivenessProbe: &Probe{}}
		assert.True(t, task.IsService())
	})
	t.Run("ReadinessProbe", func(t *testing.T) {
		task := &Task{ReadinessProbe: &Probe{}}
		assert.True(t, task.IsService())
	})
}
