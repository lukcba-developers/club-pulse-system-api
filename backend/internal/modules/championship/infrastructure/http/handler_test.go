package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	handler "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/infrastructure/http"
	clubApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
	clubDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockChampionshipRepo struct {
	mock.Mock
}

func (m *MockChampionshipRepo) CreateTournament(ctx context.Context, t *domain.Tournament) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetTournament(ctx context.Context, clubID, id string) (*domain.Tournament, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tournament), args.Error(1)
}

func (m *MockChampionshipRepo) ListTournaments(ctx context.Context, clubID string) ([]domain.Tournament, error) {
	args := m.Called(ctx, clubID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Tournament), args.Error(1)
}

func (m *MockChampionshipRepo) CreateStage(ctx context.Context, s *domain.TournamentStage) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetStage(ctx context.Context, clubID, id string) (*domain.TournamentStage, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TournamentStage), args.Error(1)
}

func (m *MockChampionshipRepo) CreateGroup(ctx context.Context, g *domain.Group) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetGroup(ctx context.Context, clubID, id string) (*domain.Group, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Group), args.Error(1)
}

func (m *MockChampionshipRepo) CreateMatch(ctx context.Context, clubID string, match *domain.TournamentMatch) error {
	args := m.Called(ctx, clubID, match)
	return args.Error(0)
}

func (m *MockChampionshipRepo) CreateMatchesBatch(ctx context.Context, clubID string, matches []domain.TournamentMatch) error {
	args := m.Called(ctx, clubID, matches)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetMatch(ctx context.Context, clubID, id string) (*domain.TournamentMatch, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TournamentMatch), args.Error(1)
}

func (m *MockChampionshipRepo) GetMatchesByGroup(ctx context.Context, clubID, groupID string) ([]domain.TournamentMatch, error) {
	args := m.Called(ctx, clubID, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.TournamentMatch), args.Error(1)
}

func (m *MockChampionshipRepo) UpdateMatchResult(ctx context.Context, clubID, matchID string, homeScore, awayScore float64) error {
	args := m.Called(ctx, clubID, matchID, homeScore, awayScore)
	return args.Error(0)
}

func (m *MockChampionshipRepo) UpdateMatchScheduling(ctx context.Context, clubID, matchID string, date time.Time, bookingID uuid.UUID) error {
	args := m.Called(ctx, clubID, matchID, date, bookingID)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetStandings(ctx context.Context, clubID, groupID string) ([]domain.Standing, error) {
	args := m.Called(ctx, clubID, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Standing), args.Error(1)
}

func (m *MockChampionshipRepo) RegisterTeam(ctx context.Context, clubID string, s *domain.Standing) error {
	args := m.Called(ctx, clubID, s)
	return args.Error(0)
}

func (m *MockChampionshipRepo) UpdateStanding(ctx context.Context, s *domain.Standing) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockChampionshipRepo) UpdateStandingsBatch(ctx context.Context, clubID string, s []domain.Standing) error {
	args := m.Called(ctx, clubID, s)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetTeamMembers(ctx context.Context, id string) ([]string, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockChampionshipRepo) CreateTeam(ctx context.Context, team *domain.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockChampionshipRepo) AddMember(ctx context.Context, teamID, userID string) error {
	args := m.Called(ctx, teamID, userID)
	return args.Error(0)
}

func (m *MockChampionshipRepo) GetMatchesByUserID(ctx context.Context, clubID, userID string) ([]domain.TournamentMatch, error) {
	args := m.Called(ctx, clubID, userID)
	return args.Get(0).([]domain.TournamentMatch), args.Error(1)
}

func (m *MockChampionshipRepo) GetUpcomingMatches(ctx context.Context, clubID string, from, to time.Time) ([]domain.TournamentMatch, error) {
	args := m.Called(ctx, clubID, from, to)
	return args.Get(0).([]domain.TournamentMatch), args.Error(1)
}

type MockVolunteerRepo struct {
	mock.Mock
}

func (m *MockVolunteerRepo) Create(ctx context.Context, v *domain.VolunteerAssignment) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}

func (m *MockVolunteerRepo) GetByMatchID(ctx context.Context, clubID string, matchID uuid.UUID) ([]domain.VolunteerAssignment, error) {
	args := m.Called(ctx, clubID, matchID)
	return args.Get(0).([]domain.VolunteerAssignment), args.Error(1)
}

func (m *MockVolunteerRepo) GetByUserID(ctx context.Context, clubID, userID string) ([]domain.VolunteerAssignment, error) {
	args := m.Called(ctx, clubID, userID)
	return args.Get(0).([]domain.VolunteerAssignment), args.Error(1)
}

func (m *MockVolunteerRepo) GetByRoleAndMatch(ctx context.Context, clubID string, matchID uuid.UUID, role domain.VolunteerRole) ([]domain.VolunteerAssignment, error) {
	args := m.Called(ctx, clubID, matchID, role)
	return args.Get(0).([]domain.VolunteerAssignment), args.Error(1)
}

func (m *MockVolunteerRepo) Update(ctx context.Context, v *domain.VolunteerAssignment) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}

func (m *MockVolunteerRepo) Delete(ctx context.Context, clubID string, id uuid.UUID) error {
	args := m.Called(ctx, clubID, id)
	return args.Error(0)
}

type MockVolunteerService struct {
	mock.Mock
}

func (m *MockVolunteerService) AssignVolunteer(ctx context.Context, clubID string, matchID uuid.UUID, userID string, role domain.VolunteerRole, assignedBy string, notes string) error {
	args := m.Called(ctx, clubID, matchID, userID, role, assignedBy, notes)
	return args.Error(0)
}

func (m *MockVolunteerService) GetMatchVolunteers(ctx context.Context, clubID string, matchID uuid.UUID) ([]domain.VolunteerAssignment, error) {
	args := m.Called(ctx, clubID, matchID)
	var res []domain.VolunteerAssignment
	if args.Get(0) != nil {
		res = args.Get(0).([]domain.VolunteerAssignment)
	}
	return res, args.Error(1)
}

func (m *MockVolunteerService) RemoveVolunteer(ctx context.Context, clubID string, assignmentID uuid.UUID) error {
	args := m.Called(ctx, clubID, assignmentID)
	return args.Error(0)
}

func (m *MockVolunteerService) GetVolunteerSummary(ctx context.Context, clubID string, matchID uuid.UUID) (*domain.VolunteerSummary, error) {
	args := m.Called(ctx, clubID, matchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.VolunteerSummary), args.Error(1)
}

func (m *MockVolunteerService) GetUserAssignments(ctx context.Context, clubID, userID string) ([]domain.VolunteerAssignment, error) {
	args := m.Called(ctx, clubID, userID)
	return args.Get(0).([]domain.VolunteerAssignment), args.Error(1)
}

func (m *MockVolunteerService) ValidateAssignment(ctx context.Context, clubID string, matchID uuid.UUID, role domain.VolunteerRole, maxPerRole int) error {
	args := m.Called(ctx, clubID, matchID, role, maxPerRole)
	return args.Error(0)
}

type MockBookingService struct {
	mock.Mock
}

func (m *MockBookingService) CreateSystemBooking(clubID, courtID string, start, end time.Time, notes string) (*uuid.UUID, error) {
	args := m.Called(clubID, courtID, start, end, notes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*uuid.UUID), args.Error(1)
}

type MockClubRepo struct {
	mock.Mock
}

func (m *MockClubRepo) Create(ctx context.Context, club *clubDomain.Club) error {
	return m.Called(ctx, club).Error(0)
}
func (m *MockClubRepo) GetByID(ctx context.Context, id string) (*clubDomain.Club, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*clubDomain.Club), args.Error(1)
}
func (m *MockClubRepo) GetBySlug(ctx context.Context, slug string) (*clubDomain.Club, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*clubDomain.Club), args.Error(1)
}
func (m *MockClubRepo) List(ctx context.Context, limit, offset int) ([]clubDomain.Club, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]clubDomain.Club), args.Error(1)
}
func (m *MockClubRepo) Update(ctx context.Context, club *clubDomain.Club) error {
	return m.Called(ctx, club).Error(0)
}
func (m *MockClubRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockClubRepo) GetMemberEmails(ctx context.Context, clubID string) ([]string, error) {
	args := m.Called(ctx, clubID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// --- Test Setup ---

func setupRouter(h *handler.ChampionshipHandler, clubID, userID, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("clubID", clubID)
		c.Set("userID", userID)
		c.Set("userRole", role)
		c.Next()
	})

	api := r.Group("/api/v1")
	auth := func(c *gin.Context) {}
	tenant := func(c *gin.Context) {}

	h.RegisterRoutes(api, auth, tenant)

	return r
}

// --- Tests ---

func TestChampionshipHandler_Tournaments(t *testing.T) {
	mockRepo := new(MockChampionshipRepo)
	uc := application.NewChampionshipUseCases(mockRepo, nil, nil)
	h := handler.NewChampionshipHandler(uc, nil, nil)
	cID := uuid.New().String()
	uID := uuid.New().String()
	r := setupRouter(h, cID, uID, userDomain.RoleAdmin)

	t.Run("Create Tournament", func(t *testing.T) {
		mockRepo.On("CreateTournament", mock.Anything, mock.Anything).Return(nil).Once()
		body, _ := json.Marshal(application.CreateTournamentInput{
			ClubID:    cID,
			Name:      "Interclub",
			Sport:     "Padel",
			StartDate: time.Now(),
		})
		req, _ := http.NewRequest("POST", "/api/v1/championships/", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("List Tournaments", func(t *testing.T) {
		mockRepo.On("ListTournaments", mock.Anything, cID).Return([]domain.Tournament{{Name: "T1"}}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/championships/", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "T1")
	})
}

func TestChampionshipHandler_StandingsAndFixture(t *testing.T) {
	mockRepo := new(MockChampionshipRepo)
	uc := application.NewChampionshipUseCases(mockRepo, nil, nil)
	h := handler.NewChampionshipHandler(uc, nil, nil)
	cID := uuid.New().String()
	uID := uuid.New().String()
	r := setupRouter(h, cID, uID, userDomain.RoleAdmin)

	t.Run("Get Standings", func(t *testing.T) {
		mockRepo.On("GetStandings", mock.Anything, cID, "g1").Return([]domain.Standing{{TeamName: "Dream Team"}}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/championships/groups/g1/standings", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Generate Fixture", func(t *testing.T) {
		mockRepo.On("GetStandings", mock.Anything, cID, "g1").Return([]domain.Standing{{TeamID: uuid.New()}, {TeamID: uuid.New()}}, nil).Once()
		mockRepo.On("GetGroup", mock.Anything, cID, "g1").Return(&domain.Group{ID: uuid.New()}, nil).Once()
		mockRepo.On("GetStage", mock.Anything, cID, mock.Anything).Return(&domain.TournamentStage{}, nil).Once()
		mockRepo.On("CreateMatchesBatch", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		req, _ := http.NewRequest("POST", "/api/v1/championships/groups/g1/fixture", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})
}

func TestChampionshipHandler_Volunteers(t *testing.T) {
	mockVolRepo := new(MockVolunteerRepo)
	volSvc := application.NewVolunteerService(mockVolRepo)
	h := handler.NewChampionshipHandler(nil, volSvc, nil)
	cID := uuid.New().String()
	uID := uuid.New().String()
	r := setupRouter(h, cID, uID, userDomain.RoleAdmin)

	t.Run("Assign Volunteer", func(t *testing.T) {
		mID := uuid.New()
		mockVolRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		body, _ := json.Marshal(map[string]string{
			"user_id": "vol-1",
			"role":    "REFREE",
		})
		req, _ := http.NewRequest("POST", "/api/v1/championships/matches/"+mID.String()+"/volunteers", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("Get Summary", func(t *testing.T) {
		mID := uuid.New()
		mockVolRepo.On("GetByMatchID", mock.Anything, cID, mID).Return([]domain.VolunteerAssignment{}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/championships/matches/"+mID.String()+"/volunteers", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestChampionshipHandler_Public(t *testing.T) {
	mockRepo := new(MockChampionshipRepo)
	mockClubRepo := new(MockClubRepo)
	uc := application.NewChampionshipUseCases(mockRepo, nil, nil)
	clubSvc := clubApp.NewClubUseCases(nil, mockClubRepo, nil, nil)
	h := handler.NewChampionshipHandler(uc, nil, clubSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	h.RegisterRoutes(api, func(c *gin.Context) {}, func(c *gin.Context) {})

	t.Run("Get Public Tournaments", func(t *testing.T) {
		clubID := uuid.New()
		mockClubRepo.On("GetBySlug", mock.Anything, "my-club").Return(&clubDomain.Club{ID: clubID.String()}, nil).Once()
		mockRepo.On("ListTournaments", mock.Anything, clubID.String()).Return([]domain.Tournament{{Name: "Open"}}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/public/clubs/my-club/championships/", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Open")
	})

	t.Run("Get Public Tournament", func(t *testing.T) {
		clubID := uuid.New()
		tID := uuid.New()
		mockClubRepo.On("GetBySlug", mock.Anything, "my-club").Return(&clubDomain.Club{ID: clubID.String()}, nil).Once()
		mockRepo.On("GetTournament", mock.Anything, clubID.String(), tID.String()).Return(&domain.Tournament{Name: "T1"}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/public/clubs/my-club/championships/"+tID.String(), nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestChampionshipHandler_StagesAndResults(t *testing.T) {
	mockRepo := new(MockChampionshipRepo)
	mockVolunteerSvc := new(MockVolunteerService)
	mockBookingSvc := new(MockBookingService)
	uc := application.NewChampionshipUseCases(mockRepo, mockBookingSvc, nil)
	h := handler.NewChampionshipHandler(uc, mockVolunteerSvc, nil)
	cID := uuid.New().String()
	uID := uuid.New().String()
	r := setupRouter(h, cID, uID, userDomain.RoleAdmin)

	t.Run("Add Stage", func(t *testing.T) {
		tID := uuid.New()
		mockRepo.On("GetTournament", mock.Anything, cID, tID.String()).Return(&domain.Tournament{}, nil).Once()
		mockRepo.On("CreateStage", mock.Anything, mock.Anything).Return(nil).Once()

		body, _ := json.Marshal(application.AddStageInput{Name: "S1", Type: "GROUP"})
		req, _ := http.NewRequest("POST", "/api/v1/championships/"+tID.String()+"/stages", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("Update Result", func(t *testing.T) {
		mID := uuid.New()
		mockRepo.On("UpdateMatchResult", mock.Anything, cID, mID.String(), 2.0, 1.0).Return(nil).Once()
		mockRepo.On("GetMatch", mock.Anything, cID, mID.String()).Return(&domain.TournamentMatch{ID: mID}, nil).Once()
		mockRepo.On("GetTeamMembers", mock.Anything, mock.Anything).Return([]string{}, nil).Maybe()

		body, _ := json.Marshal(application.UpdateMatchResultInput{MatchID: mID.String(), HomeScore: 2.0, AwayScore: 1.0})
		req, _ := http.NewRequest("POST", "/api/v1/championships/matches/result", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Generate Knockout", func(t *testing.T) {
		sID := uuid.New()
		mockRepo.On("GetStage", mock.Anything, cID, sID.String()).Return(&domain.TournamentStage{Type: domain.StageKnockout}, nil).Once()
		mockRepo.On("CreateMatchesBatch", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		body, _ := json.Marshal(application.GenerateKnockoutBracketInput{SeedOrder: []string{uuid.New().String(), uuid.New().String()}})
		req, _ := http.NewRequest("POST", "/api/v1/championships/stages/"+sID.String()+"/knockout", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("Schedule Match", func(t *testing.T) {
		mID := uuid.New()
		bID := uuid.New()
		mockBookingSvc.On("CreateSystemBooking", cID, "court-1", mock.Anything, mock.Anything, mock.Anything).Return(&bID, nil).Once()
		mockRepo.On("UpdateMatchScheduling", mock.Anything, cID, mID.String(), mock.Anything, bID).Return(nil).Once()

		body, _ := json.Marshal(application.ScheduleMatchInput{MatchID: mID.String(), CourtID: "court-1"})
		req, _ := http.NewRequest("POST", "/api/v1/championships/matches/schedule", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Get Volunteers", func(t *testing.T) {
		mID := uuid.New()
		mockVolunteerSvc.On("GetVolunteerSummary", mock.Anything, cID, mID).Return(&domain.VolunteerSummary{MatchID: mID}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/championships/matches/"+mID.String()+"/volunteers", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Remove Volunteer", func(t *testing.T) {
		vID := uuid.New()
		mockVolunteerSvc.On("RemoveVolunteer", mock.Anything, cID, vID).Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/api/v1/championships/volunteers/"+vID.String(), nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	t.Run("Add Group", func(t *testing.T) {
		sID := uuid.New()
		mockRepo.On("GetStage", mock.Anything, cID, sID.String()).Return(&domain.TournamentStage{}, nil).Once()
		mockRepo.On("CreateGroup", mock.Anything, mock.Anything).Return(nil).Once()

		body, _ := json.Marshal(application.AddGroupInput{Name: "Group A"})
		req, _ := http.NewRequest("POST", "/api/v1/championships/stages/"+sID.String()+"/groups", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})
}

func TestChampionshipHandler_DetailedErrors(t *testing.T) {
	mockRepo := new(MockChampionshipRepo)
	mockVolunteerSvc := new(MockVolunteerService)
	mockBookingSvc := new(MockBookingService)
	uc := application.NewChampionshipUseCases(mockRepo, mockBookingSvc, nil)
	h := handler.NewChampionshipHandler(uc, mockVolunteerSvc, nil)
	cID := uuid.New().String()
	uID := uuid.New().String()
	r := setupRouter(h, cID, uID, userDomain.RoleAdmin)

	// Add Stage Errors
	t.Run("AddStage_InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/championships/any/stages", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("AddStage_InternalError", func(t *testing.T) {
		tID := uuid.New()
		mockRepo.On("GetTournament", mock.Anything, cID, tID.String()).Return(nil, errors.New("db error")).Once()
		body, _ := json.Marshal(application.AddStageInput{Name: "S1"})
		req, _ := http.NewRequest("POST", "/api/v1/championships/"+tID.String()+"/stages", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	// Add Group Errors
	t.Run("AddGroup_InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/championships/stages/any/groups", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("AddGroup_UseCasesError", func(t *testing.T) {
		sID := uuid.New()
		mockRepo.On("GetStage", mock.Anything, cID, sID.String()).Return(nil, errors.New("db error")).Once()
		body, _ := json.Marshal(application.AddGroupInput{Name: "G1"})
		req, _ := http.NewRequest("POST", "/api/v1/championships/stages/"+sID.String()+"/groups", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	// Register Team Errors
	t.Run("RegisterTeam_InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/championships/groups/any/teams", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("RegisterTeam_UseCasesError", func(t *testing.T) {
		gID := uuid.New()
		mockRepo.On("GetGroup", mock.Anything, cID, gID.String()).Return(nil, errors.New("fail")).Once()
		body, _ := json.Marshal(application.RegisterTeamInput{TeamID: uuid.New().String()})
		req, _ := http.NewRequest("POST", "/api/v1/championships/groups/"+gID.String()+"/teams", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	// Generate Fixture Errors
	t.Run("GenerateFixture_Error", func(t *testing.T) {
		gID := uuid.New()
		mockRepo.On("GetStandings", mock.Anything, cID, gID.String()).Return(nil, errors.New("fail")).Once()
		req, _ := http.NewRequest("POST", "/api/v1/championships/groups/"+gID.String()+"/fixture", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	// Update Match Result Errors
	t.Run("UpdateMatchResult_InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/championships/matches/result", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("UpdateMatchResult_ServiceError", func(t *testing.T) {
		mID := uuid.New()
		mockRepo.On("UpdateMatchResult", mock.Anything, cID, mID.String(), 1.0, 1.0).Return(errors.New("fail")).Once()
		body, _ := json.Marshal(application.UpdateMatchResultInput{MatchID: mID.String(), HomeScore: 1.0, AwayScore: 1.0})
		req, _ := http.NewRequest("POST", "/api/v1/championships/matches/result", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	// Schedule Match Errors
	t.Run("ScheduleMatch_InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/championships/matches/schedule", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	// Assign Volunteer Errors
	t.Run("AssignVolunteer_InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/championships/matches/invalid/volunteers", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("AssignVolunteer_InvalidJSON", func(t *testing.T) {
		mID := uuid.New()
		req, _ := http.NewRequest("POST", "/api/v1/championships/matches/"+mID.String()+"/volunteers", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("AssignVolunteer_ServiceError", func(t *testing.T) {
		mID := uuid.New()
		uID := "user-x"
		// Corrected mock call signature:
		// AssignVolunteer(ctx, clubID, matchID, userID, role, assignedBy, notes)
		mockVolunteerSvc.On("AssignVolunteer", mock.Anything, cID, mID, uID, domain.VolunteerRole("REFREE"), mock.Anything, "test").Return(errors.New("fail")).Once()

		body := map[string]interface{}{"user_id": uID, "role": "REFREE", "notes": "test"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/championships/matches/"+mID.String()+"/volunteers", bytes.NewBuffer(jsonBody))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	// Remove Volunteer Errors
	t.Run("RemoveVolunteer_InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/championships/volunteers/invalid", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("RemoveVolunteer_ServiceError", func(t *testing.T) {
		vID := uuid.New()
		mockVolunteerSvc.On("RemoveVolunteer", mock.Anything, cID, vID).Return(errors.New("fail")).Once()
		req, _ := http.NewRequest("DELETE", "/api/v1/championships/volunteers/"+vID.String(), nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	// List Tournaments Error
	t.Run("ListTournaments_Error", func(t *testing.T) {
		mockRepo.On("ListTournaments", mock.Anything, cID).Return(nil, errors.New("fail")).Once()
		req, _ := http.NewRequest("GET", "/api/v1/championships/", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	// Get Standings Error
	t.Run("GetStandings_Error", func(t *testing.T) {
		gID := uuid.New()
		mockRepo.On("GetStandings", mock.Anything, cID, gID.String()).Return(nil, errors.New("fail")).Once()
		req, _ := http.NewRequest("GET", "/api/v1/championships/groups/"+gID.String()+"/standings", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	// Get Matches By Group Error
	t.Run("GetMatchesByGroup_Error", func(t *testing.T) {
		gID := uuid.New()
		mockRepo.On("GetMatchesByGroup", mock.Anything, cID, gID.String()).Return(nil, errors.New("fail")).Once()
		req, _ := http.NewRequest("GET", "/api/v1/championships/groups/"+gID.String()+"/matches", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestChampionshipHandler_Errors(t *testing.T) {
	h := handler.NewChampionshipHandler(nil, nil, nil)
	r := setupRouter(h, "c1", "u1", userDomain.RoleMember) // MEMBER role

	t.Run("RBAC Denied", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/championships/", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		rAdmin := setupRouter(h, "c1", "u1", userDomain.RoleAdmin)
		req, _ := http.NewRequest("POST", "/api/v1/championships/", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		rAdmin.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
