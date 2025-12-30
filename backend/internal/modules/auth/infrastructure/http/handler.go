package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/core/errors"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/application"
)

type AuthHandler struct {
	useCase *application.AuthUseCases
}

func NewAuthHandler(useCase *application.AuthUseCases) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var dto application.RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.useCase.Register(dto)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, token)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var dto application.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.useCase.Login(dto)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	type RefreshDTO struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var dto RefreshDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	token, err := h.useCase.RefreshToken(dto.RefreshToken)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	type LogoutDTO struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var dto LogoutDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	if err := h.useCase.Logout(dto.RefreshToken); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) ListSessions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sessions, err := h.useCase.ListUserSessions(userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *AuthHandler) RevokeSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID required"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.useCase.RevokeSession(sessionID, userID.(string)); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
}

// Routes registration
// Routes registration
// Routes registration
func RegisterRoutes(r *gin.RouterGroup, h *AuthHandler, authMiddleware gin.HandlerFunc) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.RefreshToken)
		authGroup.POST("/logout", h.Logout)
	}

	// Protected routes
	protected := authGroup.Group("")
	protected.Use(authMiddleware)
	{
		protected.GET("/sessions", h.ListSessions)
		protected.DELETE("/sessions/:id", h.RevokeSession)
	}
}

// Simple error handler helper (move to core/platform later if needed shared)
func handleError(c *gin.Context, err error) {
	if e, ok := err.(*errors.AppError); ok {
		c.JSON(e.Code, gin.H{"error": e.Message, "type": e.Type})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
