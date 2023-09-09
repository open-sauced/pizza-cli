package show

import (
	"context"
	"errors"
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	client "github.com/open-sauced/go-api/client"
)

// ContributorModel holds all the information related to a contributor
type ContributorModel struct {
	username      string
	userInfo      *client.DbUser
	userPrs       []client.DbPullRequest
	APIClient     *client.APIClient
	serverContext context.Context
}

type (
	// BackMsg: message to signal main model that we are back to dashboard when backspace is pressed
	BackMsg struct{}

	// ContributorErrMsg: message to signal that an error occurred when fetching contributor information
	ContributorErrMsg struct {
		name string
		err  error
	}
)

// InitContributor: initializes the contributorModel
func InitContributor(opts *Options) (tea.Model, error) {
	var model ContributorModel
	model.APIClient = opts.APIClient
	model.serverContext = opts.ServerContext

	return model, nil
}

func (m ContributorModel) Init() tea.Cmd { return nil }

func (m ContributorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case SelectMsg:
		m.username = msg.contributorName
		model, err := m.fetchUser()
		if err != nil {
			return m, func() tea.Msg { return ContributorErrMsg{name: msg.contributorName, err: err} }
		}
		return model, func() tea.Msg { return SuccessMsg{} }
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			return m, func() tea.Msg { return BackMsg{} }
		case "q", "esc", "ctrl+c", "ctrl+d":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m ContributorModel) View() string {
	return m.drawContributorView()
}

// fetchUser: fetches all the user information (general info, and pull requests)
func (m *ContributorModel) fetchUser() (tea.Model, error) {
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 2)
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		userInfo, err := m.fetchContributorInfo(m.username)
		if err != nil {
			errChan <- err
			return
		}
		m.userInfo = userInfo

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		userPRs, err := m.fetchContributorPRs(m.username)
		if err != nil {
			errChan <- err
			return
		}
		m.userPrs = userPRs
	}()

	wg.Wait()
	close(errChan)
	if len(errChan) > 0 {
		var allErrors error
		for err := range errChan {
			allErrors = errors.Join(allErrors, err)
		}
		return m, allErrors
	}

	return m, nil
}

// fetchContributorInfo: fetches the contributor info
func (m *ContributorModel) fetchContributorInfo(name string) (*client.DbUser, error) {
	resp, r, err := m.APIClient.UserServiceAPI.FindOneUserByUserame(m.serverContext, name).Execute()
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP failed: %d", r.StatusCode)
	}

	return resp, nil
}

// fetchContributorPRs: fetches the contributor pull requests
func (m *ContributorModel) fetchContributorPRs(name string) ([]client.DbPullRequest, error) {
	resp, r, err := m.APIClient.UserServiceAPI.FindContributorPullRequests(m.serverContext, name).Execute()
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP failed: %d", r.StatusCode)
	}

	return resp.Data, nil
}

// drawContributorView: view of the contributor model
func (m *ContributorModel) drawContributorView() string {
	return Viewport.Copy().Render(lipgloss.JoinVertical(lipgloss.Center, m.drawContributorInfo(), m.drawPullRequests()))
}

// drawContributorInfo: view of the contributor info (open issues, pr velocity, pr count, maintainer)
func (m *ContributorModel) drawContributorInfo() string {
	userOpenIssues := fmt.Sprintf("📄 Issues: %d", m.userInfo.OpenIssues)
	isUserMaintainer := fmt.Sprintf("🔨 Maintainer: %t", m.userInfo.GetIsMaintainer())
	prVelocity := fmt.Sprintf("🔥 PR Velocity (30d): %d%%", m.userInfo.RecentPullRequestVelocityCount)
	prCount := fmt.Sprintf("🚀 PR Count (30d): %d", m.userInfo.RecentPullRequestsCount)

	prStats := lipgloss.JoinVertical(lipgloss.Left, TextContainer.Render(prVelocity), TextContainer.Render(prCount))
	issuesAndMaintainer := lipgloss.JoinVertical(lipgloss.Center, TextContainer.Render(userOpenIssues), TextContainer.Render(isUserMaintainer))

	contributorInfo := lipgloss.JoinHorizontal(lipgloss.Center, prStats, issuesAndMaintainer)
	contributorView := lipgloss.JoinVertical(lipgloss.Center, m.userInfo.Login, contributorInfo)

	return SquareBorder.Render(contributorView)
}

// drawPullRequests: view of the contributor pull requests (draws the last 5 pull requests)
func (m *ContributorModel) drawPullRequests() string {
	if len(m.userPrs) == 0 {
		return ""
	}

	pullRequests := []string{}
	var numberOfPrs int

	if len(m.userPrs) > 5 {
		numberOfPrs = 5
	} else {
		numberOfPrs = len(m.userPrs)
	}

	for i := 0; i < numberOfPrs; i++ {
		prContainer := TextContainer.Render(fmt.Sprintf("#%d %s\n%s\n(%s)", m.userPrs[i].Number, m.userPrs[i].GetFullName(),
			m.userPrs[i].Title, m.userPrs[i].State))
		pullRequests = append(pullRequests, prContainer)
	}

	formattedPrs := lipgloss.JoinVertical(lipgloss.Left, pullRequests...)
	title := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Render("✨ Latest Pull Requests")

	pullRequestView := lipgloss.JoinVertical(lipgloss.Center, title, formattedPrs)
	return WidgetContainer.Render(pullRequestView)
}
