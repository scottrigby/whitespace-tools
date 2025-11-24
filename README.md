# Whitespace Tools

CLI tools to fix whitespace issues that correspond to Git's `core.whitespace` options:

- **`newline`** - Ensures files end with exactly one newline (blank-at-eof)
- **`trailingspace`** - Removes trailing spaces/tabs from lines (blank-at-eol)

Perfect for AI agents that consistently get whitespace wrong.

## Installation

```bash
# Go install
go install github.com/scottrigby/whitespace-tools/cmd/newline@latest
go install github.com/scottrigby/whitespace-tools/cmd/trailingspace@latest

# Or download from releases
# https://github.com/scottrigby/whitespace-tools/releases
```

## Usage

```bash
# Fix newlines
newline file.txt        # Single file
newline .              # Current directory
newline --hidden .     # Include hidden directories

# Remove trailing spaces
trailingspace file.txt  # Single file
trailingspace src/     # Directory
```

Both tools skip:
- Hidden directories (unless `--hidden` flag used or target itself is hidden)
- Build directories (`bin/`, `dist/`, `node_modules/`, `vendor/`, `target/`, `.git/`)
- Binary files (detected by null bytes)

## Building

```bash
make build       # TinyGo (~400K) if available, otherwise Go (~1.6M)
make build-full  # Force standard Go build
make build-tiny  # Force TinyGo build
```

## Setup for Maintainers

### GitHub Token for Releases

1. Go to GitHub Settings → Developer settings → Personal access tokens → Fine-grained tokens
2. Click "Generate new token"
3. Configure:
   - **Repository access**: Selected repositories → `scottrigby/whitespace-tools`, `scottrigby/homebrew-tap`
   - **Repository permissions**:
     - Contents: Read and write (for releases)
     - Metadata: Read (required)
     - Pull requests: Write (for Homebrew tap updates)
4. Copy the token and set environment variable:
   ```bash
   export GITHUB_TOKEN="your_token_here"
   ```

### Creating Releases

```bash
git tag v1.0.0
make release  # Creates GitHub release + updates Homebrew tap
```
