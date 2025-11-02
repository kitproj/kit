# coding-context - Manage Coding Agent Rules

A unified CLI tool for managing coding rules across different AI coding agents (Claude, Gemini, Cursor, GitHub Copilot, and more).

## What is coding-context?

coding-context simplifies the management of coding rules and guidelines across multiple AI coding agents. Instead of maintaining separate rule files for each agent, you can:

- **Import** rules from all agents into a normalized format
- **Export** rules to specific agent formats
- **Bootstrap** the necessary directory structure
- **View** aggregated prompts from all your rules

## Supported Agents

- **Claude** - Hierarchical concatenation with CLAUDE.md files
- **Gemini** - Strategic layer with GEMINI.md files  
- **Codex** - Uses AGENTS.md in ancestor directories
- **Cursor** - Declarative context injection with .cursor/rules/
- **Augment** - Structured rules in .augment/rules/
- **GitHub Copilot** - Uses .github/copilot-instructions.md
- **Windsurf** - Rules in .windsurf/rules/
- **Goose** - Standard AGENTS.md compatibility

## Quick Start

### Installation

Download the standalone binary from the [releases page](https://github.com/kitproj/kit/releases/latest):

```bash
# For Linux
sudo curl --fail --location --output /usr/local/bin/coding-context https://github.com/kitproj/kit/releases/download/v0.1.105/kit_v0.1.105_linux_386
sudo chmod +x /usr/local/bin/coding-context

# For Go users
go install github.com/kitproj/kit@v0.1.105
```

### Basic Usage

1. **Bootstrap** - Create initial directory structure:

```bash
coding-context bootstrap
```

2. **Import** - Gather rules from all agents:

```bash
coding-context import
```

3. **Export** - Distribute rules to a specific agent:

```bash
coding-context export Gemini
```

4. **View Prompt** - See aggregated rules:

```bash
coding-context prompt
```

## Commands

### import

Import rules from all supported agents into the normalized AGENTS.md format:

```bash
coding-context import
```

This command:
- Scans for rules from all supported agents
- Appends them to AGENTS.md with source annotations
- Converts .mdc files to .md format
- Skips duplicate imports of AGENTS.md

### export

Export rules from AGENTS.md to a specific agent's format:

```bash
coding-context export <agent>
```

Available agents: Claude, Gemini, Codex, Cursor, Augment, GitHubCopilot, Windsurf, Goose

Examples:
```bash
coding-context export Gemini    # Creates GEMINI.md
coding-context export Claude    # Creates CLAUDE.md
coding-context export Windsurf  # Creates .windsurf/rules/AGENTS.md
```

### bootstrap

Create the initial directory structure and files:

```bash
coding-context bootstrap
```

This creates:
- `AGENTS.md` - Main rules file
- `.prompts/` - Project-level rules directory
- `~/.prompts/rules/` - User-level rules directory

### prompt

Print the aggregated prompt from all rule files:

```bash
coding-context prompt
```

This displays the content of all rule files in the normalized format.

## Rule Hierarchy

Rules are organized by priority level:

1. **Project Rules (Level 0)** - Most important
   - `.prompts/` directory
   - Project-specific files like `.github/agents/`

2. **Ancestor Rules (Level 1)** - Next priority
   - Files in parent directories (e.g., AGENTS.md, CLAUDE.md)

3. **User Rules (Level 2)** - User-specific
   - `~/.prompts/rules/`
   - `~/.claude/CLAUDE.md`
   - `~/.gemini/GEMINI.md`

4. **System Rules (Level 3)** - System-wide
   - `/usr/local/prompts-rules`

## Agent-Specific Formats

### Claude
- Global: `~/.claude/CLAUDE.md`
- Ancestor: `./CLAUDE.md` (checked into Git)
- Local: `./CLAUDE.local.md` (highest precedence, typically .gitignore'd)

### Gemini
- Global: `~/.gemini/GEMINI.md`
- Ancestor: `./GEMINI.md`
- Project: `./.gemini/styleguide.md`

### Cursor
- Project: `./.cursor/rules/*.md` and `./.cursor/rules/*.mdc`
- Compatibility: `./AGENTS.md`

### GitHub Copilot
- Repository: `./.github/copilot-instructions.md`
- Agent Config: `.github/agents/`
- Compatibility: `./AGENTS.md`

### Windsurf
- Project/Ancestor: `./.windsurf/rules/*.md`

## Contributing

Contributions are welcome! Please see our [contributing guidelines](CONTRIBUTING.md) for more information.

## License

[MIT License](LICENSE)
