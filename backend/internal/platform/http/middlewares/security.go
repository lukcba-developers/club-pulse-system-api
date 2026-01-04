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

// CORSMiddleware configures Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// In production, replace "*" with specific allowed origins from config
		// origin := c.Request.Header.Get("Origin")
		// if isAllowedOrigin(origin) {
		// 	c.Header("Access-Control-Allow-Origin", origin)
		// }

		// For MVP/Dev, we can be slightly permissive but strict on methods/headers
		origin := c.Request.Header.Get("Origin")
		// fmt.Printf("DEBUG CORS: Origin=%s Method=%s Path=%s\n", origin, c.Request.Method, c.Request.URL.Path)

		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin) // Dynamic allow for localhost:3000
		} else {
			// Fallback for tools or direct access?
			// c.Header("Access-Control-Allow-Origin", "*") // Don't do this with Credentials=true
		}
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Club-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
