package offboard

import (
	"errors"
	"fmt"
	"os"
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
	outputType := "CODEOWNERS"
	var owners []byte
	var err error

	if _, err = os.Stat(path + "/CODEOWNERS"); !errors.Is(err, os.ErrNotExist) {
		fmt.Print("CODEOWNERS EXISTS")
		outputType = "CODEOWNERS"
		owners, err = os.ReadFile(path + "/CODEOWNERS")
	} else if _, err = os.Stat(path + "/OWNERS"); !errors.Is(err, os.ErrNotExist) {
		fmt.Print("OWNERS EXISTS")
		outputType = "OWNERS"
		owners, err = os.ReadFile(path + "/OWNERS")
	}

	if err != nil {
		// fmt.Errorf("failed to find existing owners: %w", err)
		fmt.Printf("WTF %v", err)
		fmt.Printf("will create a new %s file in the path %s", outputType, path)
	}

	lines := strings.Split(string(owners), "\n")
	for _, line := range lines {
		for _, name := range offboardingUsers {
			fmt.Println(name)
			strings.Cut(line, name)
		}
	}

	output := strings.Join(lines, "\n")
	file, err := os.Create(path+"/"+outputType)

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
