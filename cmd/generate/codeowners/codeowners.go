package codeowners

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/jpmcb/gopherlogs"
	"github.com/jpmcb/gopherlogs/pkg/colors"
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/pkg/config"
	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/open-sauced/pizza-cli/pkg/logging"
	"github.com/open-sauced/pizza-cli/pkg/utils"
)

// Options for the codeowners generation command
type Options struct {
	// the path to the git repository on disk to generate a codeowners file for
	path string

	// whether to generate an agnostic "OWENRS" style codeowners file.
	// The default should be to generate a GitHub style "CODEOWNERS" file.
	ownersStyleFile bool

	// the number of days to look back
	previousDays int

	logger   gopherlogs.Logger
	tty      bool
	loglevel int

	// telemetry for capturing CLI events via PostHog
	telemetry *utils.PosthogCliClient

	config *config.Spec
}

const codeownersLongDesc string = `Generates a CODEOWNERS file for a given git repository. The generated file specifies up to 3 owners for EVERY file in the git tree based on the number of lines touched in that specific file over the specified range of time.

Configuration:
The command requires a .sauced.yaml file for accurate attribution. This file maps 
commit email addresses to GitHub usernames. The command looks for this file in two locations:

1. In the root of the specified repository path
2. In the user's home directory (~/.sauced.yaml) if not found in the repository

If you run the command on a specific path, it will first look for .sauced.yaml in that 
path. If not found, it will fall back to ~/.sauced.yaml.`

func NewCodeownersCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "codeowners path/to/repo [flags]",
		Short: "Generate a CODEOWNERS file for a GitHub repository using a \"~/.sauced.yaml\" config",
		Long:  codeownersLongDesc,
		Example: `
# Generate CODEOWNERS file for the current directory
pizza generate codeowners .

# Generate CODEOWNERS file for a specific repository
pizza generate codeowners /path/to/your/repo

# Generate CODEOWNERS file analyzing the last 180 days
pizza generate codeowners . --range 180

# Generate an OWNERS style file instead of CODEOWNERS
pizza generate codeowners . --owners-style-file

# Specify a custom location for the .sauced.yaml file
pizza generate codeowners . --config /path/to/.sauced.yaml
		`,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you must provide exactly one argument: the path to the repository")
			}

			path := args[0]

			// Validate that the path is a real path on disk and accessible by the user
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			if _, err := os.Stat(absPath); os.IsNotExist(err) {
				return fmt.Errorf("the provided path does not exist: %w", err)
			}

			opts.path = absPath
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error

			disableTelem, _ := cmd.Flags().GetBool(constants.FlagNameTelemetry)

			opts.telemetry = utils.NewPosthogCliClient(!disableTelem)

			configPath, _ := cmd.Flags().GetString("config")
			opts.config, err = config.LoadConfig(configPath)
			if err != nil {
				return err
			}

			opts.ownersStyleFile, _ = cmd.Flags().GetBool("owners-style-file")
			opts.previousDays, _ = cmd.Flags().GetInt("range")
			opts.tty, _ = cmd.Flags().GetBool("tty-disable")

			loglevelS, _ := cmd.Flags().GetString("log-level")

			switch loglevelS {
			case "error":
				opts.loglevel = logging.LogError
			case "warn":
				opts.loglevel = logging.LogWarn
			case "info":
				opts.loglevel = logging.LogInfo
			case "debug":
				opts.loglevel = logging.LogDebug
			}

			err = run(opts, cmd)

			_ = opts.telemetry.Done()

			return err
		},
	}

	cmd.PersistentFlags().IntP("range", "r", 90, "The number of days to analyze commit history (default 90)")
	cmd.PersistentFlags().Bool("owners-style-file", false, "Generate an agnostic OWNERS style file instead of CODEOWNERS.")

	return cmd
}

func run(opts *Options, cmd *cobra.Command) error {
	var err error
	opts.logger, err = gopherlogs.NewLogger(
		gopherlogs.WithLogVerbosity(opts.loglevel),
		gopherlogs.WithTty(!opts.tty),
	)
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Built logger with log level: %d\n", opts.loglevel)

	repo, err := git.PlainOpen(opts.path)
	if err != nil {
		_ = opts.telemetry.CaptureFailedCodeownersGenerate()
		return fmt.Errorf("error opening repo: %w", err)
	}
	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Opened repo at: %s\n", opts.path)

	processOptions := ProcessOptions{
		repo,
		opts.previousDays,
		opts.path,
		opts.logger,
	}
	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Looking back %d days\n", opts.previousDays)

	codeowners, err := processOptions.process()
	if err != nil {
		_ = opts.telemetry.CaptureFailedCodeownersGenerate()
		return fmt.Errorf("error traversing git log: %w", err)
	}

	// Bootstrap codeowners
	var outputPath string
	if opts.ownersStyleFile {
		outputPath = filepath.Join(opts.path, "OWNERS")
	} else {
		outputPath = filepath.Join(opts.path, "CODEOWNERS")
	}

	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Processing codeowners file at: %s\n", outputPath)
	err = generateOutputFile(codeowners, outputPath, opts, cmd)
	if err != nil {
		_ = opts.telemetry.CaptureFailedCodeownersGenerate()
		return fmt.Errorf("error generating github style codeowners file: %w", err)
	}
	opts.logger.V(logging.LogInfo).Style(0, colors.FgGreen).Infof("Finished generating file: %s\n", outputPath)
	_ = opts.telemetry.CaptureCodeownersGenerate()

	opts.logger.V(logging.LogInfo).Style(0, colors.FgCyan).Infof("\nCreate an OpenSauced Contributor Insight to get metrics and insights on these codeowners:\n")
	opts.logger.V(logging.LogInfo).Style(0, colors.FgCyan).Infof("$ pizza generate insight " + opts.path + "\n")
	_ = opts.telemetry.CaptureCodeownersGenerateContributorInsight()

	return nil
}
