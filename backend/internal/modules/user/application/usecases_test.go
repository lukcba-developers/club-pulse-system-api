package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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

// --- Tests ---

func TestUserUseCases_Profiles(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := application.NewUserUseCases(mockRepo, nil)
	ctx := context.Background()

	t.Run("GetProfile Success and Errors", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "c1", "u1").Return(&domain.User{ID: "u1"}, nil).Once()
		res, _ := uc.GetProfile(ctx, "c1", "u1")
		assert.Equal(t, "u1", res.ID)

		_, err := uc.GetProfile(ctx, "c1", "")
		assert.Error(t, err)
	})

	t.Run("UpdateProfile Success and NotFound", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "c1", "u1").Return(&domain.User{ID: "u1", Name: "Old"}, nil).Once()
		mockRepo.On("Update", ctx, mock.Anything).Return(nil).Once()
		res, _ := uc.UpdateProfile(ctx, "c1", "u1", application.UpdateProfileDTO{Name: "New"})
		assert.Equal(t, "New", res.Name)

		mockRepo.On("GetByID", ctx, "c1", "u2").Return(nil, nil).Once()
		_, err := uc.UpdateProfile(ctx, "c1", "u2", application.UpdateProfileDTO{})
		assert.Error(t, err)
	})

	t.Run("DeleteUser Logic", func(t *testing.T) {
		mockRepo.On("Delete", ctx, "c1", "u2").Return(nil).Once()
		err := uc.DeleteUser(ctx, "c1", "u2", "u1")
		assert.NoError(t, err)

		err2 := uc.DeleteUser(ctx, "c1", "u1", "u1")
		assert.Error(t, err2)
	})

	t.Run("ListUsers Pagination", func(t *testing.T) {
		mockRepo.On("List", ctx, "c1", 10, 0, mock.Anything).Return([]domain.User{{ID: "1"}}, nil).Once()
		res, _ := uc.ListUsers(ctx, "c1", 0, -1, "search")
		assert.Len(t, res, 1)
	})
}

func TestUserUseCases_Gamification(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := application.NewUserUseCases(mockRepo, nil)
	ctx := context.Background()

	t.Run("UpdateMatchStats Overall", func(t *testing.T) {
		yesterday := time.Now().Add(-24 * time.Hour)
		user := &domain.User{
			ID: "u1",
			Stats: &domain.UserStats{
				Level:            1,
				Experience:       500,
				CurrentStreak:    5,
				LastActivityDate: &yesterday,
			},
		}
		mockRepo.On("GetByID", ctx, "c1", "u1").Return(user, nil).Twice()
		mockRepo.On("Update", ctx, mock.Anything).Return(nil).Twice()

		// Case 1: Win and Level Up
		err := uc.UpdateMatchStats(ctx, "c1", "u1", true, 200)
		assert.NoError(t, err)
		assert.Greater(t, user.Stats.Level, 1)
		assert.Equal(t, 6, user.Stats.CurrentStreak)

		// Case 2: Break Streak
		longAgo := time.Now().Add(-72 * time.Hour)
		user.Stats.LastActivityDate = &longAgo
		_ = uc.UpdateMatchStats(ctx, "c1", "u1", false, 10)
		assert.Equal(t, 1, user.Stats.CurrentStreak)
	})
}

func TestUserUseCases_ChildrenAndDep(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := application.NewUserUseCases(mockRepo, nil)
	ctx := context.Background()

	t.Run("ListChildren and RegisterChild", func(t *testing.T) {
		mockRepo.On("FindChildren", ctx, "c1", "p1").Return([]domain.User{{ID: "c1"}}, nil).Once()
		res, _ := uc.ListChildren(ctx, "c1", "p1")
		assert.Len(t, res, 1)

		mockRepo.On("Create", ctx, mock.Anything).Return(nil).Once()
		child, _ := uc.RegisterChild(ctx, "c1", "p1", application.RegisterChildDTO{Name: "Kid"})
		assert.Equal(t, "Kid", child.Name)

		_, err := uc.RegisterChild(ctx, "c1", "", application.RegisterChildDTO{})
		assert.Error(t, err)
	})

	t.Run("RegisterDependent", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, "dad@test.com").Return(nil, nil).Once()
		mockRepo.On("Create", ctx, mock.Anything).Return(nil).Twice()
		res, err := uc.RegisterDependent(ctx, "c1", application.RegisterDependentDTO{
			ParentEmail: "dad@test.com", ParentName: "Dad", ChildName: "Junior",
		})
		assert.NoError(t, err)
		assert.Contains(t, res.Name, "Junior")
	})
}

func TestUserUseCases_WalletAndIncidents(t *testing.T) {
	mockRepo := new(MockUserRepo)
	uc := application.NewUserUseCases(mockRepo, nil)
	ctx := context.Background()

	t.Run("Wallet Stats and Debt", func(t *testing.T) {
		user := &domain.User{ID: "u1", Stats: &domain.UserStats{Level: 2}, Wallet: &domain.Wallet{Balance: 100}}
		mockRepo.On("GetByID", ctx, "c1", "u1").Return(user, nil).Times(3)
		mockRepo.On("Update", ctx, mock.Anything).Return(nil).Once()

		s, _ := uc.GetStats(ctx, "c1", "u1")
		w, _ := uc.GetWallet(ctx, "c1", "u1")
		assert.Equal(t, 2, s.Level)
		assert.Equal(t, float64(100), w.Balance)

		_ = uc.CreateManualDebt(ctx, "c1", "u1", 30, "fees", "admin")
		assert.Equal(t, float64(70), user.Wallet.Balance)
	})

	t.Run("Incidents and Emergency", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "c1", "u1").Return(&domain.User{ID: "u1"}, nil).Once()
		mockRepo.On("Update", ctx, mock.Anything).Return(nil).Once()
		_ = uc.UpdateEmergencyInfo(ctx, "c1", "u1", "Mom", "555", "Swiss", "123")

		mockRepo.On("CreateIncident", ctx, mock.Anything).Return(nil).Once()
		inc, _ := uc.LogIncident(ctx, "c1", "u1", "Ouch", "Ice", "Witness", "admin")
		assert.Equal(t, "Ouch", inc.Description)
	})
}

func TestUserUseCases_FamilyGroups(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockFamilyRepo := new(MockFamilyGroupRepo)
	uc := application.NewUserUseCases(mockRepo, mockFamilyRepo)
	ctx := context.Background()

	t.Run("Family Lifecycle", func(t *testing.T) {
		mockFamilyRepo.On("GetByHeadUserID", ctx, "c1", "h1").Return(nil, nil).Once()
		mockFamilyRepo.On("Create", ctx, mock.Anything).Return(nil).Once()
		mockFamilyRepo.On("AddMember", ctx, "c1", mock.Anything, "h1").Return(nil).Once()
		_, _ = uc.CreateFamilyGroup(ctx, "c1", "h1", "Vader")

		mockFamilyRepo.On("GetByMemberID", ctx, "c1", "u1").Return(&domain.FamilyGroup{Name: "Vader"}, nil).Once()
		fg, _ := uc.GetMyFamilyGroup(ctx, "c1", "u1")
		assert.Equal(t, "Vader", fg.Name)

		gID := uuid.New()
		mockFamilyRepo.On("AddMember", ctx, "c1", gID, "u2").Return(nil).Once()
		_ = uc.AddFamilyMember(ctx, "c1", gID, "u2")

		// Secure Add
		mockFamilyRepo.On("GetByID", ctx, "c1", gID).Return(&domain.FamilyGroup{HeadUserID: "h1"}, nil).Once()
		mockFamilyRepo.On("AddMember", ctx, "c1", gID, "u3").Return(nil).Once()
		err := uc.AddFamilyMemberSecure(ctx, "c1", gID, "u3", "h1")
		assert.NoError(t, err)

		mockFamilyRepo.On("RemoveMember", ctx, "c1", gID, "u2").Return(nil).Once()
		_ = uc.RemoveFamilyMember(ctx, "c1", gID, "u2")
	})
}

func TestUserUseCases_GDPR(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockFamilyRepo := new(MockFamilyGroupRepo)
	uc := application.NewUserUseCases(mockRepo, mockFamilyRepo)
	ctx := context.Background()

	t.Run("GDPR Export and Anonymize", func(t *testing.T) {
		now := time.Now()
		user := &domain.User{
			ID: "u1", Name: "Neo",
			DateOfBirth:       &now,
			SportsPreferences: map[string]interface{}{"tennis": true},
			MedicalCertStatus: new(domain.MedicalCertStatus),
		}
		*user.MedicalCertStatus = domain.MedicalCertStatusValid
		user.MedicalCertExpiry = &now
		user.TermsAcceptedAt = &now
		user.PrivacyPolicyVersion = "1.0"

		mockRepo.On("GetByID", ctx, "c1", "u1").Return(user, nil).Once()
		mockRepo.On("FindChildren", ctx, "c1", "u1").Return([]domain.User{{ID: "c1", Name: "Child"}}, nil).Once()
		mockFamilyRepo.On("GetByMemberID", ctx, "c1", "u1").Return(&domain.FamilyGroup{Name: "Fam"}, nil).Once()

		exp, err := uc.ExportUserData(ctx, "c1", "u1")
		assert.NoError(t, err)
		assert.Equal(t, "Neo", exp.UserProfile["name"])
		assert.Len(t, exp.Children, 1)

		mockRepo.On("AnonymizeForGDPR", ctx, "c1", "u1").Return(nil).Once()
		err2 := uc.DeleteUserGDPR(ctx, "c1", "u1", "admin")
		assert.NoError(t, err2)
	})
}
