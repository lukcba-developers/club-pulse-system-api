package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
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

func (h *BookingHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var dto application.CreateBookingDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Override input userID with authenticated userID to prevent spoofing
	dto.UserID = userID.(string)

	clubID := c.GetString("clubID")
	booking, err := h.useCases.CreateBooking(clubID, dto)
	if err != nil {
		// Differentiate connection vs conflict vs validation errors properly in a real app
		// For MVP, if error message contains "conflict", return 409
		if err.Error() == "booking time conflict: facility is already booked for this requested time" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (h *BookingHandler) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// By default list own bookings
	// Admin logic could be added here later
	clubID := c.GetString("clubID")
	bookings, err := h.useCases.ListBookings(clubID, userID.(string))
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

	bookings, err := h.useCases.ListClubBookings(clubID, facilityID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

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
	if err := h.useCases.CancelBooking(clubID, bookingID, userID.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}

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
	availability, err := h.useCases.GetAvailability(clubID, facilityID, date)
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

func (h *BookingHandler) CreateRecurringRule(c *gin.Context) {
	// Simple RBAC check (MVP: hardcode checks or assume middleware handles basic auth)
	// In production: check if user is admin/manager

	var dto application.CreateRecurringRuleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	rule, err := h.useCases.CreateRecurringRule(clubID, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (h *BookingHandler) GenerateBookings(c *gin.Context) {
	// Admin only
	// Check query param for weeks
	// weeks := c.Query("weeks") ... parse int
	clubID := c.GetString("clubID")

	if err := h.useCases.GenerateBookingsFromRules(clubID, 4); err != nil { // Default 4 weeks
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bookings generated successfully"})
}

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

	entry, err := h.useCases.JoinWaitlist(clubID, dto)
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
		bookings.POST("/recurring", handler.CreateRecurringRule)
		bookings.POST("/generate", handler.GenerateBookings)
		bookings.POST("/waitlist", handler.JoinWaitlist)
	}
}
