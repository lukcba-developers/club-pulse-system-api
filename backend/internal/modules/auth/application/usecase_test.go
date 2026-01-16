package application_test

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/domain"
)

// --- Mocks ---

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) SaveUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthRepo) FindUserByEmail(ctx context.Context, email, clubID string) (*domain.User, error) {
	args := m.Called(ctx, email, clubID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthRepo) FindUserByID(ctx context.Context, id, clubID string) (*domain.User, error) {
	args := m.Called(ctx, id, clubID)
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

// --- Tests ---

func TestAuthUseCases_Register(t *testing.T) {
	clubID := "test-club-id"

	tests := []struct {
		name          string
		dto           application.RegisterDTO
		setupMocks    func(repo *MockAuthRepo, ts *MockTokenService)
		expectedError string
	}{
		{
			name: "Success",
			dto: application.RegisterDTO{
				Name:        "Test User",
				Email:       "test@example.com",
				Password:    "SecurePass123!",
				AcceptTerms: true,
			},
			setupMocks: func(repo *MockAuthRepo, ts *MockTokenService) {
				repo.On("FindUserByEmail", mock.Anything, "test@example.com", clubID).Return(nil, nil)
				repo.On("SaveUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.Email == "test@example.com" && u.ClubID == clubID
				})).Return(nil)
				ts.On("GenerateToken", mock.Anything).Return(&domain.Token{AccessToken: "token", RefreshToken: "refresh"}, nil)
				repo.On("SaveRefreshToken", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "Fail_WeakPassword",
			dto: application.RegisterDTO{
				Name:        "Test User",
				Email:       "test@example.com",
				Password:    "123", // Too short
				AcceptTerms: true,
			},
			setupMocks: func(repo *MockAuthRepo, ts *MockTokenService) {
				// No calls expected
			},
			expectedError: "Password must be at least 8 characters long",
		},
		{
			name: "Fail_UserAlreadyExists",
			dto: application.RegisterDTO{
				Name:        "Test User",
				Email:       "existing@example.com",
				Password:    "SecurePass123!",
				AcceptTerms: true,
			},
			setupMocks: func(repo *MockAuthRepo, ts *MockTokenService) {
				existingUser := &domain.User{ID: "existing-id", Email: "existing@example.com"}
				repo.On("FindUserByEmail", mock.Anything, "existing@example.com", clubID).Return(existingUser, nil)
			},
			expectedError: "User already exists",
		},
		{
			name: "Fail_TermsNotAccepted",
			dto: application.RegisterDTO{
				Name:        "Test User",
				Email:       "test@example.com",
				Password:    "SecurePass123!",
				AcceptTerms: false,
			},
			setupMocks: func(repo *MockAuthRepo, ts *MockTokenService) {
				// No calls expected
			},
			expectedError: "You must accept the Terms and Conditions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockAuthRepo)
			ts := new(MockTokenService)
			uc := application.NewAuthUseCases(repo, ts, nil)

			tt.setupMocks(repo, ts)

			_, err := uc.Register(context.Background(), tt.dto, clubID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
			ts.AssertExpectations(t)
		})
	}
}

func TestAuthUseCases_Login(t *testing.T) {
	clubID := "test-club-id"
	// Hashed password for "SecurePass123!" using default cost (for testing purposes, we assume bcrypt works)
	// Actually, creating a hash in test setup is better.
	// But since `bcrypt` is used inside UseCase, we can't easily inject a mocked hasher.
	// We rely on the fact that if we store a hash, the UseCase will verify it using `bcrypt`.
	// For this test, we might need a real hash.
	// We will create a hash in the test.

	// Helper to ignore errors for hash generation in test
	// hashed, _ := bcrypt.GenerateFromPassword([]byte("SecurePass123!"), bcrypt.MinCost)
	// But importing bcrypt here creates a dependency. The snippet used real bcrypt in UseCase.
	// The snippet provided in the plan didn't have Login test fully detailed but I should add it.

	tests := []struct {
		name          string
		dto           application.LoginDTO
		setupMocks    func(repo *MockAuthRepo, ts *MockTokenService)
		expectedError string
	}{
		{
			name: "Success",
			dto: application.LoginDTO{
				Email:    "test@example.com",
				Password: "SecurePass123!",
			},
			setupMocks: func(repo *MockAuthRepo, ts *MockTokenService) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("SecurePass123!"), bcrypt.MinCost)
				user := &domain.User{
					ID:       "user-id",
					Email:    "test@example.com",
					Password: string(hash),
					ClubID:   clubID,
				}
				repo.On("FindUserByEmail", mock.Anything, "test@example.com", clubID).Return(user, nil)
				ts.On("GenerateToken", user).Return(&domain.Token{AccessToken: "access", RefreshToken: "refresh"}, nil)
				repo.On("SaveRefreshToken", mock.Anything, mock.Anything).Return(nil)
				repo.On("LogAuthentication", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "Fail_UserNotFound",
			dto: application.LoginDTO{
				Email:    "notfound@example.com",
				Password: "password",
			},
			setupMocks: func(repo *MockAuthRepo, ts *MockTokenService) {
				repo.On("FindUserByEmail", mock.Anything, "notfound@example.com", clubID).Return(nil, nil)
			},
			expectedError: "Invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockAuthRepo)
			ts := new(MockTokenService)
			uc := application.NewAuthUseCases(repo, ts, nil)

			tt.setupMocks(repo, ts)

			_, err := uc.Login(context.Background(), tt.dto, clubID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
