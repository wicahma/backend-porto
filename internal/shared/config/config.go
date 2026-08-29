package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port     string
	LogLevel string

	// Whether to write logs as JSON (true) or human-readable text (false).
	// Human-readable is for local dev; set LOG_FORMAT=json when shipping logs
	// to Grafana Loki / OpenTelemetry / Vector, which ingest JSON lines.
	LogFormat string

	// Security
	APIKey          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	// Database
	DatabaseURL string
	RedisURL    string

	// Auth
	JWTSecret           string
	AdminEmail          string
	GitHubClientID      string
	GitHubClientSecret  string
	GitHubRedirectURI   string

	// Frontend
	FrontendURL string

	// S3 / Storage
	S3Region          string
	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3BucketName      string
}

// Load reads all configuration from environment variables.
// It returns an error if any required variable is missing.
func Load() (*Config, error) {
	// Load .env first; real env vars take precedence over it.
	LoadDotEnv(".env")

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		LogFormat:       getEnv("LOG_FORMAT", "text"),
		APIKey:          os.Getenv("API_KEY"),
		AccessTokenTTL:  getDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL: getDuration("REFRESH_TOKEN_TTL", 7*24*time.Hour),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		RedisURL:        getEnv("REDIS_URL", ""),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		AdminEmail:         os.Getenv("ADMIN_EMAIL"),
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GitHubRedirectURI:  os.Getenv("GITHUB_REDIRECT_URI"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		S3Region:          os.Getenv("S3_REGION"),
		S3Endpoint:        os.Getenv("S3_ENDPOINT"),
		S3AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
		S3SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
		S3BucketName:      getEnv("S3_BUCKET_NAME", "portfolio"),
	}

	required := map[string]string{
		"DATABASE_URL":         cfg.DatabaseURL,
		"API_KEY":              cfg.APIKey,
		"JWT_SECRET":           cfg.JWTSecret,
		"ADMIN_EMAIL":          cfg.AdminEmail,
		"GITHUB_CLIENT_ID":     cfg.GitHubClientID,
		"GITHUB_CLIENT_SECRET": cfg.GitHubClientSecret,
		"GITHUB_REDIRECT_URI":  cfg.GitHubRedirectURI,
		"S3_REGION":            cfg.S3Region,
		"S3_ENDPOINT":          cfg.S3Endpoint,
		"S3_ACCESS_KEY_ID":     cfg.S3AccessKeyID,
		"S3_SECRET_ACCESS_KEY": cfg.S3SecretAccessKey,
	}

	for key, val := range required {
		if val == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return d
}
