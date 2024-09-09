package docs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type Options struct {
	Path string
}

const DefaultPath = "./docs"

func GetDocsPath(path string) (string, error) {
	if path == "" {
		fmt.Printf("No path was provided. Using default path: %s\n", DefaultPath)
		path = DefaultPath
	}

	absPath, err := filepath.Abs(path)

	if err != nil {
		return "", fmt.Errorf("error resolving absolute path: %w", err)
	}

	_, err = os.Stat(absPath)

	if os.IsNotExist(err) {
		fmt.Printf("The directory %s does not exist. Creating it...\n", absPath)
		if err := os.MkdirAll(absPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("error creating directory %s: %w", absPath, err)
		}
	} else if err != nil {
		return "", fmt.Errorf("error checking directory %s: %w", absPath, err)
	}

	return absPath, nil
}

func GenerateDocumentation(rootCmd *cobra.Command, path string) error {
	fmt.Printf("Generating documentation in %s...\n", path)
	err := doc.GenMarkdownTree(rootCmd, path)

	if err != nil {
		return err
	}

	fmt.Printf("Finished generating documentation in %s\n", path)

	return nil
}

func NewDocsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "docs [path]",
		Short: "Generates the documentation for the CLI",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Parent().Root().DisableAutoGenTag = true

			var path string
			if len(args) > 0 {
				path = args[0]
			}

			resolvedPath, err := GetDocsPath(path)
			if err != nil {
				return err
			}

			return GenerateDocumentation(cmd.Parent().Root(), resolvedPath)
		},
	}
}
