package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- Test Models for SQLite Compatibility ---
// We shadow the domain models to use SQLite-friendly types (no jsonb, no gen_random_uuid)

type TestUserStats struct {
	ID               uuid.UUID `gorm:"primaryKey"`
	UserID           string    `gorm:"index"`
	MatchesPlayed    int       `gorm:"default:0"`
	MatchesWon       int       `gorm:"default:0"`
	RankingPoints    int       `gorm:"default:0"`
	Level            int       `gorm:"default:1"`
	Experience       int       `gorm:"default:0"`
	CurrentStreak    int       `gorm:"default:0"`
	LongestStreak    int       `gorm:"default:0"`
	LastActivityDate *time.Time
	TotalXP          int `gorm:"default:0"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (TestUserStats) TableName() string {
	return "user_stats"
}

type TestWallet struct {
	ID      uuid.UUID `gorm:"primaryKey"`
	UserID  string    `gorm:"index"`
	Balance float64   `gorm:"default:0.0"`
	Points  int       `gorm:"default:0"`
	// SQLite doesn't support JSONB, just use text or ignore for this test
	Transactions string `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (TestWallet) TableName() string {
	return "wallets"
}

// --------------------------------------------

type UserRepositorySuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repository.PostgresUserRepository
}

func (s *UserRepositorySuite) SetupSuite() {
	// Use SQLite :memory: for fast, isolated tests
	var err error
	s.db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	s.Require().NoError(err)

	// Migrate the schema using Test models for dependencies
	// repository.UserModel is fine if it doesn't use incompatible types directly on the table itself
	// (Relationships are fine if tables exist)
	err = s.db.AutoMigrate(&repository.UserModel{}, &TestUserStats{}, &TestWallet{})
	s.Require().NoError(err)

	s.repo = repository.NewPostgresUserRepository(s.db)
}

func (s *UserRepositorySuite) TearDownTest() {
	// Clean table after each test
	s.db.Exec("DELETE FROM users")
	s.db.Exec("DELETE FROM user_stats")
	s.db.Exec("DELETE FROM wallets")
}

func (s *UserRepositorySuite) TestCreateAndGet() {
	// Use a fixed ClubID for isolation
	clubID := "club-alpha"
	userID := uuid.New().String()

	now := time.Now().Truncate(time.Second) // Truncate to avoid microsecond diffs

	user := &domain.User{
		ID:        userID,
		Name:      "Integration User",
		Email:     "int@example.com",
		Role:      "MEMBER",
		ClubID:    clubID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 1. Create
	err := s.repo.Create(context.Background(), user)
	s.NoError(err)

	// 2. Get By ID
	fetched, err := s.repo.GetByID(context.Background(), clubID, user.ID)
	s.NoError(err)
	if fetched != nil {
		s.Equal(user.Email, fetched.Email)
		s.Equal(user.Name, fetched.Name)
		s.Equal(user.ID, fetched.ID)
	}
}

func (s *UserRepositorySuite) TestGetByEmail() {
	clubID := "club-alpha"
	user := &domain.User{
		ID:     uuid.New().String(),
		Name:   "Email User",
		Email:  "findme@example.com",
		ClubID: clubID,
	}
	s.NoError(s.repo.Create(context.Background(), user))

	// Find
	found, err := s.repo.GetByEmail(context.Background(), clubID, "findme@example.com")
	s.NoError(err)
	s.NotNil(found)
	if found != nil {
		s.Equal(user.ID, found.ID)
	}
}

func (s *UserRepositorySuite) TestUpdate() {
	clubID := "club-alpha"
	user := &domain.User{
		ID:        uuid.New().String(),
		Name:      "Old Name",
		Email:     "u@e.com",
		ClubID:    clubID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.NoError(s.repo.Create(context.Background(), user))

	// Update
	user.Name = "New Name"
	user.UpdatedAt = time.Now()
	err := s.repo.Update(context.Background(), user)
	s.NoError(err)

	// Verify
	fetched, err := s.repo.GetByID(context.Background(), clubID, user.ID)
	s.NoError(err)
	if fetched != nil {
		s.Equal("New Name", fetched.Name)
	}
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}
