package newline

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ensureSingleNewline rewrites the file so it ends with exactly one newline.
func ensureSingleNewline(path string) error {
	input, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	input = bytes.TrimRight(input, "\r\n")
	input = append(input, '\n')
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, input, info.Mode()); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewReader(input))
	os.Remove(tmp)
	return err
}

func isHidden(name string) bool {
	base := filepath.Base(name)
	return strings.HasPrefix(base, ".")
}

// Options for processing files
type Options struct {
	IncludeHidden bool
}

// Process processes a file or directory with default options.
func Process(target string) error {
	return ProcessWithOptions(target, Options{})
}

// ProcessWithOptions processes a file or directory with the given options.
func ProcessWithOptions(target string, opts Options) error {
	info, err := os.Stat(target)
	if err != nil {
		return err
	}
	if info.Mode().IsRegular() {
		return ensureSingleNewline(target)
	}
	if info.IsDir() {
		return processDirWithOptions(target, opts)
	}
	return errors.New("not a file or directory: " + target)
}

// processDir processes all files in dir, skipping hidden subdirs unless dir itself is hidden.
func processDir(dir string) error {
	return processDirWithOptions(dir, Options{})
}

// processDirWithOptions processes all files in dir with the given options.
func processDirWithOptions(dir string, opts Options) error {
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

		// Process the file
		return ensureSingleNewline(path)
	})
}
