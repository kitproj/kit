package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBootstrap(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "bootstrap-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)
	
	// Run bootstrap
	err = Bootstrap()
	if err != nil {
		t.Errorf("Bootstrap failed: %v", err)
	}
	
	// Check if AGENTS.md was created
	if _, err := os.Stat("AGENTS.md"); os.IsNotExist(err) {
		t.Error("AGENTS.md was not created")
	}
	
	// Check if .prompts directory was created
	if _, err := os.Stat(".prompts"); os.IsNotExist(err) {
		t.Error(".prompts directory was not created")
	}
}

func TestImport(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "import-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)
	
	// Create test files
	os.WriteFile("CLAUDE.md", []byte("# Claude Test"), 0644)
	os.MkdirAll(".cursor/rules", 0755)
	os.WriteFile(".cursor/rules/test.mdc", []byte("# Cursor Test"), 0644)
	
	// Create AGENTS.md
	os.WriteFile("AGENTS.md", []byte("# Initial Content\n"), 0644)
	
	// Run import
	err = Import()
	if err != nil {
		t.Errorf("Import failed: %v", err)
	}
	
	// Check if content was appended to AGENTS.md
	content, err := os.ReadFile("AGENTS.md")
	if err != nil {
		t.Fatalf("Failed to read AGENTS.md: %v", err)
	}
	
	contentStr := string(content)
	if !strings.Contains(contentStr, "Claude Test") {
		t.Error("AGENTS.md does not contain Claude content")
	}
	if !strings.Contains(contentStr, "Cursor Test") {
		t.Error("AGENTS.md does not contain Cursor content")
	}
}

func TestExport(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "export-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)
	
	// Create AGENTS.md with test content
	testContent := "# Test Agent Rules\n\nSome rules here."
	os.WriteFile("AGENTS.md", []byte(testContent), 0644)
	
	// Test export to Gemini
	err = Export("Gemini")
	if err != nil {
		t.Errorf("Export to Gemini failed: %v", err)
	}
	
	// Check if GEMINI.md was created
	content, err := os.ReadFile("GEMINI.md")
	if err != nil {
		t.Error("GEMINI.md was not created")
	}
	if string(content) != testContent {
		t.Error("GEMINI.md content does not match AGENTS.md")
	}
	
	// Test export to Claude
	err = Export("Claude")
	if err != nil {
		t.Errorf("Export to Claude failed: %v", err)
	}
	
	// Check if CLAUDE.md was created
	content, err = os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Error("CLAUDE.md was not created")
	}
	if string(content) != testContent {
		t.Error("CLAUDE.md content does not match AGENTS.md")
	}
	
	// Test export to Windsurf
	err = Export("Windsurf")
	if err != nil {
		t.Errorf("Export to Windsurf failed: %v", err)
	}
	
	// Check if .windsurf/rules/AGENTS.md was created
	content, err = os.ReadFile(filepath.Join(".windsurf", "rules", "AGENTS.md"))
	if err != nil {
		t.Error(".windsurf/rules/AGENTS.md was not created")
	}
	if string(content) != testContent {
		t.Error(".windsurf/rules/AGENTS.md content does not match AGENTS.md")
	}
	
	// Test export with unknown agent
	err = Export("UnknownAgent")
	if err == nil {
		t.Error("Expected error for unknown agent")
	}
}

func TestPrompt(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "prompt-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)
	
	// Create .prompts directory with a test file
	os.MkdirAll(".prompts", 0755)
	os.WriteFile(filepath.Join(".prompts", "test.md"), []byte("# Test Prompt"), 0644)
	
	// Run prompt command - it should not error
	err = Prompt()
	if err != nil {
		t.Errorf("Prompt failed: %v", err)
	}
}
