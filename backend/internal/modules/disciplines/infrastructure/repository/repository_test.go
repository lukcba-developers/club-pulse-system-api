package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- SQLite Compatible Models ---

type TestDiscipline struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	ClubID      string    `gorm:"index;not null"`
	Name        string    `gorm:"not null;unique;size:100"`
	Description string    `gorm:"type:text"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (TestDiscipline) TableName() string { return "disciplines" }

type TestGroup struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	ClubID       string    `gorm:"index;not null"`
	Name         string    `gorm:"not null;size:100"`
	DisciplineID uuid.UUID `gorm:"type:uuid;not null"`
	Category     string    `gorm:"not null;size:20"`
	CategoryYear int
	CoachID      string
	Schedule     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (TestGroup) TableName() string { return "training_groups" }

type TestTournament struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	ClubID       string    `gorm:"index;not null"`
	Name         string    `gorm:"not null"`
	DisciplineID uuid.UUID `gorm:"type:uuid;not null"`
	StartDate    time.Time
	EndDate      time.Time
	Status       domain.TournamentStatus
	Format       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (TestTournament) TableName() string { return "tournaments" }

type TestTeam struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	ClubID       string    `gorm:"index;not null"`
	Name         string    `gorm:"not null"`
	TournamentID uuid.UUID `gorm:"type:uuid;not null"`
	CaptainID    *string
	Members      []string `gorm:"serializer:json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (TestTeam) TableName() string { return "teams" }

type TestMatch struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	ClubID       string    `gorm:"index;not null"`
	TournamentID uuid.UUID `gorm:"type:uuid;not null"`
	HomeTeamID   uuid.UUID `gorm:"type:uuid;not null"`
	AwayTeamID   uuid.UUID `gorm:"type:uuid;not null"`
	ScoreHome    int
	ScoreAway    int
	StartTime    time.Time
	Status       domain.MatchStatus
	Round        string
	Location     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (TestMatch) TableName() string { return "matches" }

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Migrate using test models to avoid Postgres-specific defaults
	err = db.AutoMigrate(
		&TestDiscipline{},
		&TestGroup{},
		&TestTournament{},
		&TestTeam{},
		&TestMatch{},
	)
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestPostgresDisciplineRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPostgresDisciplineRepository(db)
	clubID := "club-rep-1"

	t.Run("Discipline Lifecycle", func(t *testing.T) {
		d := &domain.Discipline{
			ID:          uuid.New(),
			ClubID:      clubID,
			Name:        "Tennis",
			Description: "Court sports",
			IsActive:    true,
		}
		err := repo.CreateDiscipline(context.Background(), d)
		assert.NoError(t, err)

		saved, err := repo.GetDisciplineByID(context.Background(), clubID, d.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Tennis", saved.Name)

		list, _ := repo.ListDisciplines(context.Background(), clubID)
		assert.Len(t, list, 1)

		// Cross-tenant check
		otherClubID := "club-other"
		d2 := &domain.Discipline{ID: uuid.New(), ClubID: otherClubID, Name: "Soccer"}
		_ = repo.CreateDiscipline(context.Background(), d2)
		list, _ = repo.ListDisciplines(context.Background(), clubID)
		assert.Len(t, list, 1) // Should still be 1
	})

	t.Run("Group Lifecycle", func(t *testing.T) {
		dID := uuid.New()
		g := &domain.TrainingGroup{
			ID:           uuid.New(),
			ClubID:       clubID,
			Name:         "Advanced",
			DisciplineID: dID,
			Category:     "U18",
			CoachID:      "coach-1",
		}
		err := repo.CreateGroup(context.Background(), g)
		assert.NoError(t, err)

		saved, err := repo.GetGroupByID(context.Background(), clubID, g.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Advanced", saved.Name)

		// Filtering Tests
		// 1. By Category
		list, err := repo.ListGroups(context.Background(), clubID, map[string]interface{}{"category": "U18"})
		assert.NoError(t, err)
		assert.Len(t, list, 1)

		// 2. By Coach
		list, err = repo.ListGroups(context.Background(), clubID, map[string]interface{}{"coach_id": "coach-1"})
		assert.NoError(t, err)
		assert.Len(t, list, 1)

		// 3. No Match
		list, err = repo.ListGroups(context.Background(), clubID, map[string]interface{}{"coach_id": "coach-2"})
		assert.NoError(t, err)
		assert.Len(t, list, 0)
	})

	t.Run("Errors", func(t *testing.T) {
		// Not Found
		res, err := repo.GetDisciplineByID(context.Background(), clubID, uuid.New())
		assert.NoError(t, err)
		assert.Nil(t, res)

		resG, err := repo.GetGroupByID(context.Background(), clubID, uuid.New())
		assert.NoError(t, err)
		assert.Nil(t, resG)
	})
}

func TestPostgresTournamentRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPostgresTournamentRepository(db)
	clubID := "club-rep-2"

	t.Run("Tournament CRUD", func(t *testing.T) {
		tourney := &domain.Tournament{
			ID:     uuid.New(),
			ClubID: clubID,
			Name:   "Summer League",
		}
		err := repo.CreateTournament(context.Background(), tourney)
		assert.NoError(t, err)

		saved, _ := repo.GetTournamentByID(context.Background(), clubID, tourney.ID)
		assert.NotNil(t, saved)
		assert.Equal(t, "Summer League", saved.Name)

		tourney.Name = "Updated League"
		_ = repo.UpdateTournament(context.Background(), tourney)
		updated, _ := repo.GetTournamentByID(context.Background(), clubID, tourney.ID)
		assert.Equal(t, "Updated League", updated.Name)
	})

	t.Run("Standings Calculation", func(t *testing.T) {
		tID := uuid.New()
		team1 := &domain.Team{ID: uuid.New(), ClubID: clubID, TournamentID: tID, Name: "A"}
		team2 := &domain.Team{ID: uuid.New(), ClubID: clubID, TournamentID: tID, Name: "B"}
		_ = repo.CreateTeam(context.Background(), team1)
		_ = repo.CreateTeam(context.Background(), team2)

		match := &domain.Match{
			ID:           uuid.New(),
			ClubID:       clubID,
			TournamentID: tID,
			HomeTeamID:   team1.ID,
			AwayTeamID:   team2.ID,
			ScoreHome:    2,
			ScoreAway:    0,
			Status:       domain.MatchStatusPlayed,
		}
		_ = repo.CreateMatch(context.Background(), match)

		standings, err := repo.GetStandings(context.Background(), clubID, tID)
		assert.NoError(t, err)
		assert.Len(t, standings, 2)

		var s1, s2 domain.Standing
		for _, s := range standings {
			if s.TeamID == team1.ID {
				s1 = s
			}
			if s.TeamID == team2.ID {
				s2 = s
			}
		}

		assert.Equal(t, 3, s1.Points)
		assert.Equal(t, 1, s1.Won)
		assert.Equal(t, 0, s2.Points)
		assert.Equal(t, 1, s2.Lost)
	})

	t.Run("Extended Lists and Gets", func(t *testing.T) {
		tID := uuid.New()
		// 1. Tournaments List
		listT, err := repo.ListTournaments(context.Background(), clubID)
		assert.NoError(t, err)
		initialCount := len(listT)

		tourney := &domain.Tournament{ID: tID, ClubID: clubID, Name: "Winter Cup"}
		_ = repo.CreateTournament(context.Background(), tourney)

		listT, err = repo.ListTournaments(context.Background(), clubID)
		assert.NoError(t, err)
		assert.Len(t, listT, initialCount+1)

		// 2. Teams
		teamID := uuid.New()
		team := &domain.Team{ID: teamID, ClubID: clubID, TournamentID: tID, Name: "Tigers"}
		err = repo.CreateTeam(context.Background(), team)
		assert.NoError(t, err)

		savedTeam, err := repo.GetTeamByID(context.Background(), clubID, teamID)
		assert.NoError(t, err)
		assert.Equal(t, "Tigers", savedTeam.Name)

		listTeams, err := repo.ListTeams(context.Background(), clubID, tID)
		assert.NoError(t, err)
		assert.Len(t, listTeams, 1)

		// 3. Matches
		matchID := uuid.New()
		match := &domain.Match{
			ID:           matchID,
			ClubID:       clubID,
			TournamentID: tID,
			HomeTeamID:   teamID,
			AwayTeamID:   uuid.New(), // dummy
			StartTime:    time.Now(),
		}
		err = repo.CreateMatch(context.Background(), match)
		assert.NoError(t, err)

		savedMatch, err := repo.GetMatchByID(context.Background(), clubID, matchID)
		assert.NoError(t, err)
		assert.Equal(t, matchID, savedMatch.ID)

		match.Location = "Court 1"
		err = repo.UpdateMatch(context.Background(), match)
		assert.NoError(t, err)

		updatedMatch, _ := repo.GetMatchByID(context.Background(), clubID, matchID)
		assert.Equal(t, "Court 1", updatedMatch.Location)

		listMatches, err := repo.ListMatches(context.Background(), clubID, tID)
		assert.NoError(t, err)
		assert.Len(t, listMatches, 1)
	})
}
