package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/scottrigby/newline/newline"
)

func main() {
	var includeHidden bool
	flag.BoolVar(&includeHidden, "hidden", false, "process files in hidden directories recursively")
	flag.BoolVar(&includeHidden, "a", false, "process files in hidden directories recursively (short form)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [target]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nEnsures files end with exactly one newline.\n\n")
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "  -a, --hidden\tProcess files in hidden directories recursively\n")
		fmt.Fprintf(os.Stderr, "\nARGUMENTS:\n")
		fmt.Fprintf(os.Stderr, "  target\tFile or directory to process (default: current directory)\n\n")
		fmt.Fprintf(os.Stderr, "\nBEHAVIOR:\n")
		fmt.Fprintf(os.Stderr, "  If target is a single file: Rewrites the file to end with exactly one newline\n\n")
		fmt.Fprintf(os.Stderr, "  If target is a directory: Processes all files recursively in the directory\n\n")
		fmt.Fprintf(os.Stderr, "  Hidden directories are skipped by default unless:\n")
		fmt.Fprintf(os.Stderr, "  - The target directory itself is hidden, OR\n")
		fmt.Fprintf(os.Stderr, "  - The --hidden flag is used\n\n")
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "  %s file.txt          # Process a single file\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s src/              # Process all files in src/ (skip hidden dirs)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --hidden .        # Process all files including hidden dirs\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s .git/             # Process hidden directory (allowed)\n", os.Args[0])
	}
	flag.Parse()

	target := "."
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}
	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Error: too many arguments\n")
		flag.Usage()
		os.Exit(1)
	}

	if err := newline.ProcessWithOptions(target, newline.Options{IncludeHidden: includeHidden}); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
