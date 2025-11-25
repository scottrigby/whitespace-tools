# Whitespace Tools - Development & Testing Guide

## Requirements

### newline CLI tool

Problem: AI agents often do not ensure a single newline at EOF, even when asked.

Why this matters: a single newline at EOF corresponds to Git's default core.whitespace blank-at-eof option.

Tool requirements:
- The ability to rewrite a file so it ends with exactly one newline
- The ability to process a file or all files recursively in an entire directory
- When processing all files in dir, skip hidden subdirs unless dir itself is hidden, or unless the option to process all hidden dirs is explicitly opted-into

### trailingspace CLI tool

Problem: AI agents often leave trailing whitespace at the end of lines, even when asked not to.

Why this matters: trailing whitespace at end-of-line corresponds to Git's default core.whitespace blank-at-eol option.

Tool requirements:
- The ability to rewrite files to remove trailing spaces and tabs from the end of each line
- The ability to process a file or all files recursively in an entire directory
- When processing all files in dir, skip hidden subdirs unless dir itself is hidden, or unless the option to process all hidden dirs is explicitly opted-into
- Should preserve file content otherwise (no changes to line endings or other whitespace)

### Shared Requirements

- Part of monorepo with shared architecture and patterns
- Support same `--include-hidden` flag for processing hidden directories
- Support `--exclude` patterns for glob-based exclusion
- Both tools released together, installable as a suite
- TinyGo optimization for smaller binaries when available

## Development

### Project Structure
```
├── cmd/                    # CLI entry points
│   ├── newline/
│   └── trailingspace/
├── internal/               # Shared libraries
│   ├── whitespace/         # Core processing logic
│   └── cli/                # CLI utilities
├── .goreleaser.yml         # Cross-platform releases
└── Makefile                # Build automation
```

### Build System

```bash
# Run tests
make test

# Build binaries (uses TinyGo if available, falls back to Go)
make build

# Build with standard Go
make build-full

# Build with TinyGo (smaller binaries)
make build-tiny

# Cross-platform development build
make release-snapshot

# Format and check code
make check
```

### Version and Environment Variables

The Makefile handles version injection:
- `VERSION`: Git tag or "dev" fallback
- `COMMIT`: Git commit hash or "unknown" fallback
- `GORELEASER_PARALLELISM`: Configurable parallelism (default: 4)

Both tools support `--version` flag showing version and commit information.

## Testing

### Test Architecture

Tests are organized into three categories:

1. **File Selection Tests** (`common_test.go`):
   - Single file vs directory processing
   - Hidden directory handling (default skip vs --include-hidden)
   - Exclusion pattern matching (--exclude patterns)
   - Binary file detection and skipping

2. **Newline EOF Tests** (`newline_test.go`):
   - Empty files, content without newlines, multiple trailing newlines
   - CRLF/mixed line endings, files ending with CR only
   - Multiline content variations
   - Edge cases and integration tests

3. **Trailing Whitespace EOL Tests** (`trailingspace_test.go`):
   - Single/multiple spaces and tabs at end of lines with content
   - Lines with only whitespace (should become empty lines)
   - Whitespace at EOF (with/without newlines)
   - Preserved: leading/internal spaces and tabs
   - Complex scenarios and edge cases

### Running Tests

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run specific test categories
go test ./internal/whitespace -run TestFileSelection
go test ./internal/whitespace -run TestEnsureSingleNewline
go test ./internal/whitespace -run TestRemoveTrailingWhitespace
```

### Test Coverage

Tests cover:
- All CLI flag combinations
- Directory traversal logic with various exclusion patterns
- Text vs binary file detection
- Atomic file update safety
- Error handling and edge cases
- Cross-platform compatibility (line endings, file permissions)

### Reference Implementation

Original zsh implementations are available in `.zshrc` for reference:
- `newline`: Ensures single newline at EOF
- `trailingspace`: Uses `perl -pe 's/[ \t]+$//'` to remove trailing whitespace
- `newline_test`: Test suite for validation

The Go implementation should be at least as robust as these shell commands while providing better portability and performance.

### Process

Running Claude in a container, using this simple project: https://github.com/scottrigby/claude-container

Enabling YOLO mode.

Work until either clarification is needed, or tasks are complete.

**Important Development Rules:**
1. **Test your work**: Run `make test` after every change to ensure tests pass
2. **Validate at each step**: Ensure tests meet requirements before proceeding
3. **Use the tools**: Always run `trailingspace .` and `newline .` on the codebase after making changes (use the last known working version of the tools; if unsure which version to use, ask for clarification)
4. **No broken builds**: Never leave the codebase in a state where tests fail or builds break

A script in your PATH is provided called `notify.sh`.

When clarification is needed: Send notification when clarification needed: `notify.sh "Need clarification: reason"`

For task completion: Send audio notification: `notify.sh "completion message"`
