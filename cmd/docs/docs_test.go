package docs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDocsPath(t *testing.T) {
	t.Parallel()

	t.Run("No path provided", func(t *testing.T) {
		t.Parallel()
		actual, err := GetDocsPath("")

		if err != nil {
			t.Errorf("GetDocsPath() error = %v, wantErr false", err)
			return
		}

		expected, _ := filepath.Abs(DefaultPath)
		if actual != expected {
			t.Errorf("GetDocsPath() = %v, want %v", actual, expected)
		}
	})

	t.Run("With path provided", func(t *testing.T) {
		t.Parallel()
		inputPath := filepath.Join(os.TempDir(), "docs")
		actual, err := GetDocsPath(inputPath)

		if err != nil {
			t.Errorf("GetDocsPath() error = %v, wantErr false", err)
			return
		}

		expected, _ := filepath.Abs(inputPath)
		if actual != expected {
			t.Errorf("GetDocsPath() = %v, want %v", actual, expected)
		}

		if _, err := os.Stat(actual); os.IsNotExist(err) {
			t.Errorf("GetDocsPath() failed to create directory %s", actual)
		}
	})

	t.Run("Invalid path", func(t *testing.T) {
		t.Parallel()
		invalidPath := string([]byte{0})

		_, err := GetDocsPath(invalidPath)

		if err == nil {
			t.Errorf("GetDocsPath() error = nil, wantErr true")
		}
	})
}

func TestGetDocsPath_ExistingDirectory(t *testing.T) {
	t.Parallel()

	tempDir, err := os.MkdirTemp("", "docs_test_existing")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	actual, err := GetDocsPath(tempDir)

	if err != nil {
		t.Errorf("GetDocsPath() error = %v, wantErr false", err)
		return
	}

	expected, _ := filepath.Abs(tempDir)
	if actual != expected {
		t.Errorf("GetDocsPath() = %v, want %v", actual, expected)
	}

	if _, err := os.Stat(actual); os.IsNotExist(err) {
		t.Errorf("GetDocsPath() failed to recognize existing directory %s", actual)
	}
}
