package config

import (
	"errors"
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/pkg/config"
)

// Options for the codeowners generation command
type Options struct {
	// the path to the git repository on disk to generate a codeowners file for
	path string

	tty      bool
	loglevel int

	config *config.Spec
}

const codeownersLongDesc string = `WARNING: Proof of concept feature.

Generates a ~/.sauced.yaml configuration file. The attribution of emails to given entities
is based on the repository this command is ran in.`

func NewConfigCommand() *cobra.Command {
	opts := &Options{}
	print(opts.path);

	cmd := &cobra.Command{
		Use:   "config path/to/repo [flags]",
		Short: "Generates a \"~/.sauced.yaml\" config based on the current repository",
		Long:  codeownersLongDesc,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you must provide exactly one argument: the path to the repository")
			}

			path := args[0]
			print(path)

			return nil 
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			return nil
		},
	}

	return cmd
}
