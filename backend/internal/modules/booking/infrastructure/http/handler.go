package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"
)

type BookingHandler struct {
	useCases *application.BookingUseCases
}

func NewBookingHandler(useCases *application.BookingUseCases) *BookingHandler {
	return &BookingHandler{
		useCases: useCases,
	}
}

// mapErrorToResponse converts business errors to consistent API responses with type
// Uses snake_case to match frontend error dictionary (error-messages.ts)
func mapErrorToResponse(err error) (int, gin.H) {
	msg := err.Error()
	lowerMsg := strings.ToLower(msg)

	switch {
	case strings.Contains(lowerMsg, "conflict"):
		return http.StatusConflict, gin.H{"type": "booking_conflict", "error": msg}
	case strings.Contains(lowerMsg, "medical certificate"):
		return http.StatusBadRequest, gin.H{"type": "medical_certificate_invalid", "error": msg}
	case strings.Contains(lowerMsg, "not found"):
		return http.StatusNotFound, gin.H{"type": "not_found", "error": msg}
	case strings.Contains(lowerMsg, "unauthorized"):
		return http.StatusForbidden, gin.H{"type": "cancel_unauthorized", "error": msg}
	case strings.Contains(lowerMsg, "not active"):
		return http.StatusBadRequest, gin.H{"type": "facility_inactive", "error": msg}
	default:
		return http.StatusBadRequest, gin.H{"type": "invalid_format", "error": msg}
	}
}

// Create godoc
// @Summary      Create a new booking
// @Description  Creates a new booking for a facility. Validates health certificate and availability.
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        input body application.CreateBookingDTO true "Booking Details"
// @Success      201   {object}  domain.Booking
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      409   {object}  map[string]string "Slot conflict"
// @Router       /bookings [post]
func (h *BookingHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"type": "UNAUTHORIZED", "error": "Unauthorized"})
		return
	}

	var dto application.CreateBookingDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"type": "VALIDATION_ERROR", "error": err.Error()})
		return
	}

	// Override input userID with authenticated userID to prevent spoofing
	dto.UserID = userID.(string)

	clubID := c.GetString("clubID")
	booking, err := h.useCases.CreateBooking(c.Request.Context(), clubID, dto)
	if err != nil {
		status, resp := mapErrorToResponse(err)
		c.JSON(status, resp)
		return
	}

	// Invalidate Availability Cache
	// Key format: bookings:availability:{clubID}:{facilityID}:{YYYY-MM-DD}
	// We need date from StartTime.
	dateStr := dto.StartTime.Format("2006-01-02")
	cacheKey := fmt.Sprintf("bookings:availability:%s:%s:%s", clubID, dto.FacilityID, dateStr)
	_ = platformRedis.GetClient().Del(c.Request.Context(), cacheKey)

	c.JSON(http.StatusCreated, booking)
}

// List godoc
// @Summary      List user bookings
// @Description  Returns a list of bookings for the authenticated user.
// @Tags         bookings
// @Produce      json
// @Success      200   {object}  map[string][]domain.Booking
// @Failure      401   {object}  map[string]string
// @Router       /bookings [get]
func (h *BookingHandler) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// By default list own bookings
	// Admin logic could be added here later
	clubID := c.GetString("clubID")
	bookings, err := h.useCases.ListBookings(c.Request.Context(), clubID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

// ListAll godoc
// @Summary List all bookings (Admin only)
// @Description Lists all bookings for the club with optional date range
// @Tags bookings
// @Produce json
// @Router /bookings/all [get]
func (h *BookingHandler) ListAll(c *gin.Context) {
	// 1. Check Admin Role
	role, exists := c.Get("userRole")
	if !exists || (role != "ADMIN" && role != "SUPER_ADMIN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	clubID := c.GetString("clubID")
	facilityID := c.Query("facility_id")

	// Parse Dates
	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse("2006-01-02", fromStr); err == nil {
			from = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse("2006-01-02", toStr); err == nil {
			// Set to end of day if only date provided?
			// For simplicity, let's assume 'to' includes the whole day if we add 24h or if caller sends time.
			// Standard: if just date, assume start of day.
			// Repository logic uses strict <=, so we might want to Add(24h) if strictly date.
			// Let's rely on caller sending specific times if they want precision, or just date match.
			to = &t
		}
	}

	bookings, err := h.useCases.ListClubBookings(c.Request.Context(), clubID, facilityID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

// Cancel godoc
// @Summary      Cancel a booking
// @Description  Cancels an existing booking by its ID.
// @Tags         bookings
// @Param        id   path      string  true  "Booking ID"
// @Success      200   {object}  map[string]string "message: booking cancelled"
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /bookings/{id} [delete]
func (h *BookingHandler) Cancel(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	bookingID := c.Param("id")
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking id required"})
		return
	}

	clubID := c.GetString("clubID")
	if err := h.useCases.CancelBooking(c.Request.Context(), clubID, bookingID, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}

// GetAvailability godoc
// @Summary      Get facility availability
// @Description  Check available slots for a specific facility and date.
// @Tags         bookings
// @Produce      json
// @Param        facility_id  query     string  true  "Facility ID"
// @Param        date         query     string  true  "Date (YYYY-MM-DD)"
// @Success      200   {object}  map[string][]interface{}
// @Failure      400   {object}  map[string]string
// @Router       /bookings/availability [get]
func (h *BookingHandler) GetAvailability(c *gin.Context) {
	// 1. Parse params
	facilityID := c.Query("facility_id")
	dateStr := c.Query("date")

	if facilityID == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "facility_id and date are required"})
		return
	}

	clubID := c.GetString("clubID")
	ctx := c.Request.Context()

	// 2. Cache Check
	cacheKey := fmt.Sprintf("bookings:availability:%s:%s:%s", clubID, facilityID, dateStr)
	cached, err := platformRedis.GetClient().Get(ctx, cacheKey)
	if err == nil && cached != "" {
		c.Header("X-Cache", "HIT")
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, cached)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (expected YYYY-MM-DD)"})
		return
	}

	// 3. Call UseCase
	availability, err := h.useCases.GetAvailability(c.Request.Context(), clubID, facilityID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Set Cache (Short TTL: 1 minute)
	resp := gin.H{"data": availability}
	data, _ := json.Marshal(resp)
	_ = platformRedis.GetClient().Set(ctx, cacheKey, string(data), 1*time.Minute)

	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, resp)
}

// CreateRecurringRule godoc
// @Summary      Create a recurring booking rule
// @Description  Admin only. Sets up a pattern for automatic bookings.
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        input body application.CreateRecurringRuleDTO true "Rule Details"
// @Success      201   {object}  domain.RecurringRule
// @Failure      403   {object}  map[string]string "Requires ADMIN role"
// @Router       /bookings/recurring [post]
func (h *BookingHandler) CreateRecurringRule(c *gin.Context) {
	// RBAC: Only ADMIN or SUPER_ADMIN can create recurring rules
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	var dto application.CreateRecurringRuleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	rule, err := h.useCases.CreateRecurringRule(c.Request.Context(), clubID, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

// GenerateBookings godoc
// @Summary      Materialize recurring bookings
// @Description  Admin only. Forces generation of bookings from active recurring rules.
// @Tags         bookings
// @Success      200   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Router       /bookings/generate [post]
func (h *BookingHandler) GenerateBookings(c *gin.Context) {
	// RBAC: Only ADMIN or SUPER_ADMIN can generate bookings
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	clubID := c.GetString("clubID")

	if err := h.useCases.GenerateBookingsFromRules(c.Request.Context(), clubID, 4); err != nil { // Default 4 weeks
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bookings generated successfully"})
}

// ListRecurringRules godoc
// @Summary      List recurring rules
// @Description  Admin only. Lists all active recurring booking rules.
// @Tags         bookings
// @Produce      json
// @Success      200   {object}  map[string][]domain.RecurringRule
// @Failure      403   {object}  map[string]string "Requires ADMIN role"
// @Router       /bookings/recurring [get]
func (h *BookingHandler) ListRecurringRules(c *gin.Context) {
	// RBAC: Only ADMIN or SUPER_ADMIN can list recurring rules
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	clubID := c.GetString("clubID")
	rules, err := h.useCases.ListRecurringRules(c.Request.Context(), clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rules})
}

// JoinWaitlist godoc
// @Summary      Join a waitlist
// @Description  Adds the user to the waitlist for a specific resource and date.
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        input body application.JoinWaitlistDTO true "Waitlist Details"
// @Success      201   {object}  domain.Waitlist
// @Failure      401   {object}  map[string]string
// @Router       /bookings/waitlist [post]
func (h *BookingHandler) JoinWaitlist(c *gin.Context) {
	var dto application.JoinWaitlistDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	// Enforce user_id from token if we want strict security, similar to CreateBooking.
	userID, exists := c.Get("userID")
	if exists {
		dto.UserID = userID.(string)
	}

	entry, err := h.useCases.JoinWaitlist(c.Request.Context(), clubID, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func RegisterRoutes(r *gin.RouterGroup, handler *BookingHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	bookings := r.Group("/bookings")
	bookings.Use(authMiddleware, tenantMiddleware)
	{
		bookings.POST("", handler.Create)
		bookings.GET("", handler.List)
		bookings.GET("/all", handler.ListAll)
		bookings.GET("/availability", handler.GetAvailability)
		bookings.DELETE("/:id", handler.Cancel)
		bookings.GET("/recurring", handler.ListRecurringRules)
		bookings.POST("/recurring", handler.CreateRecurringRule)
		bookings.POST("/generate", handler.GenerateBookings)
		bookings.POST("/waitlist", handler.JoinWaitlist)
	}
}
