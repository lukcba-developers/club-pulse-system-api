package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	clubApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

type ChampionshipHandler struct {
	useCases         *application.ChampionshipUseCases
	volunteerService application.VolunteerServiceInterface
	clubUseCases     *clubApp.ClubUseCases
}

func NewChampionshipHandler(useCases *application.ChampionshipUseCases, volunteerService application.VolunteerServiceInterface, clubUseCases *clubApp.ClubUseCases) *ChampionshipHandler {
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
		group.POST("/stages/:id/knockout", h.GenerateKnockoutBracket)

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
	club, err := h.clubUseCases.GetClubBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	tournaments, err := h.useCases.ListTournaments(c.Request.Context(), club.ID)
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
	club, err := h.clubUseCases.GetClubBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	tournament, err := h.useCases.GetTournament(c.Request.Context(), club.ID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	c.JSON(http.StatusOK, tournament)
}

// ... (Existing methods remain unchanged, appending new methods)

func (h *ChampionshipHandler) AssignVolunteer(c *gin.Context) {
	// RBAC: Only ADMIN or SUPER_ADMIN can assign volunteers
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

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

	// SECURITY FIX (VUL-005): Use clubID from context, get assignerID from token
	clubID := c.GetString("clubID")
	if clubID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "club context required"})
		return
	}
	assignerID, _ := c.Get("userID")

	if err := h.volunteerService.AssignVolunteer(c.Request.Context(), clubID, matchID, input.UserID, input.Role, assignerID.(string), input.Notes); err != nil {
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

	// SECURITY FIX (VUL-005): Use clubID from context
	clubID := c.GetString("clubID")
	if clubID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "club context required"})
		return
	}

	summary, err := h.volunteerService.GetVolunteerSummary(c.Request.Context(), clubID, matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *ChampionshipHandler) RemoveVolunteer(c *gin.Context) {
	// RBAC: Only ADMIN can remove volunteers
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	assignmentIDStr := c.Param("id")
	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment ID"})
		return
	}

	// SECURITY FIX (VUL-005): Use clubID from context
	clubID := c.GetString("clubID")
	if clubID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "club context required"})
		return
	}

	if err := h.volunteerService.RemoveVolunteer(c.Request.Context(), clubID, assignmentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
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

	tournaments, err := h.useCases.ListTournaments(c.Request.Context(), clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tournaments)
}

func (h *ChampionshipHandler) GetStandings(c *gin.Context) {
	clubID := c.GetString("clubID")
	groupID := c.Param("id")
	standings, err := h.useCases.GetStandings(c.Request.Context(), clubID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standings)
}

func (h *ChampionshipHandler) GetMatchesByGroup(c *gin.Context) {
	clubID := c.GetString("clubID")
	groupID := c.Param("id")
	matches, err := h.useCases.GetMatchesByGroup(c.Request.Context(), clubID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

func (h *ChampionshipHandler) ScheduleMatch(c *gin.Context) {
	// SECURITY FIX (VUL-003): Only ADMIN can schedule matches
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	var input application.ScheduleMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ClubID = c.GetString("clubID")

	if err := h.useCases.ScheduleMatch(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "match scheduled"})
}

func (h *ChampionshipHandler) UpdateMatchResult(c *gin.Context) {
	// SECURITY FIX (VUL-003): Only ADMIN/STAFF can update match results
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin && role != "STAFF") {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN or STAFF role"})
		return
	}

	var input application.UpdateMatchResultInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ClubID = c.GetString("clubID")

	if err := h.useCases.UpdateMatchResult(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *ChampionshipHandler) CreateTournament(c *gin.Context) {
	// RBAC: Only ADMIN or SUPER_ADMIN can create tournaments
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	var input application.CreateTournamentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ClubID = c.GetString("clubID")

	tournament, err := h.useCases.CreateTournament(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tournament)
}

func (h *ChampionshipHandler) AddStage(c *gin.Context) {
	// RBAC: Only ADMIN or SUPER_ADMIN can add stages
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

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

	stage, err := h.useCases.AddStage(c.Request.Context(), input.TournamentID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, stage)
}

func (h *ChampionshipHandler) AddGroup(c *gin.Context) {
	// RBAC: Only ADMIN or SUPER_ADMIN can add groups
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

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

	group, err := h.useCases.AddGroup(c.Request.Context(), input.StageID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func (h *ChampionshipHandler) RegisterTeam(c *gin.Context) {
	// SECURITY FIX (VUL-003): Only ADMIN can register teams
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	clubID := c.GetString("clubID")
	var input application.RegisterTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	groupID := c.Param("id")
	standing, err := h.useCases.RegisterTeam(c.Request.Context(), clubID, groupID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, standing)
}

func (h *ChampionshipHandler) GenerateFixture(c *gin.Context) {
	// SECURITY FIX (VUL-003): Only ADMIN can generate fixtures
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	clubID := c.GetString("clubID")
	groupID := c.Param("id")
	matches, err := h.useCases.GenerateGroupFixture(c.Request.Context(), clubID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, matches)
}

// GenerateKnockoutBracket godoc
// @Summary      Generate knockout bracket
// @Description  Generates elimination bracket matches for a KNOCKOUT stage. Teams must be a power of 2 (2, 4, 8, 16, 32). Pairs teams by seeding: #1 vs #N, #2 vs #(N-1), etc.
// @Tags         championships
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Stage ID (must be KNOCKOUT type)"
// @Param        input body object{seed_order=[]string} true "Seed Order (Team IDs ordered by ranking/seed)"
// @Success      201   {object}  object{matches=[]domain.TournamentMatch,count=int}
// @Failure      400   {object}  map[string]string "Invalid input or stage type"
// @Failure      403   {object}  map[string]string "Requires ADMIN role"
// @Failure      500   {object}  map[string]string "Internal error"
// @Router       /championships/stages/{id}/knockout [post]
func (h *ChampionshipHandler) GenerateKnockoutBracket(c *gin.Context) {
	// RBAC: Only ADMIN can generate knockout brackets
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	var input application.GenerateKnockoutBracketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.StageID = c.Param("id")
	if input.ClubID == "" {
		input.ClubID = c.GetString("clubID")
	}

	matches, err := h.useCases.GenerateKnockoutBracket(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"matches": matches, "count": len(matches)})
}
