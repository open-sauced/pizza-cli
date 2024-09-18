package insights

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	bubblesTable "github.com/charmbracelet/bubbles/table"
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/v2/api"
	"github.com/open-sauced/pizza-cli/v2/api/services/contributors"
	apiUtils "github.com/open-sauced/pizza-cli/v2/api/utils"
	"github.com/open-sauced/pizza-cli/v2/pkg/constants"
	"github.com/open-sauced/pizza-cli/v2/pkg/utils"
)

type contributorsOptions struct {
	// APIClient is the http client for making calls to the open-sauced api
	APIClient *api.Client

	// Repos is the array of git repository urls
	Repos []string

	// FilePath is the path to yaml file containing an array of git repository urls
	FilePath string

	// RangeVal is the number of days, used for query filtering
	RangeVal int

	// Output is the formatting style for command output
	Output string

	telemetry *utils.PosthogCliClient
}

// NewContributorsCommand returns a new cobra command for 'pizza insights contributors'
func NewContributorsCommand() *cobra.Command {
	opts := &contributorsOptions{}
	cmd := &cobra.Command{
		Use:   "contributors url... [flags]",
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
	cmd.Flags().IntVarP(&opts.RangeVal, constants.FlagNameRange, "r", 30, "Number of days to look-back (7,30,90)")
	return cmd
}

func (opts *contributorsOptions) run() error {
	if !apiUtils.IsValidRange(opts.RangeVal) {
		return fmt.Errorf("invalid period: %d, accepts (7,30,90)", opts.RangeVal)
	}

	repositories, err := utils.HandleRepositoryValues(opts.Repos, opts.FilePath)
	if err != nil {
		return err
	}
	var (
		waitGroup    = new(sync.WaitGroup)
		errorChan    = make(chan error, len(repositories))
		insightsChan = make(chan contributorsInsights, len(repositories))
		doneChan     = make(chan struct{})
		insights     = make(contributorsInsightsSlice, 0, len(repositories))
		allErrors    error
	)
	go func() {
		for url := range repositories {
			repoURL := url
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				allData, err := findAllContributorsInsights(opts, repoURL)
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

type contributorsInsights struct {
	RepoURL string   `json:"repo_url" yaml:"repo_url"`
	RepoID  int      `json:"-" yaml:"-"`
	New     []string `json:"new" yaml:"new"`
	Recent  []string `json:"recent" yaml:"recent"`
	Alumni  []string `json:"alumni" yaml:"alumni"`
	Repeat  []string `json:"repeat" yaml:"repeat"`
}

type contributorsInsightsSlice []contributorsInsights

func (cis contributorsInsightsSlice) BuildOutput(format string) (string, error) {
	switch format {
	case constants.OutputTable:
		return cis.OutputTable()
	case constants.OutputJSON:
		return utils.OutputJSON(cis)
	case constants.OutputYAML:
		return utils.OutputYAML(cis)
	case constants.OuputCSV:
		return cis.OutputCSV()
	default:
		return "", fmt.Errorf("unknown output format %s", format)
	}
}

func (cis contributorsInsightsSlice) OutputCSV() (string, error) {
	if len(cis) == 0 {
		return "", errors.New("repository is either non-existent or has not been indexed yet")
	}
	b := new(bytes.Buffer)
	writer := csv.NewWriter(b)

	// write headers
	err := writer.WriteAll([][]string{{"Repository URL", "New Contributors", "Recent Contributors", "Alumni Contributors", "Repeat Contributors"}})
	if err != nil {
		return "", err
	}

	// write records
	for _, ci := range cis {
		err := writer.WriteAll([][]string{{ci.RepoURL, strconv.Itoa(len(ci.New)), strconv.Itoa(len(ci.Recent)),
			strconv.Itoa(len(ci.Alumni)), strconv.Itoa(len(ci.Repeat))}})

		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

func (cis contributorsInsightsSlice) OutputTable() (string, error) {
	tables := make([]string, 0, len(cis))
	for i := range cis {
		rows := []bubblesTable.Row{
			{
				"New contributors",
				strconv.Itoa(len(cis[i].New)),
			},
			{
				"Recent contributors",
				strconv.Itoa(len(cis[i].Recent)),
			},
			{
				"Alumni contributors",
				strconv.Itoa(len(cis[i].Alumni)),
			},
			{
				"Repeat contributors",
				strconv.Itoa(len(cis[i].Repeat)),
			},
		}
		columns := []bubblesTable.Column{
			{
				Title: "Repository URL",
				Width: utils.GetMaxTableRowWidth(rows),
			},
			{
				Title: cis[i].RepoURL,
				Width: len(cis[i].RepoURL),
			},
		}
		tables = append(tables, utils.OutputTable(rows, columns))
	}
	separator := fmt.Sprintf("\n%s\n", strings.Repeat("―", 3))
	return strings.Join(tables, separator), nil
}

func findAllContributorsInsights(opts *contributorsOptions, repoURL string) (*contributorsInsights, error) {
	var (
		waitGroup = new(sync.WaitGroup)
		errorChan = make(chan error, 4)
	)

	repo, err := findRepositoryByOwnerAndRepoName(opts.APIClient, repoURL)
	if err != nil {
		return nil, fmt.Errorf("could not get contributors insights for repository %s: %w", repoURL, err)
	}
	if repo == nil {
		return nil, nil
	}

	repoContributorsInsights := &contributorsInsights{
		RepoID:  repo.ID,
		RepoURL: repo.SvnURL,
	}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		response, err := findNewRepositoryContributors(opts.APIClient, repo.FullName, opts.RangeVal)
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
		response, err := findRecentRepositoryContributors(opts.APIClient, repo.FullName, opts.RangeVal)
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
		response, err := findAlumniRepositoryContributors(opts.APIClient, repo.FullName, opts.RangeVal)
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
		response, err := findRepeatRepositoryContributors(opts.APIClient, repo.FullName, opts.RangeVal)
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

func findNewRepositoryContributors(apiClient *api.Client, repo string, period int) (*contributors.ContribResponse, error) {
	response, _, err := apiClient.ContributorService.NewPullRequestContributors([]string{repo}, period)
	if err != nil {
		return nil, fmt.Errorf("error while calling 'ContributorsService.NewPullRequestContributors' with repository %s': %w", repo, err)
	}

	return response, nil
}

func findRecentRepositoryContributors(apiClient *api.Client, repo string, period int) (*contributors.ContribResponse, error) {
	response, _, err := apiClient.ContributorService.RecentPullRequestContributors([]string{repo}, period)
	if err != nil {
		return nil, fmt.Errorf("error while calling 'ContributorsService.RecentPullRequestContributors' with repository %s': %w", repo, err)
	}

	return response, nil
}

func findAlumniRepositoryContributors(apiClient *api.Client, repo string, period int) (*contributors.ContribResponse, error) {
	response, _, err := apiClient.ContributorService.AlumniPullRequestContributors([]string{repo}, period)
	if err != nil {
		return nil, fmt.Errorf("error while calling 'ContributorsService.AlumniPullRequestContributors' with repository %s': %w", repo, err)
	}

	return response, nil
}

func findRepeatRepositoryContributors(apiClient *api.Client, repo string, period int) (*contributors.ContribResponse, error) {
	response, _, err := apiClient.ContributorService.RepeatPullRequestContributors([]string{repo}, period)
	if err != nil {
		return nil, fmt.Errorf("error while calling 'ContributorsService.RepeatPullRequestContributors' with repository %s': %w", repo, err)
	}

	return response, nil
}
