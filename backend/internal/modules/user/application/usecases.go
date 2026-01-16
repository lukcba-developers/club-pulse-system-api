package application

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCases struct {
	repo            domain.UserRepository
	familyGroupRepo domain.FamilyGroupRepository
}

func NewUserUseCases(repo domain.UserRepository, familyGroupRepo domain.FamilyGroupRepository) *UserUseCases {
	return &UserUseCases{
		repo:            repo,
		familyGroupRepo: familyGroupRepo,
	}
}

func (uc *UserUseCases) GetProfile(ctx context.Context, clubID, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, errors.New("invalid user ID")
	}
	return uc.repo.GetByID(ctx, clubID, userID)
}

type UpdateProfileDTO struct {
	Name              string                 `json:"name"`
	DateOfBirth       *time.Time             `json:"date_of_birth"`
	SportsPreferences map[string]interface{} `json:"sports_preferences"`
}

func (uc *UserUseCases) UpdateProfile(ctx context.Context, clubID, userID string, dto UpdateProfileDTO) (*domain.User, error) {
	user, err := uc.repo.GetByID(ctx, clubID, userID)
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

	if err := uc.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCases) DeleteUser(ctx context.Context, clubID, id string, requesterID string) error {
	if id == requesterID {
		return errors.New("cannot delete yourself")
	}
	return uc.repo.Delete(ctx, clubID, id)
}

func (uc *UserUseCases) ListUsers(ctx context.Context, clubID string, limit, offset int, search string) ([]domain.User, error) {
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

	return uc.repo.List(ctx, clubID, limit, offset, filters)
}

func (uc *UserUseCases) ListChildren(ctx context.Context, clubID, parentID string) ([]domain.User, error) {
	if parentID == "" {
		return nil, errors.New("parent ID required")
	}
	return uc.repo.FindChildren(ctx, clubID, parentID)
}

type RegisterChildDTO struct {
	Name        string     `json:"name"`
	Email       string     `json:"email"` // Optional for very young children? Let's say required for uniqueness or consistency.
	DateOfBirth *time.Time `json:"date_of_birth"`
}

func (uc *UserUseCases) RegisterChild(ctx context.Context, clubID, parentID string, dto RegisterChildDTO) (*domain.User, error) {
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

	if err := uc.repo.Create(ctx, child); err != nil {
		return nil, err
	}

	return child, nil
}

type RegisterDependentDTO struct {
	ParentEmail          string                 `json:"parent_email"`
	ParentName           string                 `json:"parent_name"`
	ParentPhone          string                 `json:"parent_phone"`
	ChildName            string                 `json:"child_name"`
	ChildSurname         string                 `json:"child_surname"`
	ChildDOB             *time.Time             `json:"child_dob"`
	SportsPreferences    map[string]interface{} `json:"sports_preferences"`
	Password             string                 `json:"password" binding:"required,min=8"`
	AcceptTerms          bool                   `json:"accept_terms" binding:"required"`
	PrivacyPolicyVersion string                 `json:"privacy_policy_version"`
	ParentalConsent      bool                   `json:"parental_consent" binding:"required"`
}

func (uc *UserUseCases) RegisterDependent(ctx context.Context, clubID string, dto RegisterDependentDTO) (*domain.User, error) {
	// Input Validation
	if dto.ParentEmail == "" {
		return nil, errors.New("parent email is required")
	}
	if dto.ParentName == "" {
		return nil, errors.New("parent name is required")
	}
	if dto.ChildName == "" {
		return nil, errors.New("child name is required")
	}

	// GetByEmail now scopes by clubID for proper tenant isolation
	parent, err := uc.repo.GetByEmail(ctx, clubID, dto.ParentEmail)
	if err != nil {
		return nil, err
	}

	var parentID string
	if parent == nil {
		// Validations for new Parent
		if dto.Password == "" {
			return nil, errors.New("password is required for new registration")
		}
		// Password strength validation (matching auth module standards)
		if len(dto.Password) < 8 {
			return nil, errors.New("password must be at least 8 characters long")
		}
		hasUpper := false
		hasDigit := false
		for _, c := range dto.Password {
			if c >= 'A' && c <= 'Z' {
				hasUpper = true
			}
			if c >= '0' && c <= '9' {
				hasDigit = true
			}
		}
		if !hasUpper || !hasDigit {
			return nil, errors.New("password must contain at least one uppercase letter and one number")
		}

		if !dto.AcceptTerms {
			return nil, errors.New("terms acceptance is required")
		}
		if !dto.ParentalConsent {
			return nil, errors.New("parental consent is required")
		}

		// Hash Password
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to process password")
		}

		now := time.Now()
		privacyVersion := dto.PrivacyPolicyVersion
		if privacyVersion == "" {
			privacyVersion = "2026-01" // Default
		}

		// Create Parent
		newParent := &domain.User{
			ID:                    uuid.New().String(),
			ClubID:                clubID,
			Name:                  dto.ParentName,
			Email:                 dto.ParentEmail,
			Password:              string(hashedBytes), // Securely stored
			Role:                  domain.RoleMember,
			CreatedAt:             now,
			UpdatedAt:             now,
			EmergencyContactPhone: dto.ParentPhone,
			TermsAcceptedAt:       &now,
			PrivacyPolicyVersion:  privacyVersion,
			ParentalConsentAt:     &now, // Record when they gave consent for dependents
		}
		if err := uc.repo.Create(ctx, newParent); err != nil {
			return nil, err
		}
		parentID = newParent.ID
	} else {
		parentID = parent.ID
		// Verify parent belongs to the same club (tenant isolation)
		if parent.ClubID != clubID {
			return nil, errors.New("parent account found in a different club")
		}
		// Update ParentalConsentAt timestamp to show "Fresh" consent for this action.
		now := time.Now()
		parent.ParentalConsentAt = &now
		if err := uc.repo.Update(ctx, parent); err != nil {
			// Log error but don't fail the operation - consent update is non-critical
			// TODO: Replace with proper logger in production
			_ = err // Acknowledge error for linter
		}
	}

	// Create Child
	child := &domain.User{
		ID:                uuid.New().String(),
		ClubID:            clubID,
		Name:              dto.ChildName + " " + dto.ChildSurname,
		ParentID:          &parentID,
		SportsPreferences: dto.SportsPreferences,
		Role:              domain.RoleMember,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Email:             "child." + uuid.New().String() + "@placeholder.com",
		DateOfBirth:       dto.ChildDOB,
	}

	if err := uc.repo.Create(ctx, child); err != nil {
		return nil, err
	}

	return child, nil
}

func (uc *UserUseCases) GetStats(ctx context.Context, clubID, userID string) (*domain.UserStats, error) {
	user, err := uc.repo.GetByID(ctx, clubID, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if user.Stats != nil {
		user.Stats.NextLevelXP = user.Stats.CalculateNextLevelXP()
	}
	return user.Stats, nil
}

func (uc *UserUseCases) GetWallet(ctx context.Context, clubID, userID string) (*domain.Wallet, error) {
	user, err := uc.repo.GetByID(ctx, clubID, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user.Wallet, nil
}

func (uc *UserUseCases) CreateManualDebt(ctx context.Context, clubID, userID string, amount float64, description string, adminID string) error {
	user, err := uc.repo.GetByID(ctx, clubID, userID)
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
	return uc.repo.Update(ctx, user)
}

func (uc *UserUseCases) UpdateEmergencyInfo(ctx context.Context, clubID, userID string, contactName, contactPhone, insuranceProvider, insuranceNumber string) error {
	user, err := uc.repo.GetByID(ctx, clubID, userID)
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

	return uc.repo.Update(ctx, user)
}

func (uc *UserUseCases) LogIncident(ctx context.Context, clubID, injuredUserID, description, actionTaken, witnesses, reporterID string) (*domain.IncidentLog, error) {
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

	if err := uc.repo.CreateIncident(ctx, incident); err != nil {
		return nil, err
	}

	return incident, nil
}

func (uc *UserUseCases) UpdateMatchStats(ctx context.Context, clubID, userID string, won bool, xpGained int) error {
	user, err := uc.repo.GetByID(ctx, clubID, userID)
	if err != nil {
		return err
	}
	if user == nil {
		// If user not found, maybe just warn or skip? returning error is safer.
		return errors.New("user not found")
	}

	if user.Stats == nil {
		// Initialize stats if missing
		now := time.Now()
		user.Stats = &domain.UserStats{
			ID:            uuid.New(),
			UserID:        userID,
			MatchesPlayed: 0,
			MatchesWon:    0,
			RankingPoints: 0,
			Level:         1,
			Experience:    0,
			CurrentStreak: 0,
			LongestStreak: 0,
			TotalXP:       0,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
	}

	user.Stats.MatchesPlayed++

	// Apply streak multiplier to XP
	finalXP := domain.CalculateXPWithStreak(xpGained, user.Stats.CurrentStreak)
	user.Stats.Experience += finalXP
	user.Stats.TotalXP += finalXP

	if won {
		user.Stats.MatchesWon++
		user.Stats.RankingPoints += 3 // Example logic
		// Bonus XP for winning
		winBonus := domain.GetXPForAction(domain.XPMatchWon)
		winBonusWithStreak := domain.CalculateXPWithStreak(winBonus, user.Stats.CurrentStreak)
		user.Stats.Experience += winBonusWithStreak
		user.Stats.TotalXP += winBonusWithStreak
	} else {
		user.Stats.RankingPoints += 1 // Participation points
	}

	// Update streak (match counts as activity)
	uc.updateStreak(user.Stats)

	// Level up logic - Exponential formula: 500 * (1.15 ^ Level)
	for {
		requiredXP := CalculateRequiredXP(user.Stats.Level)
		if user.Stats.Experience >= requiredXP {
			user.Stats.Level++
			user.Stats.Experience -= requiredXP
			// Could emit LevelUp event here for notifications
		} else {
			break
		}
	}

	user.Stats.UpdatedAt = time.Now()

	return uc.repo.Update(ctx, user)
}

// CalculateRequiredXP returns XP needed to advance from the given level.
// Uses exponential formula: 500 * (1.15 ^ level)
func CalculateRequiredXP(level int) int {
	return int(500 * math.Pow(1.15, float64(level)))
}

// updateStreak updates the user's streak based on activity today.
func (uc *UserUseCases) updateStreak(stats *domain.UserStats) {
	today := time.Now().Truncate(24 * time.Hour)

	if stats.LastActivityDate == nil {
		// First activity ever
		stats.CurrentStreak = 1
		stats.LongestStreak = 1
		stats.LastActivityDate = &today
		return
	}

	lastActivity := stats.LastActivityDate.Truncate(24 * time.Hour)
	daysSinceLastActivity := int(today.Sub(lastActivity).Hours() / 24)

	switch daysSinceLastActivity {
	case 0:
		// Same day, no change to streak
		return
	case 1:
		// Consecutive day, increment streak
		stats.CurrentStreak++
		if stats.CurrentStreak > stats.LongestStreak {
			stats.LongestStreak = stats.CurrentStreak
		}
	default:
		// Streak broken, reset to 1
		stats.CurrentStreak = 1
	}

	stats.LastActivityDate = &today
}

// --- Family Group Use Cases ---

func (uc *UserUseCases) CreateFamilyGroup(ctx context.Context, clubID, headUserID, name string) (*domain.FamilyGroup, error) {
	if uc.familyGroupRepo == nil {
		return nil, errors.New("family groups not enabled")
	}
	if name == "" {
		return nil, errors.New("family group name is required")
	}

	// Check if user already has a family group as head
	existing, _ := uc.familyGroupRepo.GetByHeadUserID(ctx, clubID, headUserID)
	if existing != nil {
		return nil, errors.New("user already has a family group")
	}

	group := &domain.FamilyGroup{
		ClubID:     clubID,
		Name:       name,
		HeadUserID: headUserID,
	}

	if err := uc.familyGroupRepo.Create(ctx, group); err != nil {
		return nil, err
	}

	// Add head user to the group
	_ = uc.familyGroupRepo.AddMember(ctx, clubID, group.ID, headUserID)

	return group, nil
}

func (uc *UserUseCases) GetMyFamilyGroup(ctx context.Context, clubID, userID string) (*domain.FamilyGroup, error) {
	if uc.familyGroupRepo == nil {
		return nil, errors.New("family groups not enabled")
	}
	return uc.familyGroupRepo.GetByMemberID(ctx, clubID, userID)
}

func (uc *UserUseCases) AddFamilyMember(ctx context.Context, clubID string, groupID uuid.UUID, memberUserID string) error {
	if uc.familyGroupRepo == nil {
		return errors.New("family groups not enabled")
	}
	return uc.familyGroupRepo.AddMember(ctx, clubID, groupID, memberUserID)
}

// AddFamilyMemberSecure validates that the requesting user is the family head before adding members.
// SECURITY FIX (VUL-002): Prevents IDOR by ensuring ownership validation.
func (uc *UserUseCases) AddFamilyMemberSecure(ctx context.Context, clubID string, groupID uuid.UUID, memberUserID, requestingUserID string) error {
	if uc.familyGroupRepo == nil {
		return errors.New("family groups not enabled")
	}

	// Validate ownership - only HeadUserID can add members
	group, err := uc.familyGroupRepo.GetByID(ctx, clubID, groupID)
	if err != nil {
		return errors.New("group not found")
	}
	if group == nil {
		return errors.New("group not found")
	}

	if group.HeadUserID != requestingUserID {
		return errors.New("only the family head can add members")
	}

	return uc.familyGroupRepo.AddMember(ctx, clubID, groupID, memberUserID)
}

func (uc *UserUseCases) RemoveFamilyMember(ctx context.Context, clubID string, groupID uuid.UUID, memberUserID string) error {
	if uc.familyGroupRepo == nil {
		return errors.New("family groups not enabled")
	}
	return uc.familyGroupRepo.RemoveMember(ctx, clubID, groupID, memberUserID)
}

// --- GDPR Compliance Use Cases ---

// GDPRExportData represents the data package for right to portability (GDPR Article 20)
type GDPRExportData struct {
	ExportedAt  string                   `json:"exported_at"`
	UserProfile map[string]interface{}   `json:"user_profile"`
	Children    []map[string]interface{} `json:"children,omitempty"`
	FamilyGroup *map[string]interface{}  `json:"family_group,omitempty"`
}

// ExportUserData implements GDPR Article 20 - Right to data portability
// Returns all personal data for the user in a structured, portable format
func (uc *UserUseCases) ExportUserData(ctx context.Context, clubID, userID string) (*GDPRExportData, error) {
	user, err := uc.repo.GetByID(ctx, clubID, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Build user profile export (excluding sensitive internal fields)
	profile := map[string]interface{}{
		"id":                      user.ID,
		"name":                    user.Name,
		"email":                   user.Email,
		"role":                    user.Role,
		"club_id":                 user.ClubID,
		"created_at":              user.CreatedAt.Format(time.RFC3339),
		"updated_at":              user.UpdatedAt.Format(time.RFC3339),
		"emergency_contact_name":  user.EmergencyContactName,
		"emergency_contact_phone": user.EmergencyContactPhone,
		"insurance_provider":      user.InsuranceProvider,
		"insurance_number":        user.InsuranceNumber,
	}

	if user.DateOfBirth != nil {
		profile["date_of_birth"] = user.DateOfBirth.Format("2006-01-02")
	}
	if user.SportsPreferences != nil {
		profile["sports_preferences"] = user.SportsPreferences
	}
	if user.MedicalCertStatus != nil {
		profile["medical_cert_status"] = string(*user.MedicalCertStatus)
	}
	if user.MedicalCertExpiry != nil {
		profile["medical_cert_expiry"] = user.MedicalCertExpiry.Format("2006-01-02")
	}
	if user.TermsAcceptedAt != nil {
		profile["terms_accepted_at"] = user.TermsAcceptedAt.Format(time.RFC3339)
	}
	if user.PrivacyPolicyVersion != "" {
		profile["privacy_policy_version"] = user.PrivacyPolicyVersion
	}

	export := &GDPRExportData{
		ExportedAt:  time.Now().Format(time.RFC3339),
		UserProfile: profile,
	}

	// Export children data if any
	children, err := uc.repo.FindChildren(ctx, clubID, userID)
	if err == nil && len(children) > 0 {
		childExports := make([]map[string]interface{}, len(children))
		for i, child := range children {
			childExports[i] = map[string]interface{}{
				"id":    child.ID,
				"name":  child.Name,
				"email": child.Email,
			}
			if child.DateOfBirth != nil {
				childExports[i]["date_of_birth"] = child.DateOfBirth.Format("2006-01-02")
			}
		}
		export.Children = childExports
	}

	// Export family group if member
	if uc.familyGroupRepo != nil {
		group, err := uc.familyGroupRepo.GetByMemberID(ctx, clubID, userID)
		if err == nil && group != nil {
			fg := map[string]interface{}{
				"id":      group.ID.String(),
				"name":    group.Name,
				"is_head": group.HeadUserID == userID,
			}
			export.FamilyGroup = &fg
		}
	}

	return export, nil
}

// DeleteUserGDPR implements GDPR Article 17 - Right to erasure
// Uses anonymization instead of simple deletion
func (uc *UserUseCases) DeleteUserGDPR(ctx context.Context, clubID, id string, requesterID string) error {
	if id == requesterID {
		return errors.New("cannot delete yourself")
	}
	return uc.repo.AnonymizeForGDPR(ctx, clubID, id)
}
