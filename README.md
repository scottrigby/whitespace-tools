# Whitespace Tools

CLI tools to fix Git whitespace issues:

- `newline` - Ensures files end with exactly one newline
- `trailingspace` - Removes trailing whitespace from lines

## Usage

```bash
# Process current directory
newline .
trailingspace .

# Process single file
newline file.txt
trailingspace script.py

# Include hidden directories
newline --include-hidden .

# Exclude patterns
newline --exclude 'node_modules' --exclude '*.tmp' .
```

## Options

```bash
-i, --include-hidden    Process files in hidden directories
-e, --exclude PATTERN   Exclude files/directories matching pattern
-v, --version           Show version information
```

## Installation

```bash
# Homebrew (macOS)
brew install scottrigby/tap/whitespace-tools

# From release
curl -sLO https://github.com/scottrigby/whitespace-tools/releases/latest/download/whitespace-tools_linux_amd64.tar.gz

# From source
git clone https://github.com/scottrigby/whitespace-tools
cd whitespace-tools
make build
```
