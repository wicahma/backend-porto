package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New builds the app-wide slog logger.
// format: "text" for human-readable terminal output, "json" when logs are
// shipped to an external system (Grafana Loki, OpenTelemetry, Vector).
func New(level, format string) *slog.Logger {
	var l slog.Level
	switch strings.ToLower(level) {
	case "debug":
		l = slog.LevelDebug
	case "warn", "warning":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     l,
		AddSource: l == slog.LevelDebug,
	}
	if strings.ToLower(format) == "json" {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

func Fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}
