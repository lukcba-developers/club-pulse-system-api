package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- SQLite Compatible Models ---

type TestTournament struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	ClubID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"not null"`
	Description string
	Sport       string `gorm:"not null"`
	Category    string
	Status      string `gorm:"default:'DRAFT'"`
	StartDate   time.Time
	EndDate     *time.Time
	LogoURL     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (TestTournament) TableName() string { return "championships" }

type TestStage struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	TournamentID uuid.UUID `gorm:"type:uuid;not null;index"`
	Order        int       `gorm:"not null"`
	Name         string    `gorm:"not null"`
	Type         string    `gorm:"not null"`
	Status       string    `gorm:"default:'PENDING'"`
}

func (TestStage) TableName() string { return "tournament_stages" }

type TestGroup struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key"`
	StageID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name    string    `gorm:"not null"`
}

func (TestGroup) TableName() string { return "groups" }

type TestStanding struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	GroupID        uuid.UUID `gorm:"type:uuid;not null;index"`
	TeamID         uuid.UUID `gorm:"type:uuid;not null;index"`
	Points         int
	Played         int
	Won            int
	Drawn          int
	Lost           int
	GoalsFor       int
	GoalsAgainst   int
	GoalDifference int
	UpdatedAt      time.Time
}

func (TestStanding) TableName() string { return "standings" }

type TestMatch struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key"`
	TournamentID uuid.UUID  `gorm:"type:uuid;not null;index"`
	StageID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	GroupID      *uuid.UUID `gorm:"type:uuid;index"`
	HomeTeamID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	AwayTeamID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	HomeScore    *int
	AwayScore    *int
	BookingID    *uuid.UUID `gorm:"type:uuid;index"`
	Status       string     `gorm:"default:'SCHEDULED'"`
	Date         time.Time
}

func (TestMatch) TableName() string { return "tournament_matches" }

type TestTeam struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Name      string    `gorm:"not null"`
	LogoURL   string
	Contact   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TestTeam) TableName() string { return "teams" }

type TestVolunteerAssignment struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key"`
	ClubID     string    `gorm:"not null;index"`
	MatchID    uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID     string    `gorm:"not null;index"`
	Role       string    `gorm:"not null"`
	Notes      string
	AssignedBy string
	AssignedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (TestVolunteerAssignment) TableName() string { return "volunteer_assignments" }

type TestTeamMember struct {
	TeamID uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID string    `gorm:"primary_key"`
}

func (TestTeamMember) TableName() string { return "team_members" }

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Migrate using test models to avoid Postgres-specific defaults
	err = db.AutoMigrate(
		&TestTournament{},
		&TestStage{},
		&TestGroup{},
		&TestStanding{},
		&TestMatch{},
		&TestTeam{},
		&TestVolunteerAssignment{},
		&TestTeamMember{},
	)
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestPostgresChampionshipRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPostgresChampionshipRepository(db)
	clubID := uuid.New()

	t.Run("Tournament Lifecycle", func(t *testing.T) {
		tournament := &domain.Tournament{
			ID:     uuid.New(),
			ClubID: clubID,
			Name:   "Winter Cup",
			Sport:  "Padel",
		}

		err := repo.CreateTournament(context.TODO(), tournament)
		assert.NoError(t, err)

		saved, err := repo.GetTournament(context.TODO(), clubID.String(), tournament.ID.String())
		assert.NoError(t, err)
		assert.Equal(t, "Winter Cup", saved.Name)

		list, _ := repo.ListTournaments(context.TODO(), clubID.String())
		assert.Len(t, list, 1)
	})

	t.Run("Stages and Groups", func(t *testing.T) {
		tournament := &domain.Tournament{ID: uuid.New(), ClubID: clubID, Name: "Liga"}
		_ = repo.CreateTournament(context.TODO(), tournament)

		stage := &domain.TournamentStage{
			ID:           uuid.New(),
			TournamentID: tournament.ID,
			Name:         "Regulares",
			Type:         domain.StageGroup,
		}
		err := repo.CreateStage(context.TODO(), stage)
		assert.NoError(t, err)

		group := &domain.Group{ID: uuid.New(), StageID: stage.ID, Name: "A"}
		err = repo.CreateGroup(context.TODO(), group)
		assert.NoError(t, err)

		saved, err := repo.GetGroup(context.TODO(), clubID.String(), group.ID.String())
		assert.NoError(t, err)
		assert.Equal(t, "A", saved.Name)
	})

	t.Run("Match Operations", func(t *testing.T) {
		tournament := &domain.Tournament{ID: uuid.New(), ClubID: clubID, Name: "Cup"}
		_ = repo.CreateTournament(context.TODO(), tournament)
		stage := &domain.TournamentStage{ID: uuid.New(), TournamentID: tournament.ID, Name: "S1"}
		_ = repo.CreateStage(context.TODO(), stage)

		m1 := domain.TournamentMatch{
			ID:           uuid.New(),
			TournamentID: tournament.ID,
			StageID:      stage.ID,
			HomeTeamID:   uuid.New(),
			AwayTeamID:   uuid.New(),
			Status:       domain.MatchScheduled,
		}

		err := repo.CreateMatchesBatch(context.TODO(), []domain.TournamentMatch{m1})
		assert.NoError(t, err)

		saved, _ := repo.GetMatch(context.TODO(), clubID.String(), m1.ID.String())
		assert.NotNil(t, saved)

		err = repo.UpdateMatchResult(context.TODO(), clubID.String(), m1.ID.String(), 2, 0)
		assert.NoError(t, err)

		updated, _ := repo.GetMatch(context.TODO(), clubID.String(), m1.ID.String())
		assert.Equal(t, 2, *updated.HomeScore)
		assert.Equal(t, domain.MatchCompleted, updated.Status)

		// Create match from another club
		otherClub := uuid.New()
		m2 := domain.TournamentMatch{ID: uuid.New(), TournamentID: tournament.ID, StageID: stage.ID}
		db.Create(&m2) // Note: This doesn't strictly link to otherClub correctly if championships table isn't updated, but the Join in UpdateMatchResult will fail to find it for otherClub
		err = repo.UpdateMatchResult(context.TODO(), otherClub.String(), m2.ID.String(), 1, 1)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("UpdateMatchScheduling and GetMatchesByGroup", func(t *testing.T) {
		tournament := &domain.Tournament{ID: uuid.New(), ClubID: clubID, Name: "Schedule Cup"}
		_ = repo.CreateTournament(context.TODO(), tournament)
		stage := &domain.TournamentStage{ID: uuid.New(), TournamentID: tournament.ID, Name: "S1"}
		_ = repo.CreateStage(context.TODO(), stage)

		teamH := &TestTeam{ID: uuid.New(), Name: "Home"}
		teamA := &TestTeam{ID: uuid.New(), Name: "Away"}
		db.Create(teamH)
		db.Create(teamA)

		gID := uuid.New()
		match := domain.TournamentMatch{
			ID:           uuid.New(),
			TournamentID: tournament.ID,
			StageID:      stage.ID,
			GroupID:      &gID,
			HomeTeamID:   teamH.ID,
			AwayTeamID:   teamA.ID,
			Status:       domain.MatchScheduled,
		}
		_ = db.Create(&match)

		bID := uuid.New()
		err := repo.UpdateMatchScheduling(context.TODO(), clubID.String(), match.ID.String(), time.Now(), bID)
		assert.NoError(t, err)

		// Verify not found for different club
		err = repo.UpdateMatchScheduling(context.TODO(), uuid.New().String(), match.ID.String(), time.Now(), bID)
		assert.Error(t, err)

		matches, err := repo.GetMatchesByGroup(context.TODO(), clubID.String(), gID.String())
		assert.NoError(t, err)
		assert.Len(t, matches, 1)
		// SQLite Scan might not populate HomeTeamName/AwayTeamName if using domain model without those fields as GORM tags
		// but let's check if the scan worked for the basic fields
		assert.Equal(t, bID, *matches[0].BookingID)
	})

	t.Run("Standings and Team Registration", func(t *testing.T) {
		tournament := &domain.Tournament{ID: uuid.New(), ClubID: clubID, Name: "League"}
		_ = repo.CreateTournament(context.TODO(), tournament)
		stage := &domain.TournamentStage{ID: uuid.New(), TournamentID: tournament.ID, Name: "S1"}
		_ = repo.CreateStage(context.TODO(), stage)
		group := &domain.Group{ID: uuid.New(), StageID: stage.ID, Name: "A"}
		_ = repo.CreateGroup(context.TODO(), group)

		teamID := uuid.New()
		db.Create(&TestTeam{ID: teamID, Name: "Best Team"})

		standing := &domain.Standing{
			ID:      uuid.New(),
			GroupID: group.ID,
			TeamID:  teamID,
			Points:  3,
		}

		err := repo.RegisterTeam(context.TODO(), standing)
		assert.NoError(t, err)

		list, err := repo.GetStandings(context.TODO(), clubID.String(), group.ID.String())
		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, "Best Team", list[0].TeamName)

		standing.Won = 1
		err = repo.UpdateStandingsBatch(context.TODO(), []domain.Standing{*standing})
		assert.NoError(t, err)

		// Complex standings
		t2 := uuid.New()
		db.Create(&TestTeam{ID: t2, Name: "Team B"})
		s2 := &domain.Standing{ID: uuid.New(), GroupID: group.ID, TeamID: t2, Points: 10}
		_ = repo.RegisterTeam(context.TODO(), s2)

		list, _ = repo.GetStandings(context.TODO(), clubID.String(), group.ID.String())
		assert.Len(t, list, 2)
		assert.Equal(t, "Team B", list[0].TeamName) // Team B has more points
	})

	t.Run("GetTeamMembers", func(t *testing.T) {
		teamID := uuid.New()
		db.Create(&TestTeamMember{TeamID: teamID, UserID: "user-1"})
		db.Create(&TestTeamMember{TeamID: teamID, UserID: "user-2"})

		members, err := repo.GetTeamMembers(context.TODO(), teamID.String())
		assert.NoError(t, err)
		assert.Len(t, members, 2)
		assert.Contains(t, members, "user-1")
	})

	t.Run("CreateMatchesBatch Transaction", func(t *testing.T) {
		tournament := &domain.Tournament{ID: uuid.New(), ClubID: clubID, Name: "Batch"}
		_ = repo.CreateTournament(context.TODO(), tournament)
		stage := &domain.TournamentStage{ID: uuid.New(), TournamentID: tournament.ID, Name: "S1"}
		_ = repo.CreateStage(context.TODO(), stage)

		m1 := domain.TournamentMatch{ID: uuid.New(), TournamentID: tournament.ID, StageID: stage.ID}
		m2 := domain.TournamentMatch{ID: m1.ID, TournamentID: tournament.ID, StageID: stage.ID} // Duplicate ID for error

		err := repo.CreateMatchesBatch(context.TODO(), []domain.TournamentMatch{m1, m2})
		assert.Error(t, err) // Should fail due to duplicate ID

		var count int64
		db.Model(&TestMatch{}).Where("id = ?", m1.ID).Count(&count)
		assert.Equal(t, int64(0), count) // Should be rolled back
	})
}

func TestPostgresVolunteerRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPostgresVolunteerRepository(db)
	clubID := "club-1"
	mID := uuid.New()

	t.Run("Volunteer Lifecycle", func(t *testing.T) {
		v1 := &domain.VolunteerAssignment{
			ID:      uuid.New(),
			ClubID:  clubID,
			MatchID: mID,
			UserID:  "user-1",
			Role:    "REFREE",
		}

		err := repo.Create(context.TODO(), v1)
		assert.NoError(t, err)

		res, err := repo.GetByMatchID(context.TODO(), clubID, mID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		res, err = repo.GetByUserID(context.TODO(), clubID, "user-1")
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		res, err = repo.GetByRoleAndMatch(context.TODO(), clubID, mID, "REFREE")
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		v1.Notes = "Updated"
		err = repo.Update(context.TODO(), v1)
		assert.NoError(t, err)

		err = repo.Delete(context.TODO(), clubID, v1.ID)
		assert.NoError(t, err)

		res, _ = repo.GetByMatchID(context.TODO(), clubID, mID)
		assert.Len(t, res, 0)
	})
}
