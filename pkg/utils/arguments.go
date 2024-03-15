package utils

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// HandleUniqueValues: returns unique values
func HandleUniqueValues(args []string, filePath string) (map[string]struct{}, error) {
	uniqueValues := make(map[string]struct{})
	for _, arg := range args {
		uniqueValues[arg] = struct{}{}
	}
	if filePath != "" {
		file, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		var valuesFromYaml []string
		err = yaml.Unmarshal(file, &valuesFromYaml)
		if err != nil {
			return nil, err
		}
		for _, arg := range valuesFromYaml {
			uniqueValues[arg] = struct{}{}
		}
	}
	return uniqueValues, nil
}

// IsYAMLFile: returns true if the file given is a yaml file
func IsYAMLFile(input string) bool {
	contents := strings.Split(input, ".")
	if len(contents) != 2 {
		return false
	}
	return contents[1] == "yml" || contents[1] == "yaml"
}

// ParseFileAndCSV: parses arguments in csv form or a yaml filepath, and returns the values as a string slice
func ParseFileAndCSV(input string) ([]interface{}, error) {
	var parsedValues []interface{}
	if input == "" {
		return parsedValues, fmt.Errorf("field was not provided")
	}

	var values map[string]struct{}
	var err error
	if IsYAMLFile(input) {
		values, err = HandleUniqueValues([]string{}, input)
		if err != nil {
			return parsedValues, err
		}
	} else {
		csv := strings.Split(input, ",")
		values, err = HandleUniqueValues(csv, "")
		if err != nil {
			return parsedValues, err
		}
	}

	parsedValues = make([]interface{}, len(values))
	i := 0
	for val := range values {
		parsedValues[i] = val
		i++
	}
	return parsedValues, nil
}
