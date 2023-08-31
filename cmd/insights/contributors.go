package insights

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/open-sauced/go-api/client"
	"github.com/open-sauced/pizza-cli/pkg/api"
	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/open-sauced/pizza-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type Options struct {
	// APIClient is the http client for making calls to the open-sauced api
	APIClient *client.APIClient

	// Repos is the array of git repository urls
	Repos []string

	// FilePath is the path to yaml file containing an array of git repository urls
	FilePath string

	// Period is the number of days, used for query filtering
	Period int32
}

// NewContributorsCommand returns a new cobra command for 'pizza insights contributors'
func NewContributorsCommand() *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:   "contributors [flags]",
		Short: "Gather insights about contributors of indexed git repositories",
		Long:  "Gather insights about contributors of indexed git repositories. This command will show new, recent, alumni, repeat contributors for each git repository",
		Args: func(cmd *cobra.Command, args []string) error {
			fileFlag := cmd.Flags().Lookup(constants.FlagNameFile)
			if !fileFlag.Changed && len(args) == 0 {
				return fmt.Errorf("must specify git repository url argument(s) or provide %s flag", fileFlag.Name)
			}
			opts.Repos = append(opts.Repos, args...)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			endpointURL, _ := cmd.Flags().GetString(constants.FlagNameEndpoint)
			opts.APIClient = api.NewGoClient(endpointURL)
			return run(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.FilePath, constants.FlagNameFile, "f", "", "Path to yaml file containing an array of git repository urls")
	cmd.Flags().Int32VarP(&opts.Period, constants.FlagNamePeriod, "p", 30, "Number of days, used for query filtering")
	return cmd
}

func run(opts *Options) error {
	repositories, err := utils.HandleRepositoryValues(opts.Repos, opts.FilePath)
	if err != nil {
		return err
	}
	var (
		waitGroup = new(sync.WaitGroup)
		errorChan = make(chan error, len(repositories))
	)
	for url := range repositories {
		repoURL := url
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			repoContributorsInsights, err := findAllRepositoryContributorsInsights(context.TODO(), opts, repoURL)
			if err != nil {
				errorChan <- err
				return
			}
			repoContributorsInsights.RenderTable()
		}()
	}
	waitGroup.Wait()
	close(errorChan)
	var allErrors error
	for err = range errorChan {
		allErrors = errors.Join(allErrors, err)
	}
	return allErrors
}

type repositoryContributorsInsights struct {
	RepoID  int32
	RepoURL string
	New     []string
	Recent  []string
	Alumni  []string
	Repeat  []string
}

func (rci *repositoryContributorsInsights) RenderTable() {
	if rci == nil {
		return
	}
	rows := []table.Row{
		{"New contributors", strconv.Itoa(len(rci.New))},
		{"Recent contributors", strconv.Itoa(len(rci.Recent))},
		{"Alumni contributors", strconv.Itoa(len(rci.Alumni))},
		{"Repeat contributors", strconv.Itoa(len(rci.Repeat))},
	}
	columns := []table.Column{
		{
			Title: "Repository URL",
			Width: 20,
		},
		{
			Title: rci.RepoURL,
			Width: len(rci.RepoURL),
		},
	}
	styles := table.DefaultStyles()
	styles.Header.MarginTop(1)
	styles.Selected = lipgloss.NewStyle()
	repoTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(len(rows)),
		table.WithStyles(styles),
	)
	fmt.Println(repoTable.View())
}

func findRepositoryByOwnerAndRepoName(ctx context.Context, apiClient *client.APIClient, repoURL string) (*client.DbRepo, error) {
	owner, repoName, err := utils.GetOwnerAndRepoFromURL(repoURL)
	if err != nil {
		return nil, fmt.Errorf("could not extract owner and repo from url: %w", err)
	}
	repo, response, err := apiClient.RepositoryServiceAPI.FindOneByOwnerAndRepo(ctx, owner, repoName).Execute()
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			message := fmt.Sprintf("repository %s is either non-existent or has not been indexed yet", repoURL)
			fmt.Println("ignoring repository issue:", message)
			return nil, nil
		}
		return nil, fmt.Errorf("error while calling 'RepositoryServiceAPI.FindOneByOwnerAndRepo' with owner %q and repo %q: %w", owner, repoName, err)
	}
	return repo, nil
}

func findAllRepositoryContributorsInsights(ctx context.Context, opts *Options, repoURL string) (*repositoryContributorsInsights, error) {
	repo, err := findRepositoryByOwnerAndRepoName(ctx, opts.APIClient, repoURL)
	if err != nil {
		return nil, fmt.Errorf("could not get contributors insights for repository %s: %w", repoURL, err)
	}
	if repo == nil {
		return nil, nil
	}
	repoContributorsInsights := &repositoryContributorsInsights{
		RepoID:  repo.Id,
		RepoURL: repoURL,
	}
	var (
		waitGroup = new(sync.WaitGroup)
		errorChan = make(chan error, 4)
	)
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		response, err := findNewRepositoryContributors(ctx, opts.APIClient, repo.Id, opts.Period)
		if err != nil {
			errorChan <- err
			return
		}
		for _, data := range response.Data {
			repoContributorsInsights.New = append(repoContributorsInsights.New, data.AuthorLogin)
		}
	}()
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		response, err := findRecentRepositoryContributors(ctx, opts.APIClient, repo.Id, opts.Period)
		if err != nil {
			errorChan <- err
			return
		}
		for _, data := range response.Data {
			repoContributorsInsights.Recent = append(repoContributorsInsights.Recent, data.AuthorLogin)
		}
	}()
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		response, err := findAlumniRepositoryContributors(ctx, opts.APIClient, repo.Id, opts.Period)
		if err != nil {
			errorChan <- err
			return
		}
		for _, data := range response.Data {
			repoContributorsInsights.Alumni = append(repoContributorsInsights.Alumni, data.AuthorLogin)
		}
	}()
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		response, err := findRepeatRepositoryContributors(ctx, opts.APIClient, repo.Id, opts.Period)
		if err != nil {
			errorChan <- err
			return
		}
		for _, data := range response.Data {
			repoContributorsInsights.Repeat = append(repoContributorsInsights.Repeat, data.AuthorLogin)
		}
	}()
	waitGroup.Wait()
	close(errorChan)
	if len(errorChan) > 0 {
		var allErrors error
		for err = range errorChan {
			allErrors = errors.Join(allErrors, err)
		}
		return nil, allErrors
	}
	return repoContributorsInsights, nil
}

func findNewRepositoryContributors(ctx context.Context, apiClient *client.APIClient, repoID, period int32) (*client.SearchAllPullRequestContributors200Response, error) {
	data, _, err := apiClient.ContributorsServiceAPI.
		NewPullRequestContributors(ctx).
		RepoIds(fmt.Sprintf("%d", repoID)).
		Range_(period).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("error while calling 'ContributorsServiceAPI.NewPullRequestContributors' with repository %d': %w", repoID, err)
	}
	return data, nil
}

func findRecentRepositoryContributors(ctx context.Context, apiClient *client.APIClient, repoID, period int32) (*client.SearchAllPullRequestContributors200Response, error) {
	data, _, err := apiClient.ContributorsServiceAPI.
		FindAllRecentPullRequestContributors(ctx).
		RepoIds(fmt.Sprintf("%d", repoID)).
		Range_(period).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("error while calling 'ContributorsServiceAPI.FindAllRecentPullRequestContributors' with repository %d': %w", repoID, err)
	}
	return data, nil
}

func findAlumniRepositoryContributors(ctx context.Context, apiClient *client.APIClient, repoID, period int32) (*client.SearchAllPullRequestContributors200Response, error) {
	data, _, err := apiClient.ContributorsServiceAPI.
		FindAllChurnPullRequestContributors(ctx).
		RepoIds(fmt.Sprintf("%d", repoID)).
		Range_(period).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("error while calling 'ContributorsServiceAPI.FindAllChurnPullRequestContributors' with repository %d': %w", repoID, err)
	}
	return data, nil
}

func findRepeatRepositoryContributors(ctx context.Context, apiClient *client.APIClient, repoID, period int32) (*client.SearchAllPullRequestContributors200Response, error) {
	data, _, err := apiClient.ContributorsServiceAPI.
		FindAllRepeatPullRequestContributors(ctx).
		RepoIds(fmt.Sprintf("%d", repoID)).
		Range_(period).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("error while calling 'ContributorsServiceAPI.FindAllRepeatPullRequestContributors' with repository %d: %w", repoID, err)
	}
	return data, nil
}
