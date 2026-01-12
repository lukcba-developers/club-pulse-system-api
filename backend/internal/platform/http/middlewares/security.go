package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds OWASP recommended security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// HSTS - Enforce HTTPS for 1 year
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// Prevent MIME-type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// XSS Protection (Legacy browsers)
		c.Header("X-XSS-Protection", "1; mode=block")

		// Content Security Policy - Restrict sources
		// Adjust this based on frontend needs (e.g. allowing scripts from specific CDNs)
		// For development, allow connect-src to localhost:8080 (backend API)
		c.Header("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; font-src 'self'; script-src 'self' 'unsafe-inline'; connect-src 'self' http://localhost:8080 ws://localhost:3000")

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

// CORSMiddleware configures Cross-Origin Resource Sharing with a whitelist
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	// Build origin map for O(1) lookup
	originSet := make(map[string]bool)
	for _, o := range allowedOrigins {
		originSet[o] = true
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// SECURITY: Only allow origins in whitelist
		if originSet[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Club-ID, x-request-id")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
