package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an empty file for testing
	configFilePath := filepath.Join(tmpDir, ".sauced.yaml")

	if err := os.WriteFile(configFilePath, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	// Path for a non-existent file
	nonExistentPath := filepath.Join(tmpDir, "non_existent.yaml")

	tests := []struct {
		name          string
		path          string
		fallbackPath  string
		expectedError bool
	}{
		{"Existing file", configFilePath, "", false},
		{"Non-existent file", nonExistentPath, "", true},
		{"Non-existent file with fallback", DefaultConfigPath, configFilePath, false},
		{"Default path", DefaultConfigPath, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadConfig(tt.path, tt.fallbackPath)

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
