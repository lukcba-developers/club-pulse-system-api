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
	clubID := c.GetString("clubID")
	list, err := h.useCases.GetOrCreateList(clubID, group, date, coachID.(string))
	if err != nil {
		log.Printf("Error getting attendance list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *AttendanceHandler) GetTrainingGroupAttendance(c *gin.Context) {
	idStr := c.Param("id")
	groupID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	dateStr := c.Query("date")
	groupName := c.DefaultQuery("group_name", "Training Group")
	category := c.Query("category")

	date := time.Now()
	if dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = parsed
		}
	}

	coachID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	clubID := c.GetString("clubID")
	list, err := h.useCases.GetOrCreateListByTrainingGroup(clubID, groupID, groupName, category, date, coachID.(string))
	if err != nil {
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

	clubID := c.GetString("clubID")
	if err := h.useCases.MarkAttendance(clubID, listID, dto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func RegisterRoutes(r *gin.RouterGroup, handler *AttendanceHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	g := r.Group("/attendance")
	g.Use(authMiddleware, tenantMiddleware)
	{
		g.GET("/groups/:group", handler.GetGroupAttendance)
		g.GET("/training-groups/:id", handler.GetTrainingGroupAttendance)
		g.POST("/:listID/records", handler.SubmitAttendance)
	}
}
