package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
)

type ClubHandler struct {
	useCases *application.ClubUseCases
}

func NewClubHandler(useCases *application.ClubUseCases) *ClubHandler {
	return &ClubHandler{useCases: useCases}
}

// --- Club Handlers (Super Admin) ---

type CreateClubRequest struct {
	Name     string `json:"name" binding:"required"`
	Domain   string `json:"domain"`
	Settings string `json:"settings"`
}

func (h *ClubHandler) CreateClub(c *gin.Context) {
	var req CreateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	club, err := h.useCases.CreateClub(req.Name, req.Domain, req.Settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, club)
}

func (h *ClubHandler) ListClubs(c *gin.Context) {
	clubs, err := h.useCases.ListClubs(100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": clubs})
}

// --- Sponsor Handlers ---

type SponsorRequest struct {
	Name        string `json:"name" binding:"required"`
	ContactInfo string `json:"contact_info"`
	LogoURL     string `json:"logo_url"`
}

func (h *ClubHandler) RegisterSponsor(c *gin.Context) {
	var req SponsorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	if clubID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ClubContext required"})
		return
	}

	sponsor, err := h.useCases.RegisterSponsor(clubID, req.Name, req.ContactInfo, req.LogoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sponsor)
}

type AdPlacementRequest struct {
	SponsorID    string              `json:"sponsor_id" binding:"required"`
	LocationType domain.LocationType `json:"location_type" binding:"required"`
	Detail       string              `json:"detail"`
	EndDate      time.Time           `json:"end_date" binding:"required"`
	Amount       float64             `json:"amount" binding:"required"`
}

func (h *ClubHandler) CreateAdPlacement(c *gin.Context) {
	var req AdPlacementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ad, err := h.useCases.CreateAdPlacement(req.SponsorID, req.LocationType, req.Detail, req.EndDate, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ad)
}

func (h *ClubHandler) GetActiveAds(c *gin.Context) {
	clubID := c.GetString("clubID")
	ads, err := h.useCases.GetActiveAds(clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ads})
}

func RegisterRoutes(r *gin.RouterGroup, handler *ClubHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	// Super Admin Routes (Clubs)
	clubs := r.Group("/admin/clubs")
	clubs.Use(authMiddleware) // Tenant middleware NOT needed for creating clubs? Or maybe logical "system" tenant.
	{
		clubs.POST("", handler.CreateClub)
		clubs.GET("", handler.ListClubs)
	}

	// Club Routes (Sponsors)
	clubGroup := r.Group("/club")
	clubGroup.Use(authMiddleware, tenantMiddleware)
	{
		clubGroup.POST("/sponsors", handler.RegisterSponsor)
		clubGroup.POST("/ads", handler.CreateAdPlacement)
		clubGroup.GET("/ads", handler.GetActiveAds)
	}
}
