package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	t.Run("Existing file", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		configFilePath := filepath.Join(tmpDir, ".sauced.yaml")
		require.NoError(t, os.WriteFile(configFilePath, []byte("key: value"), 0644))

		config, err := LoadConfig(configFilePath, "")
		assert.NoError(t, err)
		assert.NotNil(t, config)
	})

	t.Run("Non-existent file", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		nonExistentPath := filepath.Join(tmpDir, ".sauced.yaml")

		config, err := LoadConfig(nonExistentPath, "")
		assert.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("Non-existent file with fallback", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		fallbackPath := filepath.Join(tmpDir, ".sauced.yaml")
		require.NoError(t, os.WriteFile(fallbackPath, []byte("key: fallback"), 0644))
		nonExistentPath := filepath.Join(tmpDir, "non-existent.yaml")

		config, err := LoadConfig(nonExistentPath, fallbackPath)
		assert.NoError(t, err)
		assert.NotNil(t, config)
	})

	//t.Run("Default path", func(t *testing.T) {
	//t.Parallel()
	//config, err := LoadConfig(DefaultConfigPath, "")
	//assert.Error(t, err)
	//assert.Nil(t, config)
	//})
}
