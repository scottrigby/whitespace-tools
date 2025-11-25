# Whitespace Tools - CLI Utilities for Git Whitespace Standards

## Overview

Two complementary CLI tools that help maintain Git's default whitespace standards:
- `newline` - Ensures files end with exactly one newline (fixes `blank-at-eof`)
- `trailingspace` - Removes trailing whitespace from lines (fixes `blank-at-eol`)

## Why These Tools Matter

Git's default `core.whitespace` setting flags two common whitespace issues:
- `blank-at-eof`: Extra blank lines at end of file
- `blank-at-eol`: Trailing spaces/tabs at end of lines

These tools fix both issues automatically, making code cleaner and preventing Git whitespace warnings.

## Problems These Tools Solve

For AI Agents:
- AI code generation often creates files without proper EOF newlines
- AI agents frequently leave trailing whitespace, even when explicitly asked not to
- Manual cleanup is tedious and error-prone

For Developers:
- Inconsistent whitespace handling across team members
- Git showing whitespace warnings during commits
- Need for automated cleanup that's predictable and safe

## Tool Capabilities

### newline
```bash
newline [OPTIONS] [target]
```

Purpose: Ensures files end with exactly one newline character

Key Features:
- Handles empty files, content-only files, files with multiple newlines
- Preserves all file content, only modifies EOF
- Supports CRLF, LF, and mixed line endings
- Safe atomic file updates

### trailingspace
```bash
trailingspace [OPTIONS] [target]
```

Purpose: Removes trailing spaces and tabs from end of each line

Key Features:
- Removes spaces and tabs from line endings
- Preserves leading/internal whitespace
- Handles lines with only whitespace (converts to empty lines)
- Safe atomic file updates

## Shared Options

Both tools support identical command-line options:

```bash
-i, --include-hidden    Process files in hidden directories recursively
-e, --exclude PATTERN   Exclude files/directories matching glob pattern
-v, --version           Show version information
```

## Directory Processing

Default Behavior:
- Processes current directory recursively if no target specified
- Skips hidden subdirectories (`.git`, `.vscode`, etc.)
- Automatically detects and skips binary files
- Only processes text files

Hidden Directory Control:
```bash
# Skip hidden dirs (default)
newline src/
trailingspace src/

# Include hidden dirs
newline --include-hidden .
trailingspace --include-hidden .
```

Exclusion Patterns:
```bash
# Exclude specific patterns
newline --exclude 'bin' --exclude '*.tmp' --exclude 'node_modules'
trailingspace --exclude '*.log' --exclude 'dist'
```

## Usage Examples

### Single File Processing
```bash
newline file.txt
trailingspace script.py
```

### Directory Processing
```bash
# Process all text files in src/
newline src/
trailingspace src/

# Process current directory, including hidden dirs
newline --include-hidden .
trailingspace --include-hidden .

# Process with exclusions
newline --exclude 'bin' --exclude '*.tmp'
trailingspace --exclude 'build' --exclude '*.log'
```

### Combined Workflow
```bash
# Clean up whitespace issues in a project
trailingspace           # Remove trailing whitespace
newline                 # Ensure proper EOF newlines
```

## AI Agent Integration

Recommended AI Agent Rules:

1. After Code Generation:
   - Run `trailingspace` on all generated files
   - Run `newline` on all generated files

2. Project Cleanup:
   - Use exclusion patterns to avoid build artifacts: `--exclude 'bin' --exclude 'dist' --exclude 'node_modules'`
   - Include hidden directories only when necessary: `--include-hidden`

3. Safe Usage:
   - Both tools only modify whitespace, never code content
   - Both tools automatically skip binary files
   - Both tools use atomic file updates

Example AI Prompts:
```
"Generate the code and then run trailingspace and newline on all created files"
"Clean up whitespace in this project using: trailingspace --exclude 'node_modules' && newline --exclude 'node_modules'"
```

## Technical Details

Architecture:
- Go monorepo with shared code
- Optimized builds with TinyGo when available
- Cross-platform releases via GoReleaser
- Comprehensive test coverage

Safety Features:
- Heuristic-based text file detection
- Atomic file updates (write to temp, then rename)
- Preserves file permissions
- No changes to file content beyond whitespace

Performance:
- Binary sizes: ~600KB (TinyGo) vs ~1.8MB (standard Go)
- Fast recursive directory processing
- Efficient pattern matching with compiled glob patterns

## Installation

```bash
# Via release binaries
curl -sLO https://github.com/scottrigby/whitespace-tools/releases/latest/download/whitespace-tools_linux_amd64.tar.gz

# From source
git clone https://github.com/scottrigby/whitespace-tools
cd whitespace-tools
make build
```

Both tools are designed to be safe, predictable, and perfect for automated use by AI agents or in CI/CD pipelines.
