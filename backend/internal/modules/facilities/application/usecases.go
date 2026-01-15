package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
)

type FacilityUseCases struct {
	repo     domain.FacilityRepository
	loanRepo domain.LoanRepository
}

func NewFacilityUseCases(repo domain.FacilityRepository, loanRepo domain.LoanRepository) *FacilityUseCases {
	return &FacilityUseCases{
		repo:     repo,
		loanRepo: loanRepo,
	}
}

type CreateFacilityDTO struct {
	Name           string                `json:"name" binding:"required"`
	Description    string                `json:"description"`
	Type           domain.FacilityType   `json:"type" binding:"required"`
	Capacity       int                   `json:"capacity" binding:"required,min=1"`
	HourlyRate     float64               `json:"hourly_rate" binding:"required,min=0"`
	OpeningTime    string                `json:"opening_time"`
	ClosingTime    string                `json:"closing_time"`
	Specifications domain.Specifications `json:"specifications"`
	Location       domain.Location       `json:"location"`
}

func (uc *FacilityUseCases) CreateFacility(ctx context.Context, clubID string, dto CreateFacilityDTO) (*domain.Facility, error) {
	// Defaults
	opening := dto.OpeningTime
	if opening == "" {
		opening = "08:00"
	}
	closing := dto.ClosingTime
	if closing == "" {
		closing = "23:00"
	}

	facility := &domain.Facility{
		ID:             uuid.New().String(),
		ClubID:         clubID,
		Name:           dto.Name,
		Description:    dto.Description,
		Type:           dto.Type,
		Status:         domain.FacilityStatusActive,
		Capacity:       dto.Capacity,
		HourlyRate:     dto.HourlyRate,
		OpeningTime:    opening,
		ClosingTime:    closing,
		Specifications: dto.Specifications,
		Location:       dto.Location,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := uc.repo.Create(ctx, facility); err != nil {
		return nil, err
	}

	return facility, nil
}

func (uc *FacilityUseCases) ListFacilities(ctx context.Context, clubID string, limit, offset int) ([]*domain.Facility, error) {
	if limit <= 0 {
		limit = 10
	}
	return uc.repo.List(ctx, clubID, limit, offset)
}

func (uc *FacilityUseCases) GetFacility(ctx context.Context, clubID, id string) (*domain.Facility, error) {
	if id == "" {
		return nil, errors.New("invalid ID")
	}
	return uc.repo.GetByID(ctx, clubID, id)
}

type UpdateFacilityDTO struct {
	Name           *string                `json:"name,omitempty"`
	Description    *string                `json:"description,omitempty"`
	Status         *domain.FacilityStatus `json:"status,omitempty"`
	OpeningTime    *string                `json:"opening_time,omitempty"`
	ClosingTime    *string                `json:"closing_time,omitempty"`
	Specifications *domain.Specifications `json:"specifications,omitempty"`
}

func (uc *FacilityUseCases) UpdateFacility(ctx context.Context, clubID, id string, dto UpdateFacilityDTO) (*domain.Facility, error) {
	facility, err := uc.repo.GetByID(ctx, clubID, id)
	if err != nil {
		return nil, err
	}
	if facility == nil {
		return nil, errors.New("facility not found")
	}

	if dto.Name != nil {
		facility.Name = *dto.Name
	}
	if dto.Description != nil {
		facility.Description = *dto.Description
	}
	if dto.Status != nil {
		facility.Status = *dto.Status
	}
	if dto.OpeningTime != nil {
		facility.OpeningTime = *dto.OpeningTime
	}
	if dto.ClosingTime != nil {
		facility.ClosingTime = *dto.ClosingTime
	}
	if dto.Specifications != nil {
		// Full replacement of specs for simplicity in MVP, or merge?
		// Let's do partial update if needed, but struct replacement is easier for now.
		// Actually, since specifications is a struct, we should probably merge if we want specific fields,
		// but given the DTO has a pointer to the full struct, we replace it.
		// Ideally we'd have meaningful merge logic, but MVP: replace.
		facility.Specifications = *dto.Specifications
	}

	if err := uc.repo.Update(ctx, facility); err != nil {
		return nil, err
	}

	return facility, nil
}

// Equipment & Loans

type AddEquipmentDTO struct {
	Name      string                    `json:"name" binding:"required"`
	Type      string                    `json:"type" binding:"required"`
	Condition domain.EquipmentCondition `json:"condition" binding:"required"`
}

func (uc *FacilityUseCases) AddEquipment(ctx context.Context, clubID, facilityID string, dto AddEquipmentDTO) (*domain.Equipment, error) {
	// 1. Verify Facility belongs to Club
	fac, err := uc.repo.GetByID(ctx, clubID, facilityID)
	if err != nil {
		return nil, err
	}
	if fac == nil {
		return nil, errors.New("facility not found")
	}

	equipment := &domain.Equipment{
		ID:         uuid.New().String(),
		FacilityID: facilityID,
		Name:       dto.Name,
		Type:       dto.Type,
		Condition:  dto.Condition,
		Status:     "available",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.repo.CreateEquipment(ctx, equipment); err != nil {
		return nil, err
	}
	return equipment, nil
}

func (uc *FacilityUseCases) ListEquipment(ctx context.Context, clubID, facilityID string) ([]*domain.Equipment, error) {
	// Verify Facility first? Or rely on ID.
	// For security, checking if facility belongs to club is better.
	fac, err := uc.repo.GetByID(ctx, clubID, facilityID)
	if err != nil {
		return nil, err
	}
	if fac == nil {
		return nil, errors.New("facility not found")
	}

	return uc.repo.ListEquipmentByFacility(ctx, facilityID)
}

func (uc *FacilityUseCases) LoanEquipment(ctx context.Context, clubID, userID, equipmentID string, expectedReturn time.Time) (*domain.EquipmentLoan, error) {
	// 1. Get Equipment
	// Note: GetEquipmentByID doesn't take ClubID, so we should verify ownership ideally.
	// But Equipment is linked to Facility. We can check Facility -> Club.
	eq, err := uc.repo.GetEquipmentByID(ctx, equipmentID)
	if err != nil {
		return nil, err
	}
	if eq == nil {
		return nil, errors.New("equipment not found")
	}

	// Verify Club (Indirectly via Facility)
	fac, err := uc.repo.GetByID(ctx, clubID, eq.FacilityID)
	if err != nil || fac == nil {
		return nil, errors.New("unauthorized access to equipment (club mismatch)")
	}

	if eq.Status != "available" {
		return nil, errors.New("equipment is not available")
	}

	loan := &domain.EquipmentLoan{
		ID:               uuid.New().String(),
		EquipmentID:      equipmentID,
		UserID:           userID,
		LoanedAt:         time.Now(),
		ExpectedReturnAt: expectedReturn,
		Status:           domain.LoanStatusActive,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Atomic Loan Creation
	if err := uc.repo.LoanEquipmentAtomic(ctx, loan, equipmentID); err != nil {
		return nil, err
	}
	// Note: We don't need to manually update loanRepo or repo.UpdateEquipment as Atomic method does it.

	return loan, nil
}

func (uc *FacilityUseCases) ReturnLoan(ctx context.Context, clubID, loanID, condition string) error {
	loan, err := uc.loanRepo.GetByID(ctx, loanID) // LoanRepo usually NO Context? Or does it need it?
	if err != nil {
		return err
	}
	if loan == nil {
		return errors.New("loan not found")
	}

	// Verify Club via Equipment -> Facility
	eq, err := uc.repo.GetEquipmentByID(ctx, loan.EquipmentID)
	if err != nil || eq == nil {
		return errors.New("equipment not found for loan")
	}
	fac, err := uc.repo.GetByID(ctx, clubID, eq.FacilityID)
	if err != nil || fac == nil {
		return errors.New("unauthorized (club mismatch)")
	}

	if loan.Status != domain.LoanStatusActive && loan.Status != domain.LoanStatusOverdue {
		return errors.New("loan is not active")
	}

	now := time.Now()
	loan.ReturnedAt = &now
	loan.Status = domain.LoanStatusReturned
	loan.ConditionOnReturn = condition
	loan.UpdatedAt = now

	if err := uc.loanRepo.Update(ctx, loan); err != nil {
		return err
	}

	// Update Equipment
	eq.Status = "available"
	if condition != "" {
		eq.Condition = domain.EquipmentCondition(condition) // assuming valid enum string
	}
	if err := uc.repo.UpdateEquipment(ctx, eq); err != nil {
		return err
	}

	return nil
}
