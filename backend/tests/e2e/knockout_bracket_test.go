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
	bookingApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	bookingRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"
	champApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	champHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/infrastructure/http"
	champRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/infrastructure/repository"
	champSvc "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/infrastructure/service"
	clubApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
	clubRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/infrastructure/repository"
	facRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"
	userApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKnockoutBracketGeneration(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	clubID := "test-club-knockout"

	// Clean up
	db.Exec("DELETE FROM tournament_matches WHERE tournament_id IN (SELECT id FROM tournaments WHERE club_id = ?)", clubID)
	db.Exec("DELETE FROM standings WHERE group_id IN (SELECT g.id FROM groups g JOIN tournament_stages s ON g.stage_id = s.id JOIN tournaments t ON s.tournament_id = t.id WHERE t.club_id = ?)", clubID)
	db.Exec("DELETE FROM groups WHERE stage_id IN (SELECT s.id FROM tournament_stages s JOIN tournaments t ON s.tournament_id = t.id WHERE t.club_id = ?)", clubID)
	db.Exec("DELETE FROM tournament_stages WHERE tournament_id IN (SELECT id FROM tournaments WHERE club_id = ?)", clubID)
	db.Exec("DELETE FROM tournaments WHERE club_id = ?", clubID)

	// 2. Setup Dependencies
	champRepository := champRepo.NewPostgresChampionshipRepository(db)
	volunteerRepository := champRepo.NewPostgresVolunteerRepository(db)
	clubRepository := clubRepo.NewPostgresClubRepository(db)

	// Booking adapter (mock)
	userR := userRepo.NewPostgresUserRepository(db)
	facR := facRepo.NewPostgresFacilityRepository(db)
	bookingR := bookingRepo.NewPostgresBookingRepository(db)
	recurringR := bookingRepo.NewPostgresRecurringRepository(db)
	familyR := userRepo.NewPostgresFamilyGroupRepository(db)
	userUC := userApp.NewUserUseCases(userR, familyR)
	bookingUC := bookingApp.NewBookingUseCases(bookingR, recurringR, facR, userR, nil, nil)
	bookingAdapter := champSvc.NewChampionshipBookingAdapter(bookingUC)

	champUC := champApp.NewChampionshipUseCases(champRepository, bookingAdapter, userUC)
	volunteerService := champApp.NewVolunteerService(volunteerRepository)
	clubUC := clubApp.NewClubUseCases(clubRepository, clubRepository, clubRepository, nil)

	handler := champHttp.NewChampionshipHandler(champUC, volunteerService, clubUC)

	// Router
	r := gin.New()
	authMiddleware := func(c *gin.Context) {
		c.Set("userID", "admin-user")
		c.Set("clubID", clubID)
		c.Set("userRole", userDomain.RoleAdmin)
		c.Next()
	}
	tenantMiddleware := func(c *gin.Context) { c.Next() }
	handler.RegisterRoutes(r.Group("/api/v1"), authMiddleware, tenantMiddleware)

	// 3. Create Tournament
	tournament, err := champUC.CreateTournament(champApp.CreateTournamentInput{
		ClubID:    clubID,
		Name:      "Knockout Cup",
		Sport:     "FUTBOL",
		Category:  "Senior",
		StartDate: time.Now().Add(7 * 24 * time.Hour),
	})
	require.NoError(t, err)

	// 4. Create KNOCKOUT Stage
	stage, err := champUC.AddStage(tournament.ID.String(), champApp.AddStageInput{
		ClubID: clubID,
		Name:   "Playoffs",
		Type:   "KNOCKOUT",
		Order:  1,
	})
	require.NoError(t, err)
	require.Equal(t, domain.StageKnockout, stage.Type)

	// 5. Create some teams for the bracket (need 4 for quarterfinals simulation)
	teamIDs := make([]string, 4)
	for i := 0; i < 4; i++ {
		teamIDs[i] = uuid.New().String()
		// Insert teams directly into DB (or use Team module if available)
		db.Exec("INSERT INTO teams (id, name) VALUES (?, ?)", teamIDs[i], "Team "+string(rune('A'+i)))
	}

	// 6. Test: Generate Knockout Bracket via API
	t.Run("Generate 4-Team Knockout Bracket", func(t *testing.T) {
		type GenerateKnockoutRequest struct {
			SeedOrder []string `json:"seed_order"`
		}
		reqBody := GenerateKnockoutRequest{
			SeedOrder: teamIDs, // 4 teams = 2 matches (semifinals)
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/championships/stages/"+stage.ID.String()+"/knockout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		// Should generate 2 matches for 4 teams
		matchCount := int(resp["count"].(float64))
		assert.Equal(t, 2, matchCount, "4 teams should generate 2 matches")

		matches := resp["matches"].([]interface{})
		assert.Len(t, matches, 2)

		// Verify pairings: #1 vs #4, #2 vs #3
		match1 := matches[0].(map[string]interface{})
		match2 := matches[1].(map[string]interface{})

		assert.Equal(t, teamIDs[0], match1["home_team_id"])
		assert.Equal(t, teamIDs[3], match1["away_team_id"])
		assert.Equal(t, teamIDs[1], match2["home_team_id"])
		assert.Equal(t, teamIDs[2], match2["away_team_id"])
	})

	// 7. Test: Invalid team count (not power of 2)
	t.Run("Reject Non-Power-of-2 Teams", func(t *testing.T) {
		// Create another stage
		stage2, _ := champUC.AddStage(tournament.ID.String(), champApp.AddStageInput{
			ClubID: clubID,
			Name:   "Playoffs 2",
			Type:   "KNOCKOUT",
			Order:  2,
		})

		type GenerateKnockoutRequest struct {
			SeedOrder []string `json:"seed_order"`
		}
		reqBody := GenerateKnockoutRequest{
			SeedOrder: []string{uuid.New().String(), uuid.New().String(), uuid.New().String()}, // 3 teams - invalid!
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/championships/stages/"+stage2.ID.String()+"/knockout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Contains(t, resp["error"], "power of 2")
	})

	// 8. Test: GROUP stage should be rejected
	t.Run("Reject GROUP Stage for Knockout", func(t *testing.T) {
		// Create GROUP stage
		groupStage, _ := champUC.AddStage(tournament.ID.String(), champApp.AddStageInput{
			ClubID: clubID,
			Name:   "Group Phase",
			Type:   "GROUP",
			Order:  3,
		})

		type GenerateKnockoutRequest struct {
			SeedOrder []string `json:"seed_order"`
		}
		reqBody := GenerateKnockoutRequest{
			SeedOrder: []string{uuid.New().String(), uuid.New().String()},
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/championships/stages/"+groupStage.ID.String()+"/knockout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Contains(t, resp["error"], "KNOCKOUT")
	})
}
