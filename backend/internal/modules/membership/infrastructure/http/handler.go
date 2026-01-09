package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/application"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/middleware"
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
		// Tenant/User routes
		memberships.POST("", h.CreateMembership)
		memberships.GET("", h.ListMemberships)
		memberships.GET("/tiers", h.ListTiers)
		memberships.GET("/:id", h.GetMembership)
		memberships.DELETE("/:id", h.CancelMembership)

		// Admin Routes
		adminOnly := memberships.Group("")
		adminOnly.Use(middleware.RequireRole(userDomain.RoleAdmin, userDomain.RoleSuperAdmin))
		{
			adminOnly.GET("/admin", h.ListAllMemberships) // Admin view
			adminOnly.POST("/process-billing", h.ProcessBilling)
			adminOnly.POST("/scholarship", h.AssignScholarship)
		}
	}
}

// CancelMembership cancels a membership by ID
func (h *MembershipHandler) CancelMembership(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"type": "invalid_format", "error": "invalid membership id"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"type": "unauthorized", "error": "unauthorized"})
		return
	}

	clubID := c.GetString("clubID")
	role := c.GetString("userRole")

	// Admin can cancel any membership, user can only cancel their own
	requestingUserID := userID.(string)
	if role == userDomain.RoleAdmin || role == userDomain.RoleSuperAdmin {
		// Admin can cancel any - pass the membership's user ID instead
		// We fetch first to get the actual owner
		membership, err := h.useCases.GetMembership(c.Request.Context(), clubID, id)
		if err != nil || membership == nil {
			c.JSON(http.StatusNotFound, gin.H{"type": "not_found", "error": "membership not found"})
			return
		}
		requestingUserID = membership.UserID.String()
	}

	cancelled, err := h.useCases.CancelMembership(c.Request.Context(), clubID, id, requestingUserID)
	if err != nil {
		if err.Error() == "membership not found" {
			c.JSON(http.StatusNotFound, gin.H{"type": "not_found", "error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"type": "cancel_unauthorized", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "membership cancelled", "data": cancelled})
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

	// SECURITY: Get authenticated user
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role := c.GetString("userRole")
	clubID := c.GetString("clubID")

	membership, err := h.useCases.GetMembership(c.Request.Context(), clubID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if membership == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "membership not found"})
		return
	}

	// SECURITY FIX (VUL-004): Validate ownership or admin role
	if membership.UserID.String() != userID.(string) &&
		role != userDomain.RoleAdmin &&
		role != userDomain.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
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
	// RBAC: Handled by middleware

	clubID := c.GetString("clubID")
	memberships, err := h.useCases.ListAllMemberships(c.Request.Context(), clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": memberships})
}

func (h *MembershipHandler) ProcessBilling(c *gin.Context) {
	// RBAC: Handled by middleware

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
	// RBAC: Handled by middleware

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
