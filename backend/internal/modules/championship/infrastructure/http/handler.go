package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/application"
)

type ChampionshipHandler struct {
	useCases *application.ChampionshipUseCases
}

func NewChampionshipHandler(useCases *application.ChampionshipUseCases) *ChampionshipHandler {
	return &ChampionshipHandler{useCases: useCases}
}

func (h *ChampionshipHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc, tenantMiddleware gin.HandlerFunc) {
	group := r.Group("/championships")
	group.Use(authMiddleware)
	// group.Use(tenantMiddleware) // Optional: If championships belong to a tenant/club context
	{
		group.GET("/", h.ListTournaments) // Might need to be public or protected? Let's protect it for Admin.
		group.POST("/", h.CreateTournament)
		group.POST("/:id/stages", h.AddStage)
		group.POST("/stages/:id/groups", h.AddGroup)
		group.POST("/groups/:id/teams", h.RegisterTeam)
		group.POST("/groups/:id/fixture", h.GenerateFixture)
		group.GET("/groups/:id/matches", h.GetMatchesByGroup)
		group.GET("/groups/:id/standings", h.GetStandings)
		group.POST("/matches/result", h.UpdateMatchResult)
		group.POST("/matches/schedule", h.ScheduleMatch)
	}
}

func (h *ChampionshipHandler) ListTournaments(c *gin.Context) {
	clubID := c.Query("club_id") // Assume passed as query param for now, or get from context if possible
	// fallback to context if query empty?
	// For now, let's require it or use a default from context if we had one.
	// But TenantMiddleware is not applied? It is applied in app.go.
	// So we can assume we might get it. But let's use Query for flexibility.
	if clubID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "club_id query parameter is required"})
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
	groupID := c.Param("id")
	standings, err := h.useCases.GetStandings(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standings)
}

func (h *ChampionshipHandler) GetMatchesByGroup(c *gin.Context) {
	groupID := c.Param("id")
	matches, err := h.useCases.GetMatchesByGroup(groupID)
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

	group, err := h.useCases.AddGroup(input.StageID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func (h *ChampionshipHandler) RegisterTeam(c *gin.Context) {
	var input application.RegisterTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	groupID := c.Param("id")
	standing, err := h.useCases.RegisterTeam(groupID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, standing)
}

func (h *ChampionshipHandler) GenerateFixture(c *gin.Context) {
	groupID := c.Param("id")
	matches, err := h.useCases.GenerateGroupFixture(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, matches)
}
