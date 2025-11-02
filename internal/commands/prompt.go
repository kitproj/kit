package commands

import (
	"fmt"
	"os"

	"github.com/kitproj/kit/internal/agents"
)

// Prompt prints the aggregated prompt from all rule files
func Prompt() error {
	fmt.Println("# Aggregated Prompt from All Rules\n")
	
	// Get paths from the default agent (normalized format)
	paths, err := agents.GetRulePaths(agents.Default)
	if err != nil {
		return fmt.Errorf("failed to get rule paths: %w", err)
	}
	
	if len(paths) == 0 {
		fmt.Println("No rule files found.")
		return nil
	}
	
	// Print each file's content
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read %s: %v\n", path, err)
			continue
		}
		
		fmt.Printf("<!-- From: %s -->\n", path)
		fmt.Println(string(content))
		fmt.Println()
	}
	
	return nil
}
