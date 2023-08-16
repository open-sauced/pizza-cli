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
	constants "github.com/open-sauced/pizza-cli/pkg"
	"github.com/spf13/cobra"
)

//go:embed success.html
var successHTML string

const loginLongDesc string = `Log into the application.

This command initiates the GitHub auth flow to log you into the application by launching your browser`

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
	codeVerifier, codeChallenge, err := pkce(constants.CodeChallengeLength)
	if err != nil {
		return fmt.Errorf("PKCE error: %v", err.Error())
	}

	supabaseAuthURL := fmt.Sprintf("https://%s.supabase.co/auth/v1/authorize", constants.SupabaseID)
	queryParams := url.Values{
		"provider":              {"github"},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
		"redirect_to":           {"http://" + constants.AuthCallbackAddr + "/"},
	}

	authenticationURL := supabaseAuthURL + "?" + queryParams.Encode()

	server := &http.Server{Addr: constants.AuthCallbackAddr}

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

func getSession(authCode, codeVerifier string) (*AccessTokenResponse, error) {
	url := fmt.Sprintf("https://%s.supabase.co/auth/v1/token?grant_type=pkce", constants.SupabaseID)

	payload := map[string]string{
		"auth_code":     authCode,
		"code_verifier": codeVerifier,
	}

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("ApiKey", constants.SupabasePublicKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't make a request with the default client: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", res.Status)
	}

	var responseData AccessTokenResponse
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

type AccessTokenResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	TokenType    string     `json:"token_type"`
	ExpiresIn    int        `json:"expires_in"`
	ExpiresAt    int        `json:"expires_at"`
	User         UserSchema `json:"user"`
}

type UserSchema struct {
	ID                     string                 `json:"id"`
	Aud                    string                 `json:"aud,omitempty"`
	Role                   string                 `json:"role"`
	Email                  string                 `json:"email"`
	EmailConfirmedAt       string                 `json:"email_confirmed_at"`
	Phone                  string                 `json:"phone"`
	PhoneConfirmedAt       string                 `json:"phone_confirmed_at"`
	ConfirmationSentAt     string                 `json:"confirmation_sent_at"`
	ConfirmedAt            string                 `json:"confirmed_at"`
	RecoverySentAt         string                 `json:"recovery_sent_at"`
	NewEmail               string                 `json:"new_email"`
	EmailChangeSentAt      string                 `json:"email_change_sent_at"`
	NewPhone               string                 `json:"new_phone"`
	PhoneChangeSentAt      string                 `json:"phone_change_sent_at"`
	ReauthenticationSentAt string                 `json:"reauthentication_sent_at"`
	LastSignInAt           string                 `json:"last_sign_in_at"`
	AppMetadata            map[string]interface{} `json:"app_metadata"`
	UserMetadata           map[string]interface{} `json:"user_metadata"`
	Factors                []MFAFactorSchema      `json:"factors"`
	Identities             []interface{}          `json:"identities"`
	BannedUntil            string                 `json:"banned_until"`
	CreatedAt              string                 `json:"created_at"`
	UpdatedAt              string                 `json:"updated_at"`
	DeletedAt              string                 `json:"deleted_at"`
}

type MFAFactorSchema struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	FriendlyName string `json:"friendly_name"`
	FactorType   string `json:"factor_type"`
}
