package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// GetOwnerAndRepoFromURL: extracts the owner and repository name
func GetOwnerAndRepoFromURL(input string) (owner, repo string, err error) {
	var repoOwner, repoName string

	// check (https://github.com/owner/repo) format
	u, err := url.Parse(input)
	if err == nil && u.Host == "github.com" {
		path := strings.Trim(u.Path, "/")
		parts := strings.Split(path, "/")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("Invalid URL: %s", input)
		}
		repoOwner = parts[0]
		repoName = parts[1]
		return repoOwner, repoName, nil
	}

	// check (owner/repo) format
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Invalid URL: %s", input)
	}
	repoOwner = parts[0]
	repoName = parts[1]

	return repoOwner, repoName, nil
}
