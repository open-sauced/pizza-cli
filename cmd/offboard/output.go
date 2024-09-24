package offboard

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
