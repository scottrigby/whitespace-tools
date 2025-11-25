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
	var includeHidden bool
	var excludePatterns cli.ArrayFlags
	var showVersion bool
	flag.BoolVar(&includeHidden, "include-hidden", false, "process files in hidden directories recursively")
	flag.BoolVar(&includeHidden, "i", false, "process files in hidden directories recursively (short form)")
	flag.Var(&excludePatterns, "exclude", "exclude files/directories matching glob pattern (can be used multiple times)")
	flag.Var(&excludePatterns, "e", "exclude files/directories matching glob pattern (short form)")
	flag.BoolVar(&showVersion, "version", false, "show version information")
	flag.BoolVar(&showVersion, "v", false, "show version information (short form)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [target]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nEnsures files end with exactly one newline.\n\n")
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "  -i, --include-hidden\tProcess files in hidden directories recursively\n")
		fmt.Fprintf(os.Stderr, "  -e, --exclude PATTERN\tExclude files/directories matching glob pattern\n")
		fmt.Fprintf(os.Stderr, "  -v, --version\t\tShow version information\n")
		fmt.Fprintf(os.Stderr, "\nARGUMENTS:\n")
		fmt.Fprintf(os.Stderr, "  target\tFile or directory to process (default: current directory)\n\n")
		fmt.Fprintf(os.Stderr, "\nBEHAVIOR:\n")
		fmt.Fprintf(os.Stderr, "  Processes all text files recursively, skipping:\n")
		fmt.Fprintf(os.Stderr, "  • Hidden directories (unless --include-hidden used)\n")
		fmt.Fprintf(os.Stderr, "  • Non-text files (detected by heuristic)\n")
		fmt.Fprintf(os.Stderr, "  • Files/directories matching --exclude patterns\n\n")
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "  %s file.txt                    # Process single file\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s src/                       # Process all text files in src/\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --include-hidden .         # Include hidden directories\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --exclude 'bin' --exclude '*.tmp' .  # Exclude patterns\n", os.Args[0])
	}
	flag.Parse()

	// Handle version flag
	if showVersion {
		if version == "" {
			version = "dev"
		}
		if commit == "" {
			commit = "unknown"
		}
		fmt.Printf("newline %s (commit: %s)\n", version, commit)
		return
	}

	target := "."
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}
	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Error: too many arguments\n")
		flag.Usage()
		os.Exit(1)
	}

	opts := whitespace.Options{
		IncludeHidden:   includeHidden,
		ExcludePatterns: []string(excludePatterns),
	}

	if err := whitespace.ProcessNewlineWithOptions(target, opts); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
