package insights

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"

	bubblesTable "github.com/charmbracelet/bubbles/table"
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/api"
	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/open-sauced/pizza-cli/pkg/utils"
)

type userContributionsOptions struct {
	// APIClient is the http client for making calls to the open-sauced api
	APIClient *api.Client

	// Repos is the array of git repository urls
	Repos []string

	// Users is the list of usernames to filter for
	Users []string

	// usersMap is a fast access set of usernames built from the Users string slice
	usersMap map[string]struct{}

	// FilePath is the path to yaml file containing an array of git repository urls
	FilePath string

	// Period is the number of days, used for query filtering
	Period int32

	// Output is the formatting style for command output
	Output string

	// Sort is the column to be used to sort user contributions (total, commits, pr, none)
	Sort string
}

// NewUserContributionsCommand returns a new user-contributions command
func NewUserContributionsCommand() *cobra.Command {
	opts := &userContributionsOptions{}

	cmd := &cobra.Command{
		Use:   "user-contributions url... [flags]",
		Short: "Gather insights on individual contributors for given repo URLs",
		Long:  "Gather insights on individual contributors given a list of repository URLs",
		Args: func(cmd *cobra.Command, args []string) error {
			fileFlag := cmd.Flags().Lookup(constants.FlagNameFile)
			if !fileFlag.Changed && len(args) == 0 {
				return fmt.Errorf("must specify git repository url argument(s) or provide %s flag", fileFlag.Name)
			}
			opts.Repos = append(opts.Repos, args...)
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			endpointURL, _ := cmd.Flags().GetString(constants.FlagNameEndpoint)
			opts.APIClient = api.NewClient(endpointURL)
			output, _ := cmd.Flags().GetString(constants.FlagNameOutput)
			opts.Output = output
			return opts.run()
		},
	}

	cmd.Flags().StringVarP(&opts.FilePath, constants.FlagNameFile, "f", "", "Path to yaml file containing an array of git repository urls")
	cmd.Flags().Int32VarP(&opts.Period, constants.FlagNameRange, "p", 30, "Number of days, used for query filtering")
	cmd.Flags().StringSliceVarP(&opts.Users, "users", "u", []string{}, "Inclusive comma separated list of GitHub usernames to filter for")
	cmd.Flags().StringVarP(&opts.Sort, "sort", "s", "none", "Sort user contributions by (total, commits, prs)")

	return cmd
}

func (opts *userContributionsOptions) run() error {
	repositories, err := utils.HandleRepositoryValues(opts.Repos, opts.FilePath)
	if err != nil {
		return err
	}

	opts.usersMap = make(map[string]struct{})
	for _, username := range opts.Users {
		// For fast access to list of users to filter out, uses an empty struct
		opts.usersMap[username] = struct{}{}
	}

	var (
		waitGroup    = new(sync.WaitGroup)
		errorChan    = make(chan error, len(repositories))
		insightsChan = make(chan *userContributionsInsightGroup, len(repositories))
		doneChan     = make(chan struct{})
		insights     = make([]*userContributionsInsightGroup, 0, len(repositories))
		allErrors    error
	)

	go func() {
		for url := range repositories {
			repoURL := url
			waitGroup.Add(1)

			go func() {
				defer waitGroup.Done()

				data, err := findAllUserContributionsInsights(opts, repoURL)
				if err != nil {
					errorChan <- err
					return
				}

				if data == nil {
					return
				}

				insightsChan <- data
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

			if opts.Sort != "none" {
				sortUserContributions(insights, opts.Sort)
			}

			for _, insight := range insights {
				output, err := insight.BuildOutput(opts.Output)
				if err != nil {
					return err
				}

				fmt.Println(output)
			}

			return nil
		}
	}
}

type userContributionsInsights struct {
	Login              string `json:"login" yaml:"login"`
	Commits            int    `json:"commits" yaml:"commits"`
	PrsCreated         int    `json:"prs_created" yaml:"prs_created"`
	TotalContributions int    `json:"total_contributions" yaml:"total_contributions"`
}

type userContributionsInsightGroup struct {
	RepoURL  string `json:"repo_url" yaml:"repo_url"`
	Insights []userContributionsInsights
}

func (ucig userContributionsInsightGroup) BuildOutput(format string) (string, error) {
	switch format {
	case constants.OutputTable:
		return ucig.OutputTable()
	case constants.OutputJSON:
		return utils.OutputJSON(ucig)
	case constants.OutputYAML:
		return utils.OutputYAML(ucig)
	case constants.OuputCSV:
		return ucig.OutputCSV()
	default:
		return "", fmt.Errorf("unknown output format %s", format)
	}
}

func (ucig userContributionsInsightGroup) OutputCSV() (string, error) {
	if len(ucig.Insights) == 0 {
		return "", errors.New("repository is either non-existent or has not been indexed yet")
	}

	b := new(bytes.Buffer)
	writer := csv.NewWriter(b)

	// write headers
	err := writer.WriteAll([][]string{
		{
			ucig.RepoURL,
		},
		{
			"User",
			"Total",
			"Commits",
			"PRs Created",
		},
	})
	if err != nil {
		return "", err
	}

	// write records
	for _, uci := range ucig.Insights {
		err := writer.WriteAll([][]string{
			{
				uci.Login,
				strconv.Itoa(uci.Commits + uci.PrsCreated),
				strconv.Itoa(uci.Commits),
				strconv.Itoa(uci.PrsCreated),
			},
		})

		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

func (ucig userContributionsInsightGroup) OutputTable() (string, error) {
	rows := []bubblesTable.Row{}

	for _, uci := range ucig.Insights {
		rows = append(rows, bubblesTable.Row{
			uci.Login,
			strconv.Itoa(uci.TotalContributions),
			strconv.Itoa(uci.Commits),
			strconv.Itoa(uci.PrsCreated),
		})
	}

	columns := []bubblesTable.Column{
		{
			Title: "User",
			Width: utils.GetMaxTableRowWidth(rows),
		},
		{
			Title: "Total",
			Width: 10,
		},
		{
			Title: "Commits",
			Width: 10,
		},
		{
			Title: "PRs Created",
			Width: 15,
		},
	}

	return fmt.Sprintf("%s\n%s\n", ucig.RepoURL, utils.OutputTable(rows, columns)), nil
}

func findAllUserContributionsInsights(opts *userContributionsOptions, repoURL string) (*userContributionsInsightGroup, error) {
	owner, name, err := utils.GetOwnerAndRepoFromURL(repoURL)
	if err != nil {
		return nil, err
	}

	repoUserContributionsInsightGroup := &userContributionsInsightGroup{
		RepoURL: repoURL,
	}

	dataPoints, _, err := opts.APIClient.RepositoryService.FindContributorsByOwnerAndRepo(owner, name, 30)

	if err != nil {
		return nil, fmt.Errorf("error while calling API RepositoryService.FindContributorsByOwnerAndRepo with repository %s/%s': %w", owner, name, err)
	}

	for _, data := range dataPoints.Data {
		_, ok := opts.usersMap[data.Login]
		if len(opts.usersMap) == 0 || ok {
			repoUserContributionsInsightGroup.Insights = append(repoUserContributionsInsightGroup.Insights, userContributionsInsights{
				Login:              data.Login,
				Commits:            data.Commits,
				PrsCreated:         data.PRsCreated,
				TotalContributions: data.Commits + data.PRsCreated,
			})
		}
	}

	return repoUserContributionsInsightGroup, nil
}

func sortUserContributions(ucig []*userContributionsInsightGroup, sortBy string) {
	for _, group := range ucig {
		if group != nil {
			sort.SliceStable(group.Insights, func(i, j int) bool {
				switch sortBy {
				case "total":
					return group.Insights[i].TotalContributions > group.Insights[j].TotalContributions
				case "prs":
					return group.Insights[i].PrsCreated > group.Insights[j].PrsCreated
				case "commits":
					return group.Insights[i].Commits > group.Insights[j].Commits
				}
				return group.Insights[i].Login < group.Insights[j].Login
			})
		}
	}
}
