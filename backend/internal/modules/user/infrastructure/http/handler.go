package http

import (
	"net/http"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
)

type UserHandler struct {
	useCases *application.UserUseCases
}

func NewUserHandler(useCases *application.UserUseCases) *UserHandler {
	return &UserHandler{
		useCases: useCases,
	}
}

type UserResponse struct {
	*domain.User
	Category string `json:"category"`
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Extract ClubID
	clubID := c.GetString("clubID")
	role := c.GetString("userRole")

	// Super Admin special handling: their user record is in "system" default club
	if role == domain.RoleSuperAdmin {
		clubID = "system"
	}

	if clubID == "" {
		// Should be handled by middleware, but safety check
		c.JSON(http.StatusBadRequest, gin.H{"error": "Club context missing"})
		return
	}

	user, err := h.useCases.GetProfile(clubID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := UserResponse{
		User:     user,
		Category: user.CalculateCategory(),
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var dto application.UpdateProfileDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	// Super Admin lives in system club
	role := c.GetString("userRole")
	if role == domain.RoleSuperAdmin {
		clubID = "system"
	}

	user, err := h.useCases.UpdateProfile(clubID, userID.(string), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	limit := 10
	offset := 0
	// Parse basic pagination query params if needed, for now defaults
	// In real implementation: strconv.Atoi(c.Query("limit")) etc.
	search := c.Query("search")

	clubID := c.GetString("clubID")
	users, err := h.useCases.ListUsers(clubID, limit, offset, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users}) // Wrap in data envelope
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	deleteID := c.Param("id")
	if deleteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	clubID := c.GetString("clubID")
	if err := h.useCases.DeleteUser(clubID, deleteID, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Bad Request for logical limits
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) GetChildren(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	clubID := c.GetString("clubID")
	role := c.GetString("userRole")
	if role == domain.RoleSuperAdmin {
		clubID = "system"
	}

	children, err := h.useCases.ListChildren(clubID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": children})
}

func (h *UserHandler) RegisterChild(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var dto application.RegisterChildDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	role := c.GetString("userRole")
	if role == domain.RoleSuperAdmin {
		clubID = "system"
	}

	child, err := h.useCases.RegisterChild(clubID, userID.(string), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": child})
}

func (h *UserHandler) GetStats(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUserID := userID.(string)

	// "me" alias
	if id == "me" {
		id = currentUserID
	} else {
		// BOLA CHECK
		roleContext, _ := c.Get("userRole")
		role, _ := roleContext.(string)
		if id != currentUserID && role != domain.RoleAdmin && role != domain.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	clubID := c.GetString("clubID")
	// For "me" of Super Admin, use system. For others, use context.
	// But Super Admin can inspect others.
	role := c.GetString("userRole")
	if role == domain.RoleSuperAdmin && id == currentUserID {
		clubID = "system"
	}

	stats, err := h.useCases.GetStats(clubID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if stats == nil {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *UserHandler) GetWallet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUserID := userID.(string)

	// "me" alias
	if id == "me" {
		id = currentUserID
	} else {
		// BOLA CHECK
		roleContext, _ := c.Get("userRole")
		role, _ := roleContext.(string)
		if id != currentUserID && role != domain.RoleAdmin && role != domain.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	clubID := c.GetString("clubID")
	role := c.GetString("userRole")
	if role == domain.RoleSuperAdmin && id == currentUserID {
		clubID = "system"
	}

	wallet, err := h.useCases.GetWallet(clubID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if wallet == nil {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": wallet})
}

func RegisterRoutes(r *gin.RouterGroup, handler *UserHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	users := r.Group("/users")
	users.Use(authMiddleware, tenantMiddleware) // Protect these routes with Auth THEN Tenant
	{
		users.GET("/me", handler.GetProfile)
		users.PUT("/me", handler.UpdateProfile)
		users.GET("/me/children", handler.GetChildren)
		users.POST("/me/children", handler.RegisterChild)
		users.GET("", handler.ListUsers)
		users.DELETE("/:id", handler.DeleteUser)
		users.GET("/:id/stats", handler.GetStats)
		users.GET("/:id/wallet", handler.GetWallet)
	}
}
