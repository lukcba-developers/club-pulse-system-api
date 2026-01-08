package application

import (
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
	Type           domain.FacilityType   `json:"type" binding:"required"`
	Capacity       int                   `json:"capacity" binding:"required,min=1"`
	HourlyRate     float64               `json:"hourly_rate" binding:"required,min=0"`
	OpeningHour    int                   `json:"opening_hour"`
	ClosingHour    int                   `json:"closing_hour"`
	Specifications domain.Specifications `json:"specifications"`
	Location       domain.Location       `json:"location"`
}

func (uc *FacilityUseCases) CreateFacility(clubID string, dto CreateFacilityDTO) (*domain.Facility, error) {
	// Defaults
	opening := dto.OpeningHour
	if opening == 0 {
		opening = 8
	}
	closing := dto.ClosingHour
	if closing == 0 {
		closing = 23
	}

	facility := &domain.Facility{
		ID:             uuid.New().String(),
		ClubID:         clubID,
		Name:           dto.Name,
		Type:           dto.Type,
		Status:         domain.FacilityStatusActive,
		Capacity:       dto.Capacity,
		HourlyRate:     dto.HourlyRate,
		OpeningHour:    opening,
		ClosingHour:    closing,
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

func (uc *FacilityUseCases) ListFacilities(clubID string, limit, offset int) ([]*domain.Facility, error) {
	if limit <= 0 {
		limit = 10
	}
	return uc.repo.List(clubID, limit, offset)
}

func (uc *FacilityUseCases) GetFacility(clubID, id string) (*domain.Facility, error) {
	if id == "" {
		return nil, errors.New("invalid ID")
	}
	return uc.repo.GetByID(clubID, id)
}

type UpdateFacilityDTO struct {
	Name           *string                `json:"name,omitempty"`
	Status         *domain.FacilityStatus `json:"status,omitempty"`
	OpeningHour    *int                   `json:"opening_hour,omitempty"`
	ClosingHour    *int                   `json:"closing_hour,omitempty"`
	Specifications *domain.Specifications `json:"specifications,omitempty"`
}

func (uc *FacilityUseCases) UpdateFacility(clubID, id string, dto UpdateFacilityDTO) (*domain.Facility, error) {
	facility, err := uc.repo.GetByID(clubID, id)
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
	if dto.OpeningHour != nil {
		facility.OpeningHour = *dto.OpeningHour
	}
	if dto.ClosingHour != nil {
		facility.ClosingHour = *dto.ClosingHour
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

// Equipment & Loans

type AddEquipmentDTO struct {
	Name      string                    `json:"name" binding:"required"`
	Type      string                    `json:"type" binding:"required"`
	Condition domain.EquipmentCondition `json:"condition" binding:"required"`
}

func (uc *FacilityUseCases) AddEquipment(clubID, facilityID string, dto AddEquipmentDTO) (*domain.Equipment, error) {
	// 1. Verify Facility belongs to Club
	fac, err := uc.repo.GetByID(clubID, facilityID)
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

	if err := uc.repo.CreateEquipment(equipment); err != nil {
		return nil, err
	}
	return equipment, nil
}

func (uc *FacilityUseCases) ListEquipment(clubID, facilityID string) ([]*domain.Equipment, error) {
	// Verify Facility first? Or rely on ID.
	// For security, checking if facility belongs to club is better.
	fac, err := uc.repo.GetByID(clubID, facilityID)
	if err != nil {
		return nil, err
	}
	if fac == nil {
		return nil, errors.New("facility not found")
	}

	return uc.repo.ListEquipmentByFacility(facilityID)
}

func (uc *FacilityUseCases) LoanEquipment(clubID, userID, equipmentID string, expectedReturn time.Time) (*domain.EquipmentLoan, error) {
	// 1. Get Equipment
	// Note: GetEquipmentByID doesn't take ClubID, so we should verify ownership ideally.
	// But Equipment is linked to Facility. We can check Facility -> Club.
	eq, err := uc.repo.GetEquipmentByID(equipmentID)
	if err != nil {
		return nil, err
	}
	if eq == nil {
		return nil, errors.New("equipment not found")
	}

	// Verify Club (Indirectly via Facility)
	fac, err := uc.repo.GetByID(clubID, eq.FacilityID)
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
	if err := uc.repo.LoanEquipmentAtomic(loan, equipmentID); err != nil {
		return nil, err
	}
	// Note: We don't need to manually update loanRepo or repo.UpdateEquipment as Atomic method does it.

	return loan, nil
}

func (uc *FacilityUseCases) ReturnLoan(clubID, loanID, condition string) error {
	loan, err := uc.loanRepo.GetByID(loanID)
	if err != nil {
		return err
	}
	if loan == nil {
		return errors.New("loan not found")
	}

	// Verify Club via Equipment -> Facility
	eq, err := uc.repo.GetEquipmentByID(loan.EquipmentID)
	if err != nil || eq == nil {
		return errors.New("equipment not found for loan")
	}
	fac, err := uc.repo.GetByID(clubID, eq.FacilityID)
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

	if err := uc.loanRepo.Update(loan); err != nil {
		return err
	}

	// Update Equipment
	eq.Status = "available"
	if condition != "" {
		eq.Condition = domain.EquipmentCondition(condition) // assuming valid enum string
	}
	if err := uc.repo.UpdateEquipment(eq); err != nil {
		return err
	}

	return nil
}
