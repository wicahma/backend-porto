package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-porto/internal/articles"
	"backend-porto/internal/auth"
	"backend-porto/internal/experiences"
	"backend-porto/internal/projects"
	"backend-porto/internal/shared/config"
	"backend-porto/internal/shared/database"
	"backend-porto/internal/shared/logger"
	"backend-porto/internal/shared/middleware"
	"backend-porto/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// ---- Load Config ----
	cfg, err := config.Load()
	if err != nil {
		// No logger yet; use the default slog logger (writes to stderr).
		// Only replace with logger.New after config loads.
		logger.Fatal("failed to load config", err)
	}
	slog.SetDefault(logger.New(cfg.LogLevel, cfg.LogFormat))

	// ---- Database ----
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to connect to postgres", err)
	}
	defer pool.Close()

	// Redis is optional
	redisClient, err := database.NewRedisClient(ctx, cfg.RedisURL)
	if err != nil {
		slog.Warn("redis not available", "error", err)
	}
	if redisClient != nil {
		defer redisClient.Close()
	}

	// ---- S3 Storage ----
	s3Client, err := storage.NewS3Client(ctx, storage.Config{
		Region:          cfg.S3Region,
		Endpoint:        cfg.S3Endpoint,
		AccessKeyID:     cfg.S3AccessKeyID,
		SecretAccessKey: cfg.S3SecretAccessKey,
		BucketName:      cfg.S3BucketName,
	})
	if err != nil {
		logger.Fatal("failed to create S3 client", err)
	}
	// Verify endpoint, credentials, and bucket before serving requests.
	if err := s3Client.CheckBucket(ctx); err != nil {
		logger.Fatal("S3 check failed (endpoint, credentials, or bucket invalid)", err)
	}
	slog.Info("S3 storage: connected and bucket verified")

	// ---- Auth Domain ----
	authQueries := auth.New(pool)
	authRepo := auth.NewRepository(authQueries)
	authSvc := auth.NewService(
		authRepo,
		redisClient,
		cfg.JWTSecret,
		cfg.AdminEmail,
		cfg.GitHubClientID,
		cfg.GitHubClientSecret,
		cfg.GitHubRedirectURI,
		cfg.FrontendURL,
		cfg.AccessTokenTTL,
		cfg.RefreshTokenTTL,
	)
	authHandler := auth.NewHandler(authSvc)

	// ---- Articles Domain ----
	articlesQueries := articles.New(pool)
	articlesRepo := articles.NewRepository(articlesQueries)
	articlesSvc := articles.NewService(articlesRepo)
	articlesHandler := articles.NewHandler(articlesSvc)

	// ---- Projects Domain ----
	projectsQueries := projects.New(pool)
	projectsRepo := projects.NewRepository(projectsQueries)
	projectsSvc := projects.NewService(projectsRepo)
	projectsHandler := projects.NewHandler(projectsSvc)

	// ---- Experiences Domain ----
	experiencesQueries := experiences.New(pool)
	experiencesRepo := experiences.NewRepository(experiencesQueries)
	experiencesSvc := experiences.NewService(experiencesRepo)
	experiencesHandler := experiences.NewHandler(experiencesSvc)

	// ---- Storage Domain ----
	storageHandler := storage.NewHandler(s3Client)

	// ---- Router Setup ----
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	// Auth routes (public)
	mux.HandleFunc("GET /api/auth/github/login", authHandler.Login)
	mux.HandleFunc("GET /api/auth/github/callback", authHandler.Callback)
	mux.HandleFunc("GET /api/auth/session", authHandler.Session)
	mux.HandleFunc("POST /api/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /api/auth/signout", authHandler.SignOut)

	// Articles — public
	mux.HandleFunc("GET /api/articles/count", articlesHandler.Count)
	mux.HandleFunc("GET /api/articles/slug/{slug}", articlesHandler.GetBySlug)
	mux.HandleFunc("GET /api/articles/related/{id}", articlesHandler.GetRelated)
	mux.HandleFunc("GET /api/articles/{id}", articlesHandler.GetByID)
	mux.HandleFunc("GET /api/articles", articlesHandler.List)

	// Articles — protected
	mux.HandleFunc("POST /api/articles", articlesHandler.Create)
	mux.HandleFunc("PUT /api/articles/{id}", articlesHandler.Update)
	mux.HandleFunc("DELETE /api/articles/{id}", articlesHandler.Delete)

	// Projects — public
	mux.HandleFunc("GET /api/projects/count", projectsHandler.Count)
	mux.HandleFunc("GET /api/projects/{id}", projectsHandler.GetByID)
	mux.HandleFunc("GET /api/projects", projectsHandler.List)

	// Projects — protected
	mux.HandleFunc("POST /api/projects", projectsHandler.Create)
	mux.HandleFunc("PUT /api/projects/{id}", projectsHandler.Update)
	mux.HandleFunc("DELETE /api/projects/{id}", projectsHandler.Delete)

	// Experiences — public
	mux.HandleFunc("GET /api/experiences/count", experiencesHandler.Count)
	mux.HandleFunc("GET /api/experiences/{id}", experiencesHandler.GetByID)
	mux.HandleFunc("GET /api/experiences", experiencesHandler.List)

	// Experiences — protected
	mux.HandleFunc("POST /api/experiences", experiencesHandler.Create)
	mux.HandleFunc("PUT /api/experiences/{id}", experiencesHandler.Update)
	mux.HandleFunc("DELETE /api/experiences/{id}", experiencesHandler.Delete)

	// Storage — public
	mux.HandleFunc("GET /api/storage/{path...}", storageHandler.Get)

	// Storage — protected
	mux.HandleFunc("POST /api/storage/upload", storageHandler.Upload)
	mux.HandleFunc("DELETE /api/storage", storageHandler.Delete)

	// ---- Middleware Chain ----
	authMiddleware := middleware.Auth(cfg.JWTSecret)
	apiKeyMiddleware := middleware.APIKey(cfg.APIKey)
	corsMiddleware := middleware.CORS(cfg.FrontendURL)

	// Apply middleware: CORS → Logging → API Key → Auth (selectively) → Handler
	var handler http.Handler = mux
	handler = middleware.Logging(handler)
	handler = corsMiddleware(handler)
	handler = apiKeyMiddleware(handler)

	// Wrap protected routes with auth middleware
	handler = selectiveAuth(handler, authMiddleware)

	// ---- Start Server ----
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		slogger := slog.With("component", "shutdown")
		slogger.Info("shutting down gracefully...")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slogger.Error("shutdown error", "error", err)
		}
	}()

	fmt.Printf("Backend Porto API starting on http://localhost%s\n", port)
	slog.Info("server listening", "addr", port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("server failed", err)
	}
}

// selectiveAuth applies auth middleware only to protected routes.
// Protected routes are those starting with POST/PUT/DELETE or specific GET paths.
func selectiveAuth(next http.Handler, authMW func(http.Handler) http.Handler) http.Handler {
	protected := authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Routes that require authentication
		if isProtectedRoute(r.Method, r.URL.Path) {
			protected.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isProtectedRoute(method, path string) bool {
	// All mutating methods on API routes require auth
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
		if len(path) >= 4 && path[:4] == "/api" {
			// Exclude auth routes from protection
			if len(path) >= 13 && path[:13] == "/api/auth/" {
				return false
			}
			return true
		}
	}

	// Session and signout are public
	return false
}

// Ensure pgxpool import is used (needed for concrete type in New function)
var _ *pgxpool.Pool
