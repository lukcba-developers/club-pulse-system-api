package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	clubApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
)

type ChampionshipHandler struct {
	useCases         *application.ChampionshipUseCases
	volunteerService *application.VolunteerService
	clubUseCases     *clubApp.ClubUseCases
}

func NewChampionshipHandler(useCases *application.ChampionshipUseCases, volunteerService *application.VolunteerService, clubUseCases *clubApp.ClubUseCases) *ChampionshipHandler {
	return &ChampionshipHandler{
		useCases:         useCases,
		volunteerService: volunteerService,
		clubUseCases:     clubUseCases,
	}
}

func (h *ChampionshipHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc, tenantMiddleware gin.HandlerFunc) {
	group := r.Group("/championships")
	group.Use(authMiddleware, tenantMiddleware)
	// protected routes
	{
		group.GET("/", h.ListTournaments)
		group.POST("/", h.CreateTournament)
		group.POST("/:id/stages", h.AddStage)
		group.POST("/stages/:id/groups", h.AddGroup)
		group.POST("/groups/:id/teams", h.RegisterTeam)
		group.POST("/groups/:id/fixture", h.GenerateFixture)
		group.GET("/groups/:id/matches", h.GetMatchesByGroup)
		group.GET("/groups/:id/standings", h.GetStandings)
		group.POST("/matches/result", h.UpdateMatchResult)
		group.POST("/matches/schedule", h.ScheduleMatch)

		group.POST("/matches/:id/volunteers", h.AssignVolunteer)
		group.GET("/matches/:id/volunteers", h.GetMatchVolunteers)
		group.DELETE("/volunteers/:id", h.RemoveVolunteer)
	}

	// Public Routes
	public := r.Group("/public/clubs/:slug/championships")
	{
		public.GET("/", h.GetPublicTournaments)
		public.GET("/:id", h.GetPublicTournament)
		public.GET("/groups/:id/standings", h.GetStandings)    // Reusing existing if possible, or wrap
		public.GET("/groups/:id/fixture", h.GetMatchesByGroup) // Reusing existing
	}
}

func (h *ChampionshipHandler) GetPublicTournaments(c *gin.Context) {
	slug := c.Param("slug")
	club, err := h.clubUseCases.GetClubBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	tournaments, err := h.useCases.ListTournaments(club.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tournaments)
}

func (h *ChampionshipHandler) GetPublicTournament(c *gin.Context) {
	id := c.Param("id")
	// Note: We should verify it belongs to the club from slug, but for MVP ID checks are ok or we assume ID is unique.
	// To be strict: Fetch tournament -> check clubID matches slug's clubID.

	// To be strict: Fetch tournament -> check clubID matches slug's clubID.
	// We might need to fetch the club from the slug (already done above in GetPublicTournaments? No, this is GetPublicTournament).
	// But we don't have slug here? We assume the path is /public/clubs/:slug/championships/:id ?
	// Yes, checking RegisterRoutes: public.GET("/:id", h.GetPublicTournament) is inside the group?
	// No: public := r.Group("/public/clubs/:slug/championships")

	// So we can get slug.
	slug := c.Param("slug")
	club, err := h.clubUseCases.GetClubBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	tournament, err := h.useCases.GetTournament(club.ID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	c.JSON(http.StatusOK, tournament)
}

// ... (Existing methods remain unchanged, appending new methods)

func (h *ChampionshipHandler) AssignVolunteer(c *gin.Context) {
	matchIDStr := c.Param("id")
	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match ID"})
		return
	}

	var input struct {
		UserID string               `json:"user_id" binding:"required"`
		Role   domain.VolunteerRole `json:"role" binding:"required"`
		Notes  string               `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get ClubID and AssignerID from context
	clubID := c.Query("club_id")
	if clubID == "" {
		clubID = "default-club-id" // Placeholder if not in context/query
	}
	assignerID := "admin-id" // Placeholder or from auth token

	if err := h.volunteerService.AssignVolunteer(c.Request.Context(), clubID, matchID, input.UserID, input.Role, assignerID, input.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "volunteer assigned"})
}

func (h *ChampionshipHandler) GetMatchVolunteers(c *gin.Context) {
	matchIDStr := c.Param("id")
	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match ID"})
		return
	}

	clubID := c.Query("club_id")
	if clubID == "" {
		clubID = "default-club-id"
	}

	summary, err := h.volunteerService.GetVolunteerSummary(c.Request.Context(), clubID, matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *ChampionshipHandler) RemoveVolunteer(c *gin.Context) {
	assignmentIDStr := c.Param("id")
	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment ID"})
		return
	}

	clubID := c.Query("club_id")
	if clubID == "" {
		clubID = "default-club-id"
	}

	if err := h.volunteerService.RemoveVolunteer(c.Request.Context(), clubID, assignmentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "volunteer removed"})
}

func (h *ChampionshipHandler) ListTournaments(c *gin.Context) {
	clubID := c.GetString("clubID")
	// Legacy fallback no longer needed if middleware is enforced, keeping safe check?
	// Middleware guarantees key existence if applied.
	if clubID == "" {
		// Fallback to Query but strictly warn/deprecate or deny?
		// If middleware is used, c.Get("clubID") is set.
		// If testing/public, logic might differ.
		// For protected routes, this IS set.
		clubID = c.Query("club_id")
	}

	if clubID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "club_id is required"})
		return
	}

	tournaments, err := h.useCases.ListTournaments(clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tournaments)
}

func (h *ChampionshipHandler) GetStandings(c *gin.Context) {
	clubID := c.GetString("clubID")
	groupID := c.Param("id")
	standings, err := h.useCases.GetStandings(clubID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standings)
}

func (h *ChampionshipHandler) GetMatchesByGroup(c *gin.Context) {
	clubID := c.GetString("clubID")
	groupID := c.Param("id")
	matches, err := h.useCases.GetMatchesByGroup(clubID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

func (h *ChampionshipHandler) ScheduleMatch(c *gin.Context) {
	var input application.ScheduleMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCases.ScheduleMatch(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "match scheduled"})
}

func (h *ChampionshipHandler) UpdateMatchResult(c *gin.Context) {
	var input application.UpdateMatchResultInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCases.UpdateMatchResult(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *ChampionshipHandler) CreateTournament(c *gin.Context) {
	var input application.CreateTournamentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournament, err := h.useCases.CreateTournament(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tournament)
}

func (h *ChampionshipHandler) AddStage(c *gin.Context) {
	var input application.AddStageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.TournamentID = c.Param("id")
	// Ensure ClubID is set (either from input or context)
	if input.ClubID == "" {
		input.ClubID = c.GetString("clubID")
	}

	stage, err := h.useCases.AddStage(input.TournamentID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, stage)
}

func (h *ChampionshipHandler) AddGroup(c *gin.Context) {
	var input application.AddGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.StageID = c.Param("id")
	// Ensure ClubID is set
	if input.ClubID == "" {
		input.ClubID = c.GetString("clubID")
	}

	group, err := h.useCases.AddGroup(input.StageID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func (h *ChampionshipHandler) RegisterTeam(c *gin.Context) {
	clubID := c.GetString("clubID")
	var input application.RegisterTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	groupID := c.Param("id")
	standing, err := h.useCases.RegisterTeam(clubID, groupID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, standing)
}

func (h *ChampionshipHandler) GenerateFixture(c *gin.Context) {
	clubID := c.GetString("clubID")
	groupID := c.Param("id")
	matches, err := h.useCases.GenerateGroupFixture(clubID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, matches)
}
