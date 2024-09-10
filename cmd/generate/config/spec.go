package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

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
		key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "next suggestion")),
		key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "prev suggestion")),
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
	ti.Placeholder = "username"
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
