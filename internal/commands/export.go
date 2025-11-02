package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kitproj/kit/internal/agents"
)

// Export exports rules from the normalized format to a specific agent
func Export(agentName string) error {
	// Parse agent name
	var targetAgent agents.Agent
	switch strings.ToLower(agentName) {
	case "claude":
		targetAgent = agents.Claude
	case "gemini":
		targetAgent = agents.Gemini
	case "codex":
		targetAgent = agents.Codex
	case "cursor":
		targetAgent = agents.Cursor
	case "augment":
		targetAgent = agents.Augment
	case "githubcopilot", "copilot":
		targetAgent = agents.GitHubCopilot
	case "windsurf":
		targetAgent = agents.Windsurf
	case "goose":
		targetAgent = agents.Goose
	default:
		return fmt.Errorf("unknown agent: %s", agentName)
	}
	
	fmt.Printf("Exporting rules to %s...\n", targetAgent)
	
	// Read the normalized AGENTS.md file
	agentsFile := "AGENTS.md"
	content, err := os.ReadFile(agentsFile)
	if err != nil {
		return fmt.Errorf("failed to read AGENTS.md: %w", err)
	}
	
	// Export to the target agent's format
	if err := exportToAgent(targetAgent, content); err != nil {
		return fmt.Errorf("failed to export to %s: %w", targetAgent, err)
	}
	
	fmt.Printf("Export to %s complete!\n", targetAgent)
	return nil
}

// exportToAgent exports content to a specific agent's format
func exportToAgent(agent agents.Agent, content []byte) error {
	switch agent {
	case agents.Claude:
		return exportToClaude(content)
	case agents.Gemini:
		return exportToGemini(content)
	case agents.Codex:
		return exportToCodex(content)
	case agents.Cursor:
		return exportToCursor(content)
	case agents.Augment:
		return exportToAugment(content)
	case agents.GitHubCopilot:
		return exportToGitHubCopilot(content)
	case agents.Windsurf:
		return exportToWindsurf(content)
	case agents.Goose:
		return exportToGoose(content)
	default:
		return fmt.Errorf("export not implemented for agent: %s", agent)
	}
}

func exportToClaude(content []byte) error {
	// Export to ./CLAUDE.md
	path := "CLAUDE.md"
	fmt.Printf("  Writing to %s\n", path)
	return os.WriteFile(path, content, 0644)
}

func exportToCodex(content []byte) error {
	// Codex uses AGENTS.md which is already the normalized format
	// No need to write - AGENTS.md is the source file
	path := "AGENTS.md"
	fmt.Printf("  Using existing %s (Codex uses the normalized format)\n", path)
	return nil
}

func exportToGemini(content []byte) error {
	// Export to ./GEMINI.md
	path := "GEMINI.md"
	fmt.Printf("  Writing to %s\n", path)
	return os.WriteFile(path, content, 0644)
}

func exportToCursor(content []byte) error {
	// Cursor uses AGENTS.md for compatibility mode
	// No need to write - AGENTS.md is the source file
	path := "AGENTS.md"
	fmt.Printf("  Using existing %s (Cursor compatibility mode)\n", path)
	return nil
}

func exportToAugment(content []byte) error {
	// Export to .augment/rules/AGENTS.md
	dir := filepath.Join(".augment", "rules")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, "AGENTS.md")
	fmt.Printf("  Writing to %s\n", path)
	return os.WriteFile(path, content, 0644)
}

func exportToGitHubCopilot(content []byte) error {
	// Export to ./AGENTS.md (compatibility) and .github/copilot-instructions.md
	paths := []string{
		"AGENTS.md",
		filepath.Join(".github", "copilot-instructions.md"),
	}
	
	for _, path := range paths {
		dir := filepath.Dir(path)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
		fmt.Printf("  Writing to %s\n", path)
		if err := os.WriteFile(path, content, 0644); err != nil {
			return err
		}
	}
	return nil
}

func exportToWindsurf(content []byte) error {
	// Export to .windsurf/rules/AGENTS.md
	dir := filepath.Join(".windsurf", "rules")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, "AGENTS.md")
	fmt.Printf("  Writing to %s\n", path)
	return os.WriteFile(path, content, 0644)
}

func exportToGoose(content []byte) error {
	// Goose uses AGENTS.md which is the normalized format
	// No need to write - AGENTS.md is the source file
	path := "AGENTS.md"
	fmt.Printf("  Using existing %s (Goose uses the normalized format)\n", path)
	return nil
}
