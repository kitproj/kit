package agents

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAllAgents(t *testing.T) {
	agents := GetAllAgents()
	if len(agents) != 9 {
		t.Errorf("Expected 9 agents, got %d", len(agents))
	}
	
	// Verify expected agents are present
	expectedAgents := map[Agent]bool{
		Default:       true,
		Claude:        true,
		Gemini:        true,
		Codex:         true,
		Cursor:        true,
		Augment:       true,
		GitHubCopilot: true,
		Windsurf:      true,
		Goose:         true,
	}
	
	for _, agent := range agents {
		if !expectedAgents[agent] {
			t.Errorf("Unexpected agent: %s", agent)
		}
	}
}

func TestGetRulePaths(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "agents-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Change to the temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)
	
	// Test Default agent with .prompts directory
	os.MkdirAll(".prompts", 0755)
	os.WriteFile(filepath.Join(".prompts", "test.md"), []byte("test"), 0644)
	
	paths, err := GetRulePaths(Default)
	if err != nil {
		t.Errorf("GetRulePaths(Default) failed: %v", err)
	}
	if len(paths) == 0 {
		t.Error("Expected at least one path for Default agent")
	}
	
	// Test Claude agent
	os.WriteFile("CLAUDE.md", []byte("test"), 0644)
	paths, err = GetRulePaths(Claude)
	if err != nil {
		t.Errorf("GetRulePaths(Claude) failed: %v", err)
	}
	if len(paths) == 0 {
		t.Error("Expected at least one path for Claude agent")
	}
	
	// Test unknown agent
	_, err = GetRulePaths("UnknownAgent")
	if err == nil {
		t.Error("Expected error for unknown agent")
	}
}

func TestExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "exists-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	testFile := filepath.Join(tmpDir, "test.txt")
	
	if exists(testFile) {
		t.Error("exists() returned true for non-existent file")
	}
	
	os.WriteFile(testFile, []byte("test"), 0644)
	
	if !exists(testFile) {
		t.Error("exists() returned false for existing file")
	}
}

func TestListMarkdownFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "md-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create some test files
	os.WriteFile(filepath.Join(tmpDir, "file1.md"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.md"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file3.txt"), []byte("test"), 0644)
	
	files, err := listMarkdownFiles(tmpDir)
	if err != nil {
		t.Fatalf("listMarkdownFiles failed: %v", err)
	}
	
	if len(files) != 2 {
		t.Errorf("Expected 2 markdown files, got %d", len(files))
	}
}

func TestListMDCFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mdc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create some test files
	os.WriteFile(filepath.Join(tmpDir, "file1.mdc"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.mdc"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file3.md"), []byte("test"), 0644)
	
	files, err := listMDCFiles(tmpDir)
	if err != nil {
		t.Fatalf("listMDCFiles failed: %v", err)
	}
	
	if len(files) != 2 {
		t.Errorf("Expected 2 mdc files, got %d", len(files))
	}
}

func TestFindAncestorFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ancestor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create nested directory structure
	subdir := filepath.Join(tmpDir, "a", "b", "c")
	os.MkdirAll(subdir, 0755)
	
	// Create test files at different levels
	os.WriteFile(filepath.Join(tmpDir, "TEST.md"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "a", "TEST.md"), []byte("test"), 0644)
	
	// Change to subdirectory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(subdir)
	
	paths := findAncestorFiles("TEST.md")
	if len(paths) < 2 {
		t.Errorf("Expected at least 2 ancestor files, got %d", len(paths))
	}
}
