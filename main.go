package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/kitproj/kit/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "import":
		if err := commands.Import(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "export":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: export command requires an agent name\n")
			fmt.Fprintf(os.Stderr, "Usage: coding-context export <agent>\n")
			fmt.Fprintf(os.Stderr, "Available agents: Claude, Gemini, Codex, Cursor, Augment, GitHubCopilot, Windsurf, Goose\n")
			os.Exit(1)
		}
		agentName := os.Args[2]
		if err := commands.Export(agentName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "bootstrap":
		if err := commands.Bootstrap(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "prompt":
		if err := commands.Prompt(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "-v", "--version":
		info, _ := debug.ReadBuildInfo()
		fmt.Printf("%v\n", info.Main.Version)
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("coding-context - Manage coding agent rules")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  coding-context <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  import                Import rules from all agents into AGENTS.md")
	fmt.Println("  export <agent>        Export rules to a specific agent format")
	fmt.Println("  bootstrap             Create initial directory structure and files")
	fmt.Println("  prompt                Print the aggregated prompt from all rules")
	fmt.Println("  -v, --version         Print version information")
	fmt.Println("  -h, --help            Print this help message")
	fmt.Println()
	fmt.Println("Supported agents:")
	fmt.Println("  Claude, Gemini, Codex, Cursor, Augment, GitHubCopilot, Windsurf, Goose")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  coding-context import")
	fmt.Println("  coding-context export Gemini")
	fmt.Println("  coding-context bootstrap")
	fmt.Println("  coding-context prompt")
}
