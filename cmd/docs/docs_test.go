package docs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDocsCommand(t *testing.T) {
	t.Parallel()

	t.Run("Docs generate using default path", func(t *testing.T) {
		t.Parallel()
		cmd := NewDocsCommand()
		err := cmd.Execute()
		require.NoError(t, err)

		// Check if a Markdown file was generated
		files, err := os.ReadDir(DefaultPath)
		assert.NoError(t, err)
		assert.NotEmpty(t, files)

		markdownFileFound := false

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".md") {
				markdownFileFound = true
				break
			}
		}

		assert.True(t, markdownFileFound, "No Markdown file was generated in the default path")

		os.RemoveAll(DefaultPath)
	})

	t.Run("Docs generate using a custom path", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		customPath := filepath.Join(tempDir, "custom_docs")

		cmd := NewDocsCommand()
		cmd.SetArgs([]string{customPath})
		err := cmd.Execute()
		require.NoError(t, err)

		// Check if a Markdown file was generated
		files, err := os.ReadDir(customPath)
		assert.NoError(t, err)
		assert.NotEmpty(t, files)

		markdownFileFound := false

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".md") {
				markdownFileFound = true
				break
			}
		}

		assert.True(t, markdownFileFound, "No Markdown file was generated in the custom path")
	})

	t.Run("Docs fail to generate when the output path is invalid", func(t *testing.T) {
		t.Parallel()
		cmd := NewDocsCommand()
		cmd.SetArgs([]string{string([]byte{0x00})})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("Docs generate using an existing directory", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()

		cmd := NewDocsCommand()
		cmd.SetArgs([]string{tempDir})
		err := cmd.Execute()
		require.NoError(t, err)

		// Check if files were generated in the existing directory
		files, err := os.ReadDir(tempDir)
		assert.NoError(t, err)
		assert.NotEmpty(t, files)

		markdownFileFound := false
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".md") {
				markdownFileFound = true
				break
			}
		}
		assert.True(t, markdownFileFound, "No Markdown file was generated in the existing directory")
	})

	t.Run("TooManyArguments", func(t *testing.T) {
		t.Parallel()
		cmd := NewDocsCommand()
		cmd.SetArgs([]string{"path1", "path2"})
		err := cmd.Execute()
		assert.Error(t, err)
	})
}
