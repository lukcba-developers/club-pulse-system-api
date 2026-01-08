package application

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	notification "github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
)

type ClubUseCases struct {
	sponsorRepo domain.SponsorRepository
	clubRepo    domain.ClubRepository
	newsRepo    domain.NewsRepository
	notifier    *notification.NotificationService
}

func NewClubUseCases(
	sponsorRepo domain.SponsorRepository,
	clubRepo domain.ClubRepository,
	newsRepo domain.NewsRepository,
	notifier *notification.NotificationService,
) *ClubUseCases {
	return &ClubUseCases{
		sponsorRepo: sponsorRepo,
		clubRepo:    clubRepo,
		newsRepo:    newsRepo,
		notifier:    notifier,
	}
}

// ... (Club methods remain)

// --- News Management ---

func (uc *ClubUseCases) PublishNews(clubID, title, content, imageURL string, notify bool) (*domain.News, error) {
	news := &domain.News{
		ClubID:    clubID,
		Title:     title,
		Content:   content,
		ImageURL:  imageURL,
		Published: true, // Auto publish for now
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.newsRepo.CreateNews(news); err != nil {
		return nil, err
	}

	if notify && uc.notifier != nil {
		// Broadcast Notification Async with restricted concurrency
		go func() {
			bgCtx := context.Background()
			emails, err := uc.clubRepo.GetMemberEmails(clubID)
			if err != nil {
				// In production, use a proper logger
				return
			}

			// Semaphore to limit concurrency (e.g., 10 concurrent sends)
			sem := make(chan struct{}, 10)
			var wg sync.WaitGroup

			for _, email := range emails {
				wg.Add(1)
				sem <- struct{}{} // Acquire token

				go func(recipient string) {
					defer wg.Done()
					defer func() { <-sem }() // Release token

					subject := "Nueva noticia: " + title
					body := "Hola,\n\nNueva noticia en tu club:\n\n" + title + "\n\n" + content

					// Fire and forget, but now controlled
					_ = uc.notifier.Send(bgCtx, notification.Notification{
						RecipientID: recipient,
						Type:        notification.NotificationTypeEmail,
						Subject:     subject,
						Message:     body,
					})
				}(email)
			}
			wg.Wait()
		}()
	}

	return news, nil
}

func (uc *ClubUseCases) GetPublicNews(slug string) ([]domain.News, error) {
	club, err := uc.clubRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	return uc.newsRepo.GetPublicNewsByClub(club.ID, 10, 0)
}

// ... (Sponsor methods below)

// --- Club Management (Super Admin) ---

func (uc *ClubUseCases) CreateClub(name, slug, domainStr, settings string) (*domain.Club, error) {
	club := &domain.Club{
		ID:        uuid.New().String(),
		Name:      name,
		Slug:      slug,
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

func (uc *ClubUseCases) GetClubBySlug(slug string) (*domain.Club, error) {
	return uc.clubRepo.GetBySlug(slug)
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
