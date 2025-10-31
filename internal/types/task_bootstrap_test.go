package types

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTask_CreateBootstrapFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	
	t.Run("Creates bootstrap file with correct naming pattern", func(t *testing.T) {
		task := &Task{
			Sh:         "echo 'Hello World'",
			WorkingDir: tmpDir,
		}
		
		filename, err := task.CreateBootstrapFile("test-task")
		if err != nil {
			t.Fatalf("CreateBootstrapFile failed: %v", err)
		}
		
		// Check that file was created
		if _, err := os.Stat(filename); err != nil {
			t.Errorf("Bootstrap file was not created: %v", err)
		}
		
		// Check filename pattern: test-task-XXXXXXXX.sh
		base := filepath.Base(filename)
		if !strings.HasPrefix(base, "test-task-") {
			t.Errorf("Bootstrap file doesn't have correct prefix. Expected 'test-task-', got: %s", base)
		}
		if !strings.HasSuffix(base, ".sh") {
			t.Errorf("Bootstrap file doesn't have .sh extension. Got: %s", base)
		}
		
		// Check that the hash part is present (should be 8 hex characters)
		parts := strings.TrimSuffix(base, ".sh")
		parts = strings.TrimPrefix(parts, "test-task-")
		if len(parts) < 8 {
			t.Errorf("Bootstrap file hash is too short. Expected at least 8 chars, got: %d (%s)", len(parts), parts)
		}
		
		// Check file content
		content, err := os.ReadFile(filename)
		if err != nil {
			t.Fatalf("Failed to read bootstrap file: %v", err)
		}
		
		expected := "echo 'Hello World'"
		if string(content) != expected {
			t.Errorf("Bootstrap file content incorrect. Expected %q, got %q", expected, string(content))
		}
		
		// Check file permissions (should be executable)
		info, err := os.Stat(filename)
		if err != nil {
			t.Fatalf("Failed to stat bootstrap file: %v", err)
		}
		
		mode := info.Mode()
		if mode&0100 == 0 {
			t.Errorf("Bootstrap file is not executable. Mode: %v", mode)
		}
		
		// Clean up
		os.Remove(filename)
	})
	
	t.Run("Different scripts produce different hashes", func(t *testing.T) {
		task1 := &Task{
			Sh:         "echo 'Hello'",
			WorkingDir: tmpDir,
		}
		task2 := &Task{
			Sh:         "echo 'World'",
			WorkingDir: tmpDir,
		}
		
		filename1, err := task1.CreateBootstrapFile("task")
		if err != nil {
			t.Fatalf("CreateBootstrapFile failed: %v", err)
		}
		defer os.Remove(filename1)
		
		filename2, err := task2.CreateBootstrapFile("task")
		if err != nil {
			t.Fatalf("CreateBootstrapFile failed: %v", err)
		}
		defer os.Remove(filename2)
		
		if filename1 == filename2 {
			t.Errorf("Different scripts should produce different filenames. Both got: %s", filename1)
		}
	})
	
	t.Run("Same script produces same hash", func(t *testing.T) {
		script := "echo 'consistent'"
		task1 := &Task{
			Sh:         script,
			WorkingDir: tmpDir,
		}
		task2 := &Task{
			Sh:         script,
			WorkingDir: tmpDir,
		}
		
		filename1, err := task1.CreateBootstrapFile("task")
		if err != nil {
			t.Fatalf("CreateBootstrapFile failed: %v", err)
		}
		defer os.Remove(filename1)
		
		filename2, err := task2.CreateBootstrapFile("task")
		if err != nil {
			t.Fatalf("CreateBootstrapFile failed: %v", err)
		}
		defer os.Remove(filename2)
		
		if filename1 != filename2 {
			t.Errorf("Same script should produce same filename. Got %s and %s", filename1, filename2)
		}
	})
	
	t.Run("Error when no shell script", func(t *testing.T) {
		task := &Task{
			WorkingDir: tmpDir,
		}
		
		_, err := task.CreateBootstrapFile("task")
		if err == nil {
			t.Errorf("Expected error when creating bootstrap file with no shell script")
		}
	})
	
	t.Run("Uses task name as prefix", func(t *testing.T) {
		task := &Task{
			Sh:         "echo 'test'",
			WorkingDir: tmpDir,
		}
		
		filename, err := task.CreateBootstrapFile("jira-bootstrap")
		if err != nil {
			t.Fatalf("CreateBootstrapFile failed: %v", err)
		}
		defer os.Remove(filename)
		
		base := filepath.Base(filename)
		if !strings.HasPrefix(base, "jira-bootstrap-") {
			t.Errorf("Bootstrap file doesn't have correct prefix. Expected 'jira-bootstrap-', got: %s", base)
		}
	})
}
