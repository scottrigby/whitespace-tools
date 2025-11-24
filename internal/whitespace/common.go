package whitespace

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

// Options for processing files
type Options struct {
	IncludeHidden   bool
	ExcludePatterns []string    // Glob patterns to exclude
	compiledGlobs   []glob.Glob // Compiled glob patterns (internal use)
}

// isHidden returns true if the file/directory name starts with a dot
func isHidden(name string) bool {
	base := filepath.Base(name)
	return strings.HasPrefix(base, ".")
}

// compileExcludePatterns compiles glob patterns for efficient matching
func compileExcludePatterns(opts *Options) error {
	if opts.compiledGlobs == nil && len(opts.ExcludePatterns) > 0 {
		opts.compiledGlobs = make([]glob.Glob, 0, len(opts.ExcludePatterns))
		for _, pattern := range opts.ExcludePatterns {
			g, err := glob.Compile(pattern)
			if err != nil {
				return err
			}
			opts.compiledGlobs = append(opts.compiledGlobs, g)
		}
	}
	return nil
}

// shouldExcludePath returns true if the path matches any exclude pattern
func shouldExcludePath(path string, opts *Options) bool {
	base := filepath.Base(path)
	for _, g := range opts.compiledGlobs {
		if g.Match(path) || g.Match(base) {
			return true
		}
	}
	return false
}

// ProcessFileFunc is a function type for processing individual files
type ProcessFileFunc func(path string) error

// processDir processes all files in a directory with the given options and file processor
func processDir(dir string, opts Options, processFile ProcessFileFunc) error {
	// Compile exclude patterns once
	if err := compileExcludePatterns(&opts); err != nil {
		return err
	}

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If this is a directory, check if we should skip it
		if d.IsDir() {
			// Don't skip the root directory
			if path == dir {
				return nil
			}

			// Check exclude patterns
			if shouldExcludePath(path, &opts) {
				return filepath.SkipDir
			}

			// Skip hidden directories based on options
			if isHidden(d.Name()) {
				// If IncludeHidden is true, never skip hidden dirs
				if opts.IncludeHidden {
					return nil
				}
				// If IncludeHidden is false (default), always skip hidden subdirectories
				return filepath.SkipDir
			}

			return nil
		}

		// Check exclude patterns for files
		if shouldExcludePath(path, &opts) {
			return nil
		}

		// Skip non-text files
		isText, err := LooksText(path)
		if err != nil {
			return err
		}
		if !isText {
			return nil
		}

		// Process the file
		return processFile(path)
	})
}

// processTarget processes a file or directory target with the given options
func processTarget(target string, opts Options, processFile ProcessFileFunc) error {
	// Compile exclude patterns once
	if err := compileExcludePatterns(&opts); err != nil {
		return err
	}

	info, err := os.Stat(target)
	if err != nil {
		return err
	}
	if info.Mode().IsRegular() {
		return processFile(target)
	}
	if info.IsDir() {
		return processDir(target, opts, processFile)
	}
	return errors.New("not a file or directory: " + target)
}