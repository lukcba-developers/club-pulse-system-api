package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- Test Models ---

type TestTrainingGroup struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	ClubID    string    `gorm:"index"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TestTrainingGroup) TableName() string { return "training_groups" }

type TestMatchEvent struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
	TrainingGroupID uuid.UUID `gorm:"type:uuid;not null;index"`
	OpponentName    string
	Location        string
	IsHomeGame      bool
	MeetupTime      time.Time
	StartTime       *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (TestMatchEvent) TableName() string { return "match_events" }

type TestPlayerAvailability struct {
	MatchEventID uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID       string    `gorm:"primary_key"`
	Status       string
	Reason       string
	UpdatedAt    time.Time
}

func (TestPlayerAvailability) TableName() string { return "player_availabilities" }

// --- Setup ---

func setupTestDB(t *testing.T) (*gorm.DB, *repository.PostgresTeamRepository) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&TestTrainingGroup{}, &TestMatchEvent{}, &TestPlayerAvailability{})
	assert.NoError(t, err)

	return db, repository.NewPostgresTeamRepository(db)
}

// --- Tests ---

func TestPostgresTeamRepository_SetPlayerAvailability_TenantIsolation(t *testing.T) {
	db, repo := setupTestDB(t)

	// Setup: Club 1 and Club 2
	club1 := "club-1"
	club2 := "club-2"

	// Create Training Group for Club 1
	tg1 := TestTrainingGroup{
		ID:        uuid.New(),
		ClubID:    club1,
		Name:      "Group 1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(&tg1)

	// Create Match Event for Club 1
	me1 := TestMatchEvent{
		ID:              uuid.New(),
		TrainingGroupID: tg1.ID,
		OpponentName:    "Opponent A",
		MeetupTime:      time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	db.Create(&me1)

	// Initial Availability
	av := &domain.PlayerAvailability{
		MatchEventID: me1.ID,
		UserID:       "user-1",
		Status:       domain.AvailabilityConfirmed,
		Reason:       "",
		UpdatedAt:    time.Now(),
	}

	t.Run("Success: Update availability for correct club", func(t *testing.T) {
		err := repo.SetPlayerAvailability(context.Background(), club1, av)
		assert.NoError(t, err)

		var stored TestPlayerAvailability
		err = db.First(&stored, "match_event_id = ? AND user_id = ?", me1.ID, "user-1").Error
		assert.NoError(t, err)
		assert.Equal(t, string(domain.AvailabilityConfirmed), stored.Status)
	})

	t.Run("Failure: Update availability for wrong club", func(t *testing.T) {
		err := repo.SetPlayerAvailability(context.Background(), club2, av)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Success: Update existing availability", func(t *testing.T) {
		av.Status = domain.AvailabilityDeclined
		err := repo.SetPlayerAvailability(context.Background(), club1, av)
		assert.NoError(t, err)

		var stored TestPlayerAvailability
		err = db.First(&stored, "match_event_id = ? AND user_id = ?", me1.ID, "user-1").Error
		assert.NoError(t, err)
		assert.Equal(t, string(domain.AvailabilityDeclined), stored.Status)
	})
}
