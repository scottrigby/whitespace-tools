package whitespace

import (
	"os"
	"testing"
)

// createTrailingspaceTestFile creates a temporary file with specified content for testing
func createTrailingspaceTestFile(t *testing.T, content string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "testtrailingspace")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()
	
	if err := os.WriteFile(tmpFile.Name(), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return tmpFile.Name()
}

// readTrailingspaceFileContent reads file content as string
func readTrailingspaceFileContent(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func TestRemoveTrailingWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single space at end of line with content",
			input:    "content \n",
			expected: "content\n",
		},
		{
			name:     "two spaces at end of line with content",
			input:    "content  \n",
			expected: "content\n",
		},
		{
			name:     "single tab at end of line with content",
			input:    "content\t\n",
			expected: "content\n",
		},
		{
			name:     "two tabs at end of line with content",
			input:    "content\t\t\n",
			expected: "content\n",
		},
		{
			name:     "mixed spaces and tabs at end of line with content",
			input:    "content \t \t\n",
			expected: "content\n",
		},
		{
			name:     "line with only spaces (should become empty line)",
			input:    "line1\n  \nline3\n",
			expected: "line1\n\nline3\n",
		},
		{
			name:     "line with only tabs (should become empty line)",
			input:    "line1\n\t\t\nline3\n",
			expected: "line1\n\nline3\n",
		},
		{
			name:     "line with only mixed spaces and tabs (should become empty line)",
			input:    "line1\n \t \t \nline3\n",
			expected: "line1\n\nline3\n",
		},
		{
			name:     "spaces at EOF with newline",
			input:    "content\n  \n",
			expected: "content\n\n",
		},
		{
			name:     "tabs at EOF with newline",
			input:    "content\n\t\t\n",
			expected: "content\n\n",
		},
		{
			name:     "spaces at EOF without newline",
			input:    "content  ",
			expected: "content",
		},
		{
			name:     "tabs at EOF without newline",
			input:    "content\t\t",
			expected: "content",
		},
		{
			name:     "preserve leading spaces and tabs (should NOT be modified)",
			input:    "  \tcontent with leading whitespace\n",
			expected: "  \tcontent with leading whitespace\n",
		},
		{
			name:     "preserve internal spaces and tabs (should NOT be modified)",
			input:    "word1  \t  word2\t\tword3\n",
			expected: "word1  \t  word2\t\tword3\n",
		},
		{
			name:     "complex case with various whitespace scenarios",
			input:    "  leading preserved\ntrailing removed  \n\t\t\nmore content \t \n",
			expected: "  leading preserved\ntrailing removed\n\nmore content\n",
		},
		{
			name:     "empty file",
			input:    "",
			expected: "",
		},
		{
			name:     "clean lines unchanged",
			input:    "clean line1\nclean line2\nclean line3\n",
			expected: "clean line1\nclean line2\nclean line3\n",
		},
		{
			name:     "file without final newline but with trailing whitespace",
			input:    "no final newline  \t",
			expected: "no final newline",
		},
		{
			name:     "multiple consecutive empty lines with whitespace",
			input:    "content\n  \n\t\n \t \nmore content\n",
			expected: "content\n\n\n\nmore content\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTrailingspaceTestFile(t, tt.input)
			defer os.Remove(path)

			if err := removeTrailingWhitespace(path); err != nil {
				t.Fatalf("removeTrailingWhitespace failed: %v", err)
			}

			result := readTrailingspaceFileContent(t, path)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Integration test for ProcessTrailingspace function
func TestProcessTrailingspaceSingleFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "file needing trailing whitespace removal",
			input:    "line with spaces   \nline with tabs\t\t\nclean line\n",
			expected: "line with spaces\nline with tabs\nclean line\n",
		},
		{
			name:     "file with mixed trailing whitespace",
			input:    "mixed \t \nanother \t\nclean\n",
			expected: "mixed\nanother\nclean\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTrailingspaceTestFile(t, tt.input)
			defer os.Remove(path)

			if err := ProcessTrailingspace(path); err != nil {
				t.Fatalf("ProcessTrailingspace failed: %v", err)
			}

			result := readTrailingspaceFileContent(t, path)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test edge cases and performance
func TestTrailingspaceEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "very long line with trailing whitespace",
			input:    "a" + string(make([]byte, 1000)) + "content   \n",
			expected: "a" + string(make([]byte, 1000)) + "content\n",
		},
		{
			name:     "file with only whitespace lines",
			input:    "   \n\t\t\n \t \n",
			expected: "\n\n\n",
		},
		{
			name:     "single character with trailing space",
			input:    "x \n",
			expected: "x\n",
		},
		{
			name:     "unicode content with trailing whitespace",
			input:    "héllo wörld  \n",
			expected: "héllo wörld\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTrailingspaceTestFile(t, tt.input)
			defer os.Remove(path)

			if err := removeTrailingWhitespace(path); err != nil {
				t.Fatalf("removeTrailingWhitespace failed: %v", err)
			}

			result := readTrailingspaceFileContent(t, path)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}