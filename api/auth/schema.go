package auth

type session struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int64       `json:"expires_in"`
	ExpiresAt    int64       `json:"expires_at"`
	User         sessionUser `json:"user"`
}

type sessionUser struct {
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
	Factors                []sessionUseFactor     `json:"factors"`
	Identities             []interface{}          `json:"identities"`
	BannedUntil            string                 `json:"banned_until"`
	CreatedAt              string                 `json:"created_at"`
	UpdatedAt              string                 `json:"updated_at"`
	DeletedAt              string                 `json:"deleted_at"`
}

type sessionUseFactor struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	FriendlyName string `json:"friendly_name"`
	FactorType   string `json:"factor_type"`
}
