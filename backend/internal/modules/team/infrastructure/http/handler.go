package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

type TeamHandler struct {
	useCases            *application.TeamUseCases
	playerStatusService *application.PlayerStatusService
	travelEventService  *application.TravelEventService
}

func NewTeamHandler(useCases *application.TeamUseCases, playerStatusService *application.PlayerStatusService, travelEventService *application.TravelEventService) *TeamHandler {
	return &TeamHandler{
		useCases:            useCases,
		playerStatusService: playerStatusService,
		travelEventService:  travelEventService,
	}
}

type ScheduleMatchRequest struct {
	TrainingGroupID string    `json:"training_group_id" binding:"required"`
	Opponent        string    `json:"opponent_name" binding:"required"`
	IsHome          bool      `json:"is_home_game"`
	MeetupTime      time.Time `json:"meetup_time" binding:"required"`
	Location        string    `json:"location"`
}

func (h *TeamHandler) ScheduleMatch(c *gin.Context) {
	// RBAC: Only COACH, ADMIN or SUPER_ADMIN can schedule matches
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleCoach && role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires COACH or ADMIN role"})
		return
	}

	clubID := c.GetString("clubID")
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

	event, err := h.useCases.ScheduleMatch(clubID, groupID, req.Opponent, req.IsHome, req.MeetupTime, req.Location)
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
	clubID := c.GetString("clubID")
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

	if err := h.useCases.RespondAvailability(clubID, req.EventID, userID.(string), req.Status, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetTeamPlayersWithStatus obtiene jugadores de un equipo con su estado unificado
// GET /teams/:teamId/players
func (h *TeamHandler) GetTeamPlayersWithStatus(c *gin.Context) {
	clubID := c.GetString("clubID")

	// TODO: Obtener userIDs del equipo desde el repositorio
	// Por ahora usamos query param como workaround
	userIDsParam := c.Query("user_ids")
	if userIDsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_ids query parameter required"})
		return
	}

	userIDs := strings.Split(userIDsParam, ",")

	players, err := h.playerStatusService.GetTeamPlayersWithStatus(c.Request.Context(), clubID, userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, players)
}

// GetPlayerStatus obtiene el estado de un jugador específico
// GET /teams/players/:playerId/status
func (h *TeamHandler) GetPlayerStatus(c *gin.Context) {
	clubID := c.GetString("clubID")
	playerID := c.Param("playerId")

	status, err := h.playerStatusService.GetPlayerStatus(c.Request.Context(), clubID, playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetInhabilitadoPlayers obtiene jugadores inhabilitados de un equipo
// GET /teams/:teamId/inhabilitados
func (h *TeamHandler) GetInhabilitadoPlayers(c *gin.Context) {
	clubID := c.GetString("clubID")

	userIDsParam := c.Query("user_ids")
	if userIDsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_ids query parameter required"})
		return
	}

	userIDs := strings.Split(userIDsParam, ",")

	inhabilitados, err := h.playerStatusService.GetInhabilitadoPlayers(c.Request.Context(), clubID, userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":         len(inhabilitados),
		"inhabilitados": inhabilitados,
	})
}

// GetPlayerIssues obtiene los problemas específicos de un jugador
// GET /teams/players/:playerId/issues
func (h *TeamHandler) GetPlayerIssues(c *gin.Context) {
	clubID := c.GetString("clubID")
	playerID := c.Param("playerId")

	issues, err := h.playerStatusService.GetPlayerIssues(c.Request.Context(), clubID, playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id": playerID,
		"issues":    issues,
	})
}

// CreateTravelEvent crea un nuevo evento de viaje
// POST /events
func (h *TeamHandler) CreateTravelEvent(c *gin.Context) {
	// RBAC: Only COACH, ADMIN or SUPER_ADMIN can create travel events
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleCoach && role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires COACH or ADMIN role"})
		return
	}

	clubID := c.GetString("clubID")
	createdBy := c.GetString("user_id")

	var event domain.TravelEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.travelEventService.CreateEvent(c.Request.Context(), clubID, event.TeamID, createdBy, &event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetTravelEvent obtiene un evento específico
// GET /events/:eventId
func (h *TeamHandler) GetTravelEvent(c *gin.Context) {
	clubID := c.GetString("clubID")
	eventID := c.Param("eventId")

	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento inválido"})
		return
	}

	event, err := h.travelEventService.GetEvent(c.Request.Context(), clubID, eventUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evento no encontrado"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// GetTeamEvents obtiene todos los eventos de un equipo
// GET /teams/:teamId/events
func (h *TeamHandler) GetTeamEvents(c *gin.Context) {
	clubID := c.GetString("clubID")
	teamID := c.Param("teamId")

	teamUUID, err := uuid.Parse(teamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de equipo inválido"})
		return
	}

	events, err := h.travelEventService.GetTeamEvents(c.Request.Context(), clubID, teamUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

type RSVPRequest struct {
	Status domain.RSVPStatus `json:"status" binding:"required"`
	Notes  string            `json:"notes"`
}

// RespondToTravelEvent registra la respuesta de un usuario a un evento
// POST /events/:eventId/rsvp
func (h *TeamHandler) RespondToTravelEvent(c *gin.Context) {
	userID := c.GetString("user_id")
	eventID := c.Param("eventId")

	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento inválido"})
		return
	}

	var req RSVPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.travelEventService.RespondToEvent(c.Request.Context(), eventUUID, userID, req.Status, req.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Respuesta registrada exitosamente"})
}

// GetEventSummary obtiene el resumen de un evento con estadísticas
// GET /events/:eventId/summary
func (h *TeamHandler) GetEventSummary(c *gin.Context) {
	clubID := c.GetString("clubID")
	eventID := c.Param("eventId")

	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento inválido"})
		return
	}

	summary, err := h.travelEventService.GetEventSummary(c.Request.Context(), clubID, eventUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func RegisterRoutes(r *gin.RouterGroup, handler *TeamHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	team := r.Group("/team")
	team.Use(authMiddleware, tenantMiddleware)
	{
		team.POST("/events", handler.ScheduleMatch)
		team.POST("/availability", handler.RespondAvailability)
	}

	// Nuevos endpoints de estado de jugadores
	teams := r.Group("/teams")
	teams.Use(authMiddleware, tenantMiddleware)
	{
		teams.GET("/:teamId/players", handler.GetTeamPlayersWithStatus)
		teams.GET("/:teamId/inhabilitados", handler.GetInhabilitadoPlayers)
		teams.GET("/players/:playerId/status", handler.GetPlayerStatus)
		teams.GET("/players/:playerId/issues", handler.GetPlayerIssues)

		// Eventos de equipo
		teams.GET("/:teamId/events", handler.GetTeamEvents)
	}

	// Endpoints de eventos de viaje
	events := r.Group("/events")
	events.Use(authMiddleware, tenantMiddleware)
	{
		events.POST("", handler.CreateTravelEvent)
		events.GET("/:eventId", handler.GetTravelEvent)
		events.POST("/:eventId/rsvp", handler.RespondToTravelEvent)
		events.GET("/:eventId/summary", handler.GetEventSummary)
	}
}
