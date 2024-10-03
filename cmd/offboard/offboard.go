package offboard

import (
	"errors"
	"fmt"
	"slices"

	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/v2/pkg/config"
	"github.com/open-sauced/pizza-cli/v2/pkg/constants"
	"github.com/open-sauced/pizza-cli/v2/pkg/utils"
)

type Options struct {
	offboardingUsers []string

	// config file path
	configPath string

	// repository path
	path string

	// from global config
	ttyDisabled bool

	// telemetry for capturing CLI events via PostHog
	telemetry *utils.PosthogCliClient
}

const offboardLongDesc string = `CAUTION: Experimental Command. Removes users from the \".sauced.yaml\" config and \"CODEOWNERS\" files.
Requires the users' name OR email.`

func NewConfigCommand() *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:   "offboard <username/email> [flags]",
		Short: "CAUTION: Experimental Command. Removes users from the \".sauced.yaml\" config and \"CODEOWNERS\" files.",
		Long:  offboardLongDesc,
		Args: func(_ *cobra.Command, args []string) error {
			if !(len(args) > 0) {
				return errors.New("you must provide at least one argument: the offboarding user's email/username")
			}

			opts.offboardingUsers = args

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.ttyDisabled, _ = cmd.Flags().GetBool("tty-disable")
			opts.configPath, _ = cmd.Flags().GetString("config")
			disableTelem, _ := cmd.Flags().GetBool(constants.FlagNameTelemetry)

			opts.telemetry = utils.NewPosthogCliClient(!disableTelem)

			opts.path, _ = cmd.Flags().GetString("path")
			err := run(opts)
			_ = opts.telemetry.Done()

			return err
		},
	}

	cmd.PersistentFlags().StringP("path", "p", "./", "the path to the repository")
	return cmd
}

func run(opts *Options) error {
	spec, _, err := config.LoadConfig(opts.configPath)

	if err != nil {
		_ = opts.telemetry.CaptureFailedOffboard()
		return fmt.Errorf("error loading config: %v", err)
	}

	var offboardingNames []string
	attributions := spec.Attributions
	for _, user := range opts.offboardingUsers {
		added := false

		// deletes if the user is a name (key)
		delete(attributions, user)

		// delete if the user is an email (value)
		for k, v := range attributions {
			if slices.Contains(v, user) {
				offboardingNames = append(offboardingNames, k)
				delete(attributions, k)
				added = true
			}
		}

		if !added {
			offboardingNames = append(offboardingNames, user)
		}
	}

	err = generateConfigFile(opts.configPath, attributions)
	if err != nil {
		_ = opts.telemetry.CaptureFailedOffboard()
		return fmt.Errorf("error generating config file: %v", err)
	}

	err = generateOwnersFile(opts.path, offboardingNames)
	if err != nil {
		_ = opts.telemetry.CaptureFailedOffboard()
		return fmt.Errorf("error generating owners file: %v", err)
	}

	_ = opts.telemetry.CaptureOffboard()
	return nil
}
