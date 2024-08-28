package codeowners

import (
		"testing"
)

func TestCleanFilename(testRunner *testing.T) {
	var tests = []struct {
		name string
		input string
		expected  string
	}{
			{"path/to/(home).go", "path/to/(home).go", `path/to/\(home\).go`},
			{"path/to/[home].go", "path/to/[home].go", `path/to/\[home\].go`},
			{"path/to/+page.go", "path/to/+page.go", `path/to/\+page.go`},
			{"path/to/go-home.go", "path/to/go-home.go", `path/to/go-home.go`},
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
