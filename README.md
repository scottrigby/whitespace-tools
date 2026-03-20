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

## Devcontainer Feature

Add to a claudeman profile or any devcontainer:

```json
{
  "features": {
    "ghcr.io/scottrigby/features/whitespace-tools:1": {}
  }
}
```

Or pin a specific version:

```json
{
  "features": {
    "ghcr.io/scottrigby/features/whitespace-tools:1": {
      "version": "v1.0.1"
    }
  }
}
```

### Publishing (manual steps)

The feature source lives in `src/whitespace-tools/`. Two workflows handle publishing:

1. **`goreleaser.yml`** — Builds binaries and creates GitHub Releases (triggered on `v*.*.*` tags).
   The install.sh downloads from GitHub Releases, so releases must exist before the feature works.

2. **`release.yml`** — Publishes the devcontainer feature to `ghcr.io/scottrigby/features/whitespace-tools`
   (manually triggered via workflow_dispatch on the main branch).

Steps to publish a new version:

1. Set up a fine-grained GitHub PAT with `Contents: write` on the tap repo
   (`scottrigby/homebrew-tap`) and add it as a repo secret named `GH_PAT`.

2. Push code to `main`, then create a semver tag to trigger GoReleaser:

   ```bash
   git tag v1.0.1
   git push origin v1.0.1
   ```

   This builds binaries and publishes GitHub Releases + Homebrew tap.

3. After the release exists, trigger the devcontainer feature publish manually:
   - GitHub → Actions → "Release dev container features" → Run workflow (main branch)
   - This pushes the feature to `ghcr.io/scottrigby/features/whitespace-tools`

4. The feature is then usable in devcontainer configs and claudeman profiles.
