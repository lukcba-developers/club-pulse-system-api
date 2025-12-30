package logger

import (
	"log/slog"
	"os"
)

// InitLogger initializes the global logger with JSON handler
func InitLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug, // Default to debug for dev; make configurable later
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize keys if needed, e.g. "time" -> "@timestamp"
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}
