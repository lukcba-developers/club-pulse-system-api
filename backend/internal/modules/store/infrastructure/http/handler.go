package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clubApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/store/application"
)

type StoreHandler struct {
	useCases     *application.StoreUseCases
	clubUseCases *clubApp.ClubUseCases
}

func NewStoreHandler(useCases *application.StoreUseCases, clubUseCases *clubApp.ClubUseCases) *StoreHandler {
	return &StoreHandler{
		useCases:     useCases,
		clubUseCases: clubUseCases,
	}
}

func (h *StoreHandler) PurchaseItems(c *gin.Context) {
	var req application.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Enforce context ClubID and UserID for security
	clubID := c.GetString("clubID")
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDVal.(string)
	req.ClubID = clubID
	req.UserID = &userID // Auth User

	order, err := h.useCases.PurchaseItems(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *StoreHandler) PublicPurchase(c *gin.Context) {
	var req application.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := c.Param("slug")
	club, err := h.clubUseCases.GetClubBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	req.ClubID = club.ID
	req.UserID = nil // Guest

	order, err := h.useCases.PurchaseItems(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *StoreHandler) GetCatalog(c *gin.Context) {
	clubID := c.GetString("clubID")
	category := c.Query("category")

	products, err := h.useCases.GetCatalog(clubID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *StoreHandler) GetPublicCatalog(c *gin.Context) {
	slug := c.Param("slug")
	category := c.Query("category")

	club, err := h.clubUseCases.GetClubBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	products, err := h.useCases.GetCatalog(club.ID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func RegisterRoutes(r *gin.RouterGroup, handler *StoreHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	store := r.Group("/store")
	store.Use(authMiddleware, tenantMiddleware)
	{
		store.POST("/purchase", handler.PurchaseItems)
		store.GET("/products", handler.GetCatalog)
	}

	// Public Routes
	public := r.Group("/public/clubs/:slug/store")
	{
		public.GET("/products", handler.GetPublicCatalog)
		public.POST("/purchase", handler.PublicPurchase)
	}
}
