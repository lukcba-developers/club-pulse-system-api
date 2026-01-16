package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/domain"
)

// AuthMiddleware creates a Gin middleware for authentication
func AuthMiddleware(tokenService domain.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Try Cookie
		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			// 2. Fallback to Header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
				return
			}

			// Expect "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
				return
			}
			tokenString = parts[1]
		}
		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"type":  "TOKEN_EXPIRED",
			})
			return
		}

		// Set userID in context for downstream handlers
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Set("userClubID", claims.ClubID)
		c.Next()
	}
}
