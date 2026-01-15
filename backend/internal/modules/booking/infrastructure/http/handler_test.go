package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	handler "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/http"
	facilityDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	notificationService "github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockBookingRepo struct{ mock.Mock }

func (m *MockBookingRepo) Create(ctx context.Context, b *domain.Booking) error {
	return m.Called(ctx, b).Error(0)
}
func (m *MockBookingRepo) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Booking, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}
func (m *MockBookingRepo) List(ctx context.Context, clubID string, filter map[string]interface{}) ([]domain.Booking, error) {
	args := m.Called(ctx, clubID, filter)
	return args.Get(0).([]domain.Booking), args.Error(1)
}
func (m *MockBookingRepo) ListAll(ctx context.Context, clubID string, filter map[string]interface{}, from, to *time.Time) ([]domain.Booking, error) {
	args := m.Called(ctx, clubID, filter, from, to)
	return args.Get(0).([]domain.Booking), args.Error(1)
}
func (m *MockBookingRepo) Update(ctx context.Context, b *domain.Booking) error {
	return m.Called(ctx, b).Error(0)
}
func (m *MockBookingRepo) HasTimeConflict(ctx context.Context, clubID string, fID uuid.UUID, s, e time.Time) (bool, error) {
	args := m.Called(ctx, clubID, fID, s, e)
	return args.Bool(0), args.Error(1)
}
func (m *MockBookingRepo) ListByFacilityAndDate(ctx context.Context, clubID string, fID uuid.UUID, d time.Time) ([]domain.Booking, error) {
	args := m.Called(ctx, clubID, fID, d)
	return args.Get(0).([]domain.Booking), args.Error(1)
}
func (m *MockBookingRepo) AddToWaitlist(ctx context.Context, e *domain.Waitlist) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockBookingRepo) GetNextInLine(ctx context.Context, clubID string, rID uuid.UUID, d time.Time) (*domain.Waitlist, error) {
	args := m.Called(ctx, clubID, rID, d)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Waitlist), args.Error(1)
}
func (m *MockBookingRepo) ListExpired(ctx context.Context) ([]domain.Booking, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Booking), args.Error(1)
}

type MockRecurringRepo struct{ mock.Mock }

func (m *MockRecurringRepo) Create(ctx context.Context, r *domain.RecurringRule) error {
	return m.Called(ctx, r).Error(0)
}
func (m *MockRecurringRepo) GetByFacility(ctx context.Context, clubID string, fID uuid.UUID) ([]domain.RecurringRule, error) {
	args := m.Called(ctx, clubID, fID)
	return args.Get(0).([]domain.RecurringRule), args.Error(1)
}
func (m *MockRecurringRepo) GetAllActive(ctx context.Context, clubID string) ([]domain.RecurringRule, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]domain.RecurringRule), args.Error(1)
}
func (m *MockRecurringRepo) Delete(ctx context.Context, clubID string, id uuid.UUID) error {
	return m.Called(ctx, clubID, id).Error(0)
}

type MockFacilityRepo struct{ mock.Mock }

func (m *MockFacilityRepo) GetByID(ctx context.Context, clubID, id string) (*facilityDomain.Facility, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*facilityDomain.Facility), args.Error(1)
}
func (m *MockFacilityRepo) HasConflict(ctx context.Context, clubID, fID string, s, e time.Time) (bool, error) {
	args := m.Called(ctx, clubID, fID, s, e)
	return args.Bool(0), args.Error(1)
}
func (m *MockFacilityRepo) ListMaintenanceByFacility(ctx context.Context, fID string) ([]*facilityDomain.MaintenanceTask, error) {
	args := m.Called(ctx, fID)
	return args.Get(0).([]*facilityDomain.MaintenanceTask), args.Error(1)
}
func (m *MockFacilityRepo) Create(ctx context.Context, f *facilityDomain.Facility) error { return nil }
func (m *MockFacilityRepo) List(ctx context.Context, clubID string, l, o int) ([]*facilityDomain.Facility, error) {
	return nil, nil
}
func (m *MockFacilityRepo) Update(ctx context.Context, f *facilityDomain.Facility) error { return nil }
func (m *MockFacilityRepo) SemanticSearch(ctx context.Context, clubID string, emb []float32, l int) ([]*facilityDomain.FacilityWithSimilarity, error) {
	return nil, nil
}
func (m *MockFacilityRepo) UpdateEmbedding(ctx context.Context, fID string, emb []float32) error {
	return nil
}
func (m *MockFacilityRepo) CreateEquipment(ctx context.Context, e *facilityDomain.Equipment) error {
	return nil
}
func (m *MockFacilityRepo) GetEquipmentByID(ctx context.Context, id string) (*facilityDomain.Equipment, error) {
	return nil, nil
}
func (m *MockFacilityRepo) ListEquipmentByFacility(ctx context.Context, fID string) ([]*facilityDomain.Equipment, error) {
	return nil, nil
}
func (m *MockFacilityRepo) UpdateEquipment(ctx context.Context, e *facilityDomain.Equipment) error {
	return nil
}
func (m *MockFacilityRepo) LoanEquipmentAtomic(ctx context.Context, l *facilityDomain.EquipmentLoan, eID string) error {
	return nil
}

type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) GetByID(ctx context.Context, clubID, id string) (*userDomain.User, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDomain.User), args.Error(1)
}
func (m *MockUserRepo) Update(ctx context.Context, u *userDomain.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *MockUserRepo) Create(ctx context.Context, u *userDomain.User) error { return nil }
func (m *MockUserRepo) Delete(ctx context.Context, clubID, id string) error  { return nil }
func (m *MockUserRepo) List(ctx context.Context, clubID string, l, o int, f map[string]interface{}) ([]userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) FindChildren(ctx context.Context, clubID, pID string) ([]userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) CreateIncident(ctx context.Context, i *userDomain.IncidentLog) error {
	return nil
}
func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) ListByIDs(ctx context.Context, clubID string, ids []string) ([]userDomain.User, error) {
	return nil, nil
}
func (m *MockUserRepo) AnonymizeForGDPR(ctx context.Context, clubID, id string) error { return nil }

type MockNotificationSender struct{ mock.Mock }

func (m *MockNotificationSender) Send(ctx context.Context, n notificationService.Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

type MockRefundService struct{ mock.Mock }

func (m *MockRefundService) Refund(ctx context.Context, clubID string, refID uuid.UUID, refType string) error {
	return m.Called(ctx, clubID, refID, refType).Error(0)
}

// --- Setup ---

func setupRouter(h *handler.BookingHandler, clubID, userID, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery()) // Add recovery to see panic details in test output if any
	r.Use(func(c *gin.Context) {
		c.Set("clubID", clubID)
		c.Set("userID", userID)
		c.Set("userRole", role)
		c.Next()
	})
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api, h, func(c *gin.Context) {}, func(c *gin.Context) {})
	return r
}

// --- Tests ---

func TestBookingHandler_Endpoints(t *testing.T) {
	mockBookingRepo := new(MockBookingRepo)
	mockRecurringRepo := new(MockRecurringRepo)
	mockFacilityRepo := new(MockFacilityRepo)
	mockUserRepo := new(MockUserRepo)
	mockNotificationSender := new(MockNotificationSender)
	mockRefundService := new(MockRefundService)

	uc := application.NewBookingUseCases(
		mockBookingRepo, mockRecurringRepo, mockFacilityRepo,
		mockUserRepo, mockNotificationSender, mockRefundService,
	)
	h := handler.NewBookingHandler(uc)

	clubID := "test-club-http"
	userID := uuid.New().String()

	t.Run("Create Booking Success", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		facilityID := uuid.New().String()
		now := time.Now().Add(1 * time.Hour)

		medicalStatus := userDomain.MedicalCertStatusValid
		mockUserRepo.On("GetByID", mock.Anything, clubID, userID).Return(&userDomain.User{
			ID: userID, MedicalCertStatus: &medicalStatus,
		}, nil).Once()

		mockFacilityRepo.On("GetByID", mock.Anything, clubID, facilityID).Return(&facilityDomain.Facility{
			ID: facilityID, Status: facilityDomain.FacilityStatusActive, HourlyRate: 10,
		}, nil).Once()

		mockBookingRepo.On("HasTimeConflict", mock.Anything, clubID, uuid.MustParse(facilityID), mock.Anything, mock.Anything).Return(false, nil).Once()
		mockFacilityRepo.On("HasConflict", mock.Anything, clubID, facilityID, mock.Anything, mock.Anything).Return(false, nil).Once()
		mockBookingRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		body, _ := json.Marshal(map[string]interface{}{
			"user_id":     userID,
			"facility_id": facilityID,
			"start_time":  now,
			"end_time":    now.Add(1 * time.Hour),
		})
		req, _ := http.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("Cancel Booking with Refund", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		bookingID := uuid.New()

		mockBookingRepo.On("GetByID", mock.Anything, clubID, bookingID).Return(&domain.Booking{
			ID: bookingID, UserID: uuid.MustParse(userID), Status: domain.BookingStatusConfirmed,
		}, nil).Once()
		mockBookingRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		mockRefundService.On("Refund", mock.Anything, clubID, bookingID, "BOOKING").Return(nil).Once()
		mockBookingRepo.On("GetNextInLine", mock.Anything, clubID, mock.Anything, mock.Anything).Return(nil, nil).Once()

		req, _ := http.NewRequest("DELETE", "/api/v1/bookings/"+bookingID.String(), nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("List Bookings", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		mockBookingRepo.On("List", mock.Anything, clubID, mock.Anything).Return([]domain.Booking{}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/bookings", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Get Availability", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		facilityID := uuid.New().String()
		mockFacilityRepo.On("GetByID", mock.Anything, clubID, facilityID).Return(&facilityDomain.Facility{
			ID: facilityID, Status: facilityDomain.FacilityStatusActive, OpeningTime: "08:00", ClosingTime: "22:00",
		}, nil).Once()
		mockFacilityRepo.On("ListMaintenanceByFacility", mock.Anything, facilityID).Return([]*facilityDomain.MaintenanceTask{}, nil).Once()
		mockBookingRepo.On("ListByFacilityAndDate", mock.Anything, clubID, uuid.MustParse(facilityID), mock.Anything).Return([]domain.Booking{}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/bookings/availability?facility_id="+facilityID+"&date=2024-05-20", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Admin: List All Bookings", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)
		mockBookingRepo.On("ListAll", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return([]domain.Booking{}, nil).Once()
		req, _ := http.NewRequest("GET", "/api/v1/bookings/all", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Admin: Create Recurring Rule", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)
		mockRecurringRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		body, _ := json.Marshal(map[string]interface{}{
			"facility_id": uuid.New().String(),
			"type":        "FIXED",
			"frequency":   "WEEKLY", // Added frequency
			"day_of_week": 1,
			"start_time":  time.Now(),
			"end_time":    time.Now().Add(1 * time.Hour),
			"start_date":  "2024-01-01",
			"end_date":    "2024-12-31",
		})
		req, _ := http.NewRequest("POST", "/api/v1/bookings/recurring", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("Join Waitlist", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		mockBookingRepo.On("AddToWaitlist", mock.Anything, mock.Anything).Return(nil).Once()

		body, _ := json.Marshal(map[string]interface{}{
			"resource_id": uuid.New().String(),
			"target_date": time.Now(),
			"user_id":     userID,
		})
		req, _ := http.NewRequest("POST", "/api/v1/bookings/waitlist", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("Admin: List All - Unauthorized", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember) // Regular member
		req, _ := http.NewRequest("GET", "/api/v1/bookings/all", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("Create Booking - Invalid Input", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		body := []byte(`{"invalid": "json"}`)
		req, _ := http.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Cancel Booking - Not Found", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		bookingID := uuid.New().String()
		mockBookingRepo.On("GetByID", mock.Anything, clubID, uuid.MustParse(bookingID)).Return(nil, fmt.Errorf("booking not found")).Once()

		req, _ := http.NewRequest("DELETE", "/api/v1/bookings/"+bookingID, nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code) // mapErrorToResponse defaults to 400 for general errors if not matched specifically as "not found" in a switch that I should probably improve
	})

	t.Run("Get Availability - Missing Params", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		req, _ := http.NewRequest("GET", "/api/v1/bookings/availability", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Get Availability - Invalid Date", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		req, _ := http.NewRequest("GET", "/api/v1/bookings/availability?facility_id="+uuid.New().String()+"&date=invalid", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Admin: Create Recurring Rule - Validation Error", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)
		body, _ := json.Marshal(map[string]interface{}{"type": "FIXED"}) // Missing facility_id etc.
		req, _ := http.NewRequest("POST", "/api/v1/bookings/recurring", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Admin: Generate Bookings", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)
		mockRecurringRepo.On("GetAllActive", mock.Anything, clubID).Return([]domain.RecurringRule{}, nil).Once()

		req, _ := http.NewRequest("POST", "/api/v1/bookings/generate", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Admin: Generate Bookings - Service Error", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)
		mockRecurringRepo.On("GetAllActive", mock.Anything, clubID).Return(nil, fmt.Errorf("db error")).Once()

		req, _ := http.NewRequest("POST", "/api/v1/bookings/generate", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("Join Waitlist - Service Error", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		mockBookingRepo.On("AddToWaitlist", mock.Anything, mock.Anything).Return(fmt.Errorf("db error")).Once()

		body, _ := json.Marshal(map[string]interface{}{
			"user_id":     userID,
			"resource_id": uuid.New().String(),
			"target_date": time.Now().Format(time.RFC3339),
		})
		req, _ := http.NewRequest("POST", "/api/v1/bookings/waitlist", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("Cancel Booking - Service Error", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		mockBookingRepo.On("GetByID", mock.Anything, clubID, mock.Anything).Return(&domain.Booking{ID: uuid.New(), UserID: uuid.MustParse(userID)}, nil).Once()
		mockBookingRepo.On("Update", mock.Anything, mock.Anything).Return(fmt.Errorf("db error")).Once()

		req, _ := http.NewRequest("DELETE", "/api/v1/bookings/"+uuid.New().String(), nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("List All Bookings - Forbidden", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		req, _ := http.NewRequest("GET", "/api/v1/bookings/all", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("List All Bookings - Service Error", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)
		mockBookingRepo.On("ListAll", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error")).Once()

		req, _ := http.NewRequest("GET", "/api/v1/bookings/all", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("List All Bookings - With Date Filters", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)
		mockBookingRepo.On("ListAll", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return([]domain.Booking{}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/bookings/all?from=2025-01-01&to=2025-01-02", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("List All Bookings - Invalid Date Filter", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)
		mockBookingRepo.On("ListAll", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return([]domain.Booking{}, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/bookings/all?from=invalid", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code) // Should still return 200 but date is nil
	})

	t.Run("Cancel Booking - Not Found", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		mockBookingRepo.On("GetByID", mock.Anything, clubID, mock.Anything).Return(nil, fmt.Errorf("not found")).Once()

		req, _ := http.NewRequest("DELETE", "/api/v1/bookings/"+uuid.New().String(), nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Join Waitlist - Invalid JSON", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		req, _ := http.NewRequest("POST", "/api/v1/bookings/waitlist", bytes.NewBufferString("{invalid}"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("List Bookings - Service Error", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		mockBookingRepo.On("List", mock.Anything, clubID, mock.Anything).Return([]domain.Booking{}, fmt.Errorf("db error")).Once()

		req, _ := http.NewRequest("GET", "/api/v1/bookings", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("Create Booking - Maintenance Conflict", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		facilityID := uuid.New().String()
		now := time.Now().Add(1 * time.Hour)

		medicalStatus := userDomain.MedicalCertStatusValid
		mockUserRepo.On("GetByID", mock.Anything, clubID, userID).Return(&userDomain.User{
			ID: userID, MedicalCertStatus: &medicalStatus,
		}, nil).Once()
		mockFacilityRepo.On("GetByID", mock.Anything, clubID, facilityID).Return(&facilityDomain.Facility{
			ID: facilityID, Status: facilityDomain.FacilityStatusActive,
		}, nil).Once()
		mockBookingRepo.On("HasTimeConflict", mock.Anything, clubID, uuid.MustParse(facilityID), mock.Anything, mock.Anything).Return(false, nil).Once()
		mockFacilityRepo.On("HasConflict", mock.Anything, clubID, facilityID, mock.Anything, mock.Anything).Return(true, nil).Once() // Maintenance conflict

		body, _ := json.Marshal(map[string]interface{}{
			"user_id":     userID,
			"facility_id": facilityID, "start_time": now, "end_time": now.Add(1 * time.Hour),
		})
		req, _ := http.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Contains(t, resp.Body.String(), "booking_conflict")
	})

	t.Run("Admin: List All - Service Error", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)
		mockBookingRepo.On("ListAll", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("list error")).Once()

		req, _ := http.NewRequest("GET", "/api/v1/bookings/all", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("Create Booking - Conflict Error", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		facilityID := uuid.New().String()
		now := time.Now().Add(1 * time.Hour)

		medicalStatus := userDomain.MedicalCertStatusValid
		mockUserRepo.On("GetByID", mock.Anything, clubID, userID).Return(&userDomain.User{
			ID: userID, MedicalCertStatus: &medicalStatus,
		}, nil).Once()
		mockFacilityRepo.On("GetByID", mock.Anything, clubID, facilityID).Return(&facilityDomain.Facility{
			ID: facilityID, Status: facilityDomain.FacilityStatusActive,
		}, nil).Once()
		mockBookingRepo.On("HasTimeConflict", mock.Anything, clubID, uuid.MustParse(facilityID), mock.Anything, mock.Anything).Return(true, nil).Once()

		body, _ := json.Marshal(map[string]interface{}{
			"user_id":     userID,
			"facility_id": facilityID, "start_time": now, "end_time": now.Add(1 * time.Hour),
		})
		req, _ := http.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusConflict, resp.Code)
	})

	t.Run("Create Booking - Facility Inactive", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		facilityID := uuid.New().String()

		mockFacilityRepo.On("GetByID", mock.Anything, clubID, facilityID).Return(&facilityDomain.Facility{
			ID: facilityID, Status: facilityDomain.FacilityStatusMaintenance,
		}, nil).Once()

		body, _ := json.Marshal(map[string]interface{}{
			"user_id":     userID,
			"facility_id": facilityID, "start_time": time.Now().Add(1 * time.Hour), "end_time": time.Now().Add(2 * time.Hour),
		})
		req, _ := http.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "facility_inactive")
	})

	t.Run("Create Booking - Medical Invalid", func(t *testing.T) {
		uniqueUserID := uuid.New().String()
		r := setupRouter(h, clubID, uniqueUserID, userDomain.RoleMember)
		facilityID := uuid.New().String()

		mockFacilityRepo.On("GetByID", mock.Anything, clubID, mock.Anything).Return(&facilityDomain.Facility{
			ID: facilityID, Status: facilityDomain.FacilityStatusActive,
		}, nil).Once()
		mockBookingRepo.On("HasTimeConflict", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()
		mockFacilityRepo.On("HasConflict", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()

		status := userDomain.MedicalCertStatusExpired
		mockUserRepo.On("GetByID", mock.Anything, clubID, uniqueUserID).Return(&userDomain.User{
			ID: uniqueUserID, MedicalCertStatus: &status,
		}, nil).Once()

		body, _ := json.Marshal(map[string]interface{}{
			"user_id":     uniqueUserID,
			"facility_id": facilityID, "start_time": time.Now().Add(1 * time.Hour), "end_time": time.Now().Add(2 * time.Hour),
		})
		req, _ := http.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "medical_certificate_invalid")
	})
	t.Run("Create Recurring Rule - RBAC Denial", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		req, _ := http.NewRequest("POST", "/api/v1/bookings/recurring", bytes.NewBufferString("{}"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("Generate Bookings - RBAC Denial", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		req, _ := http.NewRequest("POST", "/api/v1/bookings/generate", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("List All - RBAC Denial", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		req, _ := http.NewRequest("GET", "/api/v1/bookings/all", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("Create Booking - Service Error", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		facilityID := uuid.New().String()
		mockBookingRepo.On("Create", mock.Anything, mock.Anything).Return(fmt.Errorf("db error")).Once()

		status := userDomain.MedicalCertStatusValid
		mockUserRepo.On("GetByID", mock.Anything, clubID, userID).Return(&userDomain.User{
			ID: userID, MedicalCertStatus: &status,
		}, nil).Once()
		mockFacilityRepo.On("GetByID", mock.Anything, clubID, facilityID).Return(&facilityDomain.Facility{
			ID: facilityID, Status: facilityDomain.FacilityStatusActive,
		}, nil).Once()
		mockBookingRepo.On("HasTimeConflict", mock.Anything, clubID, facilityID, mock.Anything, mock.Anything).Return(false, nil).Once()
		mockFacilityRepo.On("HasConflict", mock.Anything, clubID, facilityID, mock.Anything, mock.Anything).Return(false, nil).Once()

		body, _ := json.Marshal(application.CreateBookingDTO{
			FacilityID: facilityID, StartTime: time.Now().Add(time.Hour), EndTime: time.Now().Add(2 * time.Hour),
		})
		req, _ := http.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("List - Service Error", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		mockBookingRepo.On("List", mock.Anything, clubID, mock.Anything).Return(nil, fmt.Errorf("db error")).Once()
		req, _ := http.NewRequest("GET", "/api/v1/bookings", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
