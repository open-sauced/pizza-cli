package docs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type Options struct {
	// the path to generate the documentation in
	path string
}

const DefaultPath = "./docs"

func NewDocsCommand() *cobra.Command {
	opts := &Options{}

	return &cobra.Command{
		Use:   "docs",
		Short: "Generates the documentation for the CLI",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.DisableAutoGenTag = true

			// Use default path if no argument is provided
			if len(args) == 0 {
				opts.path = DefaultPath
				fmt.Printf("No path was provided. Using default path: %s\n", DefaultPath)
			} else {
				absPath, err := filepath.Abs(args[0])
				if err != nil {
					return err
				}
				opts.path = absPath
			}

			// Create the directory if it doesn't exist
			_, err := os.Stat(opts.path)
			if os.IsNotExist(err) {
				fmt.Printf("The directory %s does not exist. Creating it...\n", opts.path)

				err := os.MkdirAll(opts.path, os.ModePerm)
				if err != nil {
					return fmt.Errorf("error creating documentation output directory %s: %s", opts.path, err)
				}
			}

			// Generate markdown documentation
			fmt.Printf("Generating documentation in %s...\n", opts.path)

			err = doc.GenMarkdownTree(cmd.Parent().Root(), opts.path)
			if err != nil {
				return err
			}

			fmt.Printf("Finished generating documentation in %s\n", opts.path)

			return nil
		},
	}
}
