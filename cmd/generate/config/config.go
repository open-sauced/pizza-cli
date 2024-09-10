package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// Options for the config generation command
type Options struct {
	// the path to the git repository on disk to generate a codeowners file for
	path string

	// where the '.sauced.yaml' file will go
	outputPath string

	// whether to use interactive mode
	isInteractive bool
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

			opts.outputPath, _ = cmd.Flags().GetString("output-path")
			opts.isInteractive, _ = cmd.Flags().GetBool("interactive")

			return run(opts)
		},
	}

	cmd.PersistentFlags().StringP("output-path", "o", "./", "Directory to create the `.sauced.yaml` file.")
	cmd.PersistentFlags().BoolP("interactive", "i", false, "Whether to be interactive")
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

	if err != nil {
		return fmt.Errorf("error opening repo commits: %w", err)
	}

	var uniqueEmails []string
	err = commitIter.ForEach(func(c *object.Commit) error {
		name := c.Author.Name
		email := c.Author.Email

		if !opts.isInteractive {
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
	if opts.isInteractive {
		program := tea.NewProgram(initialModel(opts, uniqueEmails))
		if _, err := program.Run(); err != nil {
			return fmt.Errorf("error running interactive mode: %w", err)
		}
	} else {
		// generate an output file
		// default: `./.sauced.yaml`
		// fallback for home directories
		if opts.outputPath == "~/" {
			homeDir, _ := os.UserHomeDir()
			err := generateOutputFile(filepath.Join(homeDir, ".sauced.yaml"), attributionMap)
			if err != nil {
				return fmt.Errorf("error generating output file: %w", err)
			}
		} else {
			err := generateOutputFile(filepath.Join(opts.outputPath, ".sauced.yaml"), attributionMap)
			if err != nil {
				return fmt.Errorf("error generating output file: %w", err)
			}
		}
	}

	return nil
}

// Bubbletea for Interactive Mode

type model struct {
	textInput textinput.Model
	help      help.Model
	keymap    keymap

	opts           *Options
	attributionMap map[string][]string
	uniqueEmails   []string
	currentIndex   int
}

type keymap struct{}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "next")),
		key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "prev")),
		key.NewBinding(key.WithKeys("ctrl+i"), key.WithHelp("ctrl+i", "ignore email")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func initialModel(opts *Options, uniqueEmails []string) model {
	ti := textinput.New()
	ti.Placeholder = "name"
	ti.Focus()
	ti.ShowSuggestions = true

	return model{
		textInput: ti,
		help:      help.New(),
		keymap:    keymap{},

		opts:           opts,
		attributionMap: make(map[string][]string),
		uniqueEmails:   uniqueEmails,
		currentIndex:   0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	currentEmail := m.uniqueEmails[m.currentIndex]

	existingUsers := make([]string, 0, len(m.attributionMap))
	for k := range m.attributionMap {
		existingUsers = append(existingUsers, k)
	}

	m.textInput.SetSuggestions(existingUsers)

	keyMsg, ok := msg.(tea.KeyMsg)

	if ok {
		switch keyMsg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlI:
			m.currentIndex++
			return m, nil

		case tea.KeyEnter:
			if len(strings.Trim(m.textInput.Value(), " ")) == 0 {
				return m, nil
			}
			m.attributionMap[m.textInput.Value()] = append(m.attributionMap[m.textInput.Value()], currentEmail)
			m.textInput.Reset()
			if m.currentIndex+1 >= len(m.uniqueEmails) {
				return m, runOutputGeneration(m.opts, m.attributionMap)
			}

			m.currentIndex++
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m model) View() string {
	currentEmail := ""
	if m.currentIndex < len(m.uniqueEmails) {
		currentEmail = m.uniqueEmails[m.currentIndex]
	}

	return fmt.Sprintf(
		"Found email %s - who to attribute to?: \n%s\n\n%s\n",
		currentEmail,
		m.textInput.View(),
		m.help.View(m.keymap),
	)
}

func runOutputGeneration(opts *Options, attributionMap map[string][]string) tea.Cmd {
	// generate an output file
	// default: `./.sauced.yaml`
	// fallback for home directories
	return func() tea.Msg {
		if opts.outputPath == "~/" {
			homeDir, _ := os.UserHomeDir()
			err := generateOutputFile(filepath.Join(homeDir, ".sauced.yaml"), attributionMap)
			if err != nil {
				return fmt.Errorf("error generating output file: %w", err)
			}
		} else {
			err := generateOutputFile(filepath.Join(opts.outputPath, ".sauced.yaml"), attributionMap)
			if err != nil {
				return fmt.Errorf("error generating output file: %w", err)
			}
		}

		return tea.Quit()
	}

}
