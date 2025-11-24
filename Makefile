.PHONY: test build clean install

# Build the binary
build:
	go build -o bin/newline ./cmd/newline

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

# Default target
all: check build
