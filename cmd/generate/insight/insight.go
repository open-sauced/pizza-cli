package insight

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jpmcb/gopherlogs"
	"github.com/jpmcb/gopherlogs/pkg/colors"
	"github.com/spf13/cobra"

	"github.com/open-sauced/pizza-cli/api"
	"github.com/open-sauced/pizza-cli/api/auth"
	"github.com/open-sauced/pizza-cli/api/services/workspaces"
	"github.com/open-sauced/pizza-cli/api/services/workspaces/userlists"
	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/open-sauced/pizza-cli/pkg/logging"
	"github.com/open-sauced/pizza-cli/pkg/utils"
)

// Options for the codeowners generation command
type Options struct {
	// the path to the git repository on disk to generate a codeowners file for
	path string

	logger   gopherlogs.Logger
	tty      bool
	loglevel int

	token string

	// telemetry for capturing CLI events via PostHog
	telemetry *utils.PosthogCliClient
}

const insightLongDesc string = `Generate an OpenSauced Contributor Insight based on GitHub logins in a CODEOWNERS file
to get metrics and insights on those users.

The provided path must be a local git repo with a valid CODEOWNERS file and GitHub "@login"
for each codeowner.

After logging in, the generated Contributor Insight on OpenSauced will have insights on
active contributors, contributon velocity, and more.`

const insightExamples string = `  # Use CODEOWNERS file in explicit directory
  $ pizza generate insight /path/to/repo

  # Use CODEOWNERS file in local directory
  $ pizza generate insight .`

func NewGenerateInsightCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "insight path/to/repo/with/CODEOWNERS/file [flags]",
		Short:   "Generate an OpenSauced Contributor Insight based on GitHub logins in a CODEOWNERS file",
		Long:    insightLongDesc,
		Example: insightExamples,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you must provide exactly one argument: the path to a repository with a codeowners file")
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

			err = run(opts, cmd)

			_ = opts.telemetry.Done()

			return err
		},
	}

	return cmd
}

func run(opts *Options, _ *cobra.Command) error {
	var err error
	opts.logger, err = gopherlogs.NewLogger(
		gopherlogs.WithLogVerbosity(opts.loglevel),
		gopherlogs.WithTty(!opts.tty),
	)
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Built logger with log level: %d\n", opts.loglevel)

	codeowners, err := getUniqueCodeowners(filepath.Join(opts.path, "CODEOWNERS"))
	if err != nil {
		return fmt.Errorf("could not get codeowners from file: %s - %w", opts.path, err)
	}

	// 1. Ask if they want to add users to a list
	var input string
	fmt.Print("Do you want to add these codeowners to an OpenSauced Contributor Insight? (y/n): ")
	_, err = fmt.Scanln(&input)
	if err != nil {
		return fmt.Errorf("could not scan input from terminal: %w", err)
	}

	switch input {
	case "y", "Y", "yes":
		opts.logger.V(logging.LogInfo).Style(0, colors.FgGreen).Infof("Adding codeowners to Contributor Insight\n")
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
				_ = opts.telemetry.CaptureFailedCodeownersGenerateAuth()
				opts.logger.V(logging.LogInfo).Style(0, colors.FgRed).Infof("Error logging in\n")
				return fmt.Errorf("could not log in: %w", err)
			}
			_ = opts.telemetry.CaptureCodeownersGenerateAuth(user)
			opts.logger.V(logging.LogInfo).Style(0, colors.FgGreen).Infof("Logged in as: %s\n", user)

		case "n", "N", "no":
			return nil

		default:
			return errors.New("invalid answer. Please enter y or n")
		}
	}

	opts.token, err = authenticator.GetSessionToken()
	if err != nil {
		_ = opts.telemetry.CaptureFailedCodeownersGenerateContributorInsight()
		opts.logger.V(logging.LogInfo).Style(0, colors.FgRed).Infof("Error getting session token\n")
		return fmt.Errorf("could not get session token: %w", err)
	}

	listName := filepath.Base(opts.path)

	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Looking up OpenSauced workspace: Pizza CLI\n")
	workspace, err := findCreatePizzaCliWorkspace(opts)
	if err != nil {
		_ = opts.telemetry.CaptureFailedCodeownersGenerateContributorInsight()
		opts.logger.V(logging.LogInfo).Style(0, colors.FgRed).Infof("Error finding Workspace: Pizza CLI\n")
		return fmt.Errorf("could not find Pizza CLI workspace: %w", err)
	}
	opts.logger.V(logging.LogDebug).Style(0, colors.FgGreen).Infof("Found workspace: Pizza CLI\n")

	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Looking up Contributor Insight for local repository: %s\n", listName)
	userList, err := updateCreateLocalWorkspaceUserList(opts, listName, workspace, codeowners)
	if err != nil {
		_ = opts.telemetry.CaptureFailedCodeownersGenerateContributorInsight()
		opts.logger.V(logging.LogInfo).Style(0, colors.FgRed).Infof("Error finding Workspace Contributor Insight: %s\n", listName)
		return fmt.Errorf("could not find Workspace Contributor Insight: %s - %w", listName, err)
	}
	opts.logger.V(logging.LogDebug).Style(0, colors.FgGreen).Infof("Updated Contributor Insight for local repository: %s\n", listName)
	opts.logger.V(logging.LogInfo).Style(0, colors.FgCyan).Infof("\nAccess Contributor Insight on OpenSauced:\n%s\n", fmt.Sprintf("https://app.opensauced.pizza/workspaces/%s/contributor-insights/%s", workspace.ID, userList.ID))
	_ = opts.telemetry.CaptureCodeownersGenerateContributorInsight()

	return nil
}

func getUniqueCodeowners(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a map to store unique uniqueLogins
	uniqueLogins := make(map[string]struct{})

	// Create a regular expression to match GitHub usernames
	re := regexp.MustCompile(`@(\w+)`)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Find all matches in the line
		matches := re.FindAllStringSubmatch(line, -1)

		// Add each match to the map
		for _, match := range matches {
			uniqueLogins[match[1]] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	logins := []string{}
	for login := range uniqueLogins {
		logins = append(logins, login)
	}

	fmt.Printf("%v\n", logins)
	return logins, nil
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
func updateCreateLocalWorkspaceUserList(opts *Options, listName string, workspace *workspaces.DbWorkspace, logins []string) (*userlists.DbUserList, error) {
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

	opts.logger.V(logging.LogDebug).Style(0, colors.FgBlue).Infof("Updating Contributor Insight with codeowners with GitHub aliases: %v\n", logins)
	userlist, _, err := apiClient.WorkspacesService.UserListService.PatchUserListForUser(opts.token, workspace.ID, targetUserList.ID, targetUserList.Name, logins)
	return userlist, err
}
