package application

import (
	"errors"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
)

type ClubDTO struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Domain   string            `json:"domain"`
	Status   domain.ClubStatus `json:"status"`
	Settings string            `json:"settings"`
}

type CreateClubDTO struct {
	ID       string `json:"id" binding:"required"` // Slug
	Name     string `json:"name" binding:"required"`
	Domain   string `json:"domain"`
	Settings string `json:"settings"`
}

type UpdateClubDTO struct {
	Name     string            `json:"name"`
	Domain   string            `json:"domain"`
	Status   domain.ClubStatus `json:"status"`
	Settings string            `json:"settings"`
}

type ClubUseCases struct {
	repo domain.ClubRepository
}

func NewClubUseCases(repo domain.ClubRepository) *ClubUseCases {
	return &ClubUseCases{repo: repo}
}

func (uc *ClubUseCases) CreateClub(dto CreateClubDTO) (*domain.Club, error) {
	// Check if ID already exists
	existing, err := uc.repo.GetByID(dto.ID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("club ID already exists")
	}

	club := &domain.Club{
		ID:       dto.ID,
		Name:     dto.Name,
		Domain:   dto.Domain,
		Status:   domain.ClubStatusActive,
		Settings: dto.Settings,
	}

	if err := uc.repo.Create(club); err != nil {
		return nil, err
	}
	return club, nil
}

func (uc *ClubUseCases) UpdateClub(id string, dto UpdateClubDTO) (*domain.Club, error) {
	club, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if club == nil {
		return nil, errors.New("club not found")
	}

	if dto.Name != "" {
		club.Name = dto.Name
	}
	if dto.Domain != "" {
		club.Domain = dto.Domain
	}
	if dto.Status != "" {
		club.Status = dto.Status
	}
	if dto.Settings != "" {
		club.Settings = dto.Settings
	}

	if err := uc.repo.Update(club); err != nil {
		return nil, err
	}
	return club, nil
}

func (uc *ClubUseCases) GetClub(id string) (*domain.Club, error) {
	return uc.repo.GetByID(id)
}

func (uc *ClubUseCases) ListClubs(limit, offset int) ([]domain.Club, error) {
	return uc.repo.List(limit, offset)
}
