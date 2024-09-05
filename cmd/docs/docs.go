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

func DeterminePath(args []string) (string, error) {
	if len(args) == 0 {
		fmt.Printf("No path was provided. Using default path: %s\n", DefaultPath)
		return DefaultPath, nil
	}

	absPath, err := filepath.Abs(args[0])

	if err != nil {
		return "", err
	}

	return absPath, nil
}

func EnsureDirectoryExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("The directory %s does not exist. Creating it...\n", path)
		return os.MkdirAll(path, os.ModePerm)
	}
	return err
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
		Use:   "docs",
		Short: "Generates the documentation for the CLI",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Parent().Root().DisableAutoGenTag = true

			path, err := DeterminePath(args)

			if err != nil {
				return err
			}

			if err := EnsureDirectoryExists(path); err != nil {
				return fmt.Errorf("error creating documentation output directory %s: %s", path, err)
			}

			return GenerateDocumentation(cmd.Parent().Root(), path)
		},
	}
}
