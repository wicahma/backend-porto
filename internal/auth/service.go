package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

const refreshKeyPrefix = "auth:refresh:"

// Service handles authentication business logic.
type Service struct {
	repo               *Repository
	redis              *redis.Client
	jwtSecret          []byte
	adminEmail         string
	githubClientID     string
	githubClientSecret string
	githubRedirectURI  string
	frontendURL        string
	accessTTL          time.Duration
	refreshTTL         time.Duration
	httpClient         *http.Client
}

// NewService creates a new auth Service.
func NewService(
	repo *Repository,
	redisClient *redis.Client,
	jwtSecret string,
	adminEmail string,
	githubClientID string,
	githubClientSecret string,
	githubRedirectURI string,
	frontendURL string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *Service {
	return &Service{
		repo:               repo,
		redis:              redisClient,
		jwtSecret:          []byte(jwtSecret),
		adminEmail:         adminEmail,
		githubClientID:     githubClientID,
		githubClientSecret: githubClientSecret,
		githubRedirectURI:  githubRedirectURI,
		frontendURL:        frontendURL,
		accessTTL:          accessTTL,
		refreshTTL:         refreshTTL,
		httpClient:         &http.Client{Timeout: 10 * time.Second},
	}
}

// GitHubLoginURL builds the GitHub OAuth authorize URL.
func (s *Service) GitHubLoginURL() string {
	v := url.Values{}
	v.Set("client_id", s.githubClientID)
	v.Set("redirect_uri", s.githubRedirectURI)
	v.Set("scope", "user:email")
	v.Set("allow_signup", "false")
	return "https://github.com/login/oauth/authorize?" + v.Encode()
}

// HandleGitHubCallback exchanges the OAuth code for a token, fetches the user's
// GitHub profile, validates the email, and issues a JWT.
func (s *Service) HandleGitHubCallback(ctx context.Context, code string) (*LoginResponse, error) {
	token, err := s.exchangeCodeForToken(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchanging code for token: %w", err)
	}

	ghUser, err := s.getGitHubUser(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("getting github user: %w", err)
	}

	if ghUser.Email == "" {
		return nil, fmt.Errorf("no public email on github account")
	}

	if ghUser.Email != s.adminEmail {
		return nil, fmt.Errorf("email %s is not authorized", ghUser.Email)
	}

	// Store admin user in database
	githubID := fmt.Sprintf("%d", ghUser.ID)
	if err := s.repo.UpsertAdmin(ctx, ghUser.Email, githubID); err != nil {
		return nil, fmt.Errorf("upserting admin: %w", err)
	}

	return s.issueTokenPair(ctx, githubID, ghUser.Email)
}

// Refresh rotates a refresh token and issues a new token pair.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	if s.redis == nil {
		return nil, fmt.Errorf("redis not configured")
	}

	key := refreshKeyPrefix + refreshToken
	email, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	admin, err := s.repo.GetAdminByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("admin not found")
	}

	// rotate: invalidate old refresh token, issue a new pair
	_ = s.redis.Del(ctx, key).Err()
	return s.issueTokenPair(ctx, admin.GithubID.String, email)
}

// SignOut invalidates the refresh token.
func (s *Service) SignOut(ctx context.Context, refreshToken string) error {
	if s.redis == nil {
		return nil
	}
	return s.redis.Del(ctx, refreshKeyPrefix+refreshToken).Err()
}

func (s *Service) issueTokenPair(ctx context.Context, sub, email string) (*LoginResponse, error) {
	accessExp := time.Now().Add(s.accessTTL)
	accessClaims := jwt.MapClaims{
		"sub": sub,
		"iss": email,
		"iat": time.Now().Unix(),
		"exp": accessExp.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("signing access token: %w", err)
	}

	refreshToken, err := s.newRefreshToken(ctx, sub, email)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExp,
		User: SessionResponse{
			UserID: sub,
			Email:  email,
		},
	}, nil
}

func (s *Service) newRefreshToken(ctx context.Context, sub, email string) (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generating refresh token: %w", err)
	}
	token := hex.EncodeToString(buf)

	if s.redis != nil {
		if err := s.redis.Set(ctx, refreshKeyPrefix+token, email, s.refreshTTL).Err(); err != nil {
			return "", fmt.Errorf("storing refresh token: %w", err)
		}
	}

	return token, nil
}

// ValidateSession validates a JWT token string and returns session info.
func (s *Service) ValidateSession(tokenString string) (*SessionResponse, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	userID, _ := claims.GetSubject()
	email, _ := claims.GetIssuer()

	return &SessionResponse{
		UserID: userID,
		Email:  email,
	}, nil
}

func (s *Service) exchangeCodeForToken(ctx context.Context, code string) (string, error) {
	reqBody := url.Values{}
	reqBody.Set("client_id", s.githubClientID)
	reqBody.Set("client_secret", s.githubClientSecret)
	reqBody.Set("code", code)
	reqBody.Set("redirect_uri", s.githubRedirectURI)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://github.com/login/oauth/access_token",
		stringsToReader(reqBody.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp GitHubTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("github returned no access token: %s", string(body))
	}

	return tokenResp.AccessToken, nil
}

func (s *Service) getGitHubUser(ctx context.Context, accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user GitHubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	// If no public email, fetch from /user/emails
	if user.Email == "" {
		email, err := s.getGitHubPrimaryEmail(ctx, accessToken)
		if err != nil {
			return nil, err
		}
		user.Email = email
	}

	return &user, nil
}

func (s *Service) getGitHubPrimaryEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var emails []GitHubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	return "", fmt.Errorf("no primary verified email found")
}

func stringsToReader(s string) io.Reader {
	return &stringsReader{s: s}
}

type stringsReader struct {
	s string
	i int64
}

func (r *stringsReader) Read(b []byte) (int, error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n := copy(b, r.s[r.i:])
	r.i += int64(n)
	return n, nil
}
