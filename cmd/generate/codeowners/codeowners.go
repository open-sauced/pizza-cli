package codeowners

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/jpmcb/gopherlogs"
	"github.com/jpmcb/gopherlogs/pkg/colors"
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/api"
	"github.com/open-sauced/pizza-cli/api/auth"
	"github.com/open-sauced/pizza-cli/api/services/workspaces"
	"github.com/open-sauced/pizza-cli/api/services/workspaces/userlists"
	"github.com/open-sauced/pizza-cli/pkg/config"
	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/open-sauced/pizza-cli/pkg/logging"
	"github.com/open-sauced/pizza-cli/pkg/utils"
)

// Options for the codeowners generation command
type Options struct {
	// the path to the git repository on disk to generate a codeowners file for
	path string

	// whether to generate an agnostic "OWENRS" style codeowners file.
	// The default should be to generate a GitHub style "CODEOWNERS" file.
	ownersStyleFile bool

	// the number of days to look back
	previousDays int

	// the session token adding codeowners to a workspace contributor list
	token string

	logger   gopherlogs.Logger
	tty      bool
	loglevel int

	// telemetry for capturing CLI events via PostHog
	telemetry *utils.PosthogCliClient

	config *config.Spec
}

const codeownersLongDesc string = `WARNING: Proof of concept feature.

Generates a CODEOWNERS file for a given git repository. This uses a ~/.sauced.yaml
configuration to attribute emails with given entities.

The generated file specifies up to 3 owners for EVERY file in the git tree based on the
number of lines touched in that specific file over the specified range of time.`

func NewCodeownersCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "codeowners path/to/repo [flags]",
		Short: "Generate a CODEOWNERS file for a GitHub repository using a \"~/.sauced.yaml\" config",
		Long:  codeownersLongDesc,
		Example: `
			# Generate CODEOWNERS file for the current directory
			pizza generate codeowners .

			# Generate CODEOWNERS file for a specific repository
			pizza generate codeowners /path/to/your/repo

			# Generate CODEOWNERS file analyzing the last 180 days
			pizza generate codeowners . --range 180

			# Generate an OWNERS style file instead of CODEOWNERS
			pizza generate codeowners . --owners-style-file

			# Specify a custom location for the .sauced.yaml file
			pizza generate codeowners . --config /path/to/.sauced.yaml
			`,
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
			var err error

			disableTelem, _ := cmd.Flags().GetBool(constants.FlagNameTelemetry)

			opts.telemetry = utils.NewPosthogCliClient(!disableTelem)
			defer opts.telemetry.Done()

			configPath, _ := cmd.Flags().GetString("config")
			opts.config, err = config.LoadConfig(configPath, filepath.Join(opts.path, ".sauced.yaml"))
			if err != nil {
				return err
			}

			opts.ownersStyleFile, _ = cmd.Flags().GetBool("owners-style-file")
			opts.previousDays, _ = cmd.Flags().GetInt("range")
			opts.tty, _ = cmd.Flags().GetBool("tty-disable")

			loglevelS, _ := cmd.Flags().GetString("log-level")

			switch loglevelS {
			case "error":
				opts.loglevel = logging.LogError
			case "warn":
				opts.loglevel = logging.LogWarn
			case "info":
				opts.loglevel = logging.LogInfo
			case "debug":
				opts.loglevel = logging.LogDebug
			}

			return run(opts, cmd)
		},
	}

	cmd.PersistentFlags().IntP("range", "r", 90, "The number of days to analyze commit history (default 90)")
	cmd.PersistentFlags().Bool("owners-style-file", false, "Generate an agnostic OWNERS style file instead of CODEOWNERS.")

	return cmd
}

func run(opts *Options, cmd *cobra.Command) error {
	var err error
	opts.logger, err = gopherlogs.NewLogger(
		gopherlogs.WithLogVerbosity(opts.loglevel),
		gopherlogs.WithTty(!opts.tty),
	)
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Built logger with log level: %d\n", opts.loglevel)

	repo, err := git.PlainOpen(opts.path)
	if err != nil {
		opts.telemetry.CaptureFailedCodeownersGenerate()
		return fmt.Errorf("error opening repo: %w", err)
	}
	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Opened repo at: %s\n", opts.path)

	processOptions := ProcessOptions{
		repo,
		opts.previousDays,
		opts.path,
		opts.logger,
	}
	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Looking back %d days\n", opts.previousDays)

	codeowners, err := processOptions.process()
	if err != nil {
		opts.telemetry.CaptureFailedCodeownersGenerate()
		return fmt.Errorf("error traversing git log: %w", err)
	}

	// Bootstrap codeowners
	var outputPath string
	if opts.ownersStyleFile {
		outputPath = filepath.Join(opts.path, "OWNERS")
	} else {
		outputPath = filepath.Join(opts.path, "CODEOWNERS")
	}

	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Processing codeowners file at: %s\n", outputPath)
	err = generateOutputFile(codeowners, outputPath, opts, cmd)
	if err != nil {
		opts.telemetry.CaptureFailedCodeownersGenerate()
		return fmt.Errorf("error generating github style codeowners file: %w", err)
	}
	opts.logger.V(logging.LogInfo).Style(0, colors.FgGreen).Infof("Finished generating file: %s\n", outputPath)
	opts.telemetry.CaptureCodeownersGenerate()

	// 1. Ask if they want to add users to a list
	var input string
	fmt.Print("Do you want to add these codeowners to an OpenSauced Contributor Insight? (y/n): ")
	_, err = fmt.Scanln(&input)
	if err != nil {
		return fmt.Errorf("could not scan input from terminal: %w", err)
	}

	switch input {
	case "y", "Y", "yes":
		opts.logger.V(logging.LogInfo).Style(0, colors.FgGreen).Infof("Adding codeowners to contributor insight\n")
	case "n", "N", "no":
		return nil
	default:
		return errors.New("invalid answer. Please enter y or n")
	}

	// 2. Check if user is logged in. Log them in if not.
	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Initiating log in flow\n")
	authenticator := auth.NewAuthenticator()
	err = authenticator.CheckSession()
	if err != nil {
		opts.logger.V(logging.LogInfo).Style(0, colors.FgRed).Infof("Log in session invalid: %s\n", err)
		fmt.Print("Do you want to log into OpenSauced? (y/n): ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			return fmt.Errorf("could not scan input from terminal: %w", err)
		}

		switch input {
		case "y", "Y", "yes":
			user, err := authenticator.Login()
			if err != nil {
				opts.telemetry.CaptureFailedCodeownersGenerateAuth()
				opts.logger.V(logging.LogInfo).Style(0, colors.FgRed).Infof("Error logging in\n")
				return fmt.Errorf("could not log in: %w", err)
			}
			opts.telemetry.CaptureCodeownersGenerateAuth(user)
			opts.logger.V(logging.LogInfo).Style(0, colors.FgGreen).Infof("Logged in as: %s\n", user)

		case "n", "N", "no":
			return nil

		default:
			return errors.New("invalid answer. Please enter y or n")
		}
	}

	opts.token, err = authenticator.GetSessionToken()
	if err != nil {
		opts.telemetry.CaptureFailedCodeownersGenerateContributorInsight()
		opts.logger.V(logging.LogInfo).Style(0, colors.FgRed).Infof("Error getting session token\n")
		return fmt.Errorf("could not get session token: %w", err)
	}

	listName := filepath.Base(opts.path)

	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Looking up OpenSauced workspace: Pizza CLI\n")
	workspace, err := findCreatePizzaCliWorkspace(opts)
	if err != nil {
		opts.telemetry.CaptureFailedCodeownersGenerateContributorInsight()
		opts.logger.V(logging.LogInfo).Style(0, colors.FgRed).Infof("Error finding Workspace: Pizza CLI\n")
		return fmt.Errorf("could not find Pizza CLI workspace: %w", err)
	}
	opts.logger.V(logging.LogDebug).Style(0, colors.FgGreen).Infof("Found workspace: Pizza CLI\n")

	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Looking up Contributor Insight for local repository: %s\n", listName)
	userList, err := updateCreateLocalWorkspaceUserList(opts, listName, workspace, codeowners)
	if err != nil {
		opts.telemetry.CaptureFailedCodeownersGenerateContributorInsight()
		opts.logger.V(logging.LogInfo).Style(0, colors.FgRed).Infof("Error finding Workspace Contributor Insight: %s\n", listName)
		return fmt.Errorf("could not find Workspace Contributor Insight: %s - %w", listName, err)
	}
	opts.logger.V(logging.LogDebug).Style(0, colors.FgGreen).Infof("Updated Contributor Insight for local repository: %s\n", listName)
	opts.logger.V(logging.LogInfo).Style(0, colors.FgCyan).Infof("Access list on OpenSauced:\n%s\n", fmt.Sprintf("https://app.opensauced.pizza/workspaces/%s/contributor-insights/%s", workspace.ID, userList.ID))
	opts.telemetry.CaptureCodeownersGenerateContributorInsight()

	return nil
}

// findCreatePizzaCliWorkspace finds or creates a "Pizza CLI" workspace
// for the authenticated user
func findCreatePizzaCliWorkspace(opts *Options) (*workspaces.DbWorkspace, error) {
	nextPage := true
	page := 1
	apiClient := api.NewClient("https://api.opensauced.pizza")

	for nextPage {
		opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Query user workspaces page: %d\n", page)
		workspaceResp, _, err := apiClient.WorkspacesService.GetWorkspaces(opts.token, page, 100)
		if err != nil {
			return nil, err
		}

		for _, workspace := range workspaceResp.Data {
			if workspace.Name == "Pizza CLI" {
				opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Found existing workspace named: Pizza CLI\n")
				return &workspace, nil
			}
		}

		nextPage = workspaceResp.Meta.HasNextPage
		page++
	}

	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Creating new user workspace: Pizza CLI\n")
	newWorkspace, _, err := apiClient.WorkspacesService.CreateWorkspaceForUser(opts.token, "Pizza CLI", "A workspace for the Pizza CLI", []string{})
	if err != nil {
		return nil, err
	}

	return newWorkspace, nil
}

// updateCreateLocalWorkspaceUserList updates or creates a workspace contributor list
// for the authenticated user with the given codeowners
func updateCreateLocalWorkspaceUserList(opts *Options, listName string, workspace *workspaces.DbWorkspace, codeowners FileStats) (*userlists.DbUserList, error) {
	nextPage := true
	page := 1
	apiClient := api.NewClient("https://api.opensauced.pizza")

	var targetUserListID string

	for nextPage {
		opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Query user Workspace Contributor Insight page: %d\n", page)
		userListsResp, _, err := apiClient.WorkspacesService.UserListService.GetUserLists(opts.token, workspace.ID, page, 100)
		if err != nil {
			return nil, err
		}

		nextPage = userListsResp.Meta.HasNextPage
		page++

		for _, userList := range userListsResp.Data {
			if userList.Name == listName {
				opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Found existing Workspace Contributor Insight named: %s\n", listName)
				targetUserListID = userList.ID
				nextPage = false
			}
		}
	}

	if targetUserListID == "" {
		var err error

		opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Creating new user Workspace Contributor List: %s\n", listName)
		createdUserList, _, err := apiClient.WorkspacesService.UserListService.CreateUserListForUser(opts.token, workspace.ID, listName, []string{})
		if err != nil {
			return nil, err
		}

		targetUserListID = createdUserList.UserListID
	}

	targetUserList, _, err := apiClient.WorkspacesService.UserListService.GetUserList(opts.token, workspace.ID, targetUserListID)
	if err != nil {
		return nil, err
	}

	// create a mapping of author logins to empty structs (i.e., a unique set).
	// this de-structures the { filename: author-stats } mapping that originally
	// built the codeowners
	uniqueLogins := make(map[string]struct{})
	for _, codeowner := range codeowners {
		for _, k := range codeowner {
			if k.GitHubAlias != "" {
				uniqueLogins[k.GitHubAlias] = struct{}{}
			}
		}
	}

	logins := []string{}
	for login := range uniqueLogins {
		logins = append(logins, login)
	}

	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Updating Contributor Insight with codeowners with GitHub aliases: %v\n", logins)
	userlist, _, err := apiClient.WorkspacesService.UserListService.PatchUserListForUser(opts.token, workspace.ID, targetUserList.ID, targetUserList.Name, logins)
	return userlist, err
}
