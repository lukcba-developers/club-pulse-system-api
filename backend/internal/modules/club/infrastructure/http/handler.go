package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/middleware"
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
	Slug     string `json:"slug" binding:"required"`
	Domain   string `json:"domain"`
	Settings string `json:"settings"`
}

func (h *ClubHandler) CreateClub(c *gin.Context) {
	var req CreateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	club, err := h.useCases.CreateClub(ctx, req.Name, req.Slug, req.Domain, req.Settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, club)
}

func (h *ClubHandler) ListClubs(c *gin.Context) {
	ctx := c.Request.Context()
	clubs, err := h.useCases.ListClubs(ctx, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": clubs})
}

func (h *ClubHandler) GetPublicClubBySlug(c *gin.Context) {
	slug := c.Param("slug")
	ctx := c.Request.Context()
	club, err := h.useCases.GetClubBySlug(ctx, slug)
	if err != nil {
		// Assuming error means not found or DB error.
		// For stricter 404, we'd check error type.
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}
	c.JSON(http.StatusOK, club)
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

	ctx := c.Request.Context()
	sponsor, err := h.useCases.RegisterSponsor(ctx, clubID, req.Name, req.ContactInfo, req.LogoURL)
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

	ctx := c.Request.Context()
	ad, err := h.useCases.CreateAdPlacement(ctx, req.SponsorID, req.LocationType, req.Detail, req.EndDate, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ad)
}

func (h *ClubHandler) GetActiveAds(c *gin.Context) {
	clubID := c.GetString("clubID")
	ctx := c.Request.Context()
	ads, err := h.useCases.GetActiveAds(ctx, clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ads})
}

func (h *ClubHandler) GetPublicActiveAds(c *gin.Context) {
	slug := c.Param("slug")
	ctx := c.Request.Context()
	club, err := h.useCases.GetClubBySlug(ctx, slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	ads, err := h.useCases.GetActiveAds(ctx, club.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ads})
}

func (h *ClubHandler) GetPublicNews(c *gin.Context) {
	slug := c.Param("slug")
	ctx := c.Request.Context()
	news, err := h.useCases.GetPublicNews(ctx, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": news})
}

type PublishNewsRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ImageURL string `json:"image_url"`
	Notify   bool   `json:"notify"`
}

func (h *ClubHandler) PublishNews(c *gin.Context) {
	var req PublishNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	if clubID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ClubContext required"})
		return
	}

	// Assuming we have a method PublishNews in UseCases
	ctx := c.Request.Context()
	news, err := h.useCases.PublishNews(ctx, clubID, req.Title, req.Content, req.ImageURL, req.Notify)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, news)
}

func RegisterRoutes(r *gin.RouterGroup, handler *ClubHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	// Public Routes
	public := r.Group("/public/clubs")
	{
		public.GET("/:slug", handler.GetPublicClubBySlug)
		public.GET("/:slug/ads", handler.GetPublicActiveAds)
		public.GET("/:slug/news", handler.GetPublicNews)
	}

	// Super Admin Routes (Clubs) - CRITICAL: Only SUPER_ADMIN can manage clubs
	clubs := r.Group("/admin/clubs")
	clubs.Use(authMiddleware, middleware.RequireRole(userDomain.RoleSuperAdmin))
	{
		clubs.POST("", handler.CreateClub)
		clubs.GET("", handler.ListClubs)
	}

	// Club Admin Routes (Sponsors & News Management)
	adminClubGroup := r.Group("/club")
	adminClubGroup.Use(authMiddleware, tenantMiddleware, middleware.RequireRole(userDomain.RoleAdmin, userDomain.RoleSuperAdmin))
	{
		adminClubGroup.POST("/sponsors", handler.RegisterSponsor)
		adminClubGroup.POST("/ads", handler.CreateAdPlacement)
		adminClubGroup.POST("/news", handler.PublishNews)
	}

	// Club Member Routes (View Access)
	memberClubGroup := r.Group("/club")
	memberClubGroup.Use(authMiddleware, tenantMiddleware)
	{
		memberClubGroup.GET("/ads", handler.GetActiveAds)
	}
}
