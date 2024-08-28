package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	_ "embed" // global import of embed to enable the use of the "go embed" directive
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/cli/browser"
	"github.com/go-chi/chi/v5"

	"github.com/open-sauced/pizza-cli/pkg/config"
)

// The success HTML file is embedded directly as a var.
// This allows us to include the HTML within the binary we ship to end users without
// having to ship HTML that gets served by the callback server.
//
//go:embed success.html
var successHTML string

const (
	authCallbackAddr    = "localhost:3000"
	codeChallengeLength = 87
	sessionFileName     = "session.json"

	prodSupabaseURL       = "https://fcqqkxwlntnrtjfbcioz.supabase.co"
	prodSupabasePublicKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImZjcXFreHdsbnRucnRqZmJjaW96Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2OTg0MTkyNzQsImV4cCI6MjAxMzk5NTI3NH0.ymWWYdnJC2gsnrJx4lZX2cfSOp-1xVuWFGt1Wr6zwtg"

	// TODO (jpmcb) - in the future, we'll want to encorporate the ability to
	// authenticate to our beta auth service as well
)

// Authenticator is a utility for performing authentication of the given user.
// It carries necessary metadata for spinning up a local server alongside
// channels for errors and when authentication is "done".
type Authenticator struct {
	username     string
	codeVerifier string

	errChan  chan error
	doneChan chan struct{}
}

// NewAuthenticator returns a new Authenticator for the caller with instantiated
// channels
func NewAuthenticator() *Authenticator {
	return &Authenticator{
		errChan:  make(chan error),
		doneChan: make(chan struct{}),
	}
}

// Login performs the login flow for a user. This flow uses Supabase auth and a
// local server for handling the login. Once the server has completed and received
// the session, the server is shut down and control is returned back to the CLI.
func (a *Authenticator) Login() (string, error) {
	supabaseAuthURL := fmt.Sprintf("%s/auth/v1/authorize", prodSupabaseURL)

	// 1. Generate the PKCE
	codeVerifier, codeChallenge, err := a.generatePkce(codeChallengeLength)
	if err != nil {
		return "", fmt.Errorf("PKCE error: %v", err)
	}

	a.codeVerifier = codeVerifier

	// 2. Start the local, callback login server
	r := chi.NewRouter()
	r.Get("/", a.handleLocalCallback)

	server := &http.Server{
		Addr:    authCallbackAddr,
		Handler: r,
	}

	go func() {
		a.errChan <- server.ListenAndServe()
	}()

	// 3. Open the browser to access the auth service with the necessary query params
	queryParams := url.Values{
		"provider":              {"github"},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
		"redirect_to":           {"http://" + authCallbackAddr + "/"},
	}

	authenticationURL := supabaseAuthURL + "?" + queryParams.Encode()
	err = browser.OpenURL(authenticationURL)
	if err != nil {
		fmt.Printf("Failed to open the browser: %s\nManually use authentication URL:", err)
		fmt.Println(authenticationURL)
	}

	// 4. Wait for results
	select {
	case err := <-a.errChan:
		if err != nil && err != http.ErrServerClosed {
			return "", err
		}
	case <-a.doneChan:
		a.shutdownServer(server)
	case <-time.After(60 * time.Second):
		a.shutdownServer(server)
		return "", fmt.Errorf("authentication timeout")
	}

	return a.username, nil
}

// handleLocalCallback is the callback route handler for the local server to get
// the results from the authentication service. It gets the session and saves it.
func (a *Authenticator) handleLocalCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "'code' query param not found", http.StatusBadRequest)
		a.errChan <- fmt.Errorf("'code' query param not found")
		return
	}

	sessionData, err := a.getSession(code, a.codeVerifier)
	if err != nil {
		http.Error(w, "Access token exchange failed", http.StatusInternalServerError)
		a.errChan <- fmt.Errorf("getting session failed: %w", err)
		return
	}

	if err := a.saveSession(sessionData); err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		a.errChan <- fmt.Errorf("could not save session: %w", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write([]byte(successHTML))
	if err != nil {
		a.errChan <- fmt.Errorf("error writing response writer: %s", err)
		return
	}

	a.username = sessionData.User.UserMetadata["user_name"].(string)
	a.doneChan <- struct{}{}
}

// shutdownServer shuts down the local callback server. This function panics
// if there are errors shutting down the server.
func (a *Authenticator) shutdownServer(server *http.Server) {
	err := server.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

// CheckSession checks if a session is already authenticated based on the expiration
// time for the given session on disk.
func (a *Authenticator) CheckSession() error {
	session, err := a.readSessionFile()
	if err != nil {
		return fmt.Errorf("failed to read session file: %w", err)
	}

	// Check if session is expired or about to expire (within 5 minutes)
	if time.Now().Add(5 * time.Minute).After(time.Unix(session.ExpiresAt, 0)) {
		return fmt.Errorf("session expired")
	}

	return nil
}

// readSessionFile reads a session file and returns the session struct.
func (a *Authenticator) readSessionFile() (*session, error) {
	configDir, err := config.GetConfigDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	sessionFile := path.Join(configDir, sessionFileName)

	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return nil, err
	}

	var session session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// generatePkce creates a "Proof Key for Code Exchange" (PKCE) for use in the auth
// service's authentication flow.
//
// Note on "rand.Reader" from the Go docs:
// > Reader is a global, shared instance of a cryptographically
// > secure random number generator.
func (a *Authenticator) generatePkce(length int) (string, string, error) {
	p := make([]byte, length)

	_, err := rand.Read(p)
	if err != nil {
		return "", "", fmt.Errorf("failed to read random bytes: %s", err)
	}

	verifier := base64.RawURLEncoding.EncodeToString(p)
	b := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(b[:])

	return verifier, challenge, nil
}

// getSession takes an authentication code and a verifier, using the Supabase
// auth service, to get a session
func (a *Authenticator) getSession(authCode, codeVerifier string) (*session, error) {
	url := fmt.Sprintf("%s/auth/v1/token?grant_type=pkce", prodSupabaseURL)

	payload := map[string]string{
		"auth_code":     authCode,
		"code_verifier": codeVerifier,
	}

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("ApiKey", prodSupabasePublicKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't make a request with the default client: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", res.Status)
	}

	var responseData session
	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %s", err)
	}

	return &responseData, nil
}

// saveSession saves a session to disk
func (a *Authenticator) saveSession(sessionData *session) error {
	dir, err := config.GetConfigDirectory()
	if err != nil {
		return fmt.Errorf("could not get user config directory: %w", err)
	}

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("marshaling session data failed: %w", err)
	}

	filePath := path.Join(dir, sessionFileName)
	if err := os.WriteFile(filePath, jsonData, 0600); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
