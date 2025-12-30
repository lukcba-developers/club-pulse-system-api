package http

import (
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

func (m *MockAuthRepo) SaveUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockAuthRepo) FindUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockAuthRepo) FindUserByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockAuthRepo) SaveRefreshToken(token *domain.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}
func (m *MockAuthRepo) GetRefreshToken(token string) (*domain.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}
func (m *MockAuthRepo) RevokeRefreshToken(tokenID string) error {
	args := m.Called(tokenID)
	return args.Error(0)
}
func (m *MockAuthRepo) RevokeAllUserTokens(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}
func (m *MockAuthRepo) ListUserSessions(userID string) ([]domain.RefreshToken, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.RefreshToken), args.Error(1)
}
func (m *MockAuthRepo) LogAuthentication(log *domain.AuthenticationLog) error {
	args := m.Called(log)
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
func (m *MockTokenService) ValidateToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
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
	useCase := application.NewAuthUseCases(mockRepo, mockTokenService)
	handler := NewAuthHandler(useCase)

	t.Run("Success", func(t *testing.T) {
		userID := "user123"
		mockSessions := []domain.RefreshToken{
			{ID: "session1", UserID: userID, DeviceID: "device1", CreatedAt: time.Now()},
			{ID: "session2", UserID: userID, DeviceID: "device2", CreatedAt: time.Now()},
		}

		mockRepo.On("ListUserSessions", userID).Return(mockSessions, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
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
		// Missing userID in context

		handler.ListSessions(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRevokeSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockAuthRepo)
	mockTokenService := new(MockTokenService)
	useCase := application.NewAuthUseCases(mockRepo, mockTokenService)
	handler := NewAuthHandler(useCase)

	t.Run("Success", func(t *testing.T) {
		userID := "user123"
		sessionID := "session1"

		// Note: The usecase calls RevokeRefreshToken directly.
		mockRepo.On("RevokeRefreshToken", sessionID).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{{Key: "id", Value: sessionID}}

		handler.RevokeSession(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("BadRequest_MissingID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", "user123")
		// Missing param

		handler.RevokeSession(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
