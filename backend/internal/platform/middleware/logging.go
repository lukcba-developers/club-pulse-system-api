package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

// StructuredLogger logs request details in JSON format (Canonical Log Line)
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		// Get Trace ID and Span ID from OTel if available
		span := trace.SpanFromContext(c.Request.Context())
		var traceID, spanID string
		if span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
			spanID = span.SpanContext().SpanID().String()
		} else {
			// Fallback: Check for incoming headers from Frontend/Gateway
			traceID = c.GetHeader("X-Request-ID")
			if traceID == "" {
				traceID = c.GetHeader("X-Trace-ID")
			}
			// If still empty, generate a new one to correlate logs downstream
			if traceID == "" {
				traceID = uuid.New().String()
			}
		}

		// Propagate Trace ID downstream via Response Header
		c.Header("X-Trace-ID", traceID)

		// Process Request
		c.Next()

		// Calculate Metrics
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

		// Canonical Log Line Attributes
		attrs := []any{
			slog.String("trace_id", traceID),
			slog.String("span_id", spanID),
			slog.String("user_id", ""), // Placeholder, overwritten below if present
			slog.String("club_id", ""), // Placeholder, overwritten below if present
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status_code", statusCode),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.String("ip", clientIP),
		}

		if raw != "" {
			attrs = append(attrs, slog.String("query", raw))
		}
		if errorMessage != "" {
			attrs = append(attrs, slog.String("error", errorMessage))
		}

		// Add User Metadata if extracted by auth middleware
		if userID, exists := c.Get("userID"); exists {
			// Ensure userID is a string or compatible
			if uidStr, ok := userID.(string); ok {
				attrs = append(attrs, slog.String("user_id", uidStr))
			} else {
				attrs = append(attrs, slog.Any("user_id", userID))
			}
		}
		if clubID, exists := c.Get("userClubID"); exists {
			if cidStr, ok := clubID.(string); ok {
				attrs = append(attrs, slog.String("club_id", cidStr))
			} else {
				attrs = append(attrs, slog.Any("club_id", clubID))
			}
		}

		// Context is crucial here for correlation if the logger handler supports it
		slog.Log(c.Request.Context(), level, "HTTP Request", attrs...)
	}
}
