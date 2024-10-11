package config

import (
	"fmt"
	"os"

	"github.com/open-sauced/pizza-cli/v2/pkg/config"
	"github.com/open-sauced/pizza-cli/v2/pkg/utils"
)

func generateOutputFile(outputPath string, attributionMap map[string][]string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating %s file: %w", outputPath, err)
	}
	defer file.Close()

	// write the header preamble
	_, err = file.WriteString("# Configuration for attributing commits with emails to GitHub user profiles\n# Used during codeowners generation.\n\n# List the emails associated with the given username\n# The commits associated with these emails will be attributed to\n# the username in this yaml map. Any number of emails may be listed\n\n")

	if err != nil {
		return fmt.Errorf("error writing to %s file: %w", outputPath, err)
	}

	var config config.Spec
	config.Attributions = attributionMap

	// for pretty print test
	yaml, err := utils.OutputYAML(config)

	if err != nil {
		return fmt.Errorf("failed to turn into YAML: %w", err)
	}

	_, err = file.WriteString(yaml)

	if err != nil {
		return fmt.Errorf("failed to turn into YAML: %w", err)
	}

	return nil
}
