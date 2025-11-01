package types

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvVars_Environ(t *testing.T) {
	t.Run("NoExpansion", func(t *testing.T) {
		envVars := EnvVars{
			"FOO": "bar",
			"BAZ": "qux",
		}
		environ, err := envVars.Environ()
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"FOO=bar", "BAZ=qux"}, environ)
	})

	t.Run("WithValidExpansion", func(t *testing.T) {
		// Set up environment variables for expansion
		os.Setenv("TEST_VAR", "test_value")
		defer os.Unsetenv("TEST_VAR")

		envVars := EnvVars{
			"FOO": "prefix_${TEST_VAR}_suffix",
			"BAR": "$TEST_VAR",
		}
		environ, err := envVars.Environ()
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"FOO=prefix_test_value_suffix", "BAR=test_value"}, environ)
	})

	t.Run("WithCrossReference", func(t *testing.T) {
		// Variables can reference other variables in the same EnvVars map
		envVars := EnvVars{
			"BASE":      "/home/user",
			"DATA_PATH": "${BASE}/data",
			"LOG_PATH":  "${BASE}/logs",
		}
		environ, err := envVars.Environ()
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{
			"BASE=/home/user",
			"DATA_PATH=/home/user/data",
			"LOG_PATH=/home/user/logs",
		}, environ)
	})

	t.Run("WithMissingVariable", func(t *testing.T) {
		// Make sure the variable is not set
		os.Unsetenv("MISSING_VAR")

		envVars := EnvVars{
			"FOO": "value_with_${MISSING_VAR}",
		}
		_, err := envVars.Environ()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MISSING_VAR")
	})

	t.Run("WithMultipleMissingVariables", func(t *testing.T) {
		// Make sure the variables are not set
		os.Unsetenv("MISSING_ONE")
		os.Unsetenv("MISSING_TWO")

		envVars := EnvVars{
			"FOO": "${MISSING_ONE} and ${MISSING_TWO}",
		}
		_, err := envVars.Environ()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MISSING_ONE")
		assert.Contains(t, err.Error(), "MISSING_TWO")
	})

	t.Run("MixedValidAndMissing", func(t *testing.T) {
		os.Setenv("VALID_VAR", "valid")
		defer os.Unsetenv("VALID_VAR")
		os.Unsetenv("MISSING_VAR")

		envVars := EnvVars{
			"FOO": "${VALID_VAR}",
			"BAR": "${MISSING_VAR}",
		}
		_, err := envVars.Environ()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MISSING_VAR")
	})

	t.Run("CircularReference", func(t *testing.T) {
		// Test that circular references don't cause infinite loops
		// Note: os.Expand only expands once, so this should just result in unexpanded references
		envVars := EnvVars{
			"A": "${B}",
			"B": "${A}",
		}
		environ, err := envVars.Environ()
		// This won't error because A and B are defined in the map,
		// but the values will still contain unexpanded references
		assert.NoError(t, err)
		// The result depends on which variable is processed first (maps are unordered)
		// but both will contain a reference to the other
		assert.Len(t, environ, 2)
	})

	t.Run("EmptyStringVariable", func(t *testing.T) {
		// Test that environment variables set to empty string are not treated as missing
		os.Setenv("EMPTY_VAR", "")
		defer os.Unsetenv("EMPTY_VAR")

		envVars := EnvVars{
			"FOO": "value_${EMPTY_VAR}_end",
		}
		environ, err := envVars.Environ()
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"FOO=value__end"}, environ)
	})
}
