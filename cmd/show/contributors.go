package show

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/browser"
	client "github.com/open-sauced/go-api/client"
)

// prItem: type for pull request to satisfy the list.Item interface
type prItem client.DbPullRequestGitHubEvents

func (i prItem) FilterValue() string { return i.PrTitle }
func (i prItem) GetRepoName() string {
	if i.RepoName != "" {
		return i.RepoName
	}
	return ""
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(prItem)
	if !ok {
		return
	}

	prTitle := i.PrTitle
	if len(prTitle) >= 60 {
		prTitle = fmt.Sprintf("%s...", prTitle[:60])
	}

	str := fmt.Sprintf("#%d %s\n%s\n(%s)", i.PrNumber, i.GetRepoName(), prTitle, i.PrState)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedItemStyle.Render("🍕 " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// ContributorModel holds all the information related to a contributor
type ContributorModel struct {
	username      string
	userInfo      *client.DbUser
	prList        list.Model
	prVelocity    float64
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
	case tea.WindowSizeMsg:
		WindowSize = msg
	case SelectMsg:
		m.username = msg.contributorName
		model, err := m.fetchUser()
		if err != nil {
			return m, func() tea.Msg { return ContributorErrMsg{name: msg.contributorName, err: err} }
		}
		return model, func() tea.Msg { return SuccessMsg{} }
	case tea.KeyMsg:
		switch msg.String() {
		case "B":
			if !m.prList.SettingFilter() {
				return m, func() tea.Msg { return BackMsg{} }
			}
		case "H":
			if !m.prList.SettingFilter() {
				m.prList.SetShowHelp(!m.prList.ShowHelp())
				return m, nil
			}
		case "O":
			if !m.prList.SettingFilter() {
				pr, ok := m.prList.SelectedItem().(prItem)
				if ok {
					err := browser.OpenURL(fmt.Sprintf("https://github.com/%s/pull/%d", pr.GetRepoName(), pr.PrNumber))
					if err != nil {
						fmt.Println("could not open pull request in browser")
					}
				}
			}
		case "q", "ctrl+c", "ctrl+d":
			if !m.prList.SettingFilter() {
				return m, tea.Quit
			}
		}
	}
	m.prList, cmd = m.prList.Update(msg)
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
		err := m.fetchContributorInfo(m.username)
		if err != nil {
			errChan <- err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := m.fetchContributorPRs(m.username)
		if err != nil {
			errChan <- err
			return
		}
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
func (m *ContributorModel) fetchContributorInfo(name string) error {
	resp, r, err := m.APIClient.UserServiceAPI.FindOneUserByUserame(m.serverContext, name).Execute()
	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		return fmt.Errorf("HTTP failed: %d", r.StatusCode)
	}

	m.userInfo = resp
	return nil
}

// fetchContributorPRs: fetches the contributor pull requests and creates pull request list
func (m *ContributorModel) fetchContributorPRs(name string) error {
	resp, r, err := m.APIClient.UserServiceAPI.FindContributorPullRequestGitHubEvents(m.serverContext, name).Range_(30).Execute()
	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		return fmt.Errorf("HTTP failed: %d", r.StatusCode)
	}

	// create contributor pull request list
	var items []list.Item
	var mergedPullRequests int
	for _, pr := range resp.Data {
		if pr.PrIsMerged {
			mergedPullRequests++
		}
		items = append(items, prItem(pr))
	}

	// calculate pr velocity
	if len(resp.Data) <= 0 {
		m.prVelocity = 0.0
	} else {
		m.prVelocity = (float64(mergedPullRequests) / float64(len(resp.Data))) * 100.0
	}

	l := list.New(items, itemDelegate{}, WindowSize.Width, 14)
	l.Title = "✨ Latest Pull Requests"
	l.Styles.Title = ListItemTitleStyle
	l.Styles.HelpStyle = HelpStyle
	l.Styles.NoItems = ItemStyle
	l.SetShowStatusBar(false)
	l.SetStatusBarItemName("pull request", "pull requests")
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			OpenPR,
			BackToDashboard,
			ToggleHelpMenu,
		}
	}

	m.prList = l
	return nil
}

// drawContributorView: view of the contributor model
func (m *ContributorModel) drawContributorView() string {
	contributorInfo := m.drawContributorInfo()

	contributorView := lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().PaddingLeft(2).Render(contributorInfo),
		WidgetContainer.Render(m.prList.View()))

	_, h := lipgloss.Size(contributorView)
	if WindowSize.Height < h {
		contributorView = lipgloss.JoinHorizontal(lipgloss.Center, contributorInfo, m.prList.View())
	}

	return contributorView
}

// drawContributorInfo: view of the contributor info (open issues, pr velocity, pr count, maintainer)
func (m *ContributorModel) drawContributorInfo() string {
	userOpenIssues := fmt.Sprintf("📄 Issues: %d", m.userInfo.OpenIssues)
	isUserMaintainer := fmt.Sprintf("🔨 Maintainer: %t", m.userInfo.GetIsMaintainer())
	prVelocity := fmt.Sprintf("🔥 PR Velocity (30d): %dd - %.0f%% merged", m.userInfo.RecentPullRequestVelocityCount, m.prVelocity)
	prCount := fmt.Sprintf("🚀 PR Count (30d): %d", m.userInfo.RecentPullRequestsCount)

	prStats := lipgloss.JoinVertical(lipgloss.Left, TextContainer.Render(prVelocity), TextContainer.Render(prCount))
	issuesAndMaintainer := lipgloss.JoinVertical(lipgloss.Center, TextContainer.Render(userOpenIssues), TextContainer.Render(isUserMaintainer))

	contributorInfo := lipgloss.JoinHorizontal(lipgloss.Center, prStats, issuesAndMaintainer)
	contributorView := lipgloss.JoinVertical(lipgloss.Center, m.userInfo.Login, contributorInfo)

	return SquareBorder.Render(contributorView)
}
