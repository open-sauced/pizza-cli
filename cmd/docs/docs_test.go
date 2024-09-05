package docs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeterminePath(t *testing.T) {
	t.Parallel()

	t.Run("No path passed to command", func(t *testing.T) {
		t.Parallel()
		got, err := DeterminePath([]string{})

		if err != nil {
			t.Errorf("DeterminePath() error = %v, wantErr false", err)
			return
		}

		if got != DefaultPath {
			t.Errorf("DeterminePath() = %v, want %v", got, DefaultPath)
		}
	})

	t.Run("With path passed to command", func(t *testing.T) {
		t.Parallel()
		expected := "/tmp/docs"
		got, err := DeterminePath([]string{expected})

		if err != nil {
			t.Errorf("DeterminePath() error = %v, wantErr false", err)
			return
		}

		if got != expected {
			t.Errorf("DeterminePath() = %v, want %v", got, expected)
		}
	})
}

func TestEnsureDirectoryExists(t *testing.T) {
	t.Parallel()

	t.Run("Existing directory", func(t *testing.T) {
		t.Parallel()
		tempDir, err := os.MkdirTemp(t.TempDir(), "docs_test_existing")

		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}

		err = EnsureDirectoryExists(tempDir)

		if err != nil {
			t.Errorf("EnsureDirectoryExists() error = %v, wantErr false", err)
			return
		}

		if _, err := os.Stat(tempDir); os.IsNotExist(err) {
			t.Errorf("EnsureDirectoryExists() failed to recognize existing directory %s", tempDir)
		}
	})

	t.Run("New directory", func(t *testing.T) {
		t.Parallel()
		tempDir, err := os.MkdirTemp(t.TempDir(), "docs_test_new")

		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}

		newDir := filepath.Join(tempDir, "new_dir")
		err = EnsureDirectoryExists(newDir)

		if err != nil {
			t.Errorf("EnsureDirectoryExists() error = %v, wantErr false", err)
			return
		}
		if _, err := os.Stat(newDir); os.IsNotExist(err) {
			t.Errorf("EnsureDirectoryExists() failed to create directory %s", newDir)
		}
	})
}
