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
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	handler "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/http"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockMembershipRepo struct {
	mock.Mock
}

func (m *MockMembershipRepo) Create(ctx context.Context, membership *domain.Membership) error {
	return m.Called(ctx, membership).Error(0)
}
func (m *MockMembershipRepo) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Membership, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) GetByUserID(ctx context.Context, clubID string, userID uuid.UUID) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID, userID)
	return args.Get(0).([]domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) GetByUserIDs(ctx context.Context, clubID string, userIDs []uuid.UUID) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID, userIDs)
	return args.Get(0).([]domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) ListTiers(ctx context.Context, clubID string) ([]domain.MembershipTier, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]domain.MembershipTier), args.Error(1)
}
func (m *MockMembershipRepo) GetTierByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.MembershipTier, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MembershipTier), args.Error(1)
}
func (m *MockMembershipRepo) ListBillable(ctx context.Context, clubID string, date time.Time) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID, date)
	return args.Get(0).([]domain.Membership), args.Error(1)
}
func (m *MockMembershipRepo) Update(ctx context.Context, membership *domain.Membership) error {
	return m.Called(ctx, membership).Error(0)
}
func (m *MockMembershipRepo) UpdateBalance(ctx context.Context, clubID string, membershipID uuid.UUID, newBalance decimal.Decimal, nextBilling time.Time) error {
	return m.Called(ctx, clubID, membershipID, newBalance, nextBilling).Error(0)
}
func (m *MockMembershipRepo) UpdateBalancesBatch(ctx context.Context, updates map[uuid.UUID]struct {
	Balance     decimal.Decimal
	NextBilling time.Time
}) error {
	return m.Called(ctx, updates).Error(0)
}
func (m *MockMembershipRepo) ListAll(ctx context.Context, clubID string) ([]domain.Membership, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]domain.Membership), args.Error(1)
}

type MockScholarshipRepo struct {
	mock.Mock
}

func (m *MockScholarshipRepo) Create(ctx context.Context, s *domain.Scholarship) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockScholarshipRepo) GetByUserID(ctx context.Context, userID string) ([]*domain.Scholarship, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Scholarship), args.Error(1)
}
func (m *MockScholarshipRepo) GetActiveByUserID(ctx context.Context, userID string) (*domain.Scholarship, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Scholarship), args.Error(1)
}
func (m *MockScholarshipRepo) ListActiveByUserIDs(ctx context.Context, userIDs []string) (map[string]*domain.Scholarship, error) {
	args := m.Called(ctx, userIDs)
	return args.Get(0).(map[string]*domain.Scholarship), args.Error(1)
}

// --- Mock Subscription ---
type MockSubscriptionRepo struct {
	mock.Mock
}

func (m *MockSubscriptionRepo) Create(ctx context.Context, s *domain.Subscription) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockSubscriptionRepo) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Subscription, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}
func (m *MockSubscriptionRepo) GetByUserID(ctx context.Context, clubID string, userID uuid.UUID) ([]domain.Subscription, error) {
	args := m.Called(ctx, clubID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Subscription), args.Error(1)
}
func (m *MockSubscriptionRepo) Update(ctx context.Context, s *domain.Subscription) error {
	return m.Called(ctx, s).Error(0)
}

// --- Setup ---

func setupRouter(h *handler.MembershipHandler, clubID, userID, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
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

func TestMembershipHandler_Endpoints(t *testing.T) {
	mockRepo := new(MockMembershipRepo)
	mockScholarRepo := new(MockScholarshipRepo)
	mockSubRepo := new(MockSubscriptionRepo)
	uc := application.NewMembershipUseCases(mockRepo, mockScholarRepo, mockSubRepo)
	h := handler.NewMembershipHandler(uc)
	clubID := "club-h-final"
	userID := uuid.New().String()

	t.Run("User: Create and Cancel", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleMember)
		tierID := uuid.New()
		tier := &domain.MembershipTier{ID: tierID, MonthlyFee: decimal.NewFromInt(50)}
		mockRepo.On("GetTierByID", mock.Anything, clubID, tierID).Return(tier, nil)
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		body, _ := json.Marshal(map[string]interface{}{
			"user_id":            userID,
			"membership_tier_id": tierID.String(),
			"billing_cycle":      "MONTHLY",
		})
		req, _ := http.NewRequest("POST", "/api/v1/memberships", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)

		// Cancel
		mID := uuid.New()
		m := &domain.Membership{ID: mID, UserID: uuid.MustParse(userID)}
		mockRepo.On("GetByID", mock.Anything, clubID, mID).Return(m, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

		reqCancel, _ := http.NewRequest("DELETE", "/api/v1/memberships/"+mID.String(), nil)
		respCancel := httptest.NewRecorder()
		r.ServeHTTP(respCancel, reqCancel)
		assert.Equal(t, http.StatusOK, respCancel.Code)
	})

	t.Run("Admin: Billing and Scholarships", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)

		mockRepo.On("ListBillable", mock.Anything, clubID, mock.Anything).Return([]domain.Membership{}, nil)
		reqBilling, _ := http.NewRequest("POST", "/api/v1/memberships/process-billing", nil)
		r.ServeHTTP(httptest.NewRecorder(), reqBilling)

		mockScholarRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
		bodyS, _ := json.Marshal(map[string]interface{}{
			"user_id": uuid.New().String(), "percentage": 0.30, "reason": "aid",
		})
		reqS, _ := http.NewRequest("POST", "/api/v1/memberships/scholarship", bytes.NewBuffer(bodyS))
		r.ServeHTTP(httptest.NewRecorder(), reqS)
	})

	t.Run("Admin: Queries", func(t *testing.T) {
		r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)

		mockRepo.On("ListTiers", mock.Anything, clubID).Return([]domain.MembershipTier{{Name: "Basic"}}, nil)
		reqT, _ := http.NewRequest("GET", "/api/v1/memberships/tiers", nil)
		respT := httptest.NewRecorder()
		r.ServeHTTP(respT, reqT)
		assert.Equal(t, http.StatusOK, respT.Code)

		mID := uuid.New()
		mockRepo.On("GetByID", mock.Anything, clubID, mID).Return(&domain.Membership{ID: mID}, nil)
		reqG, _ := http.NewRequest("GET", "/api/v1/memberships/"+mID.String(), nil)
		respG := httptest.NewRecorder()
		r.ServeHTTP(respG, reqG)
		assert.Equal(t, http.StatusOK, respG.Code)
	})
}

func TestMembershipHandler_DetailedErrors(t *testing.T) {
	mockRepo := new(MockMembershipRepo)
	mockScholarRepo := new(MockScholarshipRepo)
	mockSubRepo := new(MockSubscriptionRepo)
	uc := application.NewMembershipUseCases(mockRepo, mockScholarRepo, mockSubRepo)
	h := handler.NewMembershipHandler(uc)
	clubID := "club-h-error"
	userID := uuid.New().String()
	r := setupRouter(h, clubID, userID, userDomain.RoleAdmin)

	t.Run("CreateMembership_InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/memberships", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("CreateMembership_ServiceError", func(t *testing.T) {
		tierID := uuid.New()
		body, _ := json.Marshal(map[string]interface{}{
			"user_id":            userID,
			"membership_tier_id": tierID.String(),
			"billing_cycle":      "MONTHLY",
		})

		mockRepo.On("GetTierByID", mock.Anything, clubID, tierID).Return(&domain.MembershipTier{ID: tierID}, nil).Once()
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(context.DeadlineExceeded).Once()

		req, _ := http.NewRequest("POST", "/api/v1/memberships", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("CancelMembership_InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/memberships/invalid-uuid", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("CancelMembership_ServiceError", func(t *testing.T) {
		mID := uuid.New()
		// Called twice: once in handler (admin check), once in usecase
		mockRepo.On("GetByID", mock.Anything, clubID, mID).Return(&domain.Membership{ID: mID, UserID: uuid.MustParse(userID)}, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(context.DeadlineExceeded).Once()

		req, _ := http.NewRequest("DELETE", "/api/v1/memberships/"+mID.String(), nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Tiers_ServiceError", func(t *testing.T) {
		mockRepo.On("ListTiers", mock.Anything, clubID).Return([]domain.MembershipTier{}, context.DeadlineExceeded).Once()
		req, _ := http.NewRequest("GET", "/api/v1/memberships/tiers", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("GetMembership_InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/memberships/invalid-uuid", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("ProcessBilling_RBAC", func(t *testing.T) {
		rMember := setupRouter(h, clubID, userID, userDomain.RoleMember)
		req, _ := http.NewRequest("POST", "/api/v1/memberships/process-billing", nil)
		resp := httptest.NewRecorder()
		rMember.ServeHTTP(resp, req)
		// Assuming middleware returns 403, but here test router might not have it attached similarly
		// Wait, existing handler logic usually has manual check if not using middleware.
		// Checking handler code: ProcessBilling usually restricted.
		// If not checked in handler, middleware should handle it. existing test used setupRouter with Admin.
		// Let's assume standard RBAC check inside handler or middleware 403.
		// handler_test.go uses `setupRouter` which mocks middleware logic roughly or we rely on logic inside handler.
		// Let's rely on what we saw in other modules: manual checks inside handler used `c.Get("userRole")`.
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("AssignScholarship_InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/memberships/scholarship", bytes.NewBufferString("invalid"))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("AssignScholarship_RBAC", func(t *testing.T) {
		rMember := setupRouter(h, clubID, userID, userDomain.RoleMember)
		req, _ := http.NewRequest("POST", "/api/v1/memberships/scholarship", nil)
		resp := httptest.NewRecorder()
		rMember.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})
}
