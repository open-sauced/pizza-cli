package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads a configuration file at a given path. It attempts to load
// the default location of a ".sauced.yaml" in the current working directory if an
// empty path is provided. If none is found, it tries to load
// "~/.sauced.yaml" from the fallback path, which is the user's home directory.
func LoadConfig(path string) (*Spec, error) {
	println("Config path loading from -c flag", path)

	config := &Spec{}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("error resolving absolute path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		// If the file does not exist, check if the fallback path exists
		if os.IsNotExist(err) {
			// load the default file path under the user's home dir
			usr, err := user.Current()

			if err != nil {
				return nil, fmt.Errorf("could not get user home directory: %w", err)
			}

			homeDirPathConfig, err := filepath.Abs(filepath.Join(usr.HomeDir, ".sauced.yaml"))

			if err != nil {
				return nil, fmt.Errorf("error home directory absolute path: %w", err)
			}

			_, err = os.Stat(homeDirPathConfig)
			if err != nil {
				return nil, fmt.Errorf("error reading config file from %s", homeDirPathConfig)
			}

			data, err = os.ReadFile(homeDirPathConfig)
			if err != nil {
				return nil, fmt.Errorf("error reading config file from %s or %s", absPath, homeDirPathConfig)
			}
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}
