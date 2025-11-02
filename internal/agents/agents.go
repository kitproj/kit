package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Agent represents a coding agent type
type Agent string

const (
	Claude        Agent = "Claude"
	Gemini        Agent = "Gemini"
	Codex         Agent = "Codex"
	Cursor        Agent = "Cursor"
	Augment       Agent = "Augment"
	GitHubCopilot Agent = "GitHubCopilot"
	Windsurf      Agent = "Windsurf"
	Goose         Agent = "Goose"
	Default       Agent = "" // Default agent using .prompts directories
)

// RuleLevel represents the priority level of rules
type RuleLevel int

const (
	ProjectLevel  RuleLevel = 0 // Most important
	AncestorLevel RuleLevel = 1 // Next most important
	UserLevel     RuleLevel = 2 // User-specific
	SystemLevel   RuleLevel = 3 // System-wide
)

// RulePathFunc returns the list of rule file paths for a given agent at a specific level
type RulePathFunc func(level int) ([]string, error)

// AgentRulePath maps agent names to their rule path functions
var AgentRulePath = map[Agent]RulePathFunc{
	Default:       getDefaultRulePaths,
	Claude:        getClaudeRulePaths,
	Gemini:        getGeminiRulePaths,
	Codex:         getCodexRulePaths,
	Cursor:        getCursorRulePaths,
	Augment:       getAugmentRulePaths,
	GitHubCopilot: getGitHubCopilotRulePaths,
	Windsurf:      getWindsurfRulePaths,
	Goose:         getGooseRulePaths,
}

// getDefaultRulePaths returns paths for the default agent (.prompts)
func getDefaultRulePaths(level int) ([]string, error) {
	var paths []string
	
	switch RuleLevel(level) {
	case ProjectLevel:
		// .prompts directory in PWD
		if entries, err := listMarkdownFiles(".prompts"); err == nil {
			paths = append(paths, entries...)
		}
	case UserLevel:
		// ~/.prompts/rules directory
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		rulesDir := filepath.Join(home, ".prompts", "rules")
		if entries, err := listMarkdownFiles(rulesDir); err == nil {
			paths = append(paths, entries...)
		}
	}
	
	return paths, nil
}

// getClaudeRulePaths returns paths for Claude agent
func getClaudeRulePaths(level int) ([]string, error) {
	var paths []string
	
	switch RuleLevel(level) {
	case ProjectLevel:
		// ./CLAUDE.local.md
		if exists("CLAUDE.local.md") {
			paths = append(paths, "CLAUDE.local.md")
		}
	case AncestorLevel:
		// ./CLAUDE.md in ancestor directories
		ancestorPaths := findAncestorFiles("CLAUDE.md")
		paths = append(paths, ancestorPaths...)
	case UserLevel:
		// ~/.claude/CLAUDE.md
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		globalPath := filepath.Join(home, ".claude", "CLAUDE.md")
		if exists(globalPath) {
			paths = append(paths, globalPath)
		}
	}
	
	return paths, nil
}

// getGeminiRulePaths returns paths for Gemini agent
func getGeminiRulePaths(level int) ([]string, error) {
	var paths []string
	
	switch RuleLevel(level) {
	case ProjectLevel:
		// ./.gemini/styleguide.md
		styleguide := filepath.Join(".gemini", "styleguide.md")
		if exists(styleguide) {
			paths = append(paths, styleguide)
		}
	case AncestorLevel:
		// ./GEMINI.md in ancestor directories
		ancestorPaths := findAncestorFiles("GEMINI.md")
		paths = append(paths, ancestorPaths...)
	case UserLevel:
		// ~/.gemini/GEMINI.md
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		globalPath := filepath.Join(home, ".gemini", "GEMINI.md")
		if exists(globalPath) {
			paths = append(paths, globalPath)
		}
	}
	
	return paths, nil
}

// getCodexRulePaths returns paths for Codex agent
func getCodexRulePaths(level int) ([]string, error) {
	var paths []string
	
	switch RuleLevel(level) {
	case AncestorLevel:
		// AGENTS.md in ancestor directories and CWD
		ancestorPaths := findAncestorFiles("AGENTS.md")
		paths = append(paths, ancestorPaths...)
	case UserLevel:
		// ~/.codex/AGENTS.md
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		globalPath := filepath.Join(home, ".codex", "AGENTS.md")
		if exists(globalPath) {
			paths = append(paths, globalPath)
		}
	}
	
	return paths, nil
}

// getCursorRulePaths returns paths for Cursor agent
func getCursorRulePaths(level int) ([]string, error) {
	var paths []string
	
	switch RuleLevel(level) {
	case ProjectLevel:
		// ./.cursor/rules/ directory (*.mdc and *.md files)
		rulesDir := filepath.Join(".cursor", "rules")
		if entries, err := listMarkdownFiles(rulesDir); err == nil {
			paths = append(paths, entries...)
		}
		// Also check for .mdc files
		if entries, err := listMDCFiles(rulesDir); err == nil {
			paths = append(paths, entries...)
		}
		// AGENTS.md at project root
		if exists("AGENTS.md") {
			paths = append(paths, "AGENTS.md")
		}
	}
	
	return paths, nil
}

// getAugmentRulePaths returns paths for Augment agent
func getAugmentRulePaths(level int) ([]string, error) {
	var paths []string
	
	switch RuleLevel(level) {
	case ProjectLevel:
		// <workspace_root>/.augment/rules/
		rulesDir := filepath.Join(".augment", "rules")
		if entries, err := listMarkdownFiles(rulesDir); err == nil {
			paths = append(paths, entries...)
		}
		// Legacy: <workspace_root>/.augment/guidelines.md
		guidelines := filepath.Join(".augment", "guidelines.md")
		if exists(guidelines) {
			paths = append(paths, guidelines)
		}
	case AncestorLevel:
		// Compatibility: CLAUDE.md and AGENTS.md in ancestor directories
		claudePaths := findAncestorFiles("CLAUDE.md")
		paths = append(paths, claudePaths...)
		agentsPaths := findAncestorFiles("AGENTS.md")
		paths = append(paths, agentsPaths...)
	}
	
	return paths, nil
}

// getGitHubCopilotRulePaths returns paths for GitHub Copilot agent
func getGitHubCopilotRulePaths(level int) ([]string, error) {
	var paths []string
	
	switch RuleLevel(level) {
	case ProjectLevel:
		// .github/agents directory
		agentsDir := filepath.Join(".github", "agents")
		if entries, err := listMarkdownFiles(agentsDir); err == nil {
			paths = append(paths, entries...)
		}
	case AncestorLevel:
		// ./.github/copilot-instructions.md
		instructionsPath := filepath.Join(".github", "copilot-instructions.md")
		if exists(instructionsPath) {
			paths = append(paths, instructionsPath)
		}
		// ./AGENTS.md (nearest in directory tree)
		ancestorPaths := findAncestorFiles("AGENTS.md")
		if len(ancestorPaths) > 0 {
			// Only take the nearest one
			paths = append(paths, ancestorPaths[0])
		}
	}
	
	return paths, nil
}

// getWindsurfRulePaths returns paths for Windsurf agent
func getWindsurfRulePaths(level int) ([]string, error) {
	var paths []string
	
	switch RuleLevel(level) {
	case ProjectLevel, AncestorLevel:
		// ./.windsurf/rules/ (searched from workspace up to Git root)
		rulesDir := filepath.Join(".windsurf", "rules")
		if entries, err := listMarkdownFiles(rulesDir); err == nil {
			paths = append(paths, entries...)
		}
		// Also check ancestor directories for .windsurf/rules
		if RuleLevel(level) == AncestorLevel {
			ancestors := getAncestorDirs()
			for _, ancestor := range ancestors {
				ancestorRulesDir := filepath.Join(ancestor, ".windsurf", "rules")
				if entries, err := listMarkdownFiles(ancestorRulesDir); err == nil {
					paths = append(paths, entries...)
				}
			}
		}
	}
	
	return paths, nil
}

// getGooseRulePaths returns paths for Goose agent
func getGooseRulePaths(level int) ([]string, error) {
	var paths []string
	
	switch RuleLevel(level) {
	case AncestorLevel:
		// Relies on standard AGENTS.md
		ancestorPaths := findAncestorFiles("AGENTS.md")
		paths = append(paths, ancestorPaths...)
	}
	
	return paths, nil
}

// Helper functions

// exists checks if a file exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// listMarkdownFiles lists all .md files in a directory
func listMarkdownFiles(dir string) ([]string, error) {
	var files []string
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	
	return files, nil
}

// listMDCFiles lists all .mdc files in a directory
func listMDCFiles(dir string) ([]string, error) {
	var files []string
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".mdc") {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	
	return files, nil
}

// findAncestorFiles finds all instances of a file in ancestor directories
func findAncestorFiles(filename string) []string {
	var paths []string
	
	cwd, err := os.Getwd()
	if err != nil {
		return paths
	}
	
	dir := cwd
	for {
		filePath := filepath.Join(dir, filename)
		if exists(filePath) {
			paths = append(paths, filePath)
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			break
		}
		dir = parent
	}
	
	return paths
}

// getAncestorDirs returns all ancestor directories up to root
func getAncestorDirs() []string {
	var dirs []string
	
	cwd, err := os.Getwd()
	if err != nil {
		return dirs
	}
	
	dir := filepath.Dir(cwd)
	for {
		dirs = append(dirs, dir)
		
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			break
		}
		dir = parent
	}
	
	return dirs
}

// GetAllAgents returns a list of all supported agents
func GetAllAgents() []Agent {
	return []Agent{
		Default,
		Claude,
		Gemini,
		Codex,
		Cursor,
		Augment,
		GitHubCopilot,
		Windsurf,
		Goose,
	}
}

// GetRulePaths returns all rule paths for a given agent
func GetRulePaths(agent Agent) ([]string, error) {
	pathFunc, ok := AgentRulePath[agent]
	if !ok {
		return nil, fmt.Errorf("unknown agent: %s", agent)
	}
	
	var allPaths []string
	
	// Collect paths from all levels
	for level := 0; level <= int(SystemLevel); level++ {
		paths, err := pathFunc(level)
		if err != nil {
			return nil, err
		}
		allPaths = append(allPaths, paths...)
	}
	
	return allPaths, nil
}
