package e2e_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	championshipApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	championshipHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/infrastructure/http"
	championshipRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
)

func TestOddTeamsFixtureGeneration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Clean tables
	_ = db.Migrator().DropTable(&domain.Tournament{}, &domain.TournamentStage{}, &domain.Group{}, &domain.Standing{}, &domain.TournamentMatch{}, &domain.Team{})
	_ = db.AutoMigrate(&domain.Tournament{}, &domain.TournamentStage{}, &domain.Group{}, &domain.Standing{}, &domain.TournamentMatch{}, &domain.Team{})

	repo := championshipRepo.NewPostgresChampionshipRepository(db)
	uc := championshipApp.NewChampionshipUseCases(repo, nil, nil)
	h := championshipHttp.NewChampionshipHandler(uc, nil, nil)

	r := gin.New()
	clubID := uuid.New().String()
	authMw := func(c *gin.Context) {
		c.Set("userID", "test-admin")
		c.Set("clubID", clubID)
		c.Set("userRole", "ADMIN")
		c.Next()
	}

	// Data Setup
	tournamentID := uuid.New()
	db.Create(&domain.Tournament{ID: tournamentID, ClubID: uuid.MustParse(clubID), Name: "Odd Cup", Status: domain.TournamentActive})

	stageID := uuid.New()
	db.Create(&domain.TournamentStage{ID: stageID, TournamentID: tournamentID, Name: "Group Phase", Type: domain.StageGroup, Status: domain.StageActive})

	groupID := uuid.New()
	db.Create(&domain.Group{ID: groupID, StageID: stageID, Name: "Group A"})

	// Register 3 teams (ODD NUMBER)
	teams := []string{"Team A", "Team B", "Team C"}
	for range teams {
		tID := uuid.New()
		err := db.Create(&domain.Standing{
			ID:      uuid.New(),
			GroupID: groupID,
			TeamID:  tID,
			// Simplified standing, enough for fixture generation
		}).Error
		if err != nil {
			t.Fatalf("Failed to create standing: %v", err)
		}
	}

	t.Run("Generate Fixture with 3 Teams", func(t *testing.T) {
		group := r.Group("/api/v1")
		h.RegisterRoutes(group, authMw, func(c *gin.Context) { c.Next() })

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/championships/groups/"+groupID.String()+"/fixture", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify Matches in DB
		var matches []domain.TournamentMatch
		db.Find(&matches, "group_id = ?", groupID)

		// 3 teams = (3 * 2) / 2 = 3 matches (A vs B, A vs C, B vs C)
		assert.Equal(t, 3, len(matches))
	})
}
