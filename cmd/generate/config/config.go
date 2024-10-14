package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/v2/pkg/constants"
	"github.com/open-sauced/pizza-cli/v2/pkg/utils"
)

// Options for the config generation command
type Options struct {
	// the path to the git repository on disk to generate a codeowners file for
	path string

	// where the '.sauced.yaml' file will go
	outputPath string

	// whether to use interactive mode
	isInteractive bool

	// number of days to look back
	previousDays int

	// from global config
	ttyDisabled bool

	// telemetry for capturing CLI events via PostHog
	telemetry *utils.PosthogCliClient
}

const configLongDesc string = `Generates a ".sauced.yaml" configuration file for use with the Pizza CLI's codeowners command. 

This command analyzes the git history of the current repository to create a mapping 
of email addresses to GitHub usernames. `

func NewConfigCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "config path/to/repo [flags]",
		Short: "Generates a \".sauced.yaml\" config based on the current repository",
		Long:  configLongDesc,
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
			disableTelem, _ := cmd.Flags().GetBool(constants.FlagNameTelemetry)

			opts.telemetry = utils.NewPosthogCliClient(!disableTelem)
			opts.outputPath, _ = cmd.Flags().GetString("output-path")
			opts.isInteractive, _ = cmd.Flags().GetBool("interactive")
			opts.ttyDisabled, _ = cmd.Flags().GetBool("tty-disable")
			opts.previousDays, _ = cmd.Flags().GetInt("range")

			err := run(opts)
			_ = opts.telemetry.Done()

			return err
		},
	}

	cmd.PersistentFlags().StringP("output-path", "o", "./", "Directory to create the `.sauced.yaml` file.")
	cmd.PersistentFlags().BoolP("interactive", "i", false, "Whether to be interactive")
	cmd.PersistentFlags().IntP("range", "r", 90, "The number of days to analyze commit history (default 90)")
	return cmd
}

func run(opts *Options) error {
	attributionMap := make(map[string][]string)

	// Open repo
	repo, err := git.PlainOpen(opts.path)
	if err != nil {
		_ = opts.telemetry.CaptureFailedConfigGenerate()
		return fmt.Errorf("error opening repo: %w", err)
	}

	commitIter, err := repo.CommitObjects()

	if err != nil {
		_ = opts.telemetry.CaptureFailedConfigGenerate()
		return fmt.Errorf("error opening repo commits: %w", err)
	}

	now := time.Now()
	previousTime := now.AddDate(0, 0, -opts.previousDays)

	var uniqueEmails []string
	err = commitIter.ForEach(func(c *object.Commit) error {
		name := c.Author.Name
		email := c.Author.Email

		if c.Author.When.Before(previousTime) {
			return nil
		}

		if strings.Contains(name, "[bot]") {
			return nil
		}

		if opts.ttyDisabled || !opts.isInteractive {
			doesEmailExist := slices.Contains(attributionMap[name], email)
			if !doesEmailExist {
				// AUTOMATIC: set every name and associated emails
				attributionMap[name] = append(attributionMap[name], email)
			}
		} else if !slices.Contains(uniqueEmails, email) {
			uniqueEmails = append(uniqueEmails, email)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error iterating over repo commits: %w", err)
	}

	// INTERACTIVE: per unique email, set a name (existing or new or ignore)
	if opts.isInteractive && !opts.ttyDisabled {
		_ = opts.telemetry.CaptureConfigGenerateMode("interactive")
		program := tea.NewProgram(initialModel(opts, uniqueEmails))
		if _, err := program.Run(); err != nil {
			_ = opts.telemetry.CaptureFailedConfigGenerate()
			return fmt.Errorf("error running interactive mode: %w", err)
		}
	} else {
		_ = opts.telemetry.CaptureConfigGenerateMode("automatic")
		// generate an output file
		// default: `./.sauced.yaml`
		// fallback for home directories
		if opts.outputPath == "~/" {
			homeDir, _ := os.UserHomeDir()
			err := generateOutputFile(filepath.Join(homeDir, ".sauced.yaml"), attributionMap)
			if err != nil {
				_ = opts.telemetry.CaptureFailedConfigGenerate()
				return fmt.Errorf("error generating output file: %w", err)
			}
		} else {
			err := generateOutputFile(filepath.Join(opts.outputPath, ".sauced.yaml"), attributionMap)
			if err != nil {
				_ = opts.telemetry.CaptureFailedConfigGenerate()
				return fmt.Errorf("error generating output file: %w", err)
			}
		}
	}

	_ = opts.telemetry.CaptureConfigGenerate()
	return nil
}
