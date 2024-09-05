package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// Options for the config generation command
type Options struct {
	// the path to the git repository on disk to generate a codeowners file for
	path   string

	// where the '.sauced.yaml' file will go
	outputPath string
}

const configLongDesc string = `WARNING: Proof of concept feature.

Generates a ~/.sauced.yaml configuration file. The attribution of emails to given entities
is based on the repository this command is ran in.`

func NewConfigCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "config path/to/repo [flags]",
		Short: "Generates a \"~/.sauced.yaml\" config based on the current repository",
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
			// TODO: error checking based on given command

			opts.outputPath, _ = cmd.Flags().GetString("output-path");

			return run(opts)
		},
	}

	cmd.PersistentFlags().StringP("output-path", "o", "~/", "Directory to create the `.sauced.yaml` file.")
	return cmd
}

func run(opts *Options) error {
	attributionMap := make(map[string][]string)

	// Open repo
	repo, err := git.PlainOpen(opts.path)
	if err != nil {
		return fmt.Errorf("error opening repo: %w", err)
	}

	commitIter, err := repo.CommitObjects()

	commitIter.ForEach(func(c *object.Commit) error {
		name := c.Author.Name
		email := c.Author.Email

		// TODO: edge case- same email multiple names
		// eg: 'coding@zeu.dev' = 'zeudev' & 'Zeu Capua'

		// AUTOMATIC: set every name and associated emails
		doesEmailExist := slices.Contains(attributionMap[name], email)
		if !doesEmailExist {
			attributionMap[name] = append(attributionMap[name], email)
		}

		// TODO: INTERACTIVE: per unique email, set a name (existing or new)

		return nil
	})

	// generate an output file
	// default: `~/.sauced.yaml`	
	if opts.outputPath == "~/" {
		homeDir, _ := os.UserHomeDir()
		generateOutputFile(filepath.Join(homeDir, ".sauced.yaml"), attributionMap)
	} else {
		generateOutputFile(filepath.Join(opts.outputPath, ".sauced.yaml"), attributionMap)
	}


	return nil
}
