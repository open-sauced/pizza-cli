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
	"strings"

	"github.com/spf13/cobra"
)

const repoQueryURL string = "https://opensauced.tools"

type Options struct {
	// URL is the git repo URL that will be indexed
	URL string

	branch string
}

const repoQueryLongDesc string = `WARNING: Proof of concept feature.

The repo-query command takes a URL to a git repository and indexes it so
users can ask questions about it.`

func NewRepoQueryCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "repo-query url [flags]",
		Short: "Ask questions about a git repository",
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
			opts.URL = args[0]
			return run(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.branch, "branch", "b", "HEAD", "The branch to index")

	return cmd
}

func getOwnerAndRepo(url string) (owner, repo string) {
	// Remove the "https://github.com/" prefix from the URL
	url = strings.TrimPrefix(url, "https://github.com/")

	// Split the remaining URL path into segments
	segments := strings.Split(url, "/")

	// The first segment is the owner, and the second segment is the repository name
	if len(segments) >= 2 {
		owner = segments[0]
		repo = segments[1]
	} else {
		// URL is not in the expected format
		owner = ""
		repo = ""
	}

	return owner, repo
}

func run(opts *Options) error {
	// get repo name and owner name from URL
	owner, repo := getOwnerAndRepo(opts.URL)

	fmt.Printf("Checking if %s/%s is indexed by us...⏳\n", owner, repo)
	resp, err := http.Get(fmt.Sprintf("%s/collection?owner=%s&name=%s&branch=%s", repoQueryURL, owner, repo, opts.branch))
	if err != nil {
		return err
	}

	// not found or ok or error
	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Repo not found❗")
		fmt.Println("Indexing repo...⏳")
		err := indexRepo(owner, repo, opts.branch)
		if err != nil {
			return err
		}

		for {
			fmt.Printf("\nWant to ask a question about %s/%s?\n", owner, repo)
			fmt.Printf("> ")
			// read input
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input := scanner.Text()
				err := askQuestion(input, owner, repo, opts.branch)
				if err != nil {
					return err
				}
			}
		}
	} else if resp.StatusCode == http.StatusOK {
		fmt.Println("Repo found ✅")

		// this should be an infinite loop
		for {
			fmt.Printf("\nWant to ask a question about %s/%s?\n", owner, repo)
			fmt.Printf("> ")
			// read input
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input := scanner.Text()
				err := askQuestion(input, owner, repo, opts.branch)
				if err != nil {
					return err
				}
			}
		}
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Printf("An error occurred: %v\n", string(body))
	}

	return nil
}

type indexPostRequest struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
}

func indexRepo(owner string, repo string, branch string) error {
	indexPostReq := &indexPostRequest{
		Owner:  owner,
		Name:   repo,
		Branch: branch,
	}

	indexPostJSON, err := json.Marshal(indexPostReq)
	if err != nil {
		return err
	}

	responseBody := bytes.NewBuffer(indexPostJSON)
	resp, err := http.Post(fmt.Sprintf("%s/embed", repoQueryURL), "application/json", responseBody)
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
	}

	// listen for SSEs and send data,event pairs to processIndexChunk
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if line == "\n" {
			continue
		}

		if strings.HasPrefix(line, "event: ") {
			chunk := line
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					return err
				}

				if line == "\n" {
					break
				}

				chunk += line
			}

			processIndexChunk(chunk)
		}
	}
}

func processIndexChunk(chunk string) {
	chunkLines := strings.Split(chunk, "\n")
	eventLine := chunkLines[0]
	event := strings.Split(eventLine, ": ")[1]
	
	switch event {
	case "FETCH_REPO":
		fmt.Println("Fetching Repository from GitHub...")
	case "EMBED_REPO":
		fmt.Println("Embedding Repository...")
	case "SAVE_EMBEDDINGS":
		fmt.Println("Saving the embeddings to our database...")
	case "ERROR":
		fmt.Println("There was an error while indexing this repository. Redirecting to the Home Page.")
	case "DONE":
		fmt.Println("Indexing Complete. You can now ask questions about this repository! 🎉")
	default:
		break
	}
}

type queryPostRequest struct {
	Query      string `json:"query"`
	Repository struct {
		Owner  string `json:"owner"`
		Name   string `json:"name"`
		Branch string `json:"branch"`
	} `json:"repository"`
}

func askQuestion(question string, owner string, repo string, branch string) error {
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

	responseBody := bytes.NewBuffer(queryPostJSON)
	resp, err := http.Post(fmt.Sprintf("%s/query", repoQueryURL), "application/json", responseBody)
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
	}

	//  listen for SSEs and send data,event pairs to processChatChunk
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if line == "\n" {
			continue
		}

		if strings.HasPrefix(line, "event: ") {
			chunk := line
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					return err
				}

				if line == "\n" {
					break
				}

				chunk += line
			}

			processChatChunk(chunk)
		}
	}
}

func processChatChunk(chunk string) {
	chunkLines := strings.Split(chunk, "\n")
	eventLine := chunkLines[0]
	dataLine := chunkLines[1]

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
		fmt.Println("Searching the codebase for your query...🔍")
	case "SEARCH_FILE":
		fmt.Printf("Searching %s for your query...🔍\n", data.(map[string]interface{})["path"])
	case "SEARCH_PATH":
		fmt.Printf("Looking for %s in the codebase...🔍\n", data.(map[string]interface{})["path"])
	case "GENERATE_RESPONSE":
		fmt.Println("Generating a response...🧠")
	case "DONE":
		fmt.Println()
		fmt.Println(data)
	case "ERROR":
		fmt.Println("Something went wrong. Please try again.")
	default:
		break
	}
}
