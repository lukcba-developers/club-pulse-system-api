package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
)

type ClubUseCases struct {
	sponsorRepo domain.SponsorRepository
	clubRepo    domain.ClubRepository
}

func NewClubUseCases(sponsorRepo domain.SponsorRepository, clubRepo domain.ClubRepository) *ClubUseCases {
	return &ClubUseCases{
		sponsorRepo: sponsorRepo,
		clubRepo:    clubRepo,
	}
}

// --- Club Management (Super Admin) ---

func (uc *ClubUseCases) CreateClub(name, domainStr, settings string) (*domain.Club, error) {
	club := &domain.Club{
		ID:        uuid.New().String(), // OR use specific slug logic if preferred
		Name:      name,
		Domain:    domainStr,
		Status:    domain.ClubStatusActive,
		Settings:  settings,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := uc.clubRepo.Create(club); err != nil {
		return nil, err
	}
	return club, nil
}

func (uc *ClubUseCases) GetClub(id string) (*domain.Club, error) {
	return uc.clubRepo.GetByID(id)
}

func (uc *ClubUseCases) ListClubs(limit, offset int) ([]domain.Club, error) {
	return uc.clubRepo.List(limit, offset)
}

func (uc *ClubUseCases) UpdateClub(id string, name, domainStr, settings string, status domain.ClubStatus) (*domain.Club, error) {
	club, err := uc.clubRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if club == nil {
		return nil, nil // Or error
	}

	if name != "" {
		club.Name = name
	}
	if domainStr != "" {
		club.Domain = domainStr
	}
	if settings != "" {
		club.Settings = settings
	}
	if status != "" {
		club.Status = status
	}
	club.UpdatedAt = time.Now()

	if err := uc.clubRepo.Update(club); err != nil {
		return nil, err
	}
	return club, nil
}

func (uc *ClubUseCases) DeleteClub(id string) error {
	return uc.clubRepo.Delete(id)
}

// --- Sponsor Management ---

func (uc *ClubUseCases) RegisterSponsor(clubID, name, contactInfo, logoURL string) (*domain.Sponsor, error) {
	sponsor := &domain.Sponsor{
		ID:          uuid.New(),
		ClubID:      clubID,
		Name:        name,
		ContactInfo: contactInfo,
		LogoURL:     logoURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.sponsorRepo.CreateSponsor(sponsor); err != nil {
		return nil, err
	}
	return sponsor, nil
}

func (uc *ClubUseCases) CreateAdPlacement(sponsorID string, locationType domain.LocationType, detail string, endDate time.Time, amount float64) (*domain.AdPlacement, error) {
	sponsorUUID, err := uuid.Parse(sponsorID)
	if err != nil {
		return nil, err
	}

	ad := &domain.AdPlacement{
		ID:             uuid.New(),
		SponsorID:      sponsorUUID,
		LocationType:   locationType,
		LocationDetail: detail,
		ContractEnd:    endDate,
		AmountPaid:     amount,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := uc.sponsorRepo.CreateAdPlacement(ad); err != nil {
		return nil, err
	}
	return ad, nil
}

func (uc *ClubUseCases) GetActiveAds(clubID string) ([]domain.AdPlacement, error) {
	return uc.sponsorRepo.GetActiveAds(clubID)
}
