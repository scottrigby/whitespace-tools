package whitespace

import (
	"os"
	"regexp"
	"strings"
)

// removeTrailingWhitespace removes trailing spaces and tabs from each line in a file
func removeTrailingWhitespace(path string) error {
	input, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Use regexp to remove trailing spaces and tabs from each line
	trailingWhitespace := regexp.MustCompile(`[ \t]+$`)
	lines := strings.Split(string(input), "\n")

	// Process each line except handle the last one carefully to preserve EOF newlines
	for i, line := range lines {
		lines[i] = trailingWhitespace.ReplaceAllString(line, "")
	}

	output := strings.Join(lines, "\n")

	// Write back to file
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(output), info.Mode())
}

// ProcessTrailingspace processes a file or directory to remove trailing whitespace with default options.
func ProcessTrailingspace(target string) error {
	return ProcessTrailingspaceWithOptions(target, Options{})
}

// ProcessTrailingspaceWithOptions processes a file or directory to remove trailing whitespace with the given options.
func ProcessTrailingspaceWithOptions(target string, opts Options) error {
	return processTarget(target, opts, removeTrailingWhitespace)
}
