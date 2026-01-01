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

	user, err := h.useCases.GetProfile(userID.(string))
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

	user, err := h.useCases.UpdateProfile(userID.(string), dto)
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

	users, err := h.useCases.ListUsers(limit, offset, search)
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

	if err := h.useCases.DeleteUser(deleteID, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Bad Request for logical limits
		return
	}

	c.Status(http.StatusNoContent)
}

func RegisterRoutes(r *gin.RouterGroup, handler *UserHandler, authMiddleware gin.HandlerFunc) {
	users := r.Group("/users")
	users.Use(authMiddleware) // Protect these routes
	{
		users.GET("/me", handler.GetProfile)
		users.PUT("/me", handler.UpdateProfile)
		users.GET("", handler.ListUsers)
		users.DELETE("/:id", handler.DeleteUser)
	}
}
