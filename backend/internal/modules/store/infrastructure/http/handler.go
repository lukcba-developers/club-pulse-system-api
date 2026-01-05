package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/store/application"
)

type StoreHandler struct {
	useCases *application.StoreUseCases
}

func NewStoreHandler(useCases *application.StoreUseCases) *StoreHandler {
	return &StoreHandler{useCases: useCases}
}

func (h *StoreHandler) PurchaseItems(c *gin.Context) {
	var req application.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Enforce context ClubID and UserID for security
	clubID := c.GetString("clubID")
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	req.ClubID = clubID
	req.UserID = userID.(string)

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

func RegisterRoutes(r *gin.RouterGroup, handler *StoreHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	store := r.Group("/store")
	store.Use(authMiddleware, tenantMiddleware)
	{
		store.POST("/purchase", handler.PurchaseItems)
		store.GET("/products", handler.GetCatalog)
	}
}
