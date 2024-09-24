package offboard

import (
	"errors"
	"fmt"

	"github.com/open-sauced/pizza-cli/v2/pkg/config"

	"github.com/spf13/cobra"
)

type Options struct {
	offboardingUsers []string

	// config file path
	configPath string

	// CODEOWNERS file path
	ownersPath string

	// from global config
	ttyDisabled bool
}

const offboardLongDesc string = `[WIP] Removes a user from the \".sauced.yaml\" config and \"CODEOWNERS\" files.
Requires the user's name OR email.`

func NewConfigCommand() *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:   "offboard <username/email> [flags]",
		Short: "[WIP] Removes a user from the \".sauced.yaml\" config and \"CODEOWNERS\" files.",
		Long:  offboardLongDesc,
		Args: func(_ *cobra.Command, args []string) error {
			if !(len(args) > 0) {
				errors.New("you must provide at least one argument: the offboarding user's email/username")
			}

			opts.offboardingUsers = args

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ttyDisabled, _ = cmd.Flags().GetBool("tty-disable")
			opts.configPath, _ = cmd.Flags().GetString("config")

			opts.ownersPath, _ = cmd.Flags().GetString("owners-path")
			err := run(opts)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringP("owners-path", "o", "./CODEOWNERS", "the CODEOWNERS or OWNERS file to update")
	return cmd
}

func run(opts *Options) error {
	// read config spec
	spec, _, err := config.LoadConfig(opts.configPath)

	if err != nil {
		return err
	}

	attributions := spec.Attributions
	for _, user := range opts.offboardingUsers {
		// deletes if the user is a name (key)
		delete(attributions, user)
	}

	fmt.Print(attributions)

	err = generateOutputFile(opts.configPath, attributions)
	if err != nil {
		return err
	}

	return nil
}
