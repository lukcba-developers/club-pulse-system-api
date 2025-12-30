package application

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
)

type FacilityUseCases struct {
	repo domain.FacilityRepository
}

func NewFacilityUseCases(repo domain.FacilityRepository) *FacilityUseCases {
	return &FacilityUseCases{
		repo: repo,
	}
}

type CreateFacilityDTO struct {
	Name           string                `json:"name" binding:"required"`
	Type           domain.FacilityType   `json:"type" binding:"required"`
	Capacity       int                   `json:"capacity" binding:"required,min=1"`
	HourlyRate     float64               `json:"hourly_rate" binding:"required,min=0"`
	Specifications domain.Specifications `json:"specifications"`
	Location       domain.Location       `json:"location"`
}

func (uc *FacilityUseCases) CreateFacility(dto CreateFacilityDTO) (*domain.Facility, error) {
	facility := &domain.Facility{
		ID:             uuid.New().String(),
		Name:           dto.Name,
		Type:           dto.Type,
		Status:         domain.FacilityStatusActive,
		Capacity:       dto.Capacity,
		HourlyRate:     dto.HourlyRate,
		Specifications: dto.Specifications,
		Location:       dto.Location,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := uc.repo.Create(facility); err != nil {
		return nil, err
	}

	return facility, nil
}

func (uc *FacilityUseCases) ListFacilities(limit, offset int) ([]*domain.Facility, error) {
	if limit <= 0 {
		limit = 10
	}
	return uc.repo.List(limit, offset)
}

func (uc *FacilityUseCases) GetFacility(id string) (*domain.Facility, error) {
	if id == "" {
		return nil, errors.New("invalid ID")
	}
	return uc.repo.GetByID(id)
}

type UpdateFacilityDTO struct {
	Name           *string                `json:"name,omitempty"`
	Status         *domain.FacilityStatus `json:"status,omitempty"`
	Specifications *domain.Specifications `json:"specifications,omitempty"`
}

func (uc *FacilityUseCases) UpdateFacility(id string, dto UpdateFacilityDTO) (*domain.Facility, error) {
	facility, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if facility == nil {
		return nil, errors.New("facility not found")
	}

	if dto.Name != nil {
		facility.Name = *dto.Name
	}
	if dto.Status != nil {
		facility.Status = *dto.Status
	}
	if dto.Specifications != nil {
		// Full replacement of specs for simplicity in MVP, or merge?
		// Let's do partial update if needed, but struct replacement is easier for now.
		// Actually, since specifications is a struct, we should probably merge if we want specific fields,
		// but given the DTO has a pointer to the full struct, we replace it.
		// Ideally we'd have meaningful merge logic, but MVP: replace.
		facility.Specifications = *dto.Specifications
	}

	if err := uc.repo.Update(facility); err != nil {
		return nil, err
	}

	return facility, nil
}
