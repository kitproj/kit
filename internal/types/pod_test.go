package types

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestPod(t *testing.T) {
	data, err := os.ReadFile("testdata/tasks.yaml")
	assert.NoError(t, err)
	pod := &Pod{}
	err = yaml.Unmarshal(data, pod)
	assert.NoError(t, err)
	assert.Equal(t, "kit", pod.Metadata.Name)
	assert.Equal(t, map[string]string{"help": "https://github.com/kitproj/kit"}, pod.Metadata.Annotations)
	assert.Equal(t, 3*time.Second, pod.Spec.GetTerminationGracePeriod())
	assert.Len(t, pod.Spec.Tasks, 2)
	task := pod.Spec.Tasks[0]
	assert.Equal(t, []uint16{8080}, task.GetHostPorts())
	assert.Equal(t, "OnFailure", task.GetRestartPolicy())
	probe := task.GetReadinessProbe()
	assert.Equal(t, &Probe{TCPSocket: &TCPSocketAction{Port: 8080}}, probe)
	assert.Equal(t, 3*time.Second, probe.GetPeriod())
	assert.Equal(t, 3*time.Second, probe.GetInitialDelay())
	assert.Equal(t, 1, probe.GetSuccessThreshold())
	assert.Equal(t, 20, probe.GetFailureThreshold())
	assert.Nil(t, task.GetLivenessProbe())
	//
	assert.Equal(t, Strings{"sh", "-c", "echo bar"}, pod.Spec.Tasks[1].Command)
	assert.Equal(t, Strings{"baz", "qux"}, pod.Spec.Tasks[1].Dependencies)
}

func TestEnvVar_String(t *testing.T) {
	t.Run("Value", func(t *testing.T) {
		s, err := EnvVar{Name: "FOO", Value: "1"}.String()
		assert.NoError(t, err)
		assert.Equal(t, "FOO=1", s)
	})
	t.Run("ValueFrom", func(t *testing.T) {
		s, err := EnvVar{Name: "FOO", ValueFrom: &EnvVarSource{File: "testdata/six"}}.String()
		assert.NoError(t, err)
		assert.Equal(t, "FOO=6", s)
	})
}

func TestTask_AllTargetsExist(t *testing.T) {
	t.Run("Empty targets", func(t *testing.T) {
		task := &Task{}
		assert.False(t, task.AllTargetsExist())
	})
	t.Run("Missing targets", func(t *testing.T) {
		task := &Task{Targets: []string{"missing"}}
		assert.False(t, task.AllTargetsExist())
	})
	t.Run("Empty targets", func(t *testing.T) {
		task := &Task{Targets: []string{"testdata"}}
		assert.True(t, task.AllTargetsExist())
	})

}
