package application

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

type UserUseCases struct {
	repo domain.UserRepository
}

func NewUserUseCases(repo domain.UserRepository) *UserUseCases {
	return &UserUseCases{
		repo: repo,
	}
}

func (uc *UserUseCases) GetProfile(clubID, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, errors.New("invalid user ID")
	}
	return uc.repo.GetByID(clubID, userID)
}

type UpdateProfileDTO struct {
	Name              string                 `json:"name"`
	DateOfBirth       *time.Time             `json:"date_of_birth"`
	SportsPreferences map[string]interface{} `json:"sports_preferences"`
}

func (uc *UserUseCases) UpdateProfile(clubID, userID string, dto UpdateProfileDTO) (*domain.User, error) {
	user, err := uc.repo.GetByID(clubID, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update fields
	if dto.Name != "" {
		user.Name = dto.Name
	}
	if dto.DateOfBirth != nil {
		user.DateOfBirth = dto.DateOfBirth
	}
	if dto.SportsPreferences != nil {
		user.SportsPreferences = dto.SportsPreferences
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	if err := uc.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCases) DeleteUser(clubID, id string, requesterID string) error {
	if id == requesterID {
		return errors.New("cannot delete yourself")
	}
	return uc.repo.Delete(clubID, id)
}

func (uc *UserUseCases) ListUsers(clubID string, limit, offset int, search string) ([]domain.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	filters := make(map[string]interface{})
	if search != "" {
		filters["search"] = search
	}

	return uc.repo.List(clubID, limit, offset, filters)
}

func (uc *UserUseCases) ListChildren(clubID, parentID string) ([]domain.User, error) {
	if parentID == "" {
		return nil, errors.New("parent ID required")
	}
	return uc.repo.FindChildren(clubID, parentID)
}

type RegisterChildDTO struct {
	Name        string     `json:"name"`
	Email       string     `json:"email"` // Optional for very young children? Let's say required for uniqueness or consistency.
	DateOfBirth *time.Time `json:"date_of_birth"`
}

func (uc *UserUseCases) RegisterChild(clubID, parentID string, dto RegisterChildDTO) (*domain.User, error) {
	if parentID == "" {
		return nil, errors.New("parent ID required")
	}
	if dto.Name == "" {
		return nil, errors.New("name is required")
	}

	// Email Logic:
	// If email is provided, check uniqueness (repo check or constraint).
	// If not provided, maybe generate a dummy one? email column is unique not null usually.
	// For now, require email or generate one like child.UUID@placeholder.club
	email := dto.Email
	if email == "" {
		email = "child." + uuid.New().String() + "@placeholder.com"
	}

	child := &domain.User{
		ID:          uuid.New().String(),
		ClubID:      clubID,
		Name:        dto.Name,
		Email:       email,
		Role:        "USER", // Or "CHILD" if we had that role
		DateOfBirth: dto.DateOfBirth,
		ParentID:    &parentID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.repo.Create(child); err != nil {
		return nil, err
	}

	return child, nil
}

func (uc *UserUseCases) GetStats(clubID, userID string) (*domain.UserStats, error) {
	user, err := uc.repo.GetByID(clubID, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user.Stats, nil
}

func (uc *UserUseCases) GetWallet(clubID, userID string) (*domain.Wallet, error) {
	user, err := uc.repo.GetByID(clubID, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user.Wallet, nil
}

func (uc *UserUseCases) CreateManualDebt(clubID, userID string, amount float64, description string, adminID string) error {
	user, err := uc.repo.GetByID(clubID, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Assuming Wallet is loaded or we need to init it.
	// For MVP, if Wallet is nil, we might need to create it, but `gamification.go` defines it.
	// We will append to Transactions and Update User.

	if user.Wallet == nil {
		user.Wallet = &domain.Wallet{
			ID:           uuid.New(),
			UserID:       userID,
			Balance:      0,
			Transactions: []domain.Transaction{},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
	}

	// Manual Debt increases negative balance or reduces positive?
	// Usually "Debt" means they owe money. So Balance goes down (negative) or we have a Debt field.
	// Let's assume Balance represents what they have. If they have debt, it's negative.
	user.Wallet.Balance -= amount

	transaction := domain.Transaction{
		ID:          uuid.New().String(),
		Type:        "MANUAL_DEBT",
		Amount:      amount,
		Description: description + " (by Admin " + adminID + ")",
		Date:        time.Now(),
	}

	user.Wallet.Transactions = append(user.Wallet.Transactions, transaction)
	user.Wallet.UpdatedAt = time.Now()

	// Since Wallet is part of User struct in our Domain (aggregates), updating User Updates Wallet (cascade).
	return uc.repo.Update(user)
}

func (uc *UserUseCases) UpdateEmergencyInfo(clubID, userID string, contactName, contactPhone, insuranceProvider, insuranceNumber string) error {
	user, err := uc.repo.GetByID(clubID, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.EmergencyContactName = contactName
	user.EmergencyContactPhone = contactPhone
	user.InsuranceProvider = insuranceProvider
	user.InsuranceNumber = insuranceNumber
	user.UpdatedAt = time.Now()

	return uc.repo.Update(user)
}

func (uc *UserUseCases) LogIncident(clubID, injuredUserID, description, actionTaken, witnesses, reporterID string) (*domain.IncidentLog, error) {
	// injuredUserID can be empty if it's a visitor, but if provided, validate existence?
	// For now, trust input or loose coupling.

	incident := &domain.IncidentLog{
		ID:          uuid.New(),
		ClubID:      clubID,
		Description: description,
		Witnesses:   witnesses,
		ActionTaken: actionTaken,
		ReportedAt:  time.Now(),
		CreatedBy:   reporterID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if injuredUserID != "" {
		incident.InjuredUserID = &injuredUserID
	}

	if err := uc.repo.CreateIncident(incident); err != nil {
		return nil, err
	}

	return incident, nil
}
