package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
)

type BookingHandler struct {
	useCases *application.BookingUseCases
}

func NewBookingHandler(useCases *application.BookingUseCases) *BookingHandler {
	return &BookingHandler{
		useCases: useCases,
	}
}

func RegisterRoutes(r *gin.RouterGroup, handler *BookingHandler, authMiddleware gin.HandlerFunc) {
	bookings := r.Group("/bookings")
	bookings.Use(authMiddleware)
	{
		bookings.POST("", handler.Create)
		bookings.GET("", handler.List)
		bookings.GET("/availability", handler.GetAvailability)
		bookings.DELETE("/:id", handler.Cancel)
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

	booking, err := h.useCases.CreateBooking(dto)
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
	bookings, err := h.useCases.ListBookings(userID.(string))
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

	if err := h.useCases.CancelBooking(bookingID, userID.(string)); err != nil {
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

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (expected YYYY-MM-DD)"})
		return
	}

	// 2. Call UseCase
	availability, err := h.useCases.GetAvailability(facilityID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": availability})
}
