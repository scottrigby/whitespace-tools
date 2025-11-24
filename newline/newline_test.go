package newline

import (
	"os"
	"path/filepath"
	"testing"
)

// countLines counts newline characters in files within a directory (like wc -l)
func countLines(t *testing.T, dir string) map[string]int {
	t.Helper()
	result := map[string]int{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read directory %s: %v", dir, err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("failed to read file %s: %v", e.Name(), err)
		}
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

// createTestDirStructure creates a test directory structure with sample files
func createTestDirStructure(t *testing.T, tmpDir string) {
	t.Helper()
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
}

func TestProcess(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testnewline")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	createTestDirStructure(t, tmpDir)

	// Run on testdir and testdir/.hiddendir
	if err := Process(filepath.Join(tmpDir, "testdir")); err != nil {
		t.Fatal(err)
	}
	if err := Process(filepath.Join(tmpDir, "testdir/.hiddendir")); err != nil {
		t.Fatal(err)
	}

	// Check root of testdir
	lines := countLines(t, filepath.Join(tmpDir, "testdir"))
	for _, n := range lines {
		if n != 1 {
			t.Errorf("testdir: expected 1 line, got %d", n)
		}
	}

	// Check root of testdir/subdir
	lines = countLines(t, filepath.Join(tmpDir, "testdir/subdir"))
	for _, n := range lines {
		if n != 1 {
			t.Errorf("testdir/subdir: expected 1 line, got %d", n)
		}
	}

	// Check root of testdir/.hiddendir
	lines = countLines(t, filepath.Join(tmpDir, "testdir/.hiddendir"))
	for _, n := range lines {
		if n != 1 {
			t.Errorf("testdir/.hiddendir: expected 1 line, got %d", n)
		}
	}

	// Check testdir/.hiddendir/.deephidden (should NOT be processed)
	lines = countLines(t, filepath.Join(tmpDir, "testdir/.hiddendir/.deephidden"))
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
	lines = countLines(t, filepath.Join(tmpDir, "testdir/subdir/.subhidden"))
	for name, want := range expected {
		if got := lines[name]; got != want {
			t.Errorf("testdir/subdir/.subhidden: %s expected %d lines, got %d", name, want, got)
		}
	}
}

func TestProcessWithHiddenOption(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testnewlinehidden")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	createTestDirStructure(t, tmpDir)

	// Run on testdir with IncludeHidden=true
	if err := ProcessWithOptions(filepath.Join(tmpDir, "testdir"), Options{IncludeHidden: true}); err != nil {
		t.Fatal(err)
	}

	// With IncludeHidden=true, ALL directories should be processed
	// So all files should have exactly 1 line
	testDirs := []string{
		"testdir",
		"testdir/.hiddendir",
		"testdir/subdir",
		"testdir/.hiddendir/.deephidden",
		"testdir/subdir/.subhidden",
	}

	for _, dir := range testDirs {
		lines := countLines(t, filepath.Join(tmpDir, dir))
		for name, got := range lines {
			if got != 1 {
				t.Errorf("%s: %s expected 1 line (processed), got %d", dir, name, got)
			}
		}
	}
}
