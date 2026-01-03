package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

type ClubHandler struct {
	uc *application.ClubUseCases
}

func NewClubHandler(uc *application.ClubUseCases) *ClubHandler {
	return &ClubHandler{uc: uc}
}

// CreateClub godoc
// @Summary Create a new club (Super Admin only)
// @Description Creates a new tenant
// @Tags clubs
// @Accept json
// @Produce json
// @Success 201 {object} domain.Club
// @Router /clubs [post]
func (h *ClubHandler) CreateClub(c *gin.Context) {
	if !h.isSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires SUPER_ADMIN role"})
		return
	}

	var req application.CreateClubDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	club, err := h.uc.CreateClub(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, club)
}

// ListClubs godoc
// @Summary List all clubs (Super Admin only)
// @Description Lists all tenants
// @Tags clubs
// @Produce json
// @Success 200 {array} domain.Club
// @Router /clubs [get]
func (h *ClubHandler) ListClubs(c *gin.Context) {
	if !h.isSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires SUPER_ADMIN role"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	clubs, err := h.uc.ListClubs(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clubs)
}

// UpdateClub godoc
// @Summary Update a club (Super Admin only)
// @Description Updates tenant details
// @Tags clubs
// @Accept json
// @Produce json
// @Router /clubs/{id} [put]
func (h *ClubHandler) UpdateClub(c *gin.Context) {
	if !h.isSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires SUPER_ADMIN role"})
		return
	}

	id := c.Param("id")
	var req application.UpdateClubDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	club, err := h.uc.UpdateClub(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, club)
}

func (h *ClubHandler) isSuperAdmin(c *gin.Context) bool {
	// 1. Check if userRole is set in context (from AuthMiddleware)
	role, exists := c.Get("userRole")
	if !exists {
		// Fallback: Check if we can get it from DB if we have userID?
		// Ideally, AuthMiddleware should set this.
		return false
	}
	return role == userDomain.RoleSuperAdmin
}

func RegisterRoutes(r *gin.RouterGroup, h *ClubHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	clubs := r.Group("/clubs")
	clubs.Use(authMiddleware, tenantMiddleware)
	{
		clubs.POST("", h.CreateClub)
		clubs.GET("", h.ListClubs)
		clubs.PUT("/:id", h.UpdateClub)
	}
}
