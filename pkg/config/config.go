package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads a configuration file at a given path.
// If the provided path does not exist or doesn't contain a ".sauced.yaml" file,
// "~/.sauced.yaml" from the fallback path, which is the user's home directory, is used.
//
// This function returns the config Spec, the location the spec was loaded from, and an error
func LoadConfig(path string) (*Spec, string, error) {
	givenPathSpec, givenLoadedPath, givenPathErr := loadSpecAtPath(path)
	if givenPathErr == nil {
		return givenPathSpec, givenLoadedPath, nil
	}

	homePathSpec, homeLoadedPath, homePathErr := loadSpecAtHome()
	if homePathErr == nil {
		return homePathSpec, homeLoadedPath, nil
	}

	return nil, "", fmt.Errorf("could not load config at given path: %w - could not load config at home: %w", givenPathErr, homePathErr)
}

func loadSpecAtPath(path string) (*Spec, string, error) {
	config := &Spec{}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, "", fmt.Errorf("error resolving absolute path: %s - %w", path, err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, "", fmt.Errorf("error reading config file from given absolute path: %s - %w", absPath, err)
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, "", fmt.Errorf("error unmarshaling config at: %s - %w", absPath, err)
	}

	return config, absPath, nil
}

func loadSpecAtHome() (*Spec, string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, "", fmt.Errorf("could not get user home directory: %w", err)
	}

	path := filepath.Join(usr.HomeDir, ".sauced.yaml")
	conf, loadedPath, err := loadSpecAtPath(path)
	if err != nil {
		return nil, "", fmt.Errorf("could not load spec at home: %s - %w", path, err)
	}

	return conf, loadedPath, nil
}
