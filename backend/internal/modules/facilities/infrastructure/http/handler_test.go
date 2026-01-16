package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	handler "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/http"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockFacilityRepo struct {
	mock.Mock
}

func (m *MockFacilityRepo) Create(ctx context.Context, facility *domain.Facility) error {
	args := m.Called(ctx, facility)
	return args.Error(0)
}
func (m *MockFacilityRepo) GetByID(ctx context.Context, clubID, id string) (*domain.Facility, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Facility), args.Error(1)
}
func (m *MockFacilityRepo) GetByIDForUpdate(ctx context.Context, clubID, id string) (*domain.Facility, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Facility), args.Error(1)
}
func (m *MockFacilityRepo) List(ctx context.Context, clubID string, limit, offset int) ([]*domain.Facility, error) {
	args := m.Called(ctx, clubID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Facility), args.Error(1)
}
func (m *MockFacilityRepo) Update(ctx context.Context, facility *domain.Facility) error {
	args := m.Called(ctx, facility)
	return args.Error(0)
}
func (m *MockFacilityRepo) HasConflict(ctx context.Context, clubID, facilityID string, startTime, endTime time.Time) (bool, error) {
	args := m.Called(ctx, clubID, facilityID, startTime, endTime)
	return args.Bool(0), args.Error(1)
}
func (m *MockFacilityRepo) CreateEquipment(ctx context.Context, equipment *domain.Equipment) error {
	args := m.Called(ctx, equipment)
	return args.Error(0)
}
func (m *MockFacilityRepo) GetEquipmentByID(ctx context.Context, id string) (*domain.Equipment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Equipment), args.Error(1)
}
func (m *MockFacilityRepo) ListEquipmentByFacility(ctx context.Context, facilityID string) ([]*domain.Equipment, error) {
	args := m.Called(ctx, facilityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Equipment), args.Error(1)
}
func (m *MockFacilityRepo) UpdateEquipment(ctx context.Context, equipment *domain.Equipment) error {
	args := m.Called(ctx, equipment)
	return args.Error(0)
}
func (m *MockFacilityRepo) LoanEquipmentAtomic(ctx context.Context, loan *domain.EquipmentLoan, equipmentID string) error {
	args := m.Called(ctx, loan, equipmentID)
	return args.Error(0)
}
func (m *MockFacilityRepo) SemanticSearch(ctx context.Context, clubID string, embedding []float32, limit int) ([]*domain.FacilityWithSimilarity, error) {
	args := m.Called(ctx, clubID, embedding, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.FacilityWithSimilarity), args.Error(1)
}
func (m *MockFacilityRepo) UpdateEmbedding(ctx context.Context, facilityID string, embedding []float32) error {
	args := m.Called(ctx, facilityID, embedding)
	return args.Error(0)
}
func (m *MockFacilityRepo) ListMaintenanceByFacility(ctx context.Context, facilityID string) ([]*domain.MaintenanceTask, error) {
	args := m.Called(ctx, facilityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.MaintenanceTask), args.Error(1)
}

type MockLoanRepo struct {
	mock.Mock
}

func (m *MockLoanRepo) Create(ctx context.Context, loan *domain.EquipmentLoan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}
func (m *MockLoanRepo) GetByID(ctx context.Context, id string) (*domain.EquipmentLoan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EquipmentLoan), args.Error(1)
}
func (m *MockLoanRepo) ListByUser(ctx context.Context, userID string) ([]*domain.EquipmentLoan, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.EquipmentLoan), args.Error(1)
}
func (m *MockLoanRepo) ListByStatus(ctx context.Context, status domain.LoanStatus) ([]*domain.EquipmentLoan, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.EquipmentLoan), args.Error(1)
}
func (m *MockLoanRepo) Update(ctx context.Context, loan *domain.EquipmentLoan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

// --- Test Setup ---

func setupRouter(h *handler.FacilityHandler, sh *handler.SearchHandler, clubID string, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("clubID", clubID)
		c.Set("userRole", role)
		c.Next()
	})

	api := r.Group("/api/v1")
	if h != nil {
		handler.RegisterRoutes(api, h, func(c *gin.Context) {}, func(c *gin.Context) {})
	}
	if sh != nil {
		handler.RegisterSearchRoutes(api, sh)
	}

	return r
}

// --- Tests ---

func TestFacilityHandler_Create(t *testing.T) {
	mockRepo := new(MockFacilityRepo)
	mockLoanRepo := new(MockLoanRepo)
	uc := application.NewFacilityUseCases(mockRepo, mockLoanRepo)
	h := handler.NewFacilityHandler(uc)

	t.Run("Success as Admin", func(t *testing.T) {
		r := setupRouter(h, nil, "club-1", userDomain.RoleAdmin)

		input := application.CreateFacilityDTO{
			Name:       "New Court",
			Type:       domain.FacilityTypeCourt,
			Capacity:   4,
			HourlyRate: 20.0,
		}
		body, _ := json.Marshal(input)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Facility")).Return(nil).Once()

		req, _ := http.NewRequest("POST", "/api/v1/facilities", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestFacilityHandler_Equipment(t *testing.T) {
	mockRepo := new(MockFacilityRepo)
	mockLoanRepo := new(MockLoanRepo)
	uc := application.NewFacilityUseCases(mockRepo, mockLoanRepo)
	h := handler.NewFacilityHandler(uc)
	r := setupRouter(h, nil, "club-1", userDomain.RoleAdmin)

	t.Run("AddEquipment", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "club-1", "fac-1").Return(&domain.Facility{ID: "fac-1"}, nil).Once()
		mockRepo.On("CreateEquipment", mock.Anything, mock.Anything).Return(nil).Once()

		input := application.AddEquipmentDTO{
			Name:      "Ball",
			Type:      "Sports",
			Condition: domain.EquipmentConditionExcellent,
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/api/v1/facilities/fac-1/equipment", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("LoanEquipment", func(t *testing.T) {
		mockRepo.On("GetEquipmentByID", mock.Anything, "eq-1").Return(&domain.Equipment{
			ID: "eq-1", FacilityID: "fac-1", Status: "available",
		}, nil).Once()
		mockRepo.On("GetByID", mock.Anything, "club-1", "fac-1").Return(&domain.Facility{
			ID: "fac-1", ClubID: "club-1",
		}, nil).Once()
		mockRepo.On("LoanEquipmentAtomic", mock.Anything, mock.Anything, "eq-1").Return(nil).Once()

		input := map[string]string{
			"user_id":         "u1",
			"expected_return": time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/api/v1/facilities/equipment/eq-1/loan", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})
}

func TestSearchHandler(t *testing.T) {
	mockRepo := new(MockFacilityRepo)
	suc := application.NewSemanticSearchUseCase(mockRepo)
	sh := handler.NewSearchHandler(suc)
	r := setupRouter(nil, sh, "club-1", userDomain.RoleMember)

	t.Run("Search Success", func(t *testing.T) {
		mockRepo.On("SemanticSearch", mock.Anything, "club-1", mock.Anything, 10).Return([]*domain.FacilityWithSimilarity{
			{Facility: &domain.Facility{Name: "Result"}, Similarity: 0.9},
		}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/facilities/search?q=tennis", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Result")
	})
}
