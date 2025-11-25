package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/scottrigby/whitespace-tools/internal/cli"
	"github.com/scottrigby/whitespace-tools/internal/whitespace"
)

var (
	// injected by ldflags:
	// -X main.version
	// -X main.commit
	version string
	commit  string
)

func main() {
	var flags cli.CommonFlags

	cli.SetupUsage("Removes trailing whitespace from end of lines.")
	flags.SetupFlags()
	flag.Parse()

	// Handle version flag
	if cli.HandleVersion(flags.ShowVersion, "trailingspace", version, commit) {
		return
	}

	target, _ := cli.ParseTarget()

	opts := whitespace.Options{
		IncludeHidden:   flags.IncludeHidden,
		ExcludePatterns: []string(flags.ExcludePatterns),
	}

	if err := whitespace.ProcessTrailingspaceWithOptions(target, opts); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
