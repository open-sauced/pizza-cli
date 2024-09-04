package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/pkg/config"
)

// Options for the codeowners generation command
type Options struct {
	// the path to the git repository on disk to generate a codeowners file for
	path string

	previousDays int
	tty          bool
	loglevel     int

	config *config.Spec
}

const codeownersLongDesc string = `WARNING: Proof of concept feature.

Generates a ~/.sauced.yaml configuration file. The attribution of emails to given entities
is based on the repository this command is ran in.`

func NewConfigCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "config path/to/repo [flags]",
		Short: "Generates a \"~/.sauced.yaml\" config based on the current repository",
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
			// TODO: error checking based on given command

			return run(opts, cmd)
		},
	}

	return cmd
}

func run(opts *Options, cmd *cobra.Command) error {
	configuration := &config.Spec{}
	fmt.Println("CONFIG", configuration)

	// Open repo
	repo, err := git.PlainOpen(opts.path)
	if err != nil {
		return fmt.Errorf("error opening repo: %w", err)
	}

	commitIter, err := repo.CommitObjects()

	commitIter.ForEach(func(c *object.Commit) error {
		fmt.Println("COMMIT", c.Author.Email, c.Author.Name)
		return nil
	})

	return nil
}
