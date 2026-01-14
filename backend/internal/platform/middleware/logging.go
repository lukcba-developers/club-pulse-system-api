package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// StructuredLogger logs request details in JSON format (Canonical Log Line)
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		// Process Request
		c.Next()

		// Calculate Metrics
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Get Trace ID and Span ID from OTel Context (injected by otelgin middleware)
		span := trace.SpanFromContext(c.Request.Context())
		var traceID, spanID string
		if span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
			spanID = span.SpanContext().SpanID().String()
		} else {
			// Fallback: Check headers if OTel middleware didn't run or failed
			traceID = c.GetHeader("traceparent") // W3C
			if traceID == "" {
				traceID = c.GetHeader("X-Request-ID")
			}
		}

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
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status_code", statusCode),
			slog.String("latency", latency.String()), // Human readable "120ms"
			slog.Int64("latency_ns", latency.Nanoseconds()),
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
			attrs = append(attrs, slog.Any("user_id", userID))
		}
		if clubID, exists := c.Get("userClubID"); exists {
			attrs = append(attrs, slog.Any("club_id", clubID))
		}

		// Log using standard slog
		slog.Log(c.Request.Context(), level, "HTTP Request", attrs...)
	}
}
