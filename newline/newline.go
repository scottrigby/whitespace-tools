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

// Process processes a file or directory as described.
func Process(target string) error {
	info, err := os.Stat(target)
	if err != nil {
		return err
	}
	if info.Mode().IsRegular() {
		return ensureSingleNewline(target)
	}
	if info.IsDir() {
		return processDir(target)
	}
	return errors.New("not a file or directory: " + target)
}

// processDir processes all files in dir, skipping hidden subdirs unless dir itself is hidden.
func processDir(dir string) error {
	if isHidden(dir) {
		return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			return ensureSingleNewline(path)
		})
	}
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		if d.IsDir() && rel != "." && isHidden(d.Name()) {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			return ensureSingleNewline(path)
		}
		return nil
	})
}
