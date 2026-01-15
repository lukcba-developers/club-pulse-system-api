package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- SQLite Compatible Models ---

type TestFacility struct {
	ID             string                `gorm:"primaryKey"`
	Name           string                `gorm:"not null"`
	Type           string                `gorm:"not null"`
	Status         string                `gorm:"default:'active'"`
	Capacity       int                   `gorm:"not null"`
	HourlyRate     float64               `gorm:"not null"`
	OpeningTime    string                `gorm:"default:'08:00'"`
	ClosingTime    string                `gorm:"default:'23:00'"`
	GuestFee       float64               `gorm:"default:0"`
	Specifications domain.Specifications `gorm:"type:text;serializer:json"`
	Location       domain.Location       `gorm:"type:text;serializer:json"`
	ClubID         string                `gorm:"index;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (TestFacility) TableName() string { return "facilities" }

type TestMaintenanceTask struct {
	ID          string    `gorm:"primaryKey"`
	FacilityID  string    `gorm:"not null;index"`
	EquipmentID *string   `gorm:"index"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Status      string    `gorm:"not null"`
	Type        string    `gorm:"not null"`
	StartTime   time.Time `gorm:"not null;index"`
	EndTime     time.Time `gorm:"not null;index"`
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (TestMaintenanceTask) TableName() string { return "maintenance_tasks" }

type TestEquipment struct {
	ID           string `gorm:"primaryKey"`
	FacilityID   string `gorm:"not null;index"`
	Name         string `gorm:"not null"`
	Type         string `gorm:"not null"`
	Condition    string `gorm:"not null"`
	Status       string `gorm:"not null"`
	PurchaseDate *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (TestEquipment) TableName() string { return "equipment" }

type TestEquipmentLoan struct {
	ID                string    `gorm:"primaryKey"`
	EquipmentID       string    `gorm:"not null;index"`
	UserID            string    `gorm:"not null;index"`
	LoanedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ExpectedReturnAt  time.Time `gorm:"not null"`
	ReturnedAt        *time.Time
	Status            string `gorm:"default:'ACTIVE'"`
	ConditionOnReturn string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (TestEquipmentLoan) TableName() string { return "equipment_loans" }

// Need a simple Booking model for GetImpactedUsers test
type TestBooking struct {
	ID         string `gorm:"primaryKey"`
	FacilityID string `gorm:"not null"`
	UserID     string `gorm:"not null"`
	Status     string `gorm:"not null"`
	StartTime  time.Time
	EndTime    time.Time
}

func (TestBooking) TableName() string { return "bookings" }

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(
		&TestFacility{},
		&TestMaintenanceTask{},
		&TestEquipment{},
		&TestEquipmentLoan{},
		&TestBooking{},
	)
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

// --- Tests ---

func TestFacilityRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPostgresFacilityRepository(db)
	ctx := context.Background()

	t.Run("Facility CRUD", func(t *testing.T) {
		id := uuid.New().String()
		facility := &domain.Facility{ID: id, ClubID: "club-1", Name: "Test"}
		_ = repo.Create(ctx, facility)

		fetched, _ := repo.GetByID(ctx, "club-1", id)
		assert.Equal(t, "Test", fetched.Name)

		facility.Name = "Updated"
		_ = repo.Update(ctx, facility)
		fetched, _ = repo.GetByID(ctx, "club-1", id)
		assert.Equal(t, "Updated", fetched.Name)

		list, _ := repo.List(ctx, "club-1", 10, 0)
		assert.NotEmpty(t, list)
	})

	t.Run("Maintenance", func(t *testing.T) {
		facID := uuid.New().String()
		id := uuid.New().String()
		task := &domain.MaintenanceTask{
			ID:         id,
			FacilityID: facID,
			Title:      "Fix",
			Status:     domain.MaintenanceStatusScheduled,
			StartTime:  time.Now().Add(1 * time.Hour),
			EndTime:    time.Now().Add(2 * time.Hour),
		}

		_ = repo.CreateMaintenance(task)

		fetched, err := repo.GetMaintenanceByID(id)
		assert.NoError(t, err)
		assert.NotNil(t, fetched)
		assert.Equal(t, "Fix", fetched.Title)

		tasks, _ := repo.ListMaintenanceByFacility(ctx, facID)
		assert.Len(t, tasks, 1)

		conflict, _ := repo.HasConflict(ctx, "club-1", facID, time.Now().Add(90*time.Minute), time.Now().Add(3*time.Hour))
		assert.True(t, conflict)
	})

	t.Run("Equipment", func(t *testing.T) {
		facID := uuid.New().String()
		eqID := uuid.New().String()
		eq := &domain.Equipment{ID: eqID, FacilityID: facID, Name: "A", Status: "available"}
		_ = repo.CreateEquipment(ctx, eq)

		fetched, _ := repo.GetEquipmentByID(ctx, eqID)
		assert.Equal(t, "A", fetched.Name)

		eq.Name = "B"
		_ = repo.UpdateEquipment(ctx, eq)
		fetched, _ = repo.GetEquipmentByID(ctx, eqID)
		assert.Equal(t, "B", fetched.Name)

		list, _ := repo.ListEquipmentByFacility(ctx, facID)
		assert.Len(t, list, 1)

		// Atomic Loan
		loan := &domain.EquipmentLoan{ID: "l1", EquipmentID: eqID, UserID: "u1", Status: domain.LoanStatusActive}
		err := repo.LoanEquipmentAtomic(ctx, loan, eqID)
		assert.NoError(t, err)

		eq, _ = repo.GetEquipmentByID(ctx, eqID)
		assert.Equal(t, "loaned", eq.Status)
	})

	t.Run("GetImpactedUsers", func(t *testing.T) {
		facID := "fac-impact"
		start := time.Now().Add(1 * time.Hour)
		end := time.Now().Add(2 * time.Hour)

		// Create bookings
		db.Create(&TestBooking{ID: "b1", FacilityID: facID, UserID: "user-1", Status: "CONFIRMED", StartTime: start, EndTime: end})
		db.Create(&TestBooking{ID: "b2", FacilityID: facID, UserID: "user-2", Status: "PENDING", StartTime: start, EndTime: end})
		db.Create(&TestBooking{ID: "b3", FacilityID: facID, UserID: "user-3", Status: "CANCELLED", StartTime: start, EndTime: end})
		db.Create(&TestBooking{ID: "b4", FacilityID: facID, UserID: "user-1", Status: "CONFIRMED", StartTime: start, EndTime: end}) // Duplicate user

		users, err := repo.GetImpactedUsers(facID, start, end)
		assert.NoError(t, err)
		assert.Len(t, users, 2) // user-1 and user-2
		assert.Contains(t, users, "user-1")
		assert.Contains(t, users, "user-2")
	})
}

func TestLoanRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPostgresLoanRepository(db)
	ctx := context.Background()

	t.Run("Loan lifecycle", func(t *testing.T) {
		id := "loan-1"
		_ = repo.Create(ctx, &domain.EquipmentLoan{ID: id, UserID: "u1", Status: domain.LoanStatusActive})

		fetched, _ := repo.GetByID(ctx, id)
		assert.NotNil(t, fetched)

		list, _ := repo.ListByUser(ctx, "u1")
		assert.Len(t, list, 1)

		list, _ = repo.ListByStatus(ctx, domain.LoanStatusActive)
		assert.Len(t, list, 1)

		_ = repo.Update(ctx, &domain.EquipmentLoan{ID: id, Status: domain.LoanStatusReturned})
		fetched, _ = repo.GetByID(ctx, id)
		assert.Equal(t, domain.LoanStatusReturned, fetched.Status)
	})
}
