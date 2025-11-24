package whitespace

import (
	"os"
	"testing"
)

// createFileWithContent creates a temporary file with specified content for testing
func createTestFile(t *testing.T, content []byte) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "testnewline")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()
	
	if err := os.WriteFile(tmpFile.Name(), content, 0o644); err != nil {
		t.Fatal(err)
	}
	return tmpFile.Name()
}

// readFileBytes reads file content as bytes
func readFileBytes(t *testing.T, path string) []byte {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func TestEnsureSingleNewline(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "empty file",
			input:    []byte{},
			expected: []byte("\n"),
		},
		{
			name:     "content without newline",
			input:    []byte("content"),
			expected: []byte("content\n"),
		},
		{
			name:     "content with single newline",
			input:    []byte("content\n"),
			expected: []byte("content\n"),
		},
		{
			name:     "content with multiple newlines",
			input:    []byte("content\n\n\n"),
			expected: []byte("content\n"),
		},
		{
			name:     "single newline only",
			input:    []byte("\n"),
			expected: []byte("\n"),
		},
		{
			name:     "multiple newlines only",
			input:    []byte("\n\n\n"),
			expected: []byte("\n"),
		},
		{
			name:     "content with CRLF endings",
			input:    []byte("line1\r\nline2\r\n\r\n"),
			expected: []byte("line1\r\nline2\n"),
		},
		{
			name:     "content with mixed line endings",
			input:    []byte("line1\nline2\r\nline3\n\r\n\n"),
			expected: []byte("line1\nline2\r\nline3\n"),
		},
		{
			name:     "content ending with CR only",
			input:    []byte("content\r"),
			expected: []byte("content\n"),
		},
		{
			name:     "multiline content with proper ending",
			input:    []byte("line1\nline2\nline3\n"),
			expected: []byte("line1\nline2\nline3\n"),
		},
		{
			name:     "multiline content with extra newlines",
			input:    []byte("line1\nline2\nline3\n\n\n\n"),
			expected: []byte("line1\nline2\nline3\n"),
		},
		{
			name:     "content with trailing spaces then newlines",
			input:    []byte("content   \n\n\n"),
			expected: []byte("content   \n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTestFile(t, tt.input)
			defer os.Remove(path)

			if err := ensureSingleNewline(path); err != nil {
				t.Fatalf("ensureSingleNewline failed: %v", err)
			}

			result := readFileBytes(t, path)
			if string(result) != string(tt.expected) {
				t.Errorf("expected %q, got %q", string(tt.expected), string(result))
			}
		})
	}
}

// Integration test for ProcessNewline function
func TestProcessNewlineSingleFile(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "file needing newline",
			input:    []byte("content without newline"),
			expected: []byte("content without newline\n"),
		},
		{
			name:     "file with extra newlines",
			input:    []byte("content\n\n\n"),
			expected: []byte("content\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTestFile(t, tt.input)
			defer os.Remove(path)

			if err := ProcessNewline(path); err != nil {
				t.Fatalf("ProcessNewline failed: %v", err)
			}

			result := readFileBytes(t, path)
			if string(result) != string(tt.expected) {
				t.Errorf("expected %q, got %q", string(tt.expected), string(result))
			}
		})
	}
}
