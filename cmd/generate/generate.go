package generate

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/cmd/generate/codeowners"
)

const generateLongDesc string = `WARNING: Proof of concept feature.

XXX`

func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [subcommand] [flags]",
		Short: "Generates something",
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
