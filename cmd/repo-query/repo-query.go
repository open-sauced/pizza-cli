package repoquery

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/open-sauced/pizza-cli/pkg/api"
	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/open-sauced/pizza-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type Options struct {
	APIClient *api.Client

	// URL is the git repo URL that will be indexed
	URL string

	// telemetry for capturing CLI events
	telemetry *utils.PosthogCliClient

	branch string
}

const repoQueryLongDesc string = `WARNING: Proof of concept feature.

The repo-query command takes a URL to a git repository and indexes it so
users can ask questions about it.`

func NewRepoQueryCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "repo-query url [flags]",
		Short: "Ask questions about a GitHub repository",
		Long:  repoQueryLongDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("only a single url can be ingested at a time")
			}
			if len(args) == 0 {
				return errors.New("must specify the URL of a git repository")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			endpointURL, _ := cmd.Flags().GetString(constants.FlagNameEndpoint)
			useBeta, _ := cmd.Flags().GetBool(constants.FlagNameBeta)
			if useBeta {
				fmt.Printf("Warning!! Using beta API endpoint not supported for repo-query - using: %s\n", endpointURL)
			}

			// The repo-query is currently not deployed behind "api.opensauced.pizza"
			// So, if the user has not changed the desired "endpoint", use the default
			// tools URL to send SSE to the repo-query engine
			if endpointURL == constants.EndpointProd {
				endpointURL = constants.EndpointTools
			}
			opts.APIClient = api.NewClient(endpointURL)
			opts.URL = args[0]

			disableTelem, _ := cmd.Flags().GetBool(constants.FlagNameTelemetry)
			if !disableTelem {
				opts.telemetry = utils.NewPosthogCliClient()
				defer opts.telemetry.Done()

				opts.telemetry.CaptureRepoQuery(opts.URL)
			}

			return run(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.branch, "branch", "b", "HEAD", "The branch to index")

	return cmd
}

type repoQueryAgent struct {
	client *api.Client
}

func newRepoQueryAgent(apiClient *api.Client) *repoQueryAgent {
	return &repoQueryAgent{
		client: apiClient,
	}
}

func run(opts *Options) error {
	agent := newRepoQueryAgent(opts.APIClient)

	// get repo name and owner name from URL
	owner, repo, err := utils.GetOwnerAndRepoFromURL(opts.URL)
	if err != nil {
		return err
	}

	fmt.Printf("Checking if %s/%s is indexed by us ... ⏳\n", owner, repo)
	resp, err := agent.client.HTTPClient.Get(fmt.Sprintf("%s/collection?owner=%s&name=%s&branch=%s", agent.client.Endpoint, owner, repo, opts.branch))
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		// repo is not indexed
		fmt.Println("Repo not found ❗")
		fmt.Println("Indexing repo ... ⏳")
		err := agent.indexRepo(owner, repo, opts.branch)
		if err != nil {
			return err
		}

		err = agent.startQnALoop(owner, repo, opts.branch)
		if err != nil {
			return err
		}
	case http.StatusOK:
		// repo is indexed
		fmt.Println("Repo found ✅")

		err = agent.startQnALoop(owner, repo, opts.branch)
		if err != nil {
			return err
		}
	default:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Printf("An error occurred: %v\n", string(body))
		return errors.New("error while checking if repo is indexed")
	}

	return nil
}

func (rq *repoQueryAgent) startQnALoop(owner string, repo string, branch string) error {
	for {
		// if ctrl+c is pressed, exit
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		go func() {
			<-c
			fmt.Println("\n🍕 Exiting ...")
			os.Exit(0)
		}()

		fmt.Printf("\nWant to ask a question about %s/%s?\n", owner, repo)
		fmt.Printf("> ")
		// read input
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := scanner.Text()
			if input == "exit" {
				fmt.Println("🍕 Exiting ...")
				os.Exit(0)
			}
			err := rq.askQuestion(input, owner, repo, branch)
			if err != nil {
				return err
			}
		}
	}
}

type indexPostRequest struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
}

type queryPostRequest struct {
	Query      string `json:"query"`
	Repository struct {
		Owner  string `json:"owner"`
		Name   string `json:"name"`
		Branch string `json:"branch"`
	} `json:"repository"`
}

const (
	IndexChunk = iota
	ChatChunk
)

func (rq *repoQueryAgent) indexRepo(owner string, repo string, branch string) error {
	indexPostReq := &indexPostRequest{
		Owner:  owner,
		Name:   repo,
		Branch: branch,
	}

	indexPostJSON, err := json.Marshal(indexPostReq)
	if err != nil {
		return err
	}

	resp, err := rq.client.HTTPClient.Post(fmt.Sprintf("%s/embed", rq.client.Endpoint), "application/json", bytes.NewBuffer(indexPostJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error while indexing repository: %s", err.Error())
		}

		return fmt.Errorf("error while indexing repository: %s", string(body))
	}

	reader := bufio.NewReader(resp.Body)
	err = rq.listenForSSEs(reader, IndexChunk)
	if err != nil {
		return err
	}

	return nil
}

func (rq *repoQueryAgent) askQuestion(question string, owner string, repo string, branch string) error {
	queryPostReq := &queryPostRequest{
		Query: question,
		Repository: struct {
			Owner  string `json:"owner"`
			Name   string `json:"name"`
			Branch string `json:"branch"`
		}{
			Owner:  owner,
			Name:   repo,
			Branch: branch,
		},
	}

	queryPostJSON, err := json.Marshal(queryPostReq)
	if err != nil {
		return err
	}

	resp, err := rq.client.HTTPClient.Post(fmt.Sprintf("%s/query", rq.client.Endpoint), "application/json", bytes.NewBuffer(queryPostJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Printf("An error occurred: %v\n", string(body))
		return errors.New("error while asking question")
	}

	//  listen for SSEs and send data,event pairs to processChatChunk
	reader := bufio.NewReader(resp.Body)
	err = rq.listenForSSEs(reader, ChatChunk)
	if err != nil {
		return err
	}

	return nil
}

func (rq *repoQueryAgent) listenForSSEs(reader *bufio.Reader, chunkType int) error {
	// listen for SSEs and send data, event pairs to processChunk
	// we send 2 lines at a time to processChunk so it can process the event and data together.
	// the server sends empty events sometimes, so we ignore those.

	for {
		line, err := reader.ReadString('\n')
		// if we have reached the end of the stream, return
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		// ignore empty events
		if line == "\n" {
			continue
		}

		if strings.HasPrefix(line, "event: ") {
			chunk := line
			for {
				// we read the string again after getting the event, so we can get the data
				line, err := reader.ReadString('\n')
				if err != nil {
					return err
				}

				// the data data can be empty too.
				if line == "\n" {
					break
				}

				chunk += line
			}

			switch chunkType {
			case IndexChunk:
				err = rq.processIndexChunk(chunk)
			case ChatChunk:
				err = rq.processChatChunk(chunk)
			default:
				break
			}

			if err != nil {
				return err
			}
		}
	}
}

func (rq *repoQueryAgent) processIndexChunk(chunk string) error {
	// we only care about the first line of the chunk, which is the event, when indexing.
	// the data is irrelevant for now, but we still got it so we can process it later if we need to.
	// Also, for grouping the events and data together.

	chunkLines := strings.Split(chunk, "\n")
	eventLine := chunkLines[0]
	event := strings.Split(eventLine, ": ")[1]

	switch event {
	case "FETCH_REPO":
		fmt.Println("Fetching Repository from GitHub ...")
	case "EMBED_REPO":
		fmt.Println("Embedding Repository ...")
	case "SAVE_EMBEDDINGS":
		fmt.Println("Saving the embeddings to our database ...")
	case "ERROR":
		fmt.Println("There was an error while indexing this repository.")
		return errors.New("error while indexing repository")
	case "DONE":
		fmt.Println("Indexing Complete. You can now ask questions about this repository! 🎉")
	default:
		break
	}

	return nil
}

func (rq *repoQueryAgent) processChatChunk(chunk string) error {
	// The event is the first line of the chunk, and the data is the second line.
	chunkLines := strings.Split(chunk, "\n")
	eventLine := chunkLines[0]
	dataLine := chunkLines[1]

	// the event is the part after the colon
	// eg. - event: SEARCH_PATH
	//       data: {"path": "src/index.js"}
	// eg. (with a string as data) - event: DONE
	//                               data: "Here's the answer to your question"
	event := strings.Split(eventLine, ": ")[1]
	var data interface{} // the data can be a string or a JSON object

	// try to parse the data as JSON
	err := json.Unmarshal([]byte(strings.Split(dataLine, ": ")[1]), &data)
	if err != nil {
		// remove quotes from string
		data = strings.Split(dataLine, "data: ")[1][1 : len(strings.Split(dataLine, "data: ")[1])-2]
	}

	switch event {
	case "SEARCH_CODEBASE":
		fmt.Println("Searching the codebase for your query ... 🔍")
	case "SEARCH_FILE":
		fmt.Printf("Searching %s for your query ... 🔍\n", data.(map[string]interface{})["path"])
	case "SEARCH_PATH":
		fmt.Printf("Looking for %s in the codebase ... 🔍\n", data.(map[string]interface{})["path"])
	case "GENERATE_RESPONSE":
		fmt.Println("Generating a response ... 🧠")
	case "DONE":
		fmt.Println()
		// the server sends the data as a string, with \n, "", and \t escaped.
		replacer := strings.NewReplacer("\\n", "\n", "\\\"", "\"", "\\t", "\t")
		data := replacer.Replace(data.(string))
		fmt.Println(data)
	case "ERROR":
		fmt.Println("Something went wrong. Please try again.")
		return errors.New("error while asking question")
	default:
		break
	}

	return nil
}
