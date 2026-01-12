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
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	handler "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, clubID, id string) (*domain.User, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, clubID, id string) error {
	args := m.Called(ctx, clubID, id)
	return args.Error(0)
}

func (m *MockUserRepo) List(ctx context.Context, clubID string, limit, offset int, filters map[string]interface{}) ([]domain.User, error) {
	args := m.Called(ctx, clubID, limit, offset, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepo) ListByIDs(ctx context.Context, clubID string, ids []string) ([]domain.User, error) {
	args := m.Called(ctx, clubID, ids)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepo) FindChildren(ctx context.Context, clubID, parentID string) ([]domain.User, error) {
	args := m.Called(ctx, clubID, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) CreateIncident(ctx context.Context, incident *domain.IncidentLog) error {
	args := m.Called(ctx, incident)
	return args.Error(0)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) AnonymizeForGDPR(ctx context.Context, clubID, id string) error {
	args := m.Called(ctx, clubID, id)
	return args.Error(0)
}

type MockFamilyGroupRepo struct {
	mock.Mock
}

func (m *MockFamilyGroupRepo) Create(ctx context.Context, group *domain.FamilyGroup) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockFamilyGroupRepo) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.FamilyGroup, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FamilyGroup), args.Error(1)
}

func (m *MockFamilyGroupRepo) GetByHeadUserID(ctx context.Context, clubID, headUserID string) (*domain.FamilyGroup, error) {
	args := m.Called(ctx, clubID, headUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FamilyGroup), args.Error(1)
}

func (m *MockFamilyGroupRepo) GetByMemberID(ctx context.Context, clubID, userID string) (*domain.FamilyGroup, error) {
	args := m.Called(ctx, clubID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FamilyGroup), args.Error(1)
}

func (m *MockFamilyGroupRepo) AddMember(ctx context.Context, clubID string, groupID uuid.UUID, userID string) error {
	args := m.Called(ctx, clubID, groupID, userID)
	return args.Error(0)
}

func (m *MockFamilyGroupRepo) RemoveMember(ctx context.Context, clubID string, groupID uuid.UUID, userID string) error {
	args := m.Called(ctx, clubID, groupID, userID)
	return args.Error(0)
}

// --- Test Setup ---

func setupRouter(h *handler.UserHandler, clubID, userID, role string) *gin.Engine {
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

	handler.RegisterRoutes(api, h, auth, tenant)
	handler.RegisterPublicRoutes(api, h)

	return r
}

// --- Tests ---

func TestUserHandler_Basics(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := application.NewUserUseCases(mockRepo, nil)
	h := handler.NewUserHandler(uc)
	r := setupRouter(h, "club-1", "user-1", domain.RoleMember)

	t.Run("Get /me", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "club-1", "user-1").Return(&domain.User{ID: "user-1", Role: "MEMBER"}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Update /me", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "club-1", "user-1").Return(&domain.User{ID: "user-1"}, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		body, _ := json.Marshal(application.UpdateProfileDTO{Name: "Neo"})
		req, _ := http.NewRequest("PUT", "/api/v1/users/me", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestUserHandler_Children(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := application.NewUserUseCases(mockRepo, nil)
	h := handler.NewUserHandler(uc)
	r := setupRouter(h, "club-1", "user-1", domain.RoleMember)

	t.Run("Get Children", func(t *testing.T) {
		mockRepo.On("FindChildren", mock.Anything, "club-1", "user-1").Return([]domain.User{{ID: "c1"}}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/users/me/children", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Register Child", func(t *testing.T) {
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		body, _ := json.Marshal(application.RegisterChildDTO{Name: "Kid"})
		req, _ := http.NewRequest("POST", "/api/v1/users/me/children", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})
}

func TestUserHandler_StatsAndWallet(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := application.NewUserUseCases(mockRepo, nil)
	h := handler.NewUserHandler(uc)
	r := setupRouter(h, "club-1", "user-1", domain.RoleMember)

	t.Run("Get Stats", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "club-1", "user-1").Return(&domain.User{ID: "user-1", Stats: &domain.UserStats{}}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/users/me/stats", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Get Wallet", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "club-1", "user-1").Return(&domain.User{ID: "user-1", Wallet: &domain.Wallet{}}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/users/me/wallet", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestUserHandler_EmergencyAndIncidents(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := application.NewUserUseCases(mockRepo, nil)
	h := handler.NewUserHandler(uc)
	r := setupRouter(h, "club-1", "user-1", domain.RoleMember)

	t.Run("Update Emergency", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "club-1", "user-1").Return(&domain.User{ID: "user-1"}, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.EmergencyContactName == "Mom"
		})).Return(nil).Once()

		body, _ := json.Marshal(map[string]string{"contact_name": "Mom"})
		req, _ := http.NewRequest("PUT", "/api/v1/users/me/emergency", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Log Incident", func(t *testing.T) {
		mockRepo.On("CreateIncident", mock.Anything, mock.Anything).Return(nil).Once()
		body, _ := json.Marshal(map[string]string{"description": "Test Incident"})
		req, _ := http.NewRequest("POST", "/api/v1/users/incidents", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})
}

func TestUserHandler_FamilyGroups(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockFamilyRepo := new(MockFamilyGroupRepo)
	uc := application.NewUserUseCases(mockRepo, mockFamilyRepo)
	h := handler.NewUserHandler(uc)
	r := setupRouter(h, "club-1", "user-1", domain.RoleMember)

	t.Run("Create and Get Family", func(t *testing.T) {
		mockFamilyRepo.On("GetByHeadUserID", mock.Anything, "club-1", "user-1").Return(nil, nil).Once()
		mockFamilyRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		mockFamilyRepo.On("AddMember", mock.Anything, "club-1", mock.Anything, "user-1").Return(nil).Once()

		body, _ := json.Marshal(map[string]string{"name": "Vader Family"})
		req, _ := http.NewRequest("POST", "/api/v1/users/family-groups", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)

		mockFamilyRepo.On("GetByMemberID", mock.Anything, "club-1", "user-1").Return(&domain.FamilyGroup{Name: "Vader Family"}, nil).Once()
		req2, _ := http.NewRequest("GET", "/api/v1/users/family-groups/me", nil)
		resp2 := httptest.NewRecorder()
		r.ServeHTTP(resp2, req2)
		assert.Equal(t, http.StatusOK, resp2.Code)
	})

	t.Run("Add Family Member Secure", func(t *testing.T) {
		gID := uuid.New()
		mockFamilyRepo.On("GetByID", mock.Anything, "club-1", gID).Return(&domain.FamilyGroup{HeadUserID: "user-1"}, nil).Once()
		mockFamilyRepo.On("AddMember", mock.Anything, "club-1", gID, "user-2").Return(nil).Once()

		body, _ := json.Marshal(map[string]string{"user_id": "user-2"})
		req, _ := http.NewRequest("POST", "/api/v1/users/family-groups/"+gID.String()+"/members", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNoContent, resp.Code)
	})
}

func TestUserHandler_GDPR(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := application.NewUserUseCases(mockRepo, nil)
	h := handler.NewUserHandler(uc)
	r := setupRouter(h, "club-1", "user-1", domain.RoleMember)

	t.Run("Export Data", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "club-1", "user-1").Return(&domain.User{ID: "user-1"}, nil).Once()
		mockRepo.On("FindChildren", mock.Anything, "club-1", "user-1").Return([]domain.User{}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/users/me/data-export", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Erasure", func(t *testing.T) {
		mockRepo.On("AnonymizeForGDPR", mock.Anything, "club-1", "user-1").Return(nil).Once()
		req, _ := http.NewRequest("DELETE", "/api/v1/users/me/gdpr-erasure", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestUserHandler_Public(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := application.NewUserUseCases(mockRepo, nil)
	h := handler.NewUserHandler(uc)
	r := setupRouter(h, "club-1", "user-1", domain.RoleMember)

	t.Run("Register Dependent Public", func(t *testing.T) {
		mockRepo.On("GetByEmail", mock.Anything, "dad@me.com").Return(nil, nil).Once()
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Twice()
		body, _ := json.Marshal(application.RegisterDependentDTO{
			ParentEmail: "dad@me.com", ParentName: "Dad", ChildName: "Junior",
		})
		req, _ := http.NewRequest("POST", "/api/v1/users/public/register-dependent?club_id=club-1", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})
}
