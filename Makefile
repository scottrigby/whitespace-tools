.PHONY: test build build-full build-tiny clean install fmt vet check all default release-snapshot release

# Build all binaries (prefers TinyGo if available)
build:
	@if command -v tinygo >/dev/null 2>&1; then \
		tinygo build -o bin/newline ./cmd/newline; \
		tinygo build -o bin/trailingspace ./cmd/trailingspace; \
		echo "TinyGo binaries created: bin/newline bin/trailingspace"; \
	else \
		echo "TinyGo not found, using standard Go build"; \
		$(MAKE) build-full; \
	fi

# Build with standard Go (larger but full compatibility)
build-full:
	CGO_ENABLED=0 go build -ldflags="-s -w -buildid=" -o bin/newline ./cmd/newline
	CGO_ENABLED=0 go build -ldflags="-s -w -buildid=" -o bin/trailingspace ./cmd/trailingspace

# Build with tinygo (much smaller binaries)
build-tiny:
	@if command -v tinygo >/dev/null 2>&1; then \
		tinygo build -o bin/newline-tiny ./cmd/newline; \
		tinygo build -o bin/trailingspace-tiny ./cmd/trailingspace; \
		echo "TinyGo binaries created: bin/newline-tiny bin/trailingspace-tiny"; \
	else \
		echo "Error: TinyGo not found. Install from https://tinygo.org/getting-started/install/"; \
		exit 1; \
	fi

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Install the binary to GOPATH/bin
install:
	go install ./cmd/newline


# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Check code with go vet
vet:
	go vet ./...

# Run all checks (fmt, vet, test)
check: fmt vet test

# GoReleaser targets
release-snapshot:
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean; \
	else \
		echo "Error: GoReleaser not found. Install from https://goreleaser.com/install/"; \
		exit 1; \
	fi

release:
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release; \
	else \
		echo "GoReleaser not found. Install from https://goreleaser.com/install/"; \
		exit 1; \
	fi

# Default target
all: check build

# Default build target
default: build
