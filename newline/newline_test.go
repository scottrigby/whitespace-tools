package newline

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcess(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testnewline")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testDirs := []string{
		"testdir",
		"testdir/.hiddendir",
		"testdir/subdir",
		"testdir/.hiddendir/.deephidden",
		"testdir/subdir/.subhidden",
	}
	for _, d := range testDirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, d), 0o755); err != nil {
			t.Fatal(err)
		}
		files := map[string][]byte{
			"empty.txt":           {},
			"newline.txt":         []byte("\n"),
			"extra-newline.txt":   []byte("\n\n"),
			"content.txt":         []byte("content"),
			"content-newline.txt": []byte("content\n"),
		}
		for name, content := range files {
			if err := os.WriteFile(filepath.Join(tmpDir, d, name), content, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}

	// Run on testdir and testdir/.hiddendir
	if err := Process(filepath.Join(tmpDir, "testdir")); err != nil {
		t.Fatal(err)
	}
	if err := Process(filepath.Join(tmpDir, "testdir/.hiddendir")); err != nil {
		t.Fatal(err)
	}

	// Helper to count lines in all files in a dir (like wc -l)
	countLines := func(dir string) map[string]int {
		result := map[string]int{}
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			content, _ := os.ReadFile(filepath.Join(dir, e.Name()))
			lines := 0
			for _, b := range content {
				if b == '\n' {
					lines++
				}
			}
			result[e.Name()] = lines
		}
		return result
	}

	// Check root of testdir
	lines := countLines(filepath.Join(tmpDir, "testdir"))
	for _, n := range lines {
		if n != 1 {
			t.Errorf("testdir: expected 1 line, got %d", n)
		}
	}

	// Check root of testdir/subdir
	lines = countLines(filepath.Join(tmpDir, "testdir/subdir"))
	for _, n := range lines {
		if n != 1 {
			t.Errorf("testdir/subdir: expected 1 line, got %d", n)
		}
	}

	// Check root of testdir/.hiddendir
	lines = countLines(filepath.Join(tmpDir, "testdir/.hiddendir"))
	for _, n := range lines {
		if n != 1 {
			t.Errorf("testdir/.hiddendir: expected 1 line, got %d", n)
		}
	}

	// Check testdir/.hiddendir/.deephidden (should NOT be processed)
	lines = countLines(filepath.Join(tmpDir, "testdir/.hiddendir/.deephidden"))
	expected := map[string]int{
		"empty.txt":           0,
		"newline.txt":         1,
		"extra-newline.txt":   2,
		"content.txt":         0,
		"content-newline.txt": 1,
	}
	for name, want := range expected {
		if got := lines[name]; got != want {
			t.Errorf("testdir/.hiddendir/.deephidden: %s expected %d lines, got %d", name, want, got)
		}
	}

	// Check testdir/subdir/.subhidden (should NOT be processed)
	lines = countLines(filepath.Join(tmpDir, "testdir/subdir/.subhidden"))
	for name, want := range expected {
		if got := lines[name]; got != want {
			t.Errorf("testdir/subdir/.subhidden: %s expected %d lines, got %d", name, want, got)
		}
	}
}
