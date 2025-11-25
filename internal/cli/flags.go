package cli

import (
	"flag"
	"fmt"
	"os"
)

// ArrayFlags implements flag.Value for multiple string values
type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return "array flags"
}

func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// CommonFlags holds shared CLI flags and provides common setup
type CommonFlags struct {
	IncludeHidden   bool
	ExcludePatterns ArrayFlags
	ShowVersion     bool
}

// SetupFlags sets up the standard flags for both tools
func (cf *CommonFlags) SetupFlags() {
	flag.BoolVar(&cf.IncludeHidden, "include-hidden", false, "process files in hidden directories recursively")
	flag.BoolVar(&cf.IncludeHidden, "i", false, "process files in hidden directories recursively (short form)")
	flag.Var(&cf.ExcludePatterns, "exclude", "exclude files/directories matching glob pattern (can be used multiple times)")
	flag.Var(&cf.ExcludePatterns, "e", "exclude files/directories matching glob pattern (short form)")
	flag.BoolVar(&cf.ShowVersion, "version", false, "show version information")
	flag.BoolVar(&cf.ShowVersion, "v", false, "show version information (short form)")
}

// SetupUsage sets up the standard usage function for both tools
func SetupUsage(description string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [target]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n%s\n\n", description)
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "  -i, --include-hidden\t\tProcess files in hidden directories recursively\n")
		fmt.Fprintf(os.Stderr, "  -e, --exclude PATTERN\t\tExclude files/directories matching glob pattern\n")
		fmt.Fprintf(os.Stderr, "  -v, --version\t\t\tShow version information\n")
		fmt.Fprintf(os.Stderr, "\nARGUMENTS:\n")
		fmt.Fprintf(os.Stderr, "  target\t\t\tFile or directory to process (default: current directory)\n\n")
		fmt.Fprintf(os.Stderr, "\nBEHAVIOR:\n")
		fmt.Fprintf(os.Stderr, "  Processes all text files recursively, skipping:\n")
		fmt.Fprintf(os.Stderr, "  • Hidden directories (unless --include-hidden used)\n")
		fmt.Fprintf(os.Stderr, "  • Non-text files (detected by heuristic)\n")
		fmt.Fprintf(os.Stderr, "  • Files/directories matching --exclude patterns\n\n")
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "  %s file.txt\t\t\t\t# Process single file\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s src/\t\t\t\t\t# Process all text files in src/\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --include-hidden\t\t\t# Include hidden directories\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --exclude 'bin' --exclude '*.tmp'\t# Exclude patterns\n", os.Args[0])
	}
}

// HandleVersion displays version information if requested
func HandleVersion(showVersion bool, toolName, version, commit string) bool {
	if !showVersion {
		return false
	}

	if version == "" {
		version = "dev"
	}
	if commit == "" {
		commit = "unknown"
	}
	fmt.Printf("%s %s (commit: %s)\n", toolName, version, commit)
	return true
}

// ParseTarget handles target argument parsing with validation
func ParseTarget() (string, error) {
	target := "."
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}
	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Error: too many arguments\n")
		flag.Usage()
		os.Exit(1)
	}
	return target, nil
}
