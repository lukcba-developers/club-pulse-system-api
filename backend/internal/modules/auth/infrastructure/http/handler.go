package http

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/core/errors"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/middleware"
)

// isSecureCookie returns true in production (GIN_MODE=release) for Secure cookies.
func isSecureCookie() bool {
	return os.Getenv("GIN_MODE") == "release"
}

// setAuthCookie configures a secure HttpOnly cookie with SameSite protection.
func setAuthCookie(c *gin.Context, name, value string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, value, maxAge, "/", "", isSecureCookie(), true)
}

type AuthHandler struct {
	useCase *application.AuthUseCases
}

func NewAuthHandler(useCase *application.AuthUseCases) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account for the specified club.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body application.RegisterDTO true "Registration Details"
// @Success      201   {object}  domain.Token
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string "User already exists"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var dto application.RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	clubID := c.GetString("clubID")
	token, err := h.useCase.Register(c.Request.Context(), dto, clubID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, token)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return HttpOnly cookies
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body application.LoginDTO true "Login Credentials"
// @Success      200   {object}  map[string]string "message: Login successful"
// @Failure      400   {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var dto application.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	clubID := c.GetString("clubID")
	token, err := h.useCase.Login(c.Request.Context(), dto, clubID)
	if err != nil {
		handleError(c, err)
		return
	}

	// Set secure HttpOnly cookies with SameSite protection
	setAuthCookie(c, "access_token", token.AccessToken, 86400) // 24 hours
	if token.RefreshToken != "" {
		setAuthCookie(c, "refresh_token", token.RefreshToken, 86400*7) // 7 days
	}

	// Set CSRF token cookie for Double Submit Cookie protection
	_ = middleware.SetCSRFCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"access_token": token.AccessToken,
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Exchanges a refresh token for a new access token and refresh token (rotation).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body object{refresh_token=string} true "Refresh Token"
// @Success      200   {object}  domain.Token
// @Failure      401   {object}  map[string]string "Invalid or expired token"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	type RefreshDTO struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var dto RefreshDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	clubID := c.GetString("clubID")
	token, err := h.useCase.RefreshToken(c.Request.Context(), dto.RefreshToken, clubID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}

// Logout godoc
// @Summary      Logout user
// @Description  Revokes the provided refresh token and clears auth cookies.
// @Tags         auth
// @Accept       json
// @Param        input body object{refresh_token=string} true "Refresh Token"
// @Success      204   "No Content"
// @Failure      400   {object}  map[string]string
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	type LogoutDTO struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var dto LogoutDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	clubID := c.GetString("clubID")
	if err := h.useCase.Logout(c.Request.Context(), dto.RefreshToken, clubID); err != nil {
		handleError(c, err)
		return
	}

	// Clear cookies on logout (including CSRF token)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("access_token", "", -1, "/", "", isSecureCookie(), true)
	c.SetCookie("refresh_token", "", -1, "/", "", isSecureCookie(), true)
	middleware.ClearCSRFCookie(c)

	c.Status(http.StatusNoContent)
}

// ListSessions godoc
// @Summary      List user sessions
// @Description  Returns all active sessions (refresh tokens) for the authenticated user.
// @Tags         auth
// @Produce      json
// @Success      200   {array}   domain.RefreshToken
// @Failure      401   {object}  map[string]string
// @Router       /auth/sessions [get]
func (h *AuthHandler) ListSessions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sessions, err := h.useCase.ListUserSessions(c.Request.Context(), userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// RevokeSession godoc
// @Summary      Revoke a session
// @Description  Terminates a specific session by its ID.
// @Tags         auth
// @Param        id   path      string  true  "Session ID"
// @Success      200   {object}  map[string]string "message: Session revoked"
// @Failure      401   {object}  map[string]string
// @Router       /auth/sessions/{id} [delete]
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

	if err := h.useCase.RevokeSession(c.Request.Context(), sessionID, userID.(string)); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
}

// GoogleLogin godoc
// @Summary      Google OAuth Login
// @Description  Authenticates user via Google OAuth2 code.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body object{code=string} true "Google Authorization Code"
// @Success      200   {object}  map[string]string "message: Google login successful"
// @Failure      401   {object}  map[string]string
// @Router       /auth/google [post]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	type GoogleLoginDTO struct {
		Code string `json:"code" binding:"required"`
	}
	var dto GoogleLoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code is required"})
		return
	}

	clubID := c.GetString("clubID")
	token, err := h.useCase.GoogleLogin(c.Request.Context(), dto.Code, clubID)
	if err != nil {
		handleError(c, err)
		return
	}

	// Set secure HttpOnly cookies with SameSite protection
	setAuthCookie(c, "access_token", token.AccessToken, 86400)
	if token.RefreshToken != "" {
		setAuthCookie(c, "refresh_token", token.RefreshToken, 86400*7)
	}

	// Set CSRF token cookie for Double Submit Cookie protection
	_ = middleware.SetCSRFCookie(c)

	c.JSON(http.StatusOK, gin.H{"message": "Google login successful"})
}

// RegisterRoutes sets up authentication endpoints.
func RegisterRoutes(r *gin.RouterGroup, h *AuthHandler, authMiddleware gin.HandlerFunc, tenantMiddleware gin.HandlerFunc) {
	authGroup := r.Group("/auth")
	authGroup.Use(tenantMiddleware)
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.RefreshToken)
		authGroup.POST("/logout", h.Logout)
		authGroup.POST("/google", h.GoogleLogin)
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
	// TEMPORARY DEBUG: Log the actual error for CI debugging
	println("[AUTH DEBUG] Internal error:", err.Error())
	// Exposse internal error details for debugging (TODO: Use proper logger and hide in prod)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "debug": err.Error()})
}
