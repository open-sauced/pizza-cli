package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = "~/.sauced.yaml"

// LoadConfig loads a configuration file at a given path. It attempts to load
// the default location of a ".sauced.yaml" in the user's home directory if an
// empty path is not provided.
func LoadConfig(path string, fallbackPath string) (*Spec, error) {
	config := &Spec{}

	if path == DefaultConfigPath || path == "" {
		// load the default file path under the user's home dir
		usr, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("could not get user home directory: %w", err)
		}

		path = filepath.Join(usr.HomeDir, ".sauced.yaml")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("error resolving absolute path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {

		// If the file does not exist, check if the fallback path exists
		if os.IsNotExist(err) {
			_, err = os.Stat(fallbackPath)
			if err != nil {
				return nil, fmt.Errorf("error reading config file from %s or %s", absPath, fallbackPath)
			}
		}
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}
