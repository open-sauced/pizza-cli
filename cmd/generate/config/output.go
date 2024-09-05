package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-sauced/pizza-cli/pkg/config"
	"github.com/open-sauced/pizza-cli/pkg/utils"
)

func generateOutputFile(outputPath string, attributionMap map[string][]string) error {
	// Open the file for writing
	homeDir, err := os.UserHomeDir()
	file, err := os.Create(filepath.Join(homeDir, outputPath))
	if err != nil {
		return fmt.Errorf("error creating %s file: %w", outputPath, err)
	}
	defer file.Close()

	var config config.Spec
	config.Attributions = attributionMap

	// for pretty print test
	yaml, err := utils.OutputYAML(config)

	if err != nil {
		return fmt.Errorf("Failed to turn into YAML")
	}

	file.WriteString(yaml)

	return nil
}
