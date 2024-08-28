package codeowners

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/open-sauced/pizza-cli/pkg/config"
)

func generateOutputFile(fileStats FileStats, outputPath string, ownersStyleFile bool, config *config.Spec) error {
	// Open the file for writing
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating %s file: %w", outputPath, err)
	}
	defer file.Close()

	// Write the header
	_, err = file.WriteString("# This file is generated automatically by OpenSauced pizza-cli. DO NOT EDIT. Stay saucy!\n\n")
	if err != nil {
		return fmt.Errorf("error writing to %s file: %w", outputPath, err)
	}

	// Sort the filenames to ensure consistent output
	var filenames []string
	for filename := range fileStats {
		filenames = append(filenames, filename)
	}
	sort.Strings(filenames)

	// Process each file
	for _, filename := range filenames {
		authorStats := fileStats[filename]
		if ownersStyleFile {
			err = writeOwnersChunk(authorStats, config, file, filename, outputPath)
			if err != nil {
				return fmt.Errorf("error writing to %s file: %w", outputPath, err)
			}
		} else {
			_, err := writeGitHubCodeownersChunk(authorStats, config, file, filename, outputPath)
			if err != nil {
				return fmt.Errorf("error writing to %s file: %w", outputPath, err)
			}
		}
	}

	return nil
}

func writeGitHubCodeownersChunk(authorStats AuthorStats, config *config.Spec, file *os.File, srcFilename string, outputPath string) ([]string, error) {
	topContributors := getTopContributorAttributions(authorStats, 3, config)

	resultSlice := []string{}
	for _, contributor := range topContributors {
		resultSlice = append(resultSlice, contributor.GitHubAlias)
	}

	if len(topContributors) > 0 {
		_, err := fmt.Fprintf(file, "%s @%s\n", cleanFilename(srcFilename), strings.Join(resultSlice, " @"))
		if err != nil {
			return nil, fmt.Errorf("error writing to %s file: %w", outputPath, err)
		}
	} else {
		// no code owners to attribute to file
		_, err := fmt.Fprintf(file, "%s\n", cleanFilename(srcFilename))
		if err != nil {
			return nil, fmt.Errorf("error writing to %s file: %w", outputPath, err)
		}
	}

	return resultSlice, nil
}

func writeOwnersChunk(authorStats AuthorStats, config *config.Spec, file *os.File, srcFilename string, outputPath string) error {
	topContributors := getTopContributorAttributions(authorStats, 3, config)

	_, err := fmt.Fprintf(file, "%s\n", srcFilename)
	if err != nil {
		return fmt.Errorf("error writing to %s file: %w", outputPath, err)
	}

	for i := 0; i < len(topContributors) && i < 3; i++ {
		_, err = fmt.Fprintf(file, "  - %s\n", topContributors[i].Name)
		if err != nil {
			return fmt.Errorf("error writing to %s file: %w", outputPath, err)
		}

		_, err = fmt.Fprintf(file, "    - %s\n", topContributors[i].Email)
		if err != nil {
			return fmt.Errorf("error writing to %s file: %w", outputPath, err)
		}
	}

	return nil
}

func getTopContributorAttributions(authorStats AuthorStats, n int, config *config.Spec) AuthorStatSlice {
	sortedAuthorStats := authorStats.ToSortedSlice()

	// Get top n contributors (or all if less than n)
	var topContributors AuthorStatSlice

	for i := 0; i < len(sortedAuthorStats) && i < n; i++ {
		// get attributions for email / github handles
		for username, emails := range config.Attributions {
			for _, email := range emails {
				if email == sortedAuthorStats[i].Email {
					sortedAuthorStats[i].GitHubAlias = username
					topContributors = append(topContributors, sortedAuthorStats[i])
				}
			}
		}
	}

	return topContributors
}

func cleanFilename(filename string) string {
	// Split the filename in case its rename, see https://github.com/open-sauced/pizza-cli/issues/101
	parsedFilename := strings.Split(filename, " ")[0]
	// Replace anything that is not a word, period, single quote, dash, space, forward slash, or backslash with an escaped version
	re := regexp.MustCompile(`([^\w\.\'\-\s\/\\])`)
	escapedFilename := re.ReplaceAllString(parsedFilename, "\\$0")

	return escapedFilename
}
