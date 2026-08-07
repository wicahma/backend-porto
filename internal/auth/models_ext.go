package auth

import "time"

// LoginResponse is returned after successful authentication or refresh.
type LoginResponse struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	ExpiresAt    time.Time       `json:"expires_at"`
	User         SessionResponse `json:"user"`
}

// SessionResponse is returned when validating an existing session.
type SessionResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// GitHubUser is the user info returned by GitHub's /user API.
type GitHubUser struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// GitHubEmail is an email entry from GitHub's /user/emails API.
type GitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// GitHubTokenResponse is the response from GitHub's OAuth token exchange.
type GitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}
