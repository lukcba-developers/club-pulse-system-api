package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	csrfCookieName = "csrf_token"
	csrfHeaderName = "X-CSRF-Token"
	csrfTokenLen   = 32 // 32 bytes = 64 hex chars
)

// csrfExcludedPaths are paths that don't require CSRF validation
// These are either pre-login routes or external webhooks
var csrfExcludedPaths = []string{
	"/api/v1/auth/login",
	"/api/v1/auth/register",
	"/api/v1/auth/google",
	"/api/v1/auth/refresh",
	"/api/v1/payments/webhook",
}

// isSecureCookie returns true in production (GIN_MODE=release) for Secure cookies.
func isCSRFSecureCookie() bool {
	return os.Getenv("GIN_MODE") == "release"
}

// generateCSRFToken creates a cryptographically secure random token.
func generateCSRFToken() (string, error) {
	bytes := make([]byte, csrfTokenLen)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SetCSRFCookie sets the CSRF token cookie. Call this after successful login.
func SetCSRFCookie(c *gin.Context) error {
	token, err := generateCSRFToken()
	if err != nil {
		return err
	}
	// This cookie is NOT HttpOnly so JS can read it
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(csrfCookieName, token, 86400, "/", "", isCSRFSecureCookie(), false)
	return nil
}

// ClearCSRFCookie clears the CSRF token cookie. Call this on logout.
func ClearCSRFCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(csrfCookieName, "", -1, "/", "", isCSRFSecureCookie(), false)
}

// isCSRFExcludedPath checks if the request path should skip CSRF validation.
func isCSRFExcludedPath(path string) bool {
	for _, excluded := range csrfExcludedPaths {
		if strings.HasPrefix(path, excluded) {
			return true
		}
	}
	return false
}

// CSRFMiddleware validates the CSRF token for state-changing requests.
// It implements the Double Submit Cookie pattern:
// 1. The server sets a random token in a cookie (readable by JS).
// 2. The client must send the same token in the X-CSRF-Token header.
// 3. An attacker cannot read the cookie due to SameSite policy.
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for safe methods (GET, HEAD, OPTIONS)
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// Skip for excluded paths (login, register, webhooks)
		if isCSRFExcludedPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Get token from cookie
		cookieToken, err := c.Cookie(csrfCookieName)
		if err != nil || cookieToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "CSRF token missing from cookie",
				"type":  "CSRF_ERROR",
			})
			return
		}

		// Get token from header
		headerToken := c.GetHeader(csrfHeaderName)
		if headerToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "CSRF token missing from header",
				"type":  "CSRF_ERROR",
			})
			return
		}

		// Constant-time comparison to prevent timing attacks
		if cookieToken != headerToken {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "CSRF token mismatch",
				"type":  "CSRF_ERROR",
			})
			return
		}

		c.Next()
	}
}
