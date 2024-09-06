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

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

	cmd.PersistentFlags().StringP("output-path", "o", "~/", "Directory to create the `.sauced.yaml` file.")
	cmd.PersistentFlags().BoolP("interactive", "i", true, "Whether to be interactive")
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

	var uniqueEmails []string
	commitIter.ForEach(func(c *object.Commit) error {
		name := c.Author.Name
		email := c.Author.Email

		// TODO: edge case- same email multiple names
		// eg: 'coding@zeu.dev' = 'zeudev' & 'Zeu Capua'

		if !opts.isInteractive {
			doesEmailExist := slices.Contains(attributionMap[name], email)
			if !doesEmailExist {
				// AUTOMATIC: set every name and associated emails
				attributionMap[name] = append(attributionMap[name], email)
			}
		} else {
			if !slices.Contains(uniqueEmails, email) {
				uniqueEmails = append(uniqueEmails, email)
			}
		}
		return nil
	})

	// TODO: INTERACTIVE: per unique email, set a name (existing or new or ignore)
	program := tea.NewProgram(initialModel(uniqueEmails))
	if _, err := program.Run(); err != nil {
		return fmt.Errorf(err.Error())
	}

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

// Bubbletea for Interactive Mode

type model struct {
	textInput textinput.Model
	help      help.Model
	keymap    keymap

	attributionMap map[string][]string
	uniqueEmails   []string
	currentIndex   int
}

type keymap struct{}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
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

func initialModel(uniqueEmails []string) model {
	ti := textinput.New()
	ti.Placeholder = "name"
	ti.Focus()
	ti.ShowSuggestions = true

	return model{
		textInput: ti,
		help:      help.New(),
		keymap:    keymap{},

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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlI:
			m.currentIndex++
			return m, nil

		case tea.KeyEnter:
			m.attributionMap[m.textInput.Value()] = append(m.attributionMap[currentEmail], currentEmail)
			m.currentIndex++
			if m.currentIndex > len(m.attributionMap) {
				return m, tea.Quit
			}
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Found email %s - who to attribute to?: %s\n\n%s\n",
		m.uniqueEmails[m.currentIndex],
		m.textInput.View(),
		m.help.View(m.keymap),
	)
}
