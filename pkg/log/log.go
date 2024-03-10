package log

import (
	"log/slog"
	"os"
	"strings"
)

var (
	Logger *slog.Logger
)

func init() {
	logEnv := os.Getenv("BUDDY_LOG")
	logLevel := slog.LevelError

	switch strings.ToLower(logEnv) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = 9999
	}

	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
}
