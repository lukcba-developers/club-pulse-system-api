package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/middleware"
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
	// Role check handled by middleware

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
	// Role check handled by middleware

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

type UpdateEmergencyInfoRequest struct {
	ContactName       string `json:"contact_name"`
	ContactPhone      string `json:"contact_phone"`
	InsuranceProvider string `json:"insurance_provider"`
	InsuranceNumber   string `json:"insurance_number"`
}

func (h *UserHandler) UpdateEmergencyInfo(c *gin.Context) {
	var req UpdateEmergencyInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	clubID := c.GetString("clubID")

	if err := h.useCases.UpdateEmergencyInfo(clubID, userID.(string), req.ContactName, req.ContactPhone, req.InsuranceProvider, req.InsuranceNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type LogIncidentRequest struct {
	InjuredUserID string `json:"injured_user_id"`
	Description   string `json:"description" binding:"required"`
	Witnesses     string `json:"witnesses"`
	ActionTaken   string `json:"action_taken"`
}

func (h *UserHandler) LogIncident(c *gin.Context) {
	var req LogIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Reporter is the logged in user
	reporterID, _ := c.Get("userID")
	clubID := c.GetString("clubID")

	incident, err := h.useCases.LogIncident(clubID, req.InjuredUserID, req.Description, req.ActionTaken, req.Witnesses, reporterID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

func RegisterRoutes(r *gin.RouterGroup, handler *UserHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	users := r.Group("/users")
	users.Use(authMiddleware, tenantMiddleware) // Protect these routes with Auth THEN Tenant
	{
		users.GET("/me", handler.GetProfile)
		users.PUT("/me", handler.UpdateProfile)
		users.GET("/me/children", handler.GetChildren)
		users.POST("/me/children", handler.RegisterChild)
		users.GET("/:id/stats", handler.GetStats)
		users.GET("/:id/wallet", handler.GetWallet)

		// Operational
		users.PUT("/me/emergency", handler.UpdateEmergencyInfo)
		users.POST("/incidents", handler.LogIncident)

		// Family Groups
		users.POST("/family-groups", handler.CreateFamilyGroup)
		users.GET("/family-groups/me", handler.GetMyFamilyGroup)
		users.POST("/family-groups/:id/members", handler.AddFamilyMember)

		// GDPR Rights Endpoints
		users.GET("/me/data-export", handler.ExportMyData)       // Article 20 - Portability
		users.DELETE("/me/gdpr-erasure", handler.RequestErasure) // Article 17 - Right to erasure

		// Admin Routes (Protected by RBAC Middleware)
		adminOnly := users.Group("")
		adminOnly.Use(middleware.RequireRole(domain.RoleAdmin, domain.RoleSuperAdmin))
		{
			adminOnly.DELETE("/:id", handler.DeleteUser)
		}

		// Staff/Admin Routes
		staffOnly := users.Group("")
		staffOnly.Use(middleware.RequireRole(domain.RoleAdmin, domain.RoleSuperAdmin, "STAFF"))
		{
			staffOnly.GET("", handler.ListUsers)
		}
	}
}

// --- Family Group Handlers ---

type CreateFamilyGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *UserHandler) CreateFamilyGroup(c *gin.Context) {
	var req CreateFamilyGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	clubID := c.GetString("clubID")

	group, err := h.useCases.CreateFamilyGroup(clubID, userID.(string), req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": group})
}

func (h *UserHandler) GetMyFamilyGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	clubID := c.GetString("clubID")

	group, err := h.useCases.GetMyFamilyGroup(clubID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": group})
}

type AddFamilyMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

func (h *UserHandler) AddFamilyMember(c *gin.Context) {
	groupIDStr := c.Param("id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group ID required"})
		return
	}

	// SECURITY: Validate authenticated user
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req AddFamilyMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")

	// Parse group ID
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// SECURITY FIX (VUL-002): Validate ownership - only HeadUserID can add members
	if err := h.useCases.AddFamilyMemberSecure(clubID, groupID, req.UserID, currentUserID.(string)); err != nil {
		if err.Error() == "only the family head can add members" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) RegisterDependentPublic(c *gin.Context) {
	clubID := c.Query("club_id")
	if clubID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Club ID is required"})
		return
	}

	var dto application.RegisterDependentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.useCases.RegisterDependent(clubID, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func RegisterPublicRoutes(r *gin.RouterGroup, handler *UserHandler) {
	users := r.Group("/users/public")
	{
		users.POST("/register-dependent", handler.RegisterDependentPublic)
	}
}

// --- GDPR Compliance Handlers ---

// ExportMyData implements GDPR Article 20 - Right to data portability
// GET /users/me/data-export
func (h *UserHandler) ExportMyData(c *gin.Context) {
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

	exportData, err := h.useCases.ExportUserData(clubID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for file download
	c.Header("Content-Disposition", "attachment; filename=my_data_export.json")
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, exportData)
}

// RequestErasure implements GDPR Article 17 - Right to erasure (Right to be forgotten)
// DELETE /users/me/gdpr-erasure
// Note: This anonymizes the user's own data. For admin deletion, use DELETE /users/:id
func (h *UserHandler) RequestErasure(c *gin.Context) {
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

	// Users can request erasure of their own data
	// This will anonymize their data rather than delete it, preserving referential integrity
	if err := h.useCases.DeleteUserGDPR(clubID, userID.(string), "SELF_REQUEST"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your data has been anonymized. Your account is now deactivated.",
		"details": "Personal identifying information has been removed in compliance with GDPR Article 17.",
	})
}
