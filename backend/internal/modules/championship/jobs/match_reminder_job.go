package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	notificationSvc "github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
)

// MatchReminderJob sends notifications to players about their upcoming matches
type MatchReminderJob struct {
	repo                domain.ChampionshipRepository
	notificationService *notificationSvc.NotificationService
	reminderHours       int // How many hours before match to send reminder
}

func NewMatchReminderJob(
	repo domain.ChampionshipRepository,
	notificationService *notificationSvc.NotificationService,
	reminderHours int,
) *MatchReminderJob {
	if reminderHours <= 0 {
		reminderHours = 24 // Default 24 hours before
	}
	return &MatchReminderJob{
		repo:                repo,
		notificationService: notificationService,
		reminderHours:       reminderHours,
	}
}

// Run executes the job - should be called by a scheduler (e.g., cron)
func (j *MatchReminderJob) Run(ctx context.Context, clubID string) error {
	now := time.Now()
	reminderWindow := now.Add(time.Duration(j.reminderHours) * time.Hour)

	// Get all upcoming matches within the reminder window
	// For simplicity, we'll use a placeholder approach here
	// In production, you'd add a repository method like GetUpcomingMatches(clubID, from, to)

	// Placeholder: Log that the job ran
	fmt.Printf("[MatchReminderJob] Running for club %s, window: %s to %s\n",
		clubID, now.Format(time.RFC3339), reminderWindow.Format(time.RFC3339))

	// TODO: Implement actual logic when repository method is available
	matches, err := j.repo.GetUpcomingMatches(ctx, clubID, now, reminderWindow)
	if err != nil {
		fmt.Printf("[MatchReminderJob] Error fetching matches: %v\n", err)
		return err
	}

	for _, match := range matches {
		homeMembers, _ := j.repo.GetTeamMembers(ctx, match.HomeTeamID.String())
		awayMembers, _ := j.repo.GetTeamMembers(ctx, match.AwayTeamID.String())

		// Send to all
		affectedUsers := append(homeMembers, awayMembers...)
		for _, userID := range affectedUsers {
			_ = j.notificationService.Send(ctx, notificationSvc.Notification{
				RecipientID: userID,
				Type:        notificationSvc.NotificationTypePush,
				Title:       "🏆 Recordatorio de Partido",
				Body:        fmt.Sprintf("Tu partido es en menos de %d horas vs %s", j.reminderHours, match.AwayTeamName), // Simple message
			})
		}
	}

	return nil
}
