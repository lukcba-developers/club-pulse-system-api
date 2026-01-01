package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/application"
)

type AttendanceHandler struct {
	useCases *application.AttendanceUseCases
}

func NewAttendanceHandler(useCases *application.AttendanceUseCases) *AttendanceHandler {
	return &AttendanceHandler{
		useCases: useCases,
	}
}

// GetGroupAttendance
// GET /attendance/groups/:group?date=2024-01-01
func (h *AttendanceHandler) GetGroupAttendance(c *gin.Context) {
	group := c.Param("group")
	dateStr := c.Query("date")

	// Default to today if no date provided
	date := time.Now()
	if dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format (YYYY-MM-DD)"})
			return
		}
		date = parsed
	}

	coachID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Logic: Get or Create List
	list, err := h.useCases.GetOrCreateList(group, date, coachID.(string))
	if err != nil {
		log.Printf("Error getting attendance list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// SubmitAttendance
// POST /attendance/:listID/records
func (h *AttendanceHandler) SubmitAttendance(c *gin.Context) {
	listIDStr := c.Param("listID")
	listID, err := uuid.Parse(listIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list ID"})
		return
	}

	var dto application.MarkAttendanceDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCases.MarkAttendance(listID, dto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func RegisterRoutes(r *gin.RouterGroup, handler *AttendanceHandler, authMiddleware gin.HandlerFunc) {
	g := r.Group("/attendance")
	g.Use(authMiddleware)
	{
		g.GET("/groups/:group", handler.GetGroupAttendance)
		g.POST("/:listID/records", handler.SubmitAttendance)
	}
}
