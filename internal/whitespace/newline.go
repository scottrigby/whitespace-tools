package whitespace

import (
	"bytes"
	"io"
	"os"
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

// ProcessNewline processes a file or directory with default options.
func ProcessNewline(target string) error {
	return ProcessNewlineWithOptions(target, Options{})
}

// ProcessNewlineWithOptions processes a file or directory with the given options.
func ProcessNewlineWithOptions(target string, opts Options) error {
	return processTarget(target, opts, ensureSingleNewline)
}
