package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

// Bootstrap creates the necessary directory structure and files
func Bootstrap() error {
	fmt.Println("Bootstrapping coding-context...")
	
	// Create AGENTS.md if it doesn't exist
	agentsFile := "AGENTS.md"
	if _, err := os.Stat(agentsFile); os.IsNotExist(err) {
		fmt.Printf("Creating %s\n", agentsFile)
		defaultContent := []byte("# Coding Agent Rules\n\nThis file contains rules and guidelines for coding agents.\n")
		if err := os.WriteFile(agentsFile, defaultContent, 0644); err != nil {
			return fmt.Errorf("failed to create AGENTS.md: %w", err)
		}
	} else {
		fmt.Printf("%s already exists\n", agentsFile)
	}
	
	// Create .prompts directory
	promptsDir := ".prompts"
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", promptsDir, err)
	}
	fmt.Printf("Created %s directory\n", promptsDir)
	
	// Create user-level directory structure
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
	userPromptsDir := filepath.Join(home, ".prompts", "rules")
	if err := os.MkdirAll(userPromptsDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", userPromptsDir, err)
	}
	fmt.Printf("Created %s directory\n", userPromptsDir)
	
	fmt.Println("\nBootstrap complete!")
	return nil
}
