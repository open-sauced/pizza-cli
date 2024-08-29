package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	var tmpDir string

	setup := func() {
		tmpDir = t.TempDir()
	}

	teardown := func() {
		os.RemoveAll(tmpDir)
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T) (string, string)
		expectedError bool
	}{
		{
			name: "Existing file",
			setup: func(t *testing.T) (string, string) {
				configFilePath := filepath.Join(tmpDir, ".sauced.yaml")
				require.NoError(t, os.WriteFile(configFilePath, []byte("key: value"), 0644))
				return configFilePath, ""
			},
			expectedError: false,
		},
		{
			name: "Non-existent file",
			setup: func(_ *testing.T) (string, string) {
				return filepath.Join(tmpDir, ".sauced.yaml"), ""
			},
			expectedError: true,
		},
		{
			name: "Non-existent file with fallback",
			setup: func(_ *testing.T) (string, string) {
				fallbackPath := filepath.Join(tmpDir, ".sauced.yaml")
				require.NoError(t, os.WriteFile(fallbackPath, []byte("key: fallback"), 0644))
				nonExistentPath := filepath.Join(tmpDir, ".sauced.yaml")
				return nonExistentPath, fallbackPath
			},
			expectedError: false,
		},
		{
			name: "Default path",
			setup: func(_ *testing.T) (string, string) {
				return DefaultConfigPath, ""
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			setup()
			defer teardown()

			path, fallbackPath := tt.setup(t)
			config, err := LoadConfig(path, fallbackPath)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
			}
		})
	}
}
