package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/domain"
	handler "github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/infrastructure/http"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Repositories Mocks (Reusing from application_test style) ---

type MockDisciplineRepo struct {
	mock.Mock
}

func (m *MockDisciplineRepo) CreateDiscipline(ctx context.Context, d *domain.Discipline) error {
	return m.Called(ctx, d).Error(0)
}
func (m *MockDisciplineRepo) ListDisciplines(ctx context.Context, clubID string) ([]domain.Discipline, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]domain.Discipline), args.Error(1)
}
func (m *MockDisciplineRepo) GetDisciplineByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Discipline, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Discipline), args.Error(1)
}
func (m *MockDisciplineRepo) CreateGroup(ctx context.Context, g *domain.TrainingGroup) error {
	return m.Called(ctx, g).Error(0)
}
func (m *MockDisciplineRepo) ListGroups(ctx context.Context, clubID string, f map[string]interface{}) ([]domain.TrainingGroup, error) {
	args := m.Called(ctx, clubID, f)
	return args.Get(0).([]domain.TrainingGroup), args.Error(1)
}
func (m *MockDisciplineRepo) GetGroupByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.TrainingGroup, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TrainingGroup), args.Error(1)
}

type MockTournamentRepo struct {
	mock.Mock
}

func (m *MockTournamentRepo) CreateTournament(ctx context.Context, t *domain.Tournament) error {
	return m.Called(ctx, t).Error(0)
}
func (m *MockTournamentRepo) GetTournamentByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Tournament, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tournament), args.Error(1)
}
func (m *MockTournamentRepo) ListTournaments(ctx context.Context, clubID string) ([]domain.Tournament, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]domain.Tournament), args.Error(1)
}
func (m *MockTournamentRepo) UpdateTournament(ctx context.Context, t *domain.Tournament) error {
	return m.Called(ctx, t).Error(0)
}
func (m *MockTournamentRepo) CreateTeam(ctx context.Context, t *domain.Team) error {
	return m.Called(ctx, t).Error(0)
}
func (m *MockTournamentRepo) GetTeamByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Team, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Team), args.Error(1)
}
func (m *MockTournamentRepo) ListTeams(ctx context.Context, clubID string, tID uuid.UUID) ([]domain.Team, error) {
	args := m.Called(ctx, clubID, tID)
	return args.Get(0).([]domain.Team), args.Error(1)
}
func (m *MockTournamentRepo) CreateMatch(ctx context.Context, match *domain.Match) error {
	return m.Called(ctx, match).Error(0)
}
func (m *MockTournamentRepo) UpdateMatch(ctx context.Context, match *domain.Match) error {
	return m.Called(ctx, match).Error(0)
}
func (m *MockTournamentRepo) GetMatchByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Match, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Match), args.Error(1)
}
func (m *MockTournamentRepo) ListMatches(ctx context.Context, clubID string, tID uuid.UUID) ([]domain.Match, error) {
	args := m.Called(ctx, clubID, tID)
	return args.Get(0).([]domain.Match), args.Error(1)
}
func (m *MockTournamentRepo) GetStandings(ctx context.Context, clubID string, tID uuid.UUID) ([]domain.Standing, error) {
	args := m.Called(ctx, clubID, tID)
	return args.Get(0).([]domain.Standing), args.Error(1)
}

// --- Setup ---

func setupRouter(h *handler.DisciplineHandler, clubID, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("clubID", clubID)
		c.Set("userRole", role)
		c.Next()
	})
	api := r.Group("/api/v1")
	auth := func(c *gin.Context) {}
	tenant := func(c *gin.Context) {}
	handler.RegisterRoutes(api, h, auth, tenant)
	return r
}

// --- Tests ---

func TestDisciplineHandler_DisciplinesAndGroups(t *testing.T) {
	mockRepo := new(MockDisciplineRepo)
	uc := application.NewDisciplineUseCases(mockRepo, nil, nil)
	h := handler.NewDisciplineHandler(uc)
	clubID := uuid.New().String()
	r := setupRouter(h, clubID, userDomain.RoleAdmin)

	t.Run("List Disciplines", func(t *testing.T) {
		mockRepo.On("ListDisciplines", mock.Anything, clubID).Return([]domain.Discipline{{Name: "Padel"}}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/disciplines", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Padel")
	})

	t.Run("List Groups with Filters", func(t *testing.T) {
		mockRepo.On("ListGroups", mock.Anything, clubID, mock.Anything).Return([]domain.TrainingGroup{{Name: "Group A"}}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/groups?category=2012", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Group A")
	})
}

func TestDisciplineHandler_Tournaments(t *testing.T) {
	mockTourneyRepo := new(MockTournamentRepo)
	uc := application.NewDisciplineUseCases(nil, mockTourneyRepo, nil)
	h := handler.NewDisciplineHandler(uc)
	clubID := uuid.New().String()
	r := setupRouter(h, clubID, userDomain.RoleAdmin)

	t.Run("Create Tournament", func(t *testing.T) {
		mockTourneyRepo.On("CreateTournament", mock.Anything, mock.Anything).Return(nil).Once()
		body, _ := json.Marshal(map[string]string{
			"name":          "Open 2024",
			"discipline_id": uuid.New().String(),
			"start_date":    "2024-01-01",
			"end_date":      "2024-01-31",
			"format":        "LEAGUE",
		})
		req, _ := http.NewRequest("POST", "/api/v1/tournaments", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("Update Match Result", func(t *testing.T) {
		mID := uuid.New()
		mockTourneyRepo.On("GetMatchByID", mock.Anything, clubID, mID).Return(&domain.Match{ID: mID}, nil).Once()
		mockTourneyRepo.On("UpdateMatch", mock.Anything, mock.Anything).Return(nil).Once()

		body, _ := json.Marshal(map[string]int{"score_home": 3, "score_away": 1})
		req, _ := http.NewRequest("PUT", "/api/v1/matches/"+mID.String()+"/result", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Get Standings", func(t *testing.T) {
		tID := uuid.New()
		mockTourneyRepo.On("GetStandings", mock.Anything, clubID, tID).Return([]domain.Standing{{TeamName: "Eagles"}}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/tournaments/"+tID.String()+"/standings", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Eagles")
	})

	t.Run("Extended Success Paths", func(t *testing.T) {
		// 1. List Tournaments
		mockTourneyRepo.On("ListTournaments", mock.Anything, clubID).Return([]domain.Tournament{{Name: "Spring Open"}}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/tournaments", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Spring Open")

		// 2. Register Team
		tID := uuid.New()
		mockTourneyRepo.On("GetTournamentByID", mock.Anything, clubID, tID).Return(&domain.Tournament{ID: tID}, nil).Once()
		mockTourneyRepo.On("CreateTeam", mock.Anything, mock.Anything).Return(nil).Once()

		teamBody, _ := json.Marshal(map[string]string{"name": "Team Rocket"})
		req, _ = http.NewRequest("POST", "/api/v1/tournaments/"+tID.String()+"/teams", bytes.NewBuffer(teamBody))
		resp = httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)

		// 3. Schedule Match
		team1 := uuid.New()
		team2 := uuid.New()
		mockTourneyRepo.On("GetTeamByID", mock.Anything, clubID, team1).Return(&domain.Team{}, nil).Once()
		mockTourneyRepo.On("GetTeamByID", mock.Anything, clubID, team2).Return(&domain.Team{}, nil).Once()
		mockTourneyRepo.On("CreateMatch", mock.Anything, mock.Anything).Return(nil).Once() // For ScheduleMatch in usecase

		matchBody, _ := json.Marshal(map[string]interface{}{
			"home_team_id": team1.String(),
			"away_team_id": team2.String(),
			"start_time":   "2024-06-01T10:00:00Z",
		})
		req, _ = http.NewRequest("POST", "/api/v1/tournaments/"+tID.String()+"/matches", bytes.NewBuffer(matchBody))
		resp = httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)

		// 4. List Matches
		mockTourneyRepo.On("ListMatches", mock.Anything, clubID, tID).Return([]domain.Match{{Location: "Court 1"}}, nil).Once()
		req, _ = http.NewRequest("GET", "/api/v1/tournaments/"+tID.String()+"/matches", nil)
		resp = httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Court 1")
	})
}

func TestDisciplinesHandler_DetailedErrors(t *testing.T) {
	mockDisciplineRepo := new(MockDisciplineRepo)
	mockTourneyRepo := new(MockTournamentRepo)
	uc := application.NewDisciplineUseCases(mockDisciplineRepo, mockTourneyRepo, nil)
	h := handler.NewDisciplineHandler(uc)
	clubID := uuid.New().String()
	r := setupRouter(h, clubID, userDomain.RoleAdmin)

	// --- Disciplines & Groups ---

	t.Run("ListSteps_RepoError", func(t *testing.T) {
		mockDisciplineRepo.On("ListDisciplines", mock.Anything, clubID).Return([]domain.Discipline{}, context.DeadlineExceeded).Once()
		req, _ := http.NewRequest("GET", "/api/v1/disciplines", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("ListGroups_RepoError", func(t *testing.T) {
		mockDisciplineRepo.On("ListGroups", mock.Anything, clubID, mock.Anything).Return([]domain.TrainingGroup{}, context.DeadlineExceeded).Once()
		req, _ := http.NewRequest("GET", "/api/v1/groups", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("ListStudents_InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/groups/invalid/students", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	// --- Tournaments ---

	t.Run("CreateTournament_RBAC", func(t *testing.T) {
		rMember := setupRouter(h, clubID, userDomain.RoleMember)
		req, _ := http.NewRequest("POST", "/api/v1/tournaments", nil)
		resp := httptest.NewRecorder()
		rMember.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("CreateTournament_InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/tournaments", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("CreateTournament_InvalidDate", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"name":          "Open",
			"discipline_id": uuid.New().String(),
			"start_date":    "invalid-date",
			"end_date":      "2024-01-31",
			"format":        "LEAGUE",
		})
		req, _ := http.NewRequest("POST", "/api/v1/tournaments", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("ListTournaments_Error", func(t *testing.T) {
		mockTourneyRepo.On("ListTournaments", mock.Anything, clubID).Return([]domain.Tournament{}, context.DeadlineExceeded).Once()
		req, _ := http.NewRequest("GET", "/api/v1/tournaments", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("RegisterTeam_InvalidJSON", func(t *testing.T) {
		tID := uuid.New()
		req, _ := http.NewRequest("POST", "/api/v1/tournaments/"+tID.String()+"/teams", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("RegisterTeam_RepoError", func(t *testing.T) {
		tID := uuid.New().String()
		mockTourneyRepo.On("GetTournamentByID", mock.Anything, clubID, mock.Anything).Return(&domain.Tournament{ID: uuid.New()}, nil).Once()
		mockTourneyRepo.On("CreateTeam", mock.Anything, mock.Anything).Return(context.DeadlineExceeded).Once()

		body, _ := json.Marshal(map[string]string{"name": "Team A"})
		req, _ := http.NewRequest("POST", "/api/v1/tournaments/"+tID+"/teams", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	// --- Matches ---

	t.Run("ScheduleMatch_RBAC", func(t *testing.T) {
		rMember := setupRouter(h, clubID, userDomain.RoleMember)
		tID := uuid.New()
		req, _ := http.NewRequest("POST", "/api/v1/tournaments/"+tID.String()+"/matches", nil)
		resp := httptest.NewRecorder()
		rMember.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("ScheduleMatch_InvalidJSON", func(t *testing.T) {
		tID := uuid.New()
		req, _ := http.NewRequest("POST", "/api/v1/tournaments/"+tID.String()+"/matches", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("ListMatches_Error", func(t *testing.T) {
		tID := uuid.New()
		mockTourneyRepo.On("ListMatches", mock.Anything, clubID, tID).Return([]domain.Match{}, context.DeadlineExceeded).Once()
		req, _ := http.NewRequest("GET", "/api/v1/tournaments/"+tID.String()+"/matches", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("UpdateResult_RBAC", func(t *testing.T) {
		rMember := setupRouter(h, clubID, userDomain.RoleMember)
		mID := uuid.New()
		req, _ := http.NewRequest("PUT", "/api/v1/matches/"+mID.String()+"/result", nil)
		resp := httptest.NewRecorder()
		rMember.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("GetStandings_Error", func(t *testing.T) {
		tID := uuid.New()
		mockTourneyRepo.On("GetStandings", mock.Anything, clubID, tID).Return([]domain.Standing{}, context.DeadlineExceeded).Once()
		req, _ := http.NewRequest("GET", "/api/v1/tournaments/"+tID.String()+"/standings", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
