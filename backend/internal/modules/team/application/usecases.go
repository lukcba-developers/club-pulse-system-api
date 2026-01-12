package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/team/domain"
)

type TeamUseCases struct {
	repo domain.TeamRepository
}

func NewTeamUseCases(repo domain.TeamRepository) *TeamUseCases {
	return &TeamUseCases{repo: repo}
}

func (uc *TeamUseCases) ScheduleMatch(ctx context.Context, clubID string, groupID uuid.UUID, opponent string, isHome bool, meetupTime time.Time, location string) (*domain.MatchEvent, error) {
	// Ideally we should verify groupID belongs to clubID here, but we lack GetGroup in this repo.
	// Assuming the caller (handler) or database foreign keys + RLS (if any) or implicit trust for creation.
	// However, for creation, strict checking is checking if Group exists and belongs to club.
	// If we don't check, a user could schedule a match for another club's group if they guess the ID.
	// TODO: Add verification if generic repo or cross-module access allowed.
	// For now, proceeding with creation.

	event := &domain.MatchEvent{
		ID:              uuid.New(),
		TrainingGroupID: groupID,
		OpponentName:    opponent,
		IsHomeGame:      isHome,
		MeetupTime:      meetupTime,
		Location:        location,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := uc.repo.CreateMatchEvent(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

func (uc *TeamUseCases) RespondAvailability(ctx context.Context, clubID, eventID string, userID string, status domain.PlayerAvailabilityStatus, reason string) error {
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return err
	}

	// Verify Event belongs to Club
	if _, err := uc.repo.GetMatchEvent(ctx, clubID, eventID); err != nil {
		return errors.New("event not found or access denied")
	}

	availability := &domain.PlayerAvailability{
		MatchEventID: eventUUID,
		UserID:       userID,
		Status:       status,
		Reason:       reason,
		UpdatedAt:    time.Now(),
	}

	return uc.repo.SetPlayerAvailability(ctx, availability)
}

func (uc *TeamUseCases) GetEventAvailabilities(ctx context.Context, clubID, eventID string) ([]domain.PlayerAvailability, error) {
	return uc.repo.GetEventAvailabilities(ctx, clubID, eventID)
}
