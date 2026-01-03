package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

// StructuredLogger logs request details in JSON format
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		// Get Trace ID from OTel if available
		spanContext := trace.SpanContextFromContext(c.Request.Context())
		var traceID string
		if spanContext.HasTraceID() {
			traceID = spanContext.TraceID().String()
		} else {
			// Fallback (or if OTel middleware is missing/disabled)
			traceID = c.GetHeader("X-Trace-ID")
			if traceID == "" {
				traceID = uuid.New().String()
			}
		}

		// Ensure it's in header for downstream if not already handled by propagator
		c.Header("X-Trace-ID", traceID)

		// Process Request
		c.Next()

		// Calculate Latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Log Level based on status
		level := slog.LevelInfo
		if statusCode >= 500 {
			level = slog.LevelError
		} else if statusCode >= 400 {
			level = slog.LevelWarn
		}

		// Structured Log Attributes
		attrs := []any{
			slog.String("trace_id", traceID),
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
			slog.String("ip", clientIP),
		}

		if raw != "" {
			attrs = append(attrs, slog.String("query", raw))
		}
		if errorMessage != "" {
			attrs = append(attrs, slog.String("error", errorMessage))
		}

		// Add User Metadata if present
		if userID, exists := c.Get("userID"); exists {
			attrs = append(attrs, slog.Any("user_id", userID))
		}
		if clubID, exists := c.Get("userClubID"); exists {
			attrs = append(attrs, slog.Any("club_id", clubID))
		}

		slog.Log(c.Request.Context(), level, "HTTP Request", attrs...)
	}
}
