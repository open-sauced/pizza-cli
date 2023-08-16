package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/cli/browser"
	"github.com/open-sauced/pizza-cli/pkg/constants"
	"github.com/spf13/cobra"
)

//go:embed success.html
var successHTML string

const loginLongDesc string = `Log into OpenSauced.

This command initiates the GitHub auth flow to log you into the OpenSauced application by launching your browser`

func NewLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log into the CLI application via GitHub",
		Long:  loginLongDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}

	return cmd
}

func run() error {
	codeVerifier, codeChallenge, err := pkce(CodeChallengeLength)
	if err != nil {
		return fmt.Errorf("PKCE error: %v", err.Error())
	}

	supabaseAuthURL := fmt.Sprintf("https://%s.supabase.co/auth/v1/authorize", SupabaseID)
	queryParams := url.Values{
		"provider":              {"github"},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
		"redirect_to":           {"http://" + AuthCallbackAddr + "/"},
	}

	authenticationURL := supabaseAuthURL + "?" + queryParams.Encode()

	server := &http.Server{Addr: AuthCallbackAddr}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer shutdown(server)

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "'code' query param not found", http.StatusBadRequest)
			return
		}

		sessionData, err := getSession(code, codeVerifier)
		if err != nil {
			http.Error(w, "Access token exchange failed", http.StatusInternalServerError)
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			http.Error(w, "Couldn't get the Home directory", http.StatusInternalServerError)
			return
		}

		dirName := path.Join(homeDir, ".pizza")
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			http.Error(w, ".pizza directory couldn't be created", http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(sessionData)
		if err != nil {
			http.Error(w, "Marshaling session data failed", http.StatusInternalServerError)
			return
		}

		filePath := path.Join(dirName, constants.SessionFileName)
		if err := os.WriteFile(filePath, jsonData, 0o600); err != nil {
			http.Error(w, "Error writing to file", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = w.Write([]byte(successHTML))
		if err != nil {
			fmt.Println("Error writing response:", err.Error())
		}

		username := sessionData.User.UserMetadata["user_name"]
		fmt.Println("🎉 Login successful 🎉")
		fmt.Println("Welcome aboard", username, "🍕")
	})

	err = browser.OpenURL(authenticationURL)
	if err != nil {
		fmt.Println("Failed to open the browser 🤦‍♂️")
		fmt.Println("Navigate to the following URL to begin authentication:")
		fmt.Println(authenticationURL)
	}

	errCh := make(chan error)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	case <-time.After(60 * time.Second):
		shutdown(server)
		return errors.New("authentication timeout")
	case <-interruptCh:
		fmt.Println("\nAuthentication interrupted❗️")
		shutdown(server)
		os.Exit(0)
	}
	return nil
}

func getSession(authCode, codeVerifier string) (*accessTokenResponse, error) {
	url := fmt.Sprintf("https://%s.supabase.co/auth/v1/token?grant_type=pkce", SupabaseID)

	payload := map[string]string{
		"auth_code":     authCode,
		"code_verifier": codeVerifier,
	}

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("ApiKey", SupabasePublicKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't make a request with the default client: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", res.Status)
	}

	var responseData accessTokenResponse
	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %s", err.Error())
	}

	return &responseData, nil
}

func pkce(length int) (verifier, challenge string, err error) {
	p := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", "", fmt.Errorf("failed to read random bytes: %s", err.Error())
	}
	verifier = base64.RawURLEncoding.EncodeToString(p)
	b := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(b[:])
	return verifier, challenge, nil
}

func shutdown(server *http.Server) {
	go func() {
		err := server.Shutdown(context.Background())
		if err != nil {
			panic(err.Error())
		}
	}()
}
