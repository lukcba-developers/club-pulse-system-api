package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/application"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

type DisciplineHandler struct {
	useCases *application.DisciplineUseCases
}

func NewDisciplineHandler(useCases *application.DisciplineUseCases) *DisciplineHandler {
	return &DisciplineHandler{useCases: useCases}
}

func (h *DisciplineHandler) ListDisciplines(c *gin.Context) {
	clubID := c.GetString("clubID")
	disciplines, err := h.useCases.ListDisciplines(c.Request.Context(), clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, disciplines)
}

func (h *DisciplineHandler) ListGroups(c *gin.Context) {
	clubID := c.GetString("clubID")
	disciplineID := c.Query("discipline_id")
	category := c.Query("category")
	groups, err := h.useCases.ListGroups(c.Request.Context(), clubID, disciplineID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

func (h *DisciplineHandler) ListStudentsInGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	clubID := c.GetString("clubID")
	students, err := h.useCases.ListStudentsInGroup(c.Request.Context(), clubID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

// --- Championships ---

type CreateTournamentRequest struct {
	Name         string `json:"name" binding:"required"`
	DisciplineID string `json:"discipline_id" binding:"required"`
	StartDate    string `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate      string `json:"end_date" binding:"required"`   // YYYY-MM-DD
	Format       string `json:"format" binding:"required"`
}

func (h *DisciplineHandler) CreateTournament(c *gin.Context) {
	// RBAC: Only ADMIN or SUPER_ADMIN can create tournaments
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	var req CreateTournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	clubID := c.GetString("clubID")
	t, err := h.useCases.CreateTournament(c.Request.Context(), clubID, req.Name, req.DisciplineID, start, end, req.Format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *DisciplineHandler) ListTournaments(c *gin.Context) {
	clubID := c.GetString("clubID")
	tournaments, err := h.useCases.ListTournaments(c.Request.Context(), clubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tournaments)
}

type RegisterTeamRequest struct {
	Name      string   `json:"name" binding:"required"`
	CaptainID string   `json:"captain_id"`
	MemberIDs []string `json:"member_ids"`
}

func (h *DisciplineHandler) RegisterTeam(c *gin.Context) {
	tournamentID := c.Param("id")
	var req RegisterTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	team, err := h.useCases.RegisterTeam(c.Request.Context(), clubID, tournamentID, req.Name, req.CaptainID, req.MemberIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, team)
}

type ScheduleMatchRequest struct {
	HomeTeamID string    `json:"home_team_id" binding:"required"`
	AwayTeamID string    `json:"away_team_id" binding:"required"`
	StartTime  time.Time `json:"start_time" binding:"required"`
	Location   string    `json:"location"`
	Round      string    `json:"round"`
}

func (h *DisciplineHandler) ScheduleMatch(c *gin.Context) {
	// RBAC: Only ADMIN or SUPER_ADMIN can schedule matches
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	tournamentID := c.Param("id")
	var req ScheduleMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	match, err := h.useCases.ScheduleMatch(c.Request.Context(), clubID, tournamentID, req.HomeTeamID, req.AwayTeamID, req.StartTime, req.Location, req.Round)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, match)
}

func (h *DisciplineHandler) ListMatches(c *gin.Context) {
	tournamentID := c.Param("id")
	clubID := c.GetString("clubID")
	matches, err := h.useCases.ListMatches(c.Request.Context(), clubID, tournamentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

type UpdateMatchResultRequest struct {
	ScoreHome int `json:"score_home"`
	ScoreAway int `json:"score_away"`
}

func (h *DisciplineHandler) UpdateMatchResult(c *gin.Context) {
	// RBAC: Only ADMIN or SUPER_ADMIN can update match results
	role, exists := c.Get("userRole")
	if !exists || (role != userDomain.RoleAdmin && role != userDomain.RoleSuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "requires ADMIN role"})
		return
	}

	matchID := c.Param("id")
	var req UpdateMatchResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubID := c.GetString("clubID")
	match, err := h.useCases.UpdateMatchResult(c.Request.Context(), clubID, matchID, req.ScoreHome, req.ScoreAway)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, match)
}

func (h *DisciplineHandler) GetStandings(c *gin.Context) {
	tournamentID := c.Param("id")
	clubID := c.GetString("clubID")
	standings, err := h.useCases.GetStandings(c.Request.Context(), clubID, tournamentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standings)
}

func RegisterRoutes(r *gin.RouterGroup, handler *DisciplineHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
	r.Use(authMiddleware, tenantMiddleware)
	disciplines := r.Group("/disciplines")
	{
		disciplines.GET("", handler.ListDisciplines)
	}

	groups := r.Group("/groups")
	{
		groups.GET("", handler.ListGroups)
		groups.GET("/:id/students", handler.ListStudentsInGroup)
	}

	tournaments := r.Group("/tournaments")
	{
		tournaments.POST("", handler.CreateTournament)
		tournaments.GET("", handler.ListTournaments)
		tournaments.POST("/:id/teams", handler.RegisterTeam)
		tournaments.POST("/:id/matches", handler.ScheduleMatch)
		tournaments.GET("/:id/matches", handler.ListMatches)
		tournaments.GET("/:id/standings", handler.GetStandings)
	}

	matches := r.Group("/matches")
	{
		matches.PUT("/:id/result", handler.UpdateMatchResult)
	}
}
