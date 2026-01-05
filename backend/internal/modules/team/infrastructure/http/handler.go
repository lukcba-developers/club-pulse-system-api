package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/domain"
)

type TeamHandler struct {
	useCases *application.TeamUseCases
}

func NewTeamHandler(useCases *application.TeamUseCases) *TeamHandler {
	return &TeamHandler{useCases: useCases}
}

type ScheduleMatchRequest struct {
	TrainingGroupID string    `json:"training_group_id" binding:"required"`
	Opponent        string    `json:"opponent_name" binding:"required"`
	IsHome          bool      `json:"is_home_game"`
	MeetupTime      time.Time `json:"meetup_time" binding:"required"`
	Location        string    `json:"location"`
}

func (h *TeamHandler) ScheduleMatch(c *gin.Context) {
	var req ScheduleMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupID, err := uuid.Parse(req.TrainingGroupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Training Group ID"})
		return
	}

	event, err := h.useCases.ScheduleMatch(groupID, req.Opponent, req.IsHome, req.MeetupTime, req.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

type AvailabilityRequest struct {
	EventID string                          `json:"event_id" binding:"required"`
	Status  domain.PlayerAvailabilityStatus `json:"status" binding:"required"`
	Reason  string                          `json:"reason"`
}

func (h *TeamHandler) RespondAvailability(c *gin.Context) {
	var req AvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.useCases.RespondAvailability(req.EventID, userID.(string), req.Status, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func RegisterRoutes(r *gin.RouterGroup, handler *TeamHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	team := r.Group("/team")
	team.Use(authMiddleware, tenantMiddleware)
	{
		team.POST("/events", handler.ScheduleMatch)
		team.POST("/availability", handler.RespondAvailability)
	}
}
