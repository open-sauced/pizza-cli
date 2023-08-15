package show

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	dashboardView = iota
	contributorView
)

type (
	// sessionState serves as the variable to reference when looking for which model the user is in
	sessionState int
)

// MainModel: the main model is the central state manager of the TUI, decides which model is focused based on certain commands
type MainModel struct {
	state       sessionState
	dashboard   tea.Model
	contributor tea.Model
	authorName  string
}

// View: the view of the TUI
func (m MainModel) View() string {
	switch m.state {
	case contributorView:
		return m.contributor.View()
	default:
		return m.dashboard.View()
	}
}

// Init: initial IO before program start
func (m MainModel) Init() tea.Cmd { return nil }

// Update: Handle IO and Commands
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case BackMsg:
		m.state = dashboardView
	case SelectMsg:
		m.state = contributorView
		m.authorName = msg.contributorName
	}

	switch m.state {
	case dashboardView:
		newDashboard, newCmd := m.dashboard.Update(msg)
		m.dashboard = newDashboard
		cmd = newCmd
	case contributorView:
		// TODO: process cmd error if contributor fails to fetch
		var err error
		m.contributor, err = InitContributor(m.authorName)
		if err != nil {
			m.state = dashboardView
			return m, func() tea.Msg { return ContributorErrMsg{name: m.authorName, err: err} }
		}
		newContributor, newCmd := m.contributor.Update(msg)
		m.contributor = newContributor
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func pizzaTUI(opts *Options) error {
	dashboardModel, err := InitDashboard(opts)
	if err != nil {
		return err
	}

	model := MainModel{dashboard: dashboardModel}
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return fmt.Errorf("Error running program: %s", err.Error())
	}

	return nil
}
