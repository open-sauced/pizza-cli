package utils

import (
	"fmt"
	"strings"
)

func GetOwnerAndRepoFromURL(url string) (owner, repo string, err error) {
	if !strings.HasPrefix(url, "https://github.com/") {
		return "", "", fmt.Errorf("invalid URL: %s", url)
	}

	// Remove the "https://github.com/" prefix from the URL
	url = strings.TrimPrefix(url, "https://github.com/")

	// Split the remaining URL path into segments
	segments := strings.Split(url, "/")

	// The first segment is the owner, and the second segment is the repository name
	if len(segments) >= 2 {
		owner = segments[0]
		repo = segments[1]
	} else {
		return "", "", fmt.Errorf("invalid URL: %s", url)
	}

	return owner, repo, nil
}
