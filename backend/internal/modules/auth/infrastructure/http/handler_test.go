package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) SaveUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockAuthRepo) FindUserByEmail(ctx context.Context, clubID, email string) (*domain.User, error) {
	args := m.Called(ctx, clubID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockAuthRepo) FindUserByID(ctx context.Context, clubID, id string) (*domain.User, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockAuthRepo) SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
func (m *MockAuthRepo) GetRefreshToken(ctx context.Context, token, clubID string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, token, clubID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}
func (m *MockAuthRepo) RevokeRefreshToken(ctx context.Context, tokenID, userID string) error {
	args := m.Called(ctx, tokenID, userID)
	return args.Error(0)
}
func (m *MockAuthRepo) RevokeAllUserTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockAuthRepo) ListUserSessions(ctx context.Context, userID string) ([]domain.RefreshToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.RefreshToken), args.Error(1)
}
func (m *MockAuthRepo) LogAuthentication(ctx context.Context, log *domain.AuthenticationLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(user *domain.User) (*domain.Token, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Token), args.Error(1)
}
func (m *MockTokenService) ValidateToken(token string) (*domain.UserClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserClaims), args.Error(1)
}
func (m *MockTokenService) GenerateRefreshToken(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}
func (m *MockTokenService) ValidateRefreshToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func TestListSessions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockAuthRepo)
	mockTokenService := new(MockTokenService)
	useCase := application.NewAuthUseCases(mockRepo, mockTokenService, nil)
	handler := NewAuthHandler(useCase)

	t.Run("Success", func(t *testing.T) {
		userID := "user123"
		mockSessions := []domain.RefreshToken{
			{ID: "session1", UserID: userID, DeviceID: "device1", CreatedAt: time.Now()},
			{ID: "session2", UserID: userID, DeviceID: "device2", CreatedAt: time.Now()},
		}

		mockRepo.On("ListUserSessions", mock.Anything, userID).Return(mockSessions, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/sessions", nil)
		c.Set("userID", userID)

		handler.ListSessions(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []domain.RefreshToken
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, "session1", response[0].ID)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/sessions", nil)
		// Missing userID in context

		handler.ListSessions(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRevokeSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockAuthRepo)
	mockTokenService := new(MockTokenService)
	useCase := application.NewAuthUseCases(mockRepo, mockTokenService, nil)
	handler := NewAuthHandler(useCase)

	t.Run("Success", func(t *testing.T) {
		userID := "user123"
		sessionID := "session1"

		// Note: The usecase calls RevokeRefreshToken directly.
		mockRepo.On("RevokeRefreshToken", mock.Anything, sessionID, userID).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/sessions/"+sessionID, nil)
		c.Set("userID", userID)
		c.Params = gin.Params{{Key: "id", Value: sessionID}}

		handler.RevokeSession(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("BadRequest_MissingID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/sessions", nil)
		c.Set("userID", "user123")
		// Missing param

		handler.RevokeSession(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
