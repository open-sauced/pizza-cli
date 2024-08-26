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
	"github.com/open-sauced/pizza-cli/pkg/logging"
)

// Options for the codeowners generation command
type Options struct {
	// the path to the git repository on disk to generate a codeowners file for
	path string

	// whether to make a github style codeowners file or not
	githubCodeowners bool

	// the number of days to look back
	previousDays int

	tty      bool
	loglevel int

	config *config.Spec
}

const codeownersLongDesc string = `WARNING: Proof of concept feature.

Generates a CODEOWNERS file for a given git repository. This uses a .sauced.yaml
configuration to attribute emails with GitHub usernames.

The generated CODEOWNERS file specifies up to 3 owners for EVERY file based on the
number of lines touched in that specific file over the lifetime of commits.`

func NewCodeownersCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "codeowners path/to/repo [flags]",
		Short: "Generates a CODEOWNERS file for a given repository using a .sauced.yaml config",
		Long:  codeownersLongDesc,
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

			configPath, _ := cmd.Flags().GetString("config")
			opts.config, err = config.LoadConfig(configPath)
			if err != nil {
				return err
			}

			opts.githubCodeowners, _ = cmd.Flags().GetBool("github-codeowners")
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

			return run(opts)
		},
	}

	cmd.PersistentFlags().IntP("range", "r", 90, "The number of days to lookback")
	cmd.PersistentFlags().Bool("github-codeowners", false, "Whether to generate a GitHub styles CODEOWNERS file that uses attributions in a provided config")

	return cmd
}

func run(opts *Options) error {
	logger, err := gopherlogs.NewLogger(
		gopherlogs.WithLogVerbosity(opts.loglevel),
		gopherlogs.WithTty(!opts.tty),
	)
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Built logger with log level: %d\n", opts.loglevel)

	repo, err := git.PlainOpen(opts.path)
	if err != nil {
		return fmt.Errorf("error opening repo: %w", err)
	}
	logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Opened repo at: %s\n", opts.path)

	processOptions := ProcessOptions{
		repo,
		opts.previousDays,
		opts.path,
		logger,
	}
	logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Looking back %d days\n", opts.previousDays)

	codeowners, err := processOptions.process()
	if err != nil {
		return fmt.Errorf("error traversing git log: %w", err)
	}

	// Bootstrap codeowners
	outputPath := ""
	if opts.githubCodeowners {
		outputPath = filepath.Join(opts.path, "CODEOWNERS")
	} else {
		outputPath = filepath.Join(opts.path, "OWNERS")
	}

	logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Processing codeowners file at: %s\n", outputPath)
	err = generateOutputFile(codeowners, outputPath, opts.githubCodeowners, opts.config)
	if err != nil {
		return fmt.Errorf("error generating github style codeowners file: %w", err)
	}
	logger.V(logging.LogInfo).Style(0, colors.FgGreen).Infof("Finished generating file: %s\n", outputPath)

	return nil
}