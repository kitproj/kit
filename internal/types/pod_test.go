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
	assert.Equal(t, 5*time.Second, probe.GetPeriod())
	assert.Equal(t, 5*time.Second, probe.GetInitialDelay())
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

func TestPort_Unstring(t *testing.T) {

	t.Run("Invalid", func(t *testing.T) {
		p := &Port{}
		err := p.Unstring("foo")
		assert.Error(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		p := &Port{}
		err := p.Unstring("8080:80")
		assert.NoError(t, err)
		assert.Equal(t, uint16(8080), p.ContainerPort)
		assert.Equal(t, uint16(80), p.HostPort)
	})

	t.Run("NoHostPort", func(t *testing.T) {
		p := &Port{}
		err := p.Unstring("8080")
		assert.NoError(t, err)
		assert.Equal(t, uint16(8080), p.ContainerPort)
		assert.Equal(t, uint16(8080), p.HostPort)
	})
}

func TestPort_String(t *testing.T) {
	t.Run("NoHostPort", func(t *testing.T) {
		p := &Port{ContainerPort: 8080}
		assert.Equal(t, "8080", p.String())
	})

	t.Run("WithHostPort", func(t *testing.T) {
		p := &Port{ContainerPort: 8080, HostPort: 80}
		assert.Equal(t, "8080:80", p.String())
	})
}

func TestPorts_Map(t *testing.T) {
	ports := Ports{
		{ContainerPort: 8080},
		{ContainerPort: 8081, HostPort: 80},
	}
	assert.Equal(t, map[uint16]uint16{8080: 8080, 8081: 80}, ports.Map())
}
