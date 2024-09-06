package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/api/auth"
	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/open-sauced/pizza-cli/pkg/utils"
)

// Options are the persistent options for the login command
type Options struct {
	// telemetry for capturing CLI events via PostHog
	telemetry *utils.PosthogCliClient
}

const (
	loginLongDesc = `Log into the OpenSauced CLI.

This command initiates the GitHub auth flow to log you into the OpenSauced CLI
by launching your browser and logging in with GitHub.`
)

func NewLoginCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log into the CLI via GitHub",
		Long:  loginLongDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			disableTelem, _ := cmd.Flags().GetBool(constants.FlagNameTelemetry)

			opts.telemetry = utils.NewPosthogCliClient(!disableTelem)

			username, err := run()

			if err != nil {
				_ = opts.telemetry.CaptureFailedLogin()
			} else {
				_ = opts.telemetry.CaptureLogin(username)
			}

			_ = opts.telemetry.Done()

			return err
		},
	}

	return cmd
}

func run() (string, error) {
	authenticator := auth.NewAuthenticator()

	username, err := authenticator.Login()
	if err != nil {
		return "", fmt.Errorf("sad: %w", err)
	}

	fmt.Println("🎉 Login successful 🎉")
	fmt.Println("Welcome aboard", username, "🍕")

	return username, nil
}
