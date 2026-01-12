package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	clubDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

const HeaderClubID = "X-Club-ID"
const ContextClubID = "clubID"
const ContextUserRole = "userRole"

func TenantMiddleware(clubRepo clubDomain.ClubRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		clubID := c.GetHeader(HeaderClubID)

		// Public Routes Bypass (But set ClubID if provided in header)
		publicPaths := map[string]bool{
			"/api/v1/health":        true,
			"/api/v1/auth/login":    true,
			"/api/v1/auth/register": true,
			"/api/v1/auth/refresh":  true,
			"/api/v1/auth/google":   true,
		}
		if publicPaths[c.Request.URL.Path] {
			if clubID != "" {
				// Resolve slug/domain to UUID if needed (for non-UUID values)
				resolvedID := resolveClubID(c.Request.Context(), clubRepo, clubID)
				c.Set(ContextClubID, resolvedID)
			}
			c.Next()
			return
		}

		// 1. Get Auth Context
		role, _ := c.Get(ContextUserRole)
		tokenClubID, _ := c.Get("userClubID")

		// 2. Super Admin Bypass (Trust Header if provided, otherwise context optional)
		if role == userDomain.RoleSuperAdmin {
			if clubID != "" {
				c.Set(ContextClubID, clubID)
			}
			c.Next()
			return
		}

		// 3. Strict Validation for Regular Users
		// Header is optional if we trust token implicitly, BUT legacy frontend sends header.
		// CRITICAL: Ensure Header (if present) == Token Club ID
		userClubID, ok := tokenClubID.(string)
		if !ok || userClubID == "" {
			// Should be caught by AuthMiddleware, but defensive programming
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User context invalid"})
			c.Abort()
			return
		}

		if clubID != "" && clubID != userClubID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this club"})
			c.Abort()
			return
		}

		// Force Context to be the verified User Club ID
		c.Set(ContextClubID, userClubID)
		c.Next()
	}
}

// resolveClubID attempts to resolve a slug/domain to a UUID if the input is not already a UUID.
// This allows the frontend to send either the UUID directly or a slug like 'club-alpha'.
func resolveClubID(ctx context.Context, clubRepo clubDomain.ClubRepository, input string) string {
	// Check if input looks like a UUID (contains dashes and is ~36 chars)
	if len(input) == 36 && strings.Count(input, "-") == 4 {
		return input // Already a UUID
	}

	// Attempt to resolve slug to UUID
	club, err := clubRepo.GetBySlug(ctx, input)
	if err == nil && club != nil {
		return club.ID
	}

	// Fallback: return as-is (will fail gracefully in auth)
	return input
}
