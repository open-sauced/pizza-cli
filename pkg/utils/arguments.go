package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

func HandleRepositoryValues(repos []string, filePath string) (map[string]struct{}, error) {
	uniqueRepoURLs := make(map[string]struct{})
	for _, repo := range repos {
		uniqueRepoURLs[repo] = struct{}{}
	}
	if filePath != "" {
		file, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		var reposFromYaml []string
		err = yaml.Unmarshal(file, &reposFromYaml)
		if err != nil {
			return nil, err
		}
		for _, repo := range reposFromYaml {
			uniqueRepoURLs[repo] = struct{}{}
		}
	}
	return uniqueRepoURLs, nil
}
