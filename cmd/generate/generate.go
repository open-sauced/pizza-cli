package generate

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/cmd/generate/codeowners"
)

const generateLongDesc string = `The 'generate' command provides tools to automate the creation of important project documentation and derive insights from your codebase.

Currently, it supports generating CODEOWNERS files.

Available subcommands:
  - codeowners: Generate a more granular GitHub-style CODEOWNERS file based on git history.`

func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [subcommand] [flags]",
		Short: "Generates documentation and insights from your codebase",
		Long:  generateLongDesc,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you must provide a subcommand")
			}

			return nil
		},
		RunE: run,
	}

	cmd.AddCommand(codeowners.NewCodeownersCommand())

	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
