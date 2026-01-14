package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestRecurringRule struct {
	ID         uuid.UUID             `gorm:"type:text;primary_key"`
	ClubID     string                `gorm:"index;not null"`
	FacilityID uuid.UUID             `gorm:"type:text;not null"`
	Type       domain.RecurrenceType `gorm:"type:text"`
	Frequency  string                `gorm:"type:text"`
	DayOfWeek  int                   `gorm:"not null"`
	StartTime  time.Time             `gorm:"not null"`
	EndTime    time.Time             `gorm:"not null"`
	StartDate  time.Time             `gorm:"not null"`
	EndDate    time.Time             `gorm:"not null"`
	OwnerID    *uuid.UUID            `gorm:"type:text"`
	GroupID    *uuid.UUID            `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (TestRecurringRule) TableName() string { return "recurring_rules" }

type RecurringRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo domain.RecurringRepository
}

func (s *RecurringRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	s.Require().NoError(err)
	s.db = db

	err = s.db.AutoMigrate(&TestRecurringRule{})
	s.Require().NoError(err)

	s.repo = repository.NewPostgresRecurringRepository(s.db)
}

func (s *RecurringRepositoryTestSuite) TearDownTest() {
	s.db.Exec("DELETE FROM recurring_rules")
}

func (s *RecurringRepositoryTestSuite) TestCreateAndGet() {
	rule := &domain.RecurringRule{
		ID:         uuid.New(),
		ClubID:     "club-1",
		FacilityID: uuid.New(),
		Type:       domain.RecurrenceTypeFixed,
		DayOfWeek:  1,
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(1 * time.Hour),
		StartDate:  time.Now(),
		EndDate:    time.Now().AddDate(0, 0, 30),
	}

	err := s.repo.Create(context.Background(), rule)
	s.NoError(err)

	list, err := s.repo.GetByFacility(context.Background(), "club-1", rule.FacilityID)
	s.NoError(err)
	s.Len(list, 1)
	s.Equal(rule.ID, list[0].ID)
}

func (s *RecurringRepositoryTestSuite) TestGetAllActive() {
	clubID := "active-club"

	err := s.repo.Create(context.Background(), &domain.RecurringRule{
		ID: uuid.New(), ClubID: clubID, FacilityID: uuid.New(),
		EndDate: time.Now().AddDate(0, 0, 1),
	})
	s.NoError(err)

	list, err := s.repo.GetAllActive(context.Background(), clubID)
	s.NoError(err)
	s.True(len(list) >= 1)
}

func (s *RecurringRepositoryTestSuite) TestGetEmptyFacility() {
	list, err := s.repo.GetByFacility(context.Background(), "club-1", uuid.New())
	s.NoError(err)
	s.Len(list, 0)
}

func (s *RecurringRepositoryTestSuite) TestGetAllActiveExpired() {
	clubID := "expired-club"

	// Rule that ended yesterday
	err := s.repo.Create(context.Background(), &domain.RecurringRule{
		ID: uuid.New(), ClubID: clubID, FacilityID: uuid.New(),
		EndDate: time.Now().AddDate(0, 0, -1),
	})
	s.NoError(err)

	list, err := s.repo.GetAllActive(context.Background(), clubID)
	s.NoError(err)
	s.Len(list, 0)
}

func (s *RecurringRepositoryTestSuite) TestDeleteNonExistent() {
	err := s.repo.Delete(context.Background(), "club-1", uuid.New())
	s.NoError(err) // GORM Delete usually doesn't error if not found unless specified
}

func TestRecurringRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RecurringRepositoryTestSuite))
}
