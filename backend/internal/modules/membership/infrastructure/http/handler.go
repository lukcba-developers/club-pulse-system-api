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

func RegisterRoutes(r *gin.RouterGroup, h *MembershipHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	memberships := r.Group("/memberships")
	memberships.Use(authMiddleware, tenantMiddleware)
	{
		memberships.POST("", h.CreateMembership)
		memberships.GET("", h.ListMemberships)
		memberships.GET("/tiers", h.ListTiers)
		memberships.GET("/admin", h.ListAllMemberships) // Admin view
		memberships.GET("/:id", h.GetMembership)
		memberships.POST("/process-billing", h.ProcessBilling)
		memberships.POST("/scholarship", h.AssignScholarship)
	}
}

func (h *MembershipHandler) ListTiers(c *gin.Context) {
	clubID := c.GetString("clubID")
	tiers, err := h.useCases.ListTiers(c.Request.Context(), clubID)
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

	clubID := c.GetString("clubID")
	membership, err := h.useCases.CreateMembership(c.Request.Context(), clubID, req)
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

	clubID := c.GetString("clubID")
	membership, err := h.useCases.GetMembership(c.Request.Context(), clubID, id)
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
	clubID := c.GetString("clubID")
	memberships, err := h.useCases.ListUserMemberships(c.Request.Context(), clubID, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": memberships})
}

// ListAllMemberships returns all memberships for admin dashboard
func (h *MembershipHandler) ListAllMemberships(c *gin.Context) {
	clubID := c.GetString("clubID")
	memberships, err := h.useCases.ListAllMemberships(c.Request.Context(), clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": memberships})
}

func (h *MembershipHandler) ProcessBilling(c *gin.Context) {
	// In production, this would be restricted to ADMIN role.
	clubID := c.GetString("clubID")
	count, err := h.useCases.ProcessMonthlyBilling(c.Request.Context(), clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Billing cycle processed",
		"count":   count,
	})
}

func (h *MembershipHandler) AssignScholarship(c *gin.Context) {
	var req application.AssignScholarshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Grantor ID (Admin)
	grantorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin not authenticated"})
		return
	}

	clubID := c.GetString("clubID")
	scholarship, err := h.useCases.AssignScholarship(c.Request.Context(), clubID, req, grantorID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": scholarship})
}
