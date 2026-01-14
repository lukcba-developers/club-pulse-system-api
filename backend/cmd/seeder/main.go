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
	championshipDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	clubDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	disciplineDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/domain"
	facilityDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	membershipDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"

	paymentDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	storeDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/store/domain"
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
	Password          string `gorm:"column:password;not null"`
	Role              string
	DateOfBirth       *time.Time
	SportsPreferences map[string]interface{} `gorm:"serializer:json"`
	ParentID          *string                `gorm:"index"`
	ClubID            string                 `gorm:"index"`
	GoogleID          string                 `gorm:"index"`
	AvatarURL         string
	// New Business Rules
	MedicalCertStatus string `gorm:"default:'PENDING'"`
	MedicalCertExpiry *time.Time
	FamilyGroupID     *string `gorm:"type:uuid"`
	// GDPR/Legal Compliance
	TermsAcceptedAt      *time.Time
	PrivacyPolicyVersion string
	DataRetentionUntil   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (SeederUser) TableName() string {
	return "users"
}

type SeederFamilyGroup struct {
	ID         string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name       string
	HeadUserID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (SeederFamilyGroup) TableName() string {
	return "family_groups"
}

type SeederWallet struct {
	ID           string          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID       string          `gorm:"uniqueIndex"`
	Balance      decimal.Decimal `gorm:"type:decimal(10,2)"`
	Points       int
	Transactions []byte `gorm:"type:jsonb"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (SeederWallet) TableName() string {
	return "wallets"
}

func main() {
	// Initialize DB
	os.Setenv("DB_PASSWORD", "pulse_secret")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")

	database.InitDB()
	db := database.GetDB()

	sqlDB, _ := db.DB()
	var currentDBName string
	if err := sqlDB.QueryRow("SELECT current_database()").Scan(&currentDBName); err != nil {
		log.Printf("Failed to get current database name: %v", err)
	}
	log.Printf("Connected to Database: %s", currentDBName)

	log.Println("--- Starting Seeder ---")

	// Ensure UUID extension exists
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if err := db.Migrator().DropTable(
		&SeederUser{},
		&clubDom.Club{},
		&facilityDom.Facility{},
		&facilityDom.Equipment{},     // Added
		&facilityDom.EquipmentLoan{}, // Added
		&membershipDom.MembershipTier{},
		&membershipDom.Membership{},
		&membershipDom.Scholarship{}, // Added
		&disciplineDom.Discipline{},
		&disciplineDom.TrainingGroup{},
		&bookingDom.Booking{},
		&bookingDom.Waitlist{}, // Added
		&paymentDom.Payment{},
		&accessDom.AccessLog{},
		&attendanceRepo.AttendanceRecordModel{},
		&SeederFamilyGroup{},
		&SeederWallet{},
		// Championship Tables
		&championshipDom.Tournament{},
		&championshipDom.TournamentStage{},
		&championshipDom.Group{},
		&championshipDom.Team{},
		&championshipDom.TournamentMatch{},

		&championshipDom.Standing{},
		&storeDom.Product{},
		&storeDom.Order{},
	); err != nil {

		log.Printf("Error dropping tables: %v", err)
	}
	if err := db.AutoMigrate(
		&SeederUser{},
		&clubDom.Club{},
		&facilityDom.Facility{},
		&facilityDom.Equipment{},     // Added
		&facilityDom.EquipmentLoan{}, // Added
		&membershipDom.MembershipTier{},
		&membershipDom.Membership{},
		&membershipDom.Scholarship{}, // Added
		&disciplineDom.Discipline{},
		&disciplineDom.TrainingGroup{},
		&bookingDom.Booking{},
		&bookingDom.Waitlist{}, // Added
		&paymentDom.Payment{},
		&accessDom.AccessLog{},
		&attendanceRepo.AttendanceRecordModel{},
		&SeederFamilyGroup{},
		&SeederWallet{},
		// Championship Tables
		&championshipDom.Tournament{},
		&championshipDom.TournamentStage{},
		&championshipDom.Group{},
		&championshipDom.Team{},
		&championshipDom.TournamentMatch{},

		&championshipDom.Standing{},
		&storeDom.Product{},
		&storeDom.Order{},
	); err != nil {

		log.Fatalf("Failed to automigrate: %v", err)
	}

	// 0. Seed Clubs
	defaultClub := clubDom.Club{
		ID:        "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Name:      "Club Alpha",
		Slug:      "club-alpha", // Used for URL resolution and X-Club-ID header
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
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)

	// Explicitly using SeederUser to ensure table match
	admin := SeederUser{
		ID:        uuid.New().String(),
		Email:     "admin@clubpulse.com",
		Password:  string(hashedPwd),
		Name:      "System Admin",
		Role:      "ADMIN",
		ClubID:    defaultClub.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Where("email = ?", admin.Email).FirstOrCreate(&admin).Error; err != nil {
		log.Printf("Error seeding admin: %v", err)
	} else {
		log.Println("Seeded Admin User (using SeederUser struct)")
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

	validExpiry := time.Now().AddDate(1, 0, 0) // 1 year from now
	validStatus := "VALID"

	testUser := domain.User{
		ID:                uuid.New().String(),
		Email:             "testuser@example.com",
		Password:          string(hashedPwdTest),
		Name:              "Test User",
		Role:              domain.RoleMember,
		ClubID:            defaultClub.ID,
		MedicalCertStatus: &validStatus,
		MedicalCertExpiry: &validExpiry,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := db.Where("email = ?", testUser.Email).FirstOrCreate(&testUser).Error; err != nil {
		log.Printf("Error seeding test user: %v", err)
	} else {
		log.Println("Seeded Test User (With Valid Med Cert)")
	}

	// 1.5.1 Seed Coach User
	hashedPwdCoach, _ := bcrypt.GenerateFromPassword([]byte("Coach123!"), bcrypt.DefaultCost)
	coach := domain.User{
		ID:        uuid.New().String(),
		Email:     "coach@clubpulse.com",
		Password:  string(hashedPwdCoach),
		Name:      "Head Coach",
		Role:      "COACH", // Hardcoded to match enum string usually, or use domain.RoleCoach if valid
		ClubID:    defaultClub.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Where("email = ?", coach.Email).FirstOrCreate(&coach).Error; err != nil {
		log.Printf("Error seeding coach user: %v", err)
	} else {
		log.Println("Seeded Coach User")
	}

	// 1.5.2 Seed Students for Training Groups
	// Helper to get date from year
	getDateFromYear := func(year int) *time.Time {
		t := time.Date(year, time.January, 15, 0, 0, 0, 0, time.UTC)
		return &t
	}

	studentConfigs := []struct {
		CategoryYear int
		Count        int
		GroupPrefix  string
	}{
		{2012, 5, "Futbol2012"},
		{2015, 5, "Futbol2015"},
		{2014, 5, "Tenis2014"},
	}

	for _, config := range studentConfigs {
		for i := 1; i <= config.Count; i++ {
			dob := getDateFromYear(config.CategoryYear)
			student := domain.User{
				ID:                uuid.New().String(),
				Email:             config.GroupPrefix + "_student_" + uuid.New().String()[:4] + "@example.com",
				Password:          string(hashedPwdTest), // Reuse Member password
				Name:              config.GroupPrefix + " Player " + uuid.New().String()[:4],
				Role:              domain.RoleMember,
				ClubID:            defaultClub.ID,
				DateOfBirth:       dob,
				MedicalCertStatus: &validStatus,
				MedicalCertExpiry: &validExpiry,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			if err := db.Create(&student).Error; err != nil {
				log.Printf("Error seeding student %s: %v", student.Name, err)
			}
		}
		log.Printf("Seeded %d students for category %d", config.Count, config.CategoryYear)
	}

	// 1.6 Seed Family Group
	familyID := uuid.New().String()
	familyGroup := SeederFamilyGroup{
		ID:         familyID,
		Name:       "The Test Family",
		HeadUserID: testUser.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := db.FirstOrCreate(&familyGroup).Error; err != nil {
		log.Printf("Error seeding family group: %v", err)
	} else {
		log.Println("Seeded Family Group: The Test Family")
		// Update user with family ID
		db.Model(&SeederUser{}).Where("id = ?", testUser.ID).Update("family_group_id", familyID)
	}

	// 1.7 Seed Wallet
	wallet := SeederWallet{
		ID:        uuid.New().String(),
		UserID:    testUser.ID,
		Balance:   decimal.NewFromFloat(5000.00),
		Points:    150,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.FirstOrCreate(&wallet).Error; err != nil {
		log.Printf("Error seeding wallet: %v", err)
	} else {
		log.Println("Seeded Wallet for Test User ($5000)")
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

	// 6. Seed Equipment
	if len(facilities) > 0 {
		courtFac := facilities[0]
		equipments := []facilityDom.Equipment{
			{
				ID:         uuid.New().String(),
				FacilityID: courtFac.ID,
				Name:       "Tennis Racket Pro",
				Type:       "Racket",
				Condition:  facilityDom.EquipmentConditionExcellent,
				Status:     "available",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				ID:         uuid.New().String(),
				FacilityID: courtFac.ID,
				Name:       "Tennis Ball Basket",
				Type:       "Accessory",
				Condition:  facilityDom.EquipmentConditionGood,
				Status:     "available",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}
		for _, eq := range equipments {
			db.Create(&eq)
		}
		log.Println("Seeded Equipment")
	}

	// 7. Seed Championship (Basic)
	endDate := time.Now().AddDate(0, 1, 0)
	tournament := championshipDom.Tournament{
		ID:        uuid.New(),
		ClubID:    uuid.MustParse(defaultClub.ID),
		Name:      "Copa Verano 2026",
		Sport:     "TENNIS",
		Status:    championshipDom.TournamentActive, // Fixed constant name
		StartDate: time.Now(),
		EndDate:   &endDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(&tournament)
	log.Println("Seeded Tournament: Copa Verano 2026")

	// 8. Seed Store Products
	products := []storeDom.Product{
		{
			ID:            uuid.New(),
			ClubID:        defaultClub.ID,
			Name:          "Club Pulse T-Shirt",
			Description:   "Official Club T-Shirt",
			Price:         25.00,
			StockQuantity: 100,
			Category:      "Merch",
			IsActive:      true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            uuid.New(),
			ClubID:        defaultClub.ID,
			Name:          "Wilson US Open Balls (3-Pack)",
			Description:   "Official Tournament Balls",
			Price:         9.99,
			StockQuantity: 500,
			Category:      "Equipment",
			IsActive:      true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	for _, p := range products {
		if err := db.Create(&p).Error; err != nil {
			log.Printf("Error creating product %s: %v", p.Name, err)
		}
	}
	log.Println("Seeded Store Products")

	log.Println("--- Seeding Completed Successfully ---")

	var userCount int64
	db.Model(&SeederUser{}).Count(&userCount)
	log.Printf("Final Verification: Found %d users in database", userCount)
}
