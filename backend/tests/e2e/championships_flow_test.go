package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/application"
	authRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/repository"
	authToken "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/token"
	disciplineApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/domain"
	disciplineHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/infrastructure/http"
	disciplineRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/infrastructure/repository"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChampionshipsFlow(t *testing.T) {
	// 1. Setup Environment
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Ensure clean state
	_ = db.Migrator().DropTable(&domain.Discipline{}, &domain.Tournament{}, &domain.Match{}, &domain.Team{}, &domain.Standing{}, &userDomain.UserStats{}, &userDomain.Wallet{}, &userRepo.UserModel{}) // Tables likely created
	_ = db.AutoMigrate(&domain.Discipline{}, &domain.Tournament{}, &domain.Match{}, &domain.Team{}, &domain.Standing{}, &userRepo.UserModel{})

	// Clear PostgreSQL cached prepared statements after schema change
	db.Exec("DISCARD ALL")

	// 2. Setup Dependencies
	userR := userRepo.NewPostgresUserRepository(db)
	dRepo := disciplineRepo.NewPostgresDisciplineRepository(db)
	tRepo := disciplineRepo.NewPostgresTournamentRepository(db)
	dUC := disciplineApp.NewDisciplineUseCases(dRepo, tRepo, userR)
	dHandler := disciplineHttp.NewDisciplineHandler(dUC)

	// Auth Setup for Helper
	tokenService := authToken.NewJWTService("secret")
	authR := authRepo.NewPostgresAuthRepository(db)
	_ = authApp.NewAuthUseCases(authR, tokenService, nil)

	r := gin.New()
	clubID := "test-club-championships"
	authMiddleware := func(c *gin.Context) {
		c.Set("userID", "admin-user")
		c.Set("clubID", clubID)
		c.Set("userRole", userDomain.RoleAdmin) // Fix for RBAC
		c.Next()
	}
	disciplineHttp.RegisterRoutes(r.Group("/api/v1"), dHandler, authMiddleware, func(c *gin.Context) { c.Next() })

	// 3. Create Discipline (Prerequisite)
	discID := uuid.New()
	disc := &domain.Discipline{
		ID:     discID,
		ClubID: clubID,
		Name:   "Football 5",
	}
	// Direct DB insert for prerequisite
	db.Create(disc)

	// 4. Test: Create Tournament
	t.Run("Create Tournament", func(t *testing.T) {
		start := time.Now()
		// Define DTO locally to match Handler expectation
		type CreateTournamentRequest struct {
			Name         string `json:"name"`
			DisciplineID string `json:"discipline_id"`
			StartDate    string `json:"start_date"`
			EndDate      string `json:"end_date"`
			Format       string `json:"format"`
		}
		reqBody := CreateTournamentRequest{
			Name:         "Summer Cup 2024",
			DisciplineID: discID.String(),
			StartDate:    start.Format("2006-01-02"),
			EndDate:      start.AddDate(0, 1, 0).Format("2006-01-02"),
			Format:       "LEAGUE",
		}
		w := httptest.NewRecorder()
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/tournaments", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		tID, ok := resp["id"].(string)
		require.True(t, ok)
		require.NotEmpty(t, tID)
	})

	// Get the created Tournament ID (Assuming single one for test simplicity or query DB)
	var tournament domain.Tournament
	db.Where("club_id = ?", clubID).First(&tournament)
	require.NotNil(t, tournament.ID)

	// 5. Test: Register Teams
	var teamAID, teamBID uuid.UUID

	type RegisterTeamRequest struct {
		Name      string   `json:"name"`
		CaptainID string   `json:"captain_id"`
		MemberIDs []string `json:"member_ids"`
	}

	t.Run("Register Teams", func(t *testing.T) {
		// Register Team A
		w := httptest.NewRecorder()
		reqBody := RegisterTeamRequest{
			Name: "Team Alpha",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/tournaments/"+tournament.ID.String()+"/teams", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		teamAID = uuid.MustParse(resp["id"].(string))

		// Register Team B
		w2 := httptest.NewRecorder()
		reqBody2 := RegisterTeamRequest{
			Name: "Team Beta",
		}
		body2, _ := json.Marshal(reqBody2)
		req2, _ := http.NewRequest("POST", "/api/v1/tournaments/"+tournament.ID.String()+"/teams", bytes.NewBuffer(body2))
		r.ServeHTTP(w2, req2)
		require.Equal(t, http.StatusCreated, w2.Code)
		var resp2 map[string]interface{}
		_ = json.Unmarshal(w2.Body.Bytes(), &resp2)
		teamBID = uuid.MustParse(resp2["id"].(string))
	})

	// 6. Test: Schedule Match
	var matchID uuid.UUID

	type ScheduleMatchRequest struct {
		HomeTeamID string    `json:"home_team_id"`
		AwayTeamID string    `json:"away_team_id"`
		StartTime  time.Time `json:"start_time"`
		Location   string    `json:"location"`
		Round      string    `json:"round"`
	}

	t.Run("Schedule Match", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := ScheduleMatchRequest{
			HomeTeamID: teamAID.String(),
			AwayTeamID: teamBID.String(),
			StartTime:  time.Now().Add(24 * time.Hour),
			Location:   "Pitch 1",
			Round:      "1",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/tournaments/"+tournament.ID.String()+"/matches", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		matchID = uuid.MustParse(resp["id"].(string))
	})

	// 7. Test: Update Match Result
	t.Run("Update Match Result", func(t *testing.T) {
		w := httptest.NewRecorder()

		type UpdateMatchResultRequest struct {
			ScoreHome int `json:"score_home"`
			ScoreAway int `json:"score_away"`
		}

		reqBody := UpdateMatchResultRequest{
			ScoreHome: 3,
			ScoreAway: 1,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/api/v1/matches/"+matchID.String()+"/result", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	// 8. Get Standings
	t.Run("Get Standings", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/tournaments/"+tournament.ID.String()+"/standings", nil)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var standings []domain.Standing
		_ = json.Unmarshal(w.Body.Bytes(), &standings)

		// Verify
		// Team A should have 3 points (Won)
		// Team B should have 0 points (Lost)
		foundA := false
		for _, s := range standings {
			if s.TeamID == teamAID {
				assert.Equal(t, 3, s.Points)
				assert.Equal(t, 1, s.Won)
				foundA = true
			}
			if s.TeamID == teamBID {
				assert.Equal(t, 0, s.Points)
				assert.Equal(t, 1, s.Lost)
			}
		}
		assert.True(t, foundA, "Team A standing should be present")
	})
}
