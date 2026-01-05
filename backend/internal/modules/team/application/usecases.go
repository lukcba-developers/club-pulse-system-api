package application

import (
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

func (uc *TeamUseCases) ScheduleMatch(groupID uuid.UUID, opponent string, isHome bool, meetupTime time.Time, location string) (*domain.MatchEvent, error) {
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

	if err := uc.repo.CreateMatchEvent(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (uc *TeamUseCases) RespondAvailability(eventID string, userID string, status domain.PlayerAvailabilityStatus, reason string) error {
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return err
	}

	availability := &domain.PlayerAvailability{
		MatchEventID: eventUUID,
		UserID:       userID,
		Status:       status,
		Reason:       reason,
		UpdatedAt:    time.Now(),
	}

	return uc.repo.SetPlayerAvailability(availability)
}

func (uc *TeamUseCases) GetEventAvailabilities(eventID string) ([]domain.PlayerAvailability, error) {
	return uc.repo.GetEventAvailabilities(eventID)
}
