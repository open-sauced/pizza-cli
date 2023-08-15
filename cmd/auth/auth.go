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

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

const (
	codeChallengeLength = 87
	supabaseID          = "ibcwmlhcimymasokhgvn"
	supabasePublicKey   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYW5vbiIsImlhdCI6MTYyOTkzMDc3OCwiZXhwIjoxOTQ1NTA2Nzc4fQ.zcdbd7kDhk7iNSMo8SjsTaXi0wlLNNQcSZkzZ84NUDg"
	authCallbackAddr    = "localhost:3000"
	sessionFileName     = "session.json"
)

//go:embed success.html
var successHTML string

func NewLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log into the CLI application via GitHub",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}

	return cmd
}

func run() error {
	codeVerifier, codeChallenge, err := pkce(codeChallengeLength)
	if err != nil {
		return fmt.Errorf("PKCE error: %v", err)
	}

	supabaseAuthURL := fmt.Sprintf("https://%s.supabase.co/auth/v1/authorize", supabaseID)
	queryParams := url.Values{
		"provider":              {"github"},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
		"redirect_to":           {"http://" + authCallbackAddr + "/"},
	}

	authenticationURL := supabaseAuthURL + "?" + queryParams.Encode()

	server := &http.Server{Addr: authCallbackAddr}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer shutdown(server)

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Error: auth-code query param not found", http.StatusBadRequest)
			return
		}

		sessionData, err := getSession(code, codeVerifier)
		if err != nil {
			http.Error(w, "Error: Access token exchange failed", http.StatusInternalServerError)
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}

		dirName := path.Join(homeDir, ".pizza")
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			http.Error(w, ".pizza directory couldn't be created", http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(sessionData)
		if err != nil {
			http.Error(w, "Error marshaling session data", http.StatusInternalServerError)
			return
		}

		filePath := path.Join(dirName, sessionFileName)
		if err := os.WriteFile(filePath, jsonData, 0o600); err != nil {
			http.Error(w, "Error writing to file", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = w.Write([]byte(successHTML))
		if err != nil {
			fmt.Println("Error writing response:", err.Error())
		}

		username := sessionData["user"].(map[string]interface{})["user_metadata"].(map[string]interface{})["user_name"].(string)
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

func getSession(authCode, codeVerifier string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s.supabase.co/auth/v1/token?grant_type=pkce", supabaseID)

	payload := map[string]string{
		"auth_code":     authCode,
		"code_verifier": codeVerifier,
	}

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("ApiKey", supabasePublicKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", res.Status)
	}

	var responseData map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	return responseData, nil
}

func pkce(length int) (verifier, challenge string, err error) {
	p := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", "", err
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
			fmt.Println("Graceful shutdown failed", err)
		}
	}()
}
