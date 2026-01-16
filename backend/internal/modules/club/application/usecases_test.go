package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	notificationSvc "github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockClubRepo struct {
	mock.Mock
}

func (m *MockClubRepo) Create(ctx context.Context, club *domain.Club) error {
	args := m.Called(ctx, club)
	return args.Error(0)
}
func (m *MockClubRepo) GetByID(ctx context.Context, id string) (*domain.Club, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Club), args.Error(1)
}
func (m *MockClubRepo) GetBySlug(ctx context.Context, slug string) (*domain.Club, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Club), args.Error(1)
}
func (m *MockClubRepo) List(ctx context.Context, limit, offset int) ([]domain.Club, error) {
	return nil, nil
}
func (m *MockClubRepo) Update(ctx context.Context, club *domain.Club) error {
	args := m.Called(ctx, club)
	return args.Error(0)
}
func (m *MockClubRepo) Delete(ctx context.Context, id string) error {
	return nil
}
func (m *MockClubRepo) GetMemberEmails(ctx context.Context, clubID string) ([]string, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]string), args.Error(1)
}

type MockNewsRepo struct {
	mock.Mock
}

func (m *MockNewsRepo) CreateNews(ctx context.Context, news *domain.News) error {
	args := m.Called(ctx, news)
	return args.Error(0)
}
func (m *MockNewsRepo) GetPublicNewsByClub(ctx context.Context, clubID string, limit, offset int) ([]domain.News, error) {
	return nil, nil
}
func (m *MockNewsRepo) GetNewsByClub(ctx context.Context, clubID string, limit, offset int) ([]domain.News, error) {
	return nil, nil
}
func (m *MockNewsRepo) GetNewsByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.News, error) { // Note: using uuid in interface? checking file news.go id is uuid.UUID
	return nil, nil
}

// Checking news.go again. id is uuid.UUID.
// I will implement correctly.

func (m *MockNewsRepo) UpdateNews(ctx context.Context, clubID string, news *domain.News) error {
	return nil
}
func (m *MockNewsRepo) DeleteNews(ctx context.Context, clubID string, id uuid.UUID) error {
	return nil
}

type MockSponsorRepo struct{ mock.Mock }

// Implement stubs if needed, mostly unused in these tests
func (m *MockSponsorRepo) CreateSponsor(ctx context.Context, s *domain.Sponsor) error { return nil }
func (m *MockSponsorRepo) CreateAdPlacement(ctx context.Context, ad *domain.AdPlacement) error {
	return nil
}
func (m *MockSponsorRepo) GetActiveAds(ctx context.Context, clubID string) ([]domain.AdPlacement, error) {
	return nil, nil
}

// Notification Mocks
type MockEmailProvider struct{ mock.Mock }

func (m *MockEmailProvider) SendEmail(ctx context.Context, to, subject, body string) (*notificationSvc.DeliveryResult, error) {
	args := m.Called(ctx, to, subject, body)
	return &notificationSvc.DeliveryResult{Success: true}, args.Error(1)
}

type MockSMSProvider struct{ mock.Mock }

func (m *MockSMSProvider) SendSMS(ctx context.Context, to, body string) (*notificationSvc.DeliveryResult, error) {
	return nil, nil
}

// --- Tests ---

func TestClubUseCases_UpdateClub(t *testing.T) {
	clubRepo := new(MockClubRepo)
	uc := application.NewClubUseCases(nil, clubRepo, nil, nil)
	clubID := "c1"

	t.Run("Success: Update Validation", func(t *testing.T) {
		existing := &domain.Club{
			ID:           clubID,
			Name:         "Old Name",
			PrimaryColor: "#000000",
		}

		clubRepo.On("GetByID", mock.Anything, clubID).Return(existing, nil).Once()

		clubRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *domain.Club) bool {
			return c.Name == "New Name" && c.PrimaryColor == "#FFFFFF" && c.SecondaryColor == "#FF0000"
		})).Return(nil).Once()

		updated, err := uc.UpdateClub(context.TODO(), clubID, "New Name", "", "", "#FFFFFF", "#FF0000", "", "", "", "", domain.ClubStatusActive)
		assert.NoError(t, err)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, "#FFFFFF", updated.PrimaryColor)
		clubRepo.AssertExpectations(t)
	})

	t.Run("Fail: Club Not Found", func(t *testing.T) {
		clubRepo.On("GetByID", mock.Anything, "unknown").Return(nil, nil).Once()
		res, err := uc.UpdateClub(context.TODO(), "unknown", "Name", "", "", "", "", "", "", "", "", domain.ClubStatusActive)
		assert.NoError(t, err) // Current logic returns nil, nil
		assert.Nil(t, res)
	})
}

func TestClubUseCases_PublishNews(t *testing.T) {
	clubRepo := new(MockClubRepo)
	newsRepo := new(MockNewsRepo)
	emailMock := new(MockEmailProvider)

	notifier := notificationSvc.NewNotificationService(emailMock, nil)
	uc := application.NewClubUseCases(nil, clubRepo, newsRepo, notifier)
	clubID := "c1"

	t.Run("Success: Create News and Notify", func(t *testing.T) {
		newsRepo.On("CreateNews", mock.Anything, mock.Anything).Return(nil).Once()
		clubRepo.On("GetMemberEmails", mock.Anything, clubID).Return([]string{"u1@example.com"}, nil).Once()

		// Setup Email Mock expectation.
		// Note: Since notification is async (goroutine), asserting validation requires waiting or mocked channels.
		// However, testified 'Called' might race.
		// For robustness, we can use a WaitGroup or just assert that CreateNews passed
		// and rely on integration tests for full async verification.
		// But let's try to match call. Update: Mock is thread-safe mostly.
		emailMock.On("SendEmail", mock.Anything, "u1@example.com", mock.MatchedBy(func(s string) bool {
			return s == "Nueva noticia: Title"
		}), mock.Anything).Return(nil, nil).Once()

		_, err := uc.PublishNews(context.TODO(), clubID, "Title", "Content", "", true)
		assert.NoError(t, err)

		// Simple sleep to allow goroutine to fire
		time.Sleep(50 * time.Millisecond)
		emailMock.AssertExpectations(t)
	})

	t.Run("Success: Create News No Notify", func(t *testing.T) {
		newsRepo.On("CreateNews", mock.Anything, mock.Anything).Return(nil).Once()
		_, err := uc.PublishNews(context.TODO(), clubID, "Title 2", "Content", "", false)
		assert.NoError(t, err)
		// No GetMemberEmails call
	})
}
