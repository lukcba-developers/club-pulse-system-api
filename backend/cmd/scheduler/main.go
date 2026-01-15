package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	championshipRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/infrastructure/repository"
	championshipJobs "github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/jobs"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/repository"
	notificationSvc "github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func main() {
	log.Println("🚀 Starting Billing Scheduler Service...")

	// 1. Initialize Database
	database.InitDB()
	db := database.GetDB()

	// 2. Check for one-shot mode
	if len(os.Args) > 1 && os.Args[1] == "--run-once" {
		log.Println("Running in one-shot mode...")
		if err := processAllClubs(db); err != nil {
			log.Fatalf("Job Failed: %v", err)
		}
		log.Println("✅ Billing Job Completed Successfully.")
		return
	}

	// 3. Setup Cron Scheduler
	c := cron.New(cron.WithSeconds())

	// Schedule billing job to run daily at 2:00 AM (low traffic time)
	// Format: Second Minute Hour DayOfMonth Month DayOfWeek
	cronSchedule := os.Getenv("BILLING_CRON_SCHEDULE")
	if cronSchedule == "" {
		cronSchedule = "0 0 2 * * *" // Default: 2:00 AM daily
	}

	_, err := c.AddFunc(cronSchedule, func() {
		log.Printf("⏰ [%s] Starting scheduled billing job...", time.Now().Format(time.RFC3339))
		if err := processAllClubs(db); err != nil {
			log.Printf("❌ Billing job failed: %v", err)
		} else {
			log.Printf("✅ [%s] Billing job completed successfully", time.Now().Format(time.RFC3339))
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule cron job: %v", err)
	}

	log.Printf("📅 Scheduled billing job with pattern: %s", cronSchedule)

	// 4. Schedule Match Reminder Job (hourly)
	matchReminderSchedule := os.Getenv("MATCH_REMINDER_CRON_SCHEDULE")
	if matchReminderSchedule == "" {
		matchReminderSchedule = "0 0 * * * *" // Default: Every hour at :00
	}

	champRepo := championshipRepo.NewPostgresChampionshipRepository(db)
	notifService := notificationSvc.NewNotificationService(nil, nil) // Console fallback
	matchReminderJob := championshipJobs.NewMatchReminderJob(champRepo, notifService, 24)

	_, err = c.AddFunc(matchReminderSchedule, func() {
		log.Printf("🏆 [%s] Starting match reminder job...", time.Now().Format(time.RFC3339))
		// Process all clubs
		var clubIDs []string
		db.Table("clubs").Select("id").Find(&clubIDs)
		for _, clubID := range clubIDs {
			if err := matchReminderJob.Run(context.Background(), clubID); err != nil {
				log.Printf("⚠️ Match reminder failed for club %s: %v", clubID, err)
			}
		}
		log.Printf("✅ [%s] Match reminder job completed", time.Now().Format(time.RFC3339))
	})
	if err != nil {
		log.Printf("⚠️ Failed to schedule match reminder job: %v", err)
	} else {
		log.Printf("📅 Scheduled match reminder job with pattern: %s", matchReminderSchedule)
	}

	// 5. Start scheduler
	c.Start()

	// 5. Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down scheduler...")
	ctx := c.Stop()
	<-ctx.Done()
	log.Println("👋 Scheduler stopped gracefully")
}

func processAllClubs(db *gorm.DB) error {
	var clubIDs []string
	// Get all unique club_ids from memberships to process
	result := db.Table("memberships").Select("DISTINCT club_id").Find(&clubIDs)
	if result.Error != nil {
		return result.Error
	}

	if len(clubIDs) == 0 {
		log.Println("ℹ️ No clubs found with active memberships")
		return nil
	}

	membershipRepo := repository.NewPostgresMembershipRepository(db)
	scholarshipRepo := repository.NewPostgresScholarshipRepository(db)
	useCases := application.NewMembershipUseCases(membershipRepo, scholarshipRepo)

	ctx := context.Background()
	totalProcessed := 0
	failedClubs := []string{}

	for _, clubID := range clubIDs {
		log.Printf("  📋 Processing Club: %s", clubID)
		count, err := useCases.ProcessMonthlyBilling(ctx, clubID)
		if err != nil {
			log.Printf("  ⚠️ Error processing club %s: %v", clubID, err)
			failedClubs = append(failedClubs, clubID)
			continue // Continue with next club
		}
		log.Printf("  ✅ Processed %d memberships for club %s", count, clubID)
		totalProcessed += count
	}

	// Summary
	log.Printf("📊 Billing Summary: %d memberships processed across %d clubs", totalProcessed, len(clubIDs)-len(failedClubs))
	if len(failedClubs) > 0 {
		log.Printf("⚠️ Failed clubs: %v", failedClubs)
	}

	return nil
}
