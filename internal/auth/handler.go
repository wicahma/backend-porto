package auth

import (
	"encoding/json"
	"net/http"
	"net/url"

	"backend-porto/internal/shared/response"
)

// Handler handles HTTP requests for authentication.
type Handler struct {
	svc *Service
}

// NewHandler creates a new auth Handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Login redirects the user to GitHub OAuth.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.svc.GitHubLoginURL(), http.StatusTemporaryRedirect)
}

// Callback handles the GitHub OAuth callback and redirects to the frontend with tokens.
func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		response.Error(w, http.StatusBadRequest, "missing authorization code")
		return
	}

	result, err := h.svc.HandleGitHubCallback(r.Context(), code)
	if err != nil {
		http.Redirect(w, r, h.svc.frontendURL+"/auth/callback?error="+url.QueryEscape(err.Error()), http.StatusTemporaryRedirect)
		return
	}

	q := url.Values{}
	q.Set("access_token", result.AccessToken)
	q.Set("refresh_token", result.RefreshToken)
	http.Redirect(w, r, h.svc.frontendURL+"/auth/callback?"+q.Encode(), http.StatusTemporaryRedirect)
}

// Session validates the current session token and returns user info.
func (h *Handler) Session(w http.ResponseWriter, r *http.Request) {
	tokenString := extractTokenFromRequest(r)
	if tokenString == "" {
		response.Error(w, http.StatusUnauthorized, "missing token")
		return
	}

	session, err := h.svc.ValidateSession(tokenString)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid token")
		return
	}

	response.Success(w, session)
}

// Refresh rotates a refresh token and returns a new token pair.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "missing refresh token")
		return
	}

	result, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	response.Success(w, result)
}

// SignOut invalidates the refresh token.
func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.RefreshToken != "" {
		_ = h.svc.SignOut(r.Context(), req.RefreshToken)
	}

	response.Success(w, map[string]string{"message": "signed out"})
}

func extractTokenFromRequest(r *http.Request) string {
	cookie, err := r.Cookie("porto_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	return ""
}
