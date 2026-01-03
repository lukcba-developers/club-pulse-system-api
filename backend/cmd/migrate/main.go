package main

import (
	"log"
	"os"

	accessDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/domain"
	attendanceRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/infrastructure/repository"
	bookingDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	disciplineDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/domain"
	membershipDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	paymentDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
)

func main() {
	log.Println("Starting Comprehensive Migration...")

	// Default envs for local if not set
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_PASSWORD", "pulse_secret")
		os.Setenv("DB_USER", "postgres")
		os.Setenv("DB_NAME", "club_pulse")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
	}

	database.InitDB()
	db := database.GetDB()

	log.Println("Migrating All Models...")

	err := db.AutoMigrate(
		&membershipDomain.MembershipTier{},
		&membershipDomain.Membership{},
		&bookingDomain.Booking{},
		&paymentDomain.Payment{},
		&accessDomain.AccessLog{},
		&attendanceRepo.AttendanceListModel{},
		&attendanceRepo.AttendanceRecordModel{},
		&disciplineDomain.Discipline{},
		&disciplineDomain.TrainingGroup{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration successful!")
}
