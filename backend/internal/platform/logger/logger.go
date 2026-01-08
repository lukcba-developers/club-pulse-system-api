package logger

import (
	"log/slog"
	"os"
	"strings"
)

// sensitiveKeys are fields that should be redacted in logs
var sensitiveKeys = map[string]bool{
	"password":      true,
	"token":         true,
	"access_token":  true,
	"refresh_token": true,
	"secret":        true,
	"authorization": true,
	"api_key":       true,
	"apikey":        true,
	"credit_card":   true,
	"ssn":           true,
}

// InitLogger initializes the global logger with JSON handler and sensitive data redaction
func InitLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug, // Default to debug for dev; make configurable later
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// SECURITY: Redact sensitive fields
			if sensitiveKeys[strings.ToLower(a.Key)] {
				return slog.Attr{Key: a.Key, Value: slog.StringValue("[REDACTED]")}
			}
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
