package offboard

import (
	"fmt"
	"errors"

	"github.com/spf13/cobra"
)

type Options struct {
	offboardingUsers []string
}

const configLongDesc string = `[WIP] Removes a user from the \".sauced.yaml\" config and \"CODEOWNERS\" files.
Requires the user's name OR email.`

func NewConfigCommand() *cobra.Command {
	options := &Options{}
	cmd := &cobra.Command{
		Use:   "offboard <username/email> [flags]",
		Short: "[WIP] Removes a user from the \".sauced.yaml\" config and \"CODEOWNERS\" files.",
		Long:  configLongDesc,
		Args: func(_ *cobra.Command, args []string) error {
			if !(len(args) > 0) {
				errors.New("you must provide at least one argument: the offboarding user's email/username")
			}

			options.offboardingUsers = args

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}
