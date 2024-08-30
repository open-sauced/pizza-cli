package config

import (
	"fmt"
	"os"
	"path"
)

const (
	configDir = ".pizza-cli"
)

// GetConfigDirectory gets the config directory path for the Pizza CLI.
// This function should be used to ensure consistency among commands for loading
// and modifying the config.
func GetConfigDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldn't get user home directory: %w", err)
	}

	dirName := path.Join(homeDir, configDir)

	// Check if the directory already exists. If not, create it.
	_, err = os.Stat(dirName)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			return "", fmt.Errorf(".pizza-cli directory could not be created: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("error checking ~/.pizza-cli directory: %w", err)
	}

	return dirName, nil
}
