package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/application"
)

type MembershipHandler struct {
	useCases *application.MembershipUseCases
}

func NewMembershipHandler(useCases *application.MembershipUseCases) *MembershipHandler {
	return &MembershipHandler{useCases: useCases}
}

func RegisterRoutes(r *gin.RouterGroup, h *MembershipHandler, authMiddleware gin.HandlerFunc) {
	memberships := r.Group("/memberships")
	memberships.Use(authMiddleware)
	{
		memberships.POST("", h.CreateMembership)
		memberships.GET("", h.ListMemberships)
		memberships.GET("/tiers", h.ListTiers)
		memberships.GET("/:id", h.GetMembership)
	}
}

func (h *MembershipHandler) ListTiers(c *gin.Context) {
	tiers, err := h.useCases.ListTiers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tiers})
}

func (h *MembershipHandler) CreateMembership(c *gin.Context) {
	var req application.CreateMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Override UserID from token if needed, or validate it matches
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Force UserID from token to ensure security
	req.UserID = uuid.MustParse(userID.(string))

	membership, err := h.useCases.CreateMembership(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": membership})
}

func (h *MembershipHandler) GetMembership(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid membership id"})
		return
	}

	membership, err := h.useCases.GetMembership(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": membership})
}

func (h *MembershipHandler) ListMemberships(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	uid := uuid.MustParse(userID.(string))
	memberships, err := h.useCases.ListUserMemberships(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": memberships})
}
