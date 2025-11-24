package whitespace

import (
	"bufio"
	"io"
	"os"
	"unicode/utf8"
)

// LooksText reports whether the file appears to be text based on a small prefix.
// - Reads up to maxBytes (default 4096).
// - Returns false for NUL bytes or high ratio of non-text bytes.
// - Directories are false. Symlinks are followed once.
func LooksText(path string) (bool, error) {
	const maxBytes = 4096

	fi, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	if fi.IsDir() {
		return false, nil
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		t, err := os.Readlink(path)
		if err != nil {
			return false, err
		}
		path = t
	}

	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	buf, err := r.Peek(maxBytes)
	if err != nil && err != bufio.ErrBufferFull && err != io.EOF {
		// For short files, EOF is expected; ErrBufferFull means we read maxBytes.
		return false, err
	}
	if len(buf) == 0 {
		// Empty file: treat as text (safe to add newline/trim).
		return true, nil
	}

	// Immediate binary indicator: any NUL byte.
	for _, b := range buf {
		if b == 0x00 {
			return false, nil
		}
	}

	// Count suspicious bytes: invalid UTF-8 or disallowed controls.
	var nonText int
	for i := 0; i < len(buf); {
		b := buf[i]
		if b < 0x80 {
			switch b {
			case '\n', '\r', '\t':
				// common controls allowed
			default:
				if b < 0x20 && b != '\n' && b != '\r' && b != '\t' {
					nonText++
				}
			}
			i++
			continue
		}
		_, size := utf8.DecodeRune(buf[i:])
		if size == 1 && b >= 0x80 {
			nonText++
			i++
		} else {
			i += size
		}
	}

	// If more than ~30% of sampled bytes are non-text, treat as binary.
	return float64(nonText) <= 0.30*float64(len(buf)), nil
}
