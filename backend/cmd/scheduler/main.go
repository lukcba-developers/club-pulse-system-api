package main

import (
	"context"
	"log"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"gorm.io/gorm"
)

func main() {
	log.Println("Starting Billing Scheduler Job...")

	// 1. Initialize Database
	database.InitDB()
	db := database.GetDB()

	// 2. Setup Dependencies
	// Note: We need to iterate over ALL clubs.
	// For MVP, we'll assume a single tenant or iterate if we had a ClubRepository.
	// The current ProcessMonthlyBilling requires a ClubID.
	// We will query distinct ClubIDs from the memberships table to be robust.

	if err := processAllClubs(db); err != nil {
		log.Fatalf("Job Failed: %v", err)
	}

	log.Println("Billing Scheduler Job Completed Successfully.")
}

func processAllClubs(db *gorm.DB) error {
	var clubIDs []string
	// Get all unique club_ids from memberships to process
	result := db.Table("memberships").Select("DISTINCT club_id").Find(&clubIDs)
	if result.Error != nil {
		return result.Error
	}

	membershipRepo := repository.NewPostgresMembershipRepository(db)
	scholarshipRepo := repository.NewPostgresScholarshipRepository(db)
	useCases := application.NewMembershipUseCases(membershipRepo, scholarshipRepo)

	ctx := context.Background()

	for _, clubID := range clubIDs {
		log.Printf("Processing Club: %s", clubID)
		count, err := useCases.ProcessMonthlyBilling(ctx, clubID)
		if err != nil {
			log.Printf("Error processing club %s: %v", clubID, err)
			continue // Continue with next club
		}
		log.Printf("Processed %d memberships for club %s", count, clubID)
	}

	return nil
}
