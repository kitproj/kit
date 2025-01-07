package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_sortUnstructureds(t *testing.T) {
	t.Run("CRD before CR", func(t *testing.T) {
		unstructureds := []*unstructured.Unstructured{
			{Object: map[string]any{}},
			{Object: map[string]any{"kind": "CustomResourceDefinition"}},
		}
		sortUnstructureds(unstructureds)
		assert.Equal(t, "CustomResourceDefinition", unstructureds[0].GetKind())
		assert.Equal(t, "", unstructureds[1].GetKind())
	})
	t.Run("Namespace before other resources", func(t *testing.T) {
		unstructureds := []*unstructured.Unstructured{
			{Object: map[string]any{"kind": "Deployment"}},
			{Object: map[string]any{"kind": "Namespace"}},
		}
		sortUnstructureds(unstructureds)
		assert.Equal(t, "Namespace", unstructureds[0].GetKind())
		assert.Equal(t, "Deployment", unstructureds[1].GetKind())
	})
}
