package whitespace

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// mockProcessor records which files would be processed without modifying them
type mockProcessor struct {
	processedFiles []string
}

func (m *mockProcessor) process(path string) error {
	m.processedFiles = append(m.processedFiles, path)
	return nil
}

func (m *mockProcessor) reset() {
	m.processedFiles = nil
}

// createTestFileStructure creates a comprehensive test directory structure
func createTestFileStructure(t *testing.T, tmpDir string) {
	t.Helper()
	
	// Directory structure to create
	testDirs := []string{
		"testdir",
		"testdir/.hiddendir",
		"testdir/subdir", 
		"testdir/.hiddendir/.deephidden",
		"testdir/subdir/.subhidden",
		"testdir/bin",           // for exclude pattern testing
		"testdir/build",         // for exclude pattern testing
	}

	for _, d := range testDirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	// Files to create in each directory
	files := map[string][]byte{
		"text.txt":      []byte("text content\n"),
		"code.go":       []byte("package main\n"),
		"readme.md":     []byte("# README\n"), 
		"config.json":   []byte("{}\n"),
		"script.sh":     []byte("#!/bin/bash\n"),
		"data.tmp":      []byte("temp data\n"),    // for exclude testing
		"build.log":     []byte("build output\n"), // for exclude testing
		"binary":        []byte("\x00\x01\x02\x03"), // binary file (should be skipped)
	}

	// Create files in all directories
	for _, d := range testDirs {
		for name, content := range files {
			path := filepath.Join(tmpDir, d, name)
			if err := os.WriteFile(path, content, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestFileSelection_SingleFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testfileselection")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	createTestFileStructure(t, tmpDir)
	mock := &mockProcessor{}

	// Test processing single text file
	textFile := filepath.Join(tmpDir, "testdir", "text.txt")
	if err := processTarget(textFile, Options{}, mock.process); err != nil {
		t.Fatal(err)
	}

	expected := []string{textFile}
	if len(mock.processedFiles) != 1 || mock.processedFiles[0] != expected[0] {
		t.Errorf("expected %v, got %v", expected, mock.processedFiles)
	}

	// Test processing binary file (should be processed since we're targeting it directly)
	mock.reset()
	binaryFile := filepath.Join(tmpDir, "testdir", "binary")
	if err := processTarget(binaryFile, Options{}, mock.process); err != nil {
		t.Fatal(err)
	}

	if len(mock.processedFiles) != 1 || mock.processedFiles[0] != binaryFile {
		t.Errorf("direct file processing should work even for binary files, got %v", mock.processedFiles)
	}
}

func TestFileSelection_DirectoryDefault(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testfileselection")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	createTestFileStructure(t, tmpDir)
	mock := &mockProcessor{}

	// Process directory with default options (no hidden dirs)
	testDir := filepath.Join(tmpDir, "testdir")
	if err := processTarget(testDir, Options{}, mock.process); err != nil {
		t.Fatal(err)
	}

	// Should process files in testdir/ and testdir/subdir/ but NOT hidden directories
	// Should skip binary files automatically
	expectedDirs := []string{
		"testdir",
		"testdir/subdir",
		"testdir/bin", 
		"testdir/build",
	}

	expectedFiles := []string{
		"text.txt", "code.go", "readme.md", "config.json", "script.sh", "data.tmp", "build.log",
	}

	var expectedPaths []string
	for _, dir := range expectedDirs {
		for _, file := range expectedFiles {
			expectedPaths = append(expectedPaths, filepath.Join(tmpDir, dir, file))
		}
	}

	sort.Strings(expectedPaths)
	sort.Strings(mock.processedFiles)

	if len(mock.processedFiles) != len(expectedPaths) {
		t.Errorf("expected %d files, got %d", len(expectedPaths), len(mock.processedFiles))
		t.Logf("expected: %v", expectedPaths)
		t.Logf("got: %v", mock.processedFiles)
		return
	}

	for i, expected := range expectedPaths {
		if mock.processedFiles[i] != expected {
			t.Errorf("file %d: expected %s, got %s", i, expected, mock.processedFiles[i])
		}
	}

	// Verify hidden directories were NOT processed
	hiddenFiles := []string{
		filepath.Join(tmpDir, "testdir/.hiddendir/text.txt"),
		filepath.Join(tmpDir, "testdir/.hiddendir/.deephidden/text.txt"),
		filepath.Join(tmpDir, "testdir/subdir/.subhidden/text.txt"),
	}

	for _, hiddenFile := range hiddenFiles {
		for _, processedFile := range mock.processedFiles {
			if processedFile == hiddenFile {
				t.Errorf("hidden file should not be processed: %s", hiddenFile)
			}
		}
	}
}

func TestFileSelection_DirectoryWithHidden(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testfileselection")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	createTestFileStructure(t, tmpDir)
	mock := &mockProcessor{}

	// Process directory with IncludeHidden=true
	testDir := filepath.Join(tmpDir, "testdir")
	opts := Options{IncludeHidden: true}
	if err := processTarget(testDir, opts, mock.process); err != nil {
		t.Fatal(err)
	}

	// Should process files in ALL directories including hidden ones
	expectedDirs := []string{
		"testdir",
		"testdir/.hiddendir",
		"testdir/subdir",
		"testdir/.hiddendir/.deephidden",
		"testdir/subdir/.subhidden",
		"testdir/bin",
		"testdir/build",
	}

	expectedFiles := []string{
		"text.txt", "code.go", "readme.md", "config.json", "script.sh", "data.tmp", "build.log",
	}

	var expectedPaths []string
	for _, dir := range expectedDirs {
		for _, file := range expectedFiles {
			expectedPaths = append(expectedPaths, filepath.Join(tmpDir, dir, file))
		}
	}

	sort.Strings(expectedPaths)
	sort.Strings(mock.processedFiles)

	if len(mock.processedFiles) != len(expectedPaths) {
		t.Errorf("expected %d files, got %d", len(expectedPaths), len(mock.processedFiles))
		t.Logf("expected: %v", expectedPaths)
		t.Logf("got: %v", mock.processedFiles)
		return
	}

	// Verify all expected files were processed
	for i, expected := range expectedPaths {
		if mock.processedFiles[i] != expected {
			t.Errorf("file %d: expected %s, got %s", i, expected, mock.processedFiles[i])
		}
	}
}

func TestFileSelection_ExcludePatterns(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testfileselection")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	createTestFileStructure(t, tmpDir)
	mock := &mockProcessor{}

	// Process directory with exclude patterns
	testDir := filepath.Join(tmpDir, "testdir")
	opts := Options{
		ExcludePatterns: []string{"*.tmp", "*.log", "bin", "build"},
	}
	if err := processTarget(testDir, opts, mock.process); err != nil {
		t.Fatal(err)
	}

	// Should process files in testdir/ and testdir/subdir/ but exclude:
	// - *.tmp files
	// - *.log files  
	// - bin directory
	// - build directory
	expectedDirs := []string{
		"testdir",
		"testdir/subdir",
	}

	expectedFiles := []string{
		"text.txt", "code.go", "readme.md", "config.json", "script.sh",
	}

	var expectedPaths []string
	for _, dir := range expectedDirs {
		for _, file := range expectedFiles {
			expectedPaths = append(expectedPaths, filepath.Join(tmpDir, dir, file))
		}
	}

	sort.Strings(expectedPaths)
	sort.Strings(mock.processedFiles)

	if len(mock.processedFiles) != len(expectedPaths) {
		t.Errorf("expected %d files, got %d", len(expectedPaths), len(mock.processedFiles))
		t.Logf("expected: %v", expectedPaths)
		t.Logf("got: %v", mock.processedFiles)
		return
	}

	for i, expected := range expectedPaths {
		if mock.processedFiles[i] != expected {
			t.Errorf("file %d: expected %s, got %s", i, expected, mock.processedFiles[i])
		}
	}

	// Verify excluded patterns were NOT processed
	excludedPatterns := []string{
		"data.tmp", "build.log",
	}
	excludedDirs := []string{
		"testdir/bin", "testdir/build",
	}

	for _, pattern := range excludedPatterns {
		for _, processedFile := range mock.processedFiles {
			if filepath.Base(processedFile) == pattern {
				t.Errorf("excluded file should not be processed: %s", pattern)
			}
		}
	}

	for _, dir := range excludedDirs {
		fullDir := filepath.Join(tmpDir, dir)
		for _, processedFile := range mock.processedFiles {
			if filepath.Dir(processedFile) == fullDir {
				t.Errorf("excluded directory should not be processed: %s", dir)
			}
		}
	}
}

func TestFileSelection_CombinedOptions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testfileselection")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	createTestFileStructure(t, tmpDir)
	mock := &mockProcessor{}

	// Process directory with both IncludeHidden and exclude patterns
	testDir := filepath.Join(tmpDir, "testdir")
	opts := Options{
		IncludeHidden:   true,
		ExcludePatterns: []string{"*.tmp", "bin"},
	}
	if err := processTarget(testDir, opts, mock.process); err != nil {
		t.Fatal(err)
	}

	// Should process files in ALL directories (including hidden) but exclude:
	// - *.tmp files everywhere
	// - bin directory
	expectedDirs := []string{
		"testdir",
		"testdir/.hiddendir", 
		"testdir/subdir",
		"testdir/.hiddendir/.deephidden",
		"testdir/subdir/.subhidden",
		"testdir/build", // build not excluded, only bin
	}

	expectedFiles := []string{
		"text.txt", "code.go", "readme.md", "config.json", "script.sh", "build.log",
		// data.tmp excluded by pattern
	}

	var expectedPaths []string
	for _, dir := range expectedDirs {
		for _, file := range expectedFiles {
			expectedPaths = append(expectedPaths, filepath.Join(tmpDir, dir, file))
		}
	}

	sort.Strings(expectedPaths)
	sort.Strings(mock.processedFiles)

	if len(mock.processedFiles) != len(expectedPaths) {
		t.Errorf("expected %d files, got %d", len(expectedPaths), len(mock.processedFiles))
		t.Logf("expected: %v", expectedPaths)
		t.Logf("got: %v", mock.processedFiles)
	}
}

func TestFileSelection_BinaryFileDetection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testfileselection")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create text and binary files
	textFile := filepath.Join(tmpDir, "text.txt")
	binaryFile := filepath.Join(tmpDir, "binary.bin")
	
	if err := os.WriteFile(textFile, []byte("hello world\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binaryFile, []byte("\x00\x01\x02\x03\xFF\xFE"), 0o644); err != nil {
		t.Fatal(err)
	}

	mock := &mockProcessor{}

	// Process directory - should only process text files
	if err := processTarget(tmpDir, Options{}, mock.process); err != nil {
		t.Fatal(err)
	}

	// Should only process text file, binary should be skipped
	if len(mock.processedFiles) != 1 {
		t.Errorf("expected 1 file, got %d: %v", len(mock.processedFiles), mock.processedFiles)
		return
	}

	if mock.processedFiles[0] != textFile {
		t.Errorf("expected %s, got %s", textFile, mock.processedFiles[0])
	}

	// Verify binary file was not processed
	for _, processedFile := range mock.processedFiles {
		if processedFile == binaryFile {
			t.Errorf("binary file should not be processed: %s", binaryFile)
		}
	}
}