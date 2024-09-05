package codeowners

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/open-sauced/pizza-cli/pkg/config"
)

func TestCleanFilename(testRunner *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected string
	}{
		{"path/to/(home).go", "path/to/(home).go", `path/to/\(home\).go`},
		{"path/to/[home].go", "path/to/[home].go", `path/to/\[home\].go`},
		{"path/to/+page.go", "path/to/+page.go", `path/to/\+page.go`},
		{"path/to/go-home.go", "path/to/go-home.go", `path/to/go-home.go`},
		{`path\to\(home).go`, `path\to\(home).go`, `path\to\\(home\).go`},
		{`path\to\[home].go`, `path\to\[home].go`, `path\to\\[home\].go`},
		{`path\to\+page.go`, `path\to\+page.go`, `path\to\\+page.go`},
		{`path\to\go-home.go`, `path\to\go-home.go`, `path\to\go-home.go`},
	}

	for _, testItem := range tests {
		testRunner.Run(testItem.name, func(tester *testing.T) {
			ans := cleanFilename(testItem.input)
			if ans != testItem.expected {
				tester.Errorf("got %s, expected %s", ans, testItem.expected)
			}
		})
	}
}

func TestGetTopContributorAttributions(testRunner *testing.T) {
	configSpec := config.Spec{
		Attributions: map[string][]string{
			"brandonroberts": {"brandon@opensauced.pizza"},
		},
		AttributionFallback: []string{"open-sauced/engineering"},
	}

	var authorStats = AuthorStats{
		"brandon": {GitHubAlias: "brandon", Email: "brandon@opensauced.pizza", Lines: 20},
		"john":    {GitHubAlias: "john", Email: "john@opensauced.pizza", Lines: 15},
	}

	results := getTopContributorAttributions(authorStats, 3, &configSpec)

	assert.Equal(testRunner, len(results), 1, "Expected 1 result")
	assert.Equal(testRunner, results[0].GitHubAlias, "brandonroberts", "Expected brandonroberts")
}

func TestGetFallbackAttributions(testRunner *testing.T) {
	configSpec := config.Spec{
		Attributions: map[string][]string{
			"jpmcb":          {"jpmcb@opensauced.pizza"},
			"brandonroberts": {"brandon@opensauced.pizza"},
		},
		AttributionFallback: []string{"open-sauced/engineering"},
	}

	results := getTopContributorAttributions(AuthorStats{}, 3, &configSpec)

	assert.Equal(testRunner, len(results), 1, "Expected 1 result")
	assert.Equal(testRunner, results[0].GitHubAlias, "open-sauced/engineering", "Expected open-sauced/engineering")
}
