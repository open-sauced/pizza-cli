package offboard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-sauced/pizza-cli/v2/pkg/config"
	"github.com/open-sauced/pizza-cli/v2/pkg/utils"
)

func generateConfigFile(outputPath string, attributionMap map[string][]string) error {
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

func generateOwnersFile(path string, offboardingUsers []string) error {
	outputType := "/CODEOWNERS"
	var owners []byte
	var err error

	var ownersPath string

	if _, err = os.Stat(filepath.Join(path, "/CODEOWNERS")); !errors.Is(err, os.ErrNotExist) {
		fmt.Print("CODEOWNERS EXISTS")
		outputType = "CODEOWNERS"
		ownersPath = filepath.Join(path, "/CODEOWNERS")
		owners, err = os.ReadFile(ownersPath)
	} else if _, err = os.Stat(filepath.Join(path, "OWNERS")); !errors.Is(err, os.ErrNotExist) {
		fmt.Print("OWNERS EXISTS")
		outputType = "OWNERS"
		ownersPath = filepath.Join(path, "/OWNERS")
		owners, err = os.ReadFile(ownersPath)
	}

	if err != nil {
		fmt.Printf("will create a new %s file in the path %s", outputType, path)
	}

	lines := strings.Split(string(owners), "\n")
	var newLines []string
	for _, line := range lines {
		newLine := line
		for _, name := range offboardingUsers {
			result, _, found := strings.Cut(newLine, "@"+name)
			if found {
				newLine = result
			}
		}
		newLines = append(newLines, newLine)
	}

	output := strings.Join(newLines, "\n")
	file, err := os.Create(ownersPath)
	if err != nil {
		return fmt.Errorf("error creating %s file: %w", outputType, err)
	}
	defer file.Close()

	_, err = file.WriteString(output)
	if err != nil {
		return fmt.Errorf("failed writing file %s: %w", path+outputType, err)
	}

	return nil
}
