package main

import (
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/domain"
	facilityDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	membershipDom "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

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
	db.Migrator().DropTable(&domain.User{}, &facilityDom.Facility{}, &membershipDom.MembershipTier{})
	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&facilityDom.Facility{})
	db.AutoMigrate(&membershipDom.MembershipTier{})

	// 1. Seed Users
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := domain.User{
		ID:        uuid.New().String(),
		Email:     "admin@clubpulse.com",
		Password:  string(hashedPwd),
		Name:      "System Admin",
		Role:      "ADMIN",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Where("email = ?", admin.Email).FirstOrCreate(&admin).Error; err != nil {
		log.Printf("Error seeding admin: %v", err)
	} else {
		log.Println("Seeded Admin User")
	}

	// 1.5 Seed Test User (Member)
	hashedPwdTest, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := domain.User{
		ID:        uuid.New().String(),
		Email:     "testuser@example.com",
		Password:  string(hashedPwdTest),
		Name:      "Test User",
		Role:      "MEMBER",
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
			Name:        "Bronze",
			Description: "Entry level access for casual players",
			MonthlyFee:  decimal.NewFromFloat(29.99),
			Colors:      "bg-amber-100 text-amber-800",
			Benefits:    pq.StringArray{"Access to Gym", "10% Off Court Booking"},
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			Name:        "Silver",
			Description: "Perfect for regular members",
			MonthlyFee:  decimal.NewFromFloat(59.99),
			Colors:      "bg-slate-200 text-slate-800",
			Benefits:    pq.StringArray{"Access to Gym", "Unlimited Sauna", "20% Off Court Booking", "1 Free Guest Pass/mo"},
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
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

	log.Println("--- Seeding Completed Successfully ---")
}
