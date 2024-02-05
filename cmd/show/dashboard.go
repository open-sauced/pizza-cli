package show

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	client "github.com/open-sauced/go-api/client"
)

const (
	newContributorsView = iota
	alumniContributorsView
)

// DashboardModel holds all the information related to the repository queried (issues, stars, new contributors, alumni contributors)
type DashboardModel struct {
	newContributorsTable    table.Model
	alumniContributorsTable table.Model
	RepositoryInfo          *client.DbRepo
	contributorErr          string
	tableView               int
	queryOptions            [3]int
	APIClient               *client.APIClient
	serverContext           context.Context
}

// SelectMsg: message to signal the main model that we want to go to the contributor model when 'enter' is pressed
type SelectMsg struct {
	contributorName string
}

// FetchRepoInfo: initializes the dashboard model
func InitDashboard(opts *Options) (tea.Model, error) {
	var model DashboardModel
	err := validateShowQuery(opts)
	if err != nil {
		return model, err
	}

	resp, r, err := opts.APIClient.RepositoryServiceAPI.FindOneByOwnerAndRepo(opts.ServerContext, opts.Owner, opts.RepoName).Execute()
	if err != nil {
		return model, err
	}

	if r.StatusCode != 200 {
		return model, fmt.Errorf("HTTP status: %d", r.StatusCode)
	}

	// configuring the dashboardModel
	model.RepositoryInfo = resp
	model.queryOptions = [3]int{opts.Page, opts.Limit, opts.Range}
	model.APIClient = opts.APIClient
	model.serverContext = opts.ServerContext

	// fetching the contributor tables
	err = model.FetchAllContributors()
	if err != nil {
		return model, err
	}

	return model, nil
}

func (m DashboardModel) Init() tea.Cmd {
	return nil
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		WindowSize = msg

	case ErrMsg:
		fmt.Printf("Failed to retrieve contributors table data: %s", msg.err.Error())
		return m, tea.Quit

	case ContributorErrMsg:
		m.contributorErr = fmt.Sprintf("🚧 could not fetch %s: %s", msg.name, msg.err.Error())
	default:
		m.contributorErr = ""

	case tea.KeyMsg:
		switch msg.String() {
		case "right", "l":
			m.tableView = (m.tableView + 1) % 2
		case "left", "h":
			if m.tableView-1 <= 0 {
				m.tableView = 0
			} else {
				m.tableView--
			}
		case "q", "esc", "ctrl+c", "ctrl+d":
			return m, tea.Quit
		case "enter":
			switch m.tableView {
			case newContributorsView:
				if len(m.newContributorsTable.Rows()) > 0 {
					return m, func() tea.Msg { return SelectMsg{contributorName: m.newContributorsTable.SelectedRow()[1]} }
				}
			case alumniContributorsView:
				if len(m.alumniContributorsTable.Rows()) > 0 {
					return m, func() tea.Msg { return SelectMsg{contributorName: m.alumniContributorsTable.SelectedRow()[1]} }
				}
			}
		}

		switch m.tableView {
		case newContributorsView:
			m.newContributorsTable, cmd = m.newContributorsTable.Update(msg)
		case alumniContributorsView:
			m.alumniContributorsTable, cmd = m.alumniContributorsTable.Update(msg)
		}
	}

	return m, cmd
}

func (m DashboardModel) View() string {
	return m.drawDashboardView()
}

// drawTitle: view of PIZZA
func (m *DashboardModel) drawTitle() string {
	titleRunes1 := []rune{'█', '█', '█', '█', '█', '█', '╗', ' ', '█', '█', '╗', '█', '█', '█', '█', '█', '█', '█', '╗', '█', '█', '█', '█', '█', '█', '█', '╗', ' ', '█', '█', '█', '█', '█', '╗', ' '}
	titleRunes2 := []rune{'█', '█', '╔', '═', '═', '█', '█', '╗', '█', '█', '║', '╚', '═', '═', '█', '█', '█', '╔', '╝', '╚', '═', '═', '█', '█', '█', '╔', '╝', '█', '█', '╔', '═', '═', '█', '█', '╗'}
	titleRunes3 := []rune{'█', '█', '█', '█', '█', '█', '╔', '╝', '█', '█', '║', ' ', ' ', '█', '█', '█', '╔', '╝', ' ', ' ', ' ', '█', '█', '█', '╔', '╝', ' ', '█', '█', '█', '█', '█', '█', '█', '║'}
	titleRunes4 := []rune{'█', '█', '╔', '═', '═', '═', '╝', ' ', '█', '█', '║', ' ', '█', '█', '█', '╔', '╝', ' ', ' ', ' ', '█', '█', '█', '╔', '╝', ' ', ' ', '█', '█', '╔', '═', '═', '█', '█', '║'}
	titleRunes5 := []rune{'█', '█', '║', ' ', ' ', ' ', ' ', ' ', '█', '█', '║', '█', '█', '█', '█', '█', '█', '█', '╗', '█', '█', '█', '█', '█', '█', '█', '╗', '█', '█', '║', ' ', ' ', '█', '█', '║'}
	titleRunes6 := []rune{'╚', '═', '╝', ' ', ' ', ' ', ' ', ' ', '╚', '═', '╝', '╚', '═', '═', '═', '═', '═', '═', '╝', '╚', '═', '═', '═', '═', '═', '═', '╝', '╚', '═', '╝', ' ', ' ', '╚', '═', '╝'}

	title1 := lipgloss.JoinHorizontal(lipgloss.Left, string(titleRunes1))
	title2 := lipgloss.JoinHorizontal(lipgloss.Left, string(titleRunes2))
	title3 := lipgloss.JoinHorizontal(lipgloss.Left, string(titleRunes3))
	title4 := lipgloss.JoinHorizontal(lipgloss.Left, string(titleRunes4))
	title5 := lipgloss.JoinHorizontal(lipgloss.Left, string(titleRunes5))
	title6 := lipgloss.JoinHorizontal(lipgloss.Left, string(titleRunes6))
	title := lipgloss.JoinVertical(lipgloss.Center, title1, title2, title3, title4, title5, title6)
	titleView := lipgloss.NewStyle().Foreground(Color).Render(title)

	return titleView
}

// drawRepositoryInfo: view of the repository info (name, stars, size, and issues)
func (m *DashboardModel) drawRepositoryInfo() string {
	repoName := lipgloss.NewStyle().Bold(true).AlignHorizontal(lipgloss.Center).Render(fmt.Sprintf("Repository: %s", m.RepositoryInfo.FullName))
	repoStars := fmt.Sprintf("🌟 stars: %d", m.RepositoryInfo.Stars)
	repoSize := fmt.Sprintf("💾 size: %dB", m.RepositoryInfo.Size)
	repoIssues := fmt.Sprintf("📄 issues: %d", m.RepositoryInfo.Issues)
	repoForks := fmt.Sprintf("⑂ forks: %d", m.RepositoryInfo.Forks)

	issuesAndForks := lipgloss.JoinVertical(lipgloss.Center, TextContainer.Render(repoIssues), TextContainer.Render(repoForks))
	sizeAndStars := lipgloss.JoinVertical(lipgloss.Left, TextContainer.Render(repoSize), TextContainer.Render(repoStars))

	repoGeneralSection := lipgloss.JoinHorizontal(lipgloss.Center, issuesAndForks, sizeAndStars)
	repositoryInfoView := lipgloss.JoinVertical(lipgloss.Center, repoName, repoGeneralSection)
	frame := SquareBorder.Render(repositoryInfoView)

	return frame
}

// drawMetrics: view of metrics includes.
// - new contributors table
// - alumni contributors table
func (m *DashboardModel) drawMetrics() string {
	var newContributorsDisplay, alumniContributorsDisplay string

	switch m.tableView {
	case newContributorsView:
		newContributorsDisplay = lipgloss.JoinVertical(lipgloss.Center, TableTitle.Render("🍕 New Contributors"), ActiveStyle.Render(m.newContributorsTable.View()))
		alumniContributorsDisplay = lipgloss.JoinVertical(lipgloss.Center, TableTitle.Render("🍁 Alumni Contributors"), InactiveStyle.Render(m.alumniContributorsTable.View()))
	case alumniContributorsView:
		newContributorsDisplay = lipgloss.JoinVertical(lipgloss.Center, TableTitle.Render("🍕 New Contributors"), InactiveStyle.Render(m.newContributorsTable.View()))
		alumniContributorsDisplay = lipgloss.JoinVertical(lipgloss.Center, TableTitle.Render("🍁 Alumni Contributors"), ActiveStyle.Render(m.alumniContributorsTable.View()))
	}

	contributorsMetrics := lipgloss.JoinHorizontal(lipgloss.Center, WidgetContainer.Render(newContributorsDisplay), WidgetContainer.Render(alumniContributorsDisplay))

	return contributorsMetrics
}

// drawDashboardView: this is the main model view (shows repository info and tables)
func (m *DashboardModel) drawDashboardView() string {
	if WindowSize.Width == 0 {
		return "Loading..."
	}
	titleView, repoInfoView, metricsView := m.drawTitle(), m.drawRepositoryInfo(), m.drawMetrics()
	mainView := lipgloss.JoinVertical(lipgloss.Center, titleView, repoInfoView, metricsView, m.contributorErr)

	_, h := lipgloss.Size(mainView)
	if WindowSize.Height < h {
		contentLeft := lipgloss.JoinVertical(lipgloss.Center, titleView, repoInfoView)
		contentRight := lipgloss.JoinVertical(lipgloss.Center, metricsView, m.contributorErr)
		mainView = lipgloss.JoinHorizontal(lipgloss.Center, contentLeft, contentRight)
	}
	frame := Viewport.Render(mainView)
	return frame
}

// validateShowQuery: validates fields set to query the contributor tables
func validateShowQuery(opts *Options) error {
	if opts.Limit < 1 {
		return errors.New("--limit flag must be a positive integer value")
	}
	if opts.Range < 1 {
		return errors.New("--range flag must be a positive integer value")
	}

	if opts.Page < 1 {
		return errors.New("--page flag must be a positive integer value")
	}

	return nil
}

// FetchAllContributors: fetchs and sets all the contributors (new, alumni)
func (m *DashboardModel) FetchAllContributors() error {
	var (
		errorChan = make(chan error, 2)
		wg        sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		newContributors, err := m.FetchNewContributors()
		if err != nil {
			errorChan <- err
			return
		}
		m.newContributorsTable = setupContributorsTable(newContributors)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		alumniContributors, err := m.FetchAlumniContributors()
		if err != nil {
			errorChan <- err
			return
		}
		m.alumniContributorsTable = setupContributorsTable(alumniContributors)
	}()

	wg.Wait()
	close(errorChan)
	if len(errorChan) > 0 {
		var allErrors error
		for err := range errorChan {
			allErrors = errors.Join(allErrors, err)
		}
		return allErrors
	}

	return nil
}

// FetchNewContributors: Returns all the new contributors
func (m *DashboardModel) FetchNewContributors() ([]client.DbPullRequestContributor, error) {
	resp, r, err := m.APIClient.ContributorsServiceAPI.NewPullRequestContributors(m.serverContext).Page(int32(m.queryOptions[0])).
		Limit(int32(m.queryOptions[1])).Repos(m.RepositoryInfo.FullName).Execute()
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status: %d", r.StatusCode)
	}

	return resp.Data, nil

}

// FetchAlumniContributors: Returns all alumni contributors
func (m *DashboardModel) FetchAlumniContributors() ([]client.DbPullRequestContributor, error) {
	resp, r, err := m.APIClient.ContributorsServiceAPI.FindAllChurnPullRequestContributors(m.serverContext).
		Page(int32(m.queryOptions[0])).Limit(int32(m.queryOptions[1])).
		Range_(int32(m.queryOptions[2])).Repos(m.RepositoryInfo.FullName).Execute()
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status: %d", r.StatusCode)
	}

	return resp.Data, nil
}

// setupContributorsTable: sets the contributor table UI
func setupContributorsTable(contributors []client.DbPullRequestContributor) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Name", Width: 20},
	}

	rows := []table.Row{}

	for i, contributor := range contributors {
		rows = append(rows, table.Row{strconv.Itoa(i), contributor.AuthorLogin})
	}

	contributorTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF4500"))

	contributorTable.SetStyles(s)
	return contributorTable
}
