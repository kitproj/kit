package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kitproj/kit/internal/agents"
)

// Import imports rules from all agents into the normalized format
func Import() error {
	fmt.Println("Importing rules from all agents...")
	
	// Open AGENTS.md for appending
	agentsFile := "AGENTS.md"
	
	// Get absolute path to AGENTS.md to avoid importing it
	agentsFilePath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	agentsFilePath = agentsFilePath + "/AGENTS.md"
	
	f, err := os.OpenFile(agentsFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open AGENTS.md: %w", err)
	}
	defer f.Close()
	
	// Process each agent
	for _, agent := range agents.GetAllAgents() {
		if agent == agents.Default {
			// Skip the default agent (normalized format)
			continue
		}
		
		fmt.Printf("\nProcessing agent: %s\n", agent)
		
		paths, err := agents.GetRulePaths(agent)
		if err != nil {
			fmt.Printf("  Warning: failed to get paths for %s: %v\n", agent, err)
			continue
		}
		
		if len(paths) == 0 {
			fmt.Printf("  No rules found for %s\n", agent)
			continue
		}
		
		// Process each rule file, skip AGENTS.md to avoid recursion
		for _, path := range paths {
			absPath := path
			if !filepath.IsAbs(path) {
				cwd, _ := os.Getwd()
				absPath = filepath.Join(cwd, path)
			}
			
			// Skip AGENTS.md to avoid importing itself
			if absPath == agentsFilePath {
				continue
			}
			
			if err := importFile(f, agent, path); err != nil {
				fmt.Printf("  Warning: failed to import %s: %v\n", path, err)
			}
		}
	}
	
	fmt.Println("\nImport complete!")
	return nil
}

// importFile imports a single file into AGENTS.md
func importFile(w io.Writer, agent agents.Agent, path string) error {
	fmt.Printf("  Importing: %s\n", path)
	
	// Read the file content
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	
	// Write to AGENTS.md with a header
	_, err = fmt.Fprintf(w, "\n<!-- Imported from %s: %s -->\n", agent, path)
	if err != nil {
		return err
	}
	
	_, err = w.Write(content)
	if err != nil {
		return err
	}
	
	_, err = w.Write([]byte("\n"))
	return err
}
