package main

import (
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	accessDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/domain"
	attendanceRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/domain"
	bookingDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	clubDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	disciplineDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/domain"
	facilityDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	membershipDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	paymentDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// SeederUser helps us create the full table schema required by all modules
type SeederUser struct {
	ID                string `gorm:"primaryKey"`
	Name              string
	Email             string `gorm:"uniqueIndex"`
	Password          string
	Role              string
	DateOfBirth       *time.Time
	SportsPreferences map[string]interface{} `gorm:"serializer:json"`
	ParentID          *string                `gorm:"index"`
	ClubID            string                 `gorm:"index"`
	GoogleID          string                 `gorm:"index"`
	AvatarURL         string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (SeederUser) TableName() string {
	return "users"
}

func main() {
	// Initialize DB
	os.Setenv("DB_PASSWORD", "pulse_secret")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")

	database.InitDB()
	db := database.GetDB()

	log.Println("--- Starting Seeder ---")

	// Reset Tables for clean seed (MVP)
	log.Println("Resetting Tables...")
	// Reset Tables logic removed to preserve App Schema

	// Ensure UUID extension exists
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if err := db.Migrator().DropTable(
		&SeederUser{},
		&clubDom.Club{},
		&facilityDom.Facility{},
		&membershipDom.MembershipTier{},
		&membershipDom.Membership{},
		&disciplineDom.Discipline{},
		&disciplineDom.TrainingGroup{},
		&bookingDom.Booking{},
		&paymentDom.Payment{},
		&accessDom.AccessLog{},
		&attendanceRepo.AttendanceRecordModel{},
		// Add others if we seed them
	); err != nil {
		log.Printf("Error dropping tables: %v", err)
	}
	if err := db.AutoMigrate(
		&SeederUser{},
		&clubDom.Club{},
		&facilityDom.Facility{},
		&membershipDom.MembershipTier{},
		&membershipDom.Membership{},
		&disciplineDom.Discipline{},
		&disciplineDom.TrainingGroup{},
		&bookingDom.Booking{},
		&paymentDom.Payment{},
		&accessDom.AccessLog{},
		&attendanceRepo.AttendanceRecordModel{},
		// Add others if we seed them
	); err != nil {
		log.Fatalf("Failed to automigrate: %v", err)
	}

	// 0. Seed Clubs
	defaultClub := clubDom.Club{
		ID:        "club-alpha",
		Name:      "Club Alpha",
		Domain:    "club-alpha.com",
		Status:    clubDom.ClubStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.FirstOrCreate(&defaultClub, "id = ?", defaultClub.ID).Error; err != nil {
		log.Printf("Error creating club: %v", err)
	} else {
		log.Println("Seeded Club: Club Alpha")
	}

	// 1. Seed Users
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := domain.User{
		ID:        uuid.New().String(),
		Email:     "admin@clubpulse.com",
		Password:  string(hashedPwd),
		Name:      "System Admin",
		Role:      domain.RoleAdmin,
		ClubID:    defaultClub.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Where("email = ?", admin.Email).FirstOrCreate(&admin).Error; err != nil {
		log.Printf("Error seeding admin: %v", err)
	} else {
		log.Println("Seeded Admin User")
	}

	// 1.2 Seed Super Admin
	superAdmin := domain.User{
		ID:        uuid.New().String(),
		Email:     "superadmin@clubpulse.com",
		Password:  string(hashedPwd), // Same password for MVP convenience
		Name:      "Super Administrator",
		Role:      domain.RoleSuperAdmin,
		ClubID:    "system",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Where("email = ?", superAdmin.Email).FirstOrCreate(&superAdmin).Error; err != nil {
		log.Printf("Error seeding super admin: %v", err)
	} else {
		log.Println("Seeded Super Admin User")
	}

	// 1.5 Seed Test User (Member)
	hashedPwdTest, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := domain.User{
		ID:        uuid.New().String(),
		Email:     "testuser@example.com",
		Password:  string(hashedPwdTest),
		Name:      "Test User",
		Role:      domain.RoleMember,
		ClubID:    defaultClub.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Where("email = ?", testUser.Email).FirstOrCreate(&testUser).Error; err != nil {
		log.Printf("Error seeding test user: %v", err)
	} else {
		log.Println("Seeded Test User")
	}

	// 2. Seed Facilities
	surfaceHard := "Hard Court"
	surfaceTurf := "Artificial Turf"

	facilities := []facilityDom.Facility{
		{
			ID:         uuid.New().String(),
			ClubID:     defaultClub.ID, // Add ClubID
			Name:       "Centre Court (Tennis)",
			Type:       facilityDom.FacilityTypeCourt,
			Status:     facilityDom.FacilityStatusActive,
			HourlyRate: 50.00,
			Capacity:   4,
			Location: facilityDom.Location{
				Name:        "Main Complex",
				Description: "North Wing",
			},
			Specifications: facilityDom.Specifications{
				SurfaceType: &surfaceHard,
				Lighting:    true,
				Covered:     false,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:         uuid.New().String(),
			ClubID:     defaultClub.ID, // Add ClubID
			Name:       "Padel Court 1",
			Type:       facilityDom.FacilityTypeCourt, // Padel is a court type
			Status:     facilityDom.FacilityStatusActive,
			HourlyRate: 35.00,
			Capacity:   4,
			Location: facilityDom.Location{
				Name:        "Padel Zone",
				Description: "South Wing",
			},
			Specifications: facilityDom.Specifications{
				SurfaceType: &surfaceTurf,
				Lighting:    true,
				Covered:     true, // Glass
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:         uuid.New().String(),
			ClubID:     defaultClub.ID, // Add ClubID
			Name:       "Main Gym",
			Type:       facilityDom.FacilityTypeGym,
			Status:     facilityDom.FacilityStatusActive,
			HourlyRate: 10.00,
			Capacity:   50,
			Location: facilityDom.Location{
				Name:        "Fitness Building",
				Description: "Floor 1",
			},
			Specifications: facilityDom.Specifications{
				// Gym might not have surface type relevant here, leaving defaults
				Lighting: true,
				Covered:  true,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, f := range facilities {
		if err := db.FirstOrCreate(&f, "name = ?", f.Name).Error; err != nil {
			log.Printf("Error seeding facility %s: %v", f.Name, err)
		} else {
			log.Printf("Seeded Facility: %s", f.Name)
		}
	}

	// 3. Seed Membership Tiers
	tiers := []membershipDom.MembershipTier{
		{
			ID:          uuid.New(),
			ClubID:      defaultClub.ID, // Add ClubID
			Name:        "Bronze",
			Description: "Entry level access for casual players",
			MonthlyFee:  decimal.NewFromFloat(29.99),
			Colors:      "bg-amber-100 text-amber-800",
			Benefits:    pq.StringArray{"Access to Gym", "10% Off Court Booking"},
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			ClubID:      defaultClub.ID, // Add ClubID
			Name:        "Silver",
			Description: "Perfect for regular members",
			MonthlyFee:  decimal.NewFromFloat(59.99),
			Colors:      "bg-slate-200 text-slate-800",
			Benefits:    pq.StringArray{"Access to Gym", "Unlimited Sauna", "20% Off Court Booking", "1 Free Guest Pass/mo"},
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			ClubID:      defaultClub.ID, // Add ClubID
			Name:        "Gold",
			Description: "The ultimate VIP experience",
			MonthlyFee:  decimal.NewFromFloat(99.99),
			Colors:      "bg-yellow-100 text-yellow-800",
			Benefits:    pq.StringArray{"All Facilities Access", "Priority Booking", "Free Court Rentals (Off-peak)", "5 Free Guest Passes/mo"},
			IsActive:    true,
		},
	}

	for _, t := range tiers {
		if err := db.FirstOrCreate(&t, "name = ?", t.Name).Error; err != nil {
			log.Printf("Error seeding tier %s: %v", t.Name, err)
		} else {
			log.Printf("Seeded Tier: %s", t.Name)
		}
	}

	// 4. Seed Disciplines
	futbol := disciplineDom.Discipline{
		ID:          uuid.New(),
		ClubID:      defaultClub.ID, // Add ClubID
		Name:        "Fútbol",
		Description: "Escuela de fútbol infantil y juvenil",
		IsActive:    true,
	}
	if err := db.FirstOrCreate(&futbol, "name = ?", futbol.Name).Error; err != nil {
		log.Printf("Error creating discipline %s: %v", futbol.Name, err)
	}

	tennis := disciplineDom.Discipline{
		ID:          uuid.New(),
		ClubID:      defaultClub.ID, // Add ClubID
		Name:        "Tenis",
		Description: "Clases grupales e individuales",
		IsActive:    true,
	}
	if err := db.FirstOrCreate(&tennis, "name = ?", tennis.Name).Error; err != nil {
		log.Printf("Error creating discipline %s: %v", tennis.Name, err)
	}

	// 5. Seed Training Groups
	groups := []disciplineDom.TrainingGroup{
		{
			ID:           uuid.New(),
			ClubID:       defaultClub.ID, // Add ClubID
			Name:         "Fútbol 2012",
			DisciplineID: futbol.ID,
			Category:     "2012",
			CoachID:      admin.ID,
			Schedule:     "Lun/Mie 18:00",
		},
		{
			ID:           uuid.New(),
			ClubID:       defaultClub.ID, // Add ClubID
			Name:         "Fútbol 2015",
			DisciplineID: futbol.ID,
			Category:     "2015",
			CoachID:      admin.ID,
			Schedule:     "Mar/Jue 17:30",
		},
		{
			ID:           uuid.New(),
			ClubID:       defaultClub.ID, // Add ClubID
			Name:         "Tenis Inicial - Niños",
			DisciplineID: tennis.ID,
			Category:     "2014",
			CoachID:      admin.ID,
			Schedule:     "Sab 10:00",
		},
	}

	for _, g := range groups {
		if err := db.FirstOrCreate(&g, "name = ?", g.Name).Error; err != nil {
			log.Printf("Error seeding group %s: %v", g.Name, err)
		}
	}

	log.Println("Seeded Disciplines and Training Groups")

	log.Println("--- Seeding Completed Successfully ---")
}
