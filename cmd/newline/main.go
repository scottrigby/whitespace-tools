package main

import (
	"fmt"
	"os"

	"github.com/scottrigby/newline/newline"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [target]\n", os.Args[0])
		os.Exit(1)
	}
	target := os.Args[1]
	if err := newline.Process(target); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
