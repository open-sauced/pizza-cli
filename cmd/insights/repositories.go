package insights

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	bubblesTable "github.com/charmbracelet/bubbles/table"
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/v2/api"
	"github.com/open-sauced/pizza-cli/v2/api/services/contributors"
	"github.com/open-sauced/pizza-cli/v2/api/services/histogram"
	"github.com/open-sauced/pizza-cli/v2/pkg/constants"
	"github.com/open-sauced/pizza-cli/v2/pkg/utils"
)

type repositoriesOptions struct {
	// APIClient is the http client for making calls to the open-sauced api
	APIClient *api.Client

	// Repos is the array of git repository urls
	Repos []string

	// FilePath is the path to yaml file containing an array of git repository urls
	FilePath string

	// RangeVal is the number of days, used for query filtering
	// Constrained to either 30 or 60
	RangeVal int

	// Output is the formatting style for command output
	Output string

	telemetry *utils.PosthogCliClient
}

// NewRepositoriesCommand returns a new cobra command for 'pizza insights repositories'
func NewRepositoriesCommand() *cobra.Command {
	opts := &repositoriesOptions{}
	cmd := &cobra.Command{
		Use:     "repositories url... [flags]",
		Aliases: []string{"repos"},
		Short:   "Gather insights about indexed git repositories",
		Long:    "Gather insights about indexed git repositories. This command will show info about contributors, pull requests, etc.",
		Args: func(cmd *cobra.Command, args []string) error {
			fileFlag := cmd.Flags().Lookup(constants.FlagNameFile)
			if !fileFlag.Changed && len(args) == 0 {
				return fmt.Errorf("must specify git repository url argument(s) or provide %s flag", fileFlag.Name)
			}
			opts.Repos = append(opts.Repos, args...)
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			disableTelem, _ := cmd.Flags().GetBool(constants.FlagNameTelemetry)

			opts.telemetry = utils.NewPosthogCliClient(!disableTelem)

			endpointURL, _ := cmd.Flags().GetString(constants.FlagNameEndpoint)
			opts.APIClient = api.NewClient(endpointURL)
			output, _ := cmd.Flags().GetString(constants.FlagNameOutput)
			opts.Output = output

			err := opts.run()

			if err != nil {
				_ = opts.telemetry.CaptureInsights()
			} else {
				_ = opts.telemetry.CaptureFailedInsights()
			}

			_ = opts.telemetry.Done()

			return err
		},
	}
	cmd.Flags().StringVarP(&opts.FilePath, constants.FlagNameFile, "f", "", "Path to yaml file containing an array of git repository urls")
	cmd.Flags().IntVarP(&opts.RangeVal, constants.FlagNameRange, "p", 30, "Number of days to look-back")
	return cmd
}

func (opts *repositoriesOptions) run() error {
	repositories, err := utils.HandleRepositoryValues(opts.Repos, opts.FilePath)
	if err != nil {
		return err
	}
	var (
		waitGroup    = new(sync.WaitGroup)
		errorChan    = make(chan error, len(repositories))
		insightsChan = make(chan repositoryInsights, len(repositories))
		doneChan     = make(chan struct{})
		insights     = make(repositoryInsightsSlice, 0, len(repositories))
		allErrors    error
	)
	go func() {
		for url := range repositories {
			repoURL := url
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				allData, err := findAllRepositoryInsights(opts, repoURL)
				if err != nil {
					errorChan <- err
					return
				}
				if allData == nil {
					return
				}
				insightsChan <- *allData
			}()
		}
		waitGroup.Wait()
		close(doneChan)
	}()
	for {
		select {
		case err = <-errorChan:
			allErrors = errors.Join(allErrors, err)
		case data := <-insightsChan:
			insights = append(insights, data)
		case <-doneChan:
			if allErrors != nil {
				return allErrors
			}
			output, err := insights.BuildOutput(opts.Output)
			if err != nil {
				return err
			}
			fmt.Println(output)
			return nil
		}
	}
}

type repositoryInsights struct {
	RepoURL              string   `json:"repo_url" yaml:"repo_url"`
	RepoID               int      `json:"-" yaml:"-"`
	AllPullRequests      int      `json:"all_pull_requests" yaml:"all_pull_requests"`
	AcceptedPullRequests int      `json:"accepted_pull_requests" yaml:"accepted_pull_requests"`
	SpamPullRequests     int      `json:"spam_pull_requests" yaml:"spam_pull_requests"`
	Contributors         []string `json:"contributors" yaml:"contributors"`
}

type repositoryInsightsSlice []repositoryInsights

func (ris repositoryInsightsSlice) BuildOutput(format string) (string, error) {
	switch format {
	case constants.OutputTable:
		return ris.OutputTable()
	case constants.OutputJSON:
		return utils.OutputJSON(ris)
	case constants.OutputYAML:
		return utils.OutputYAML(ris)
	default:
		return "", fmt.Errorf("unknown output format %s", format)
	}
}

func (ris repositoryInsightsSlice) OutputTable() (string, error) {
	tables := make([]string, 0, len(ris))
	for i := range ris {
		rows := []bubblesTable.Row{
			{
				"All pull requests",
				strconv.Itoa(ris[i].AllPullRequests),
			},
			{
				"Accepted pull requests",
				strconv.Itoa(ris[i].AcceptedPullRequests),
			},
			{
				"Spam pull requests",
				strconv.Itoa(ris[i].SpamPullRequests),
			},
			{
				"Contributors",
				strconv.Itoa(len(ris[i].Contributors)),
			},
		}
		columns := []bubblesTable.Column{
			{
				Title: "Repository URL",
				Width: utils.GetMaxTableRowWidth(rows),
			},
			{
				Title: ris[i].RepoURL,
				Width: len(ris[i].RepoURL),
			},
		}
		tables = append(tables, utils.OutputTable(rows, columns))
	}
	separator := fmt.Sprintf("\n%s\n", strings.Repeat("―", 3))
	return strings.Join(tables, separator), nil
}

func findAllRepositoryInsights(opts *repositoriesOptions, repoURL string) (*repositoryInsights, error) {
	repo, err := findRepositoryByOwnerAndRepoName(opts.APIClient, repoURL)
	if err != nil {
		return nil, fmt.Errorf("could not get repository insights for repository %s: %w", repoURL, err)
	}
	if repo == nil {
		return nil, nil
	}
	repoInsights := &repositoryInsights{
		RepoID:  repo.ID,
		RepoURL: repo.SvnURL,
	}

	var (
		waitGroup = new(sync.WaitGroup)
		errorChan = make(chan error, 4)
	)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		response, err := getPullRequestInsights(opts.APIClient, repo.FullName, opts.RangeVal)
		if err != nil {
			errorChan <- err
			return
		}

		for _, bucket := range response {
			repoInsights.AllPullRequests += bucket.PrCount
			repoInsights.AcceptedPullRequests += bucket.AcceptedPrs
			repoInsights.SpamPullRequests += bucket.SpamPrs
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		response, err := searchAllPullRequestContributors(opts.APIClient, []string{repo.FullName}, opts.RangeVal)
		if err != nil {
			errorChan <- err
			return
		}
		var contributors []string
		for _, contributor := range response.Data {
			contributors = append(contributors, contributor.AuthorLogin)
		}
		repoInsights.Contributors = contributors
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

	return repoInsights, nil
}

func getPullRequestInsights(apiClient *api.Client, repo string, rangeVal int) ([]histogram.PrHistogramData, error) {
	data, _, err := apiClient.HistogramService.PrsHistogram(repo, rangeVal)
	if err != nil {
		return nil, fmt.Errorf("error while calling 'PullRequestsServiceAPI.GetPullRequestInsights' with repository %s': %w", repo, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("could not find pull request insights for repository %s with interval %d", repo, rangeVal)
	}

	return data, nil
}

func searchAllPullRequestContributors(apiClient *api.Client, repos []string, rangeVal int) (*contributors.ContribResponse, error) {
	data, _, err := apiClient.ContributorService.SearchPullRequestContributors(repos, rangeVal)
	if err != nil {
		return nil, fmt.Errorf("error while calling 'ContributorService.SearchPullRequestContributors' with repository %v': %w", repos, err)
	}

	return data, nil
}
