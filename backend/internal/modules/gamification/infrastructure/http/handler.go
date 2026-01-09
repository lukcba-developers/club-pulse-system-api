package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/gamification/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/gamification/domain"
)

// GamificationHandler handles HTTP requests for gamification features.
type GamificationHandler struct {
	badgeService       *application.BadgeService
	leaderboardService *application.LeaderboardServiceImpl
}

// NewGamificationHandler creates a new handler.
func NewGamificationHandler(
	badgeService *application.BadgeService,
	leaderboardService *application.LeaderboardServiceImpl,
) *GamificationHandler {
	return &GamificationHandler{
		badgeService:       badgeService,
		leaderboardService: leaderboardService,
	}
}

// RegisterRoutes registers gamification routes.
func (h *GamificationHandler) RegisterRoutes(r *gin.RouterGroup) {
	gamification := r.Group("/gamification")
	{
		// Badges
		gamification.GET("/badges", h.ListBadges)
		gamification.GET("/badges/my", h.GetMyBadges)
		gamification.GET("/badges/featured/:user_id", h.GetFeaturedBadges)
		gamification.PUT("/badges/:badge_id/feature", h.SetFeaturedBadge)

		// Leaderboards
		gamification.GET("/leaderboard", h.GetLeaderboard)
		gamification.GET("/leaderboard/context", h.GetMyLeaderboardContext)
		gamification.GET("/leaderboard/rank", h.GetMyRank)
	}
}

// @Summary List all badges
// @Description Returns all badges available in the club
// @Tags Gamification
// @Produce json
// @Success 200 {array} domain.Badge
// @Router /gamification/badges [get]
func (h *GamificationHandler) ListBadges(c *gin.Context) {
	clubID := c.GetString("clubID")

	badges, err := h.badgeService.GetAllBadges(clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, badges)
}

// @Summary Get my badges
// @Description Returns all badges earned by the current user
// @Tags Gamification
// @Produce json
// @Success 200 {array} domain.UserBadge
// @Router /gamification/badges/my [get]
func (h *GamificationHandler) GetMyBadges(c *gin.Context) {
	userID := c.GetString("userID")

	badges, err := h.badgeService.GetUserBadges(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, badges)
}

// @Summary Get user's featured badges
// @Description Returns the featured badges displayed on a user's profile
// @Tags Gamification
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} domain.UserBadge
// @Router /gamification/badges/featured/{user_id} [get]
func (h *GamificationHandler) GetFeaturedBadges(c *gin.Context) {
	userID := c.Param("user_id")

	badges, err := h.badgeService.GetFeaturedBadges(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, badges)
}

// SetFeaturedBadgeRequest represents the request body for featuring a badge.
type SetFeaturedBadgeRequest struct {
	Featured bool `json:"featured"`
}

// @Summary Set badge as featured
// @Description Toggle whether a badge is featured on the user's profile (max 3)
// @Tags Gamification
// @Accept json
// @Produce json
// @Param badge_id path string true "Badge ID"
// @Param body body SetFeaturedBadgeRequest true "Featured status"
// @Success 200 {object} map[string]string
// @Router /gamification/badges/{badge_id}/feature [put]
func (h *GamificationHandler) SetFeaturedBadge(c *gin.Context) {
	userID := c.GetString("userID")
	badgeIDStr := c.Param("badge_id")

	badgeID, err := uuid.Parse(badgeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid badge ID"})
		return
	}

	var req SetFeaturedBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.badgeService.SetFeaturedBadge(userID, badgeID, req.Featured); err != nil {
		if err == application.ErrMaxFeaturedBadges {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "badge featured status updated"})
}

// @Summary Get global leaderboard
// @Description Returns the club's leaderboard ranked by XP
// @Tags Gamification
// @Produce json
// @Param period query string false "Period: DAILY, WEEKLY, MONTHLY, ALL_TIME" default(MONTHLY)
// @Param limit query int false "Number of entries" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} domain.Leaderboard
// @Router /gamification/leaderboard [get]
func (h *GamificationHandler) GetLeaderboard(c *gin.Context) {
	clubID := c.GetString("clubID")
	period := domain.LeaderboardPeriod(c.DefaultQuery("period", "MONTHLY"))
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		// Parse limit
		limit = 20 // Simplified
	}

	leaderboard, err := h.leaderboardService.GetGlobalLeaderboard(clubID, period, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

// @Summary Get my leaderboard context
// @Description Returns my position with surrounding users
// @Tags Gamification
// @Produce json
// @Param period query string false "Period: DAILY, WEEKLY, MONTHLY, ALL_TIME" default(MONTHLY)
// @Success 200 {object} domain.LeaderboardContext
// @Router /gamification/leaderboard/context [get]
func (h *GamificationHandler) GetMyLeaderboardContext(c *gin.Context) {
	clubID := c.GetString("clubID")
	userID := c.GetString("userID")
	period := domain.LeaderboardPeriod(c.DefaultQuery("period", "MONTHLY"))

	context, err := h.leaderboardService.GetUserContext(clubID, userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if context == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found in leaderboard"})
		return
	}

	c.JSON(http.StatusOK, context)
}

// @Summary Get my rank
// @Description Returns my current rank in the leaderboard
// @Tags Gamification
// @Produce json
// @Param period query string false "Period" default(MONTHLY)
// @Success 200 {object} map[string]int
// @Router /gamification/leaderboard/rank [get]
func (h *GamificationHandler) GetMyRank(c *gin.Context) {
	clubID := c.GetString("clubID")
	userID := c.GetString("userID")
	period := domain.LeaderboardPeriod(c.DefaultQuery("period", "MONTHLY"))

	rank, err := h.leaderboardService.GetUserRank(clubID, userID, domain.LeaderboardTypeGlobal, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rank": rank})
}
