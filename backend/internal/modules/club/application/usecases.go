package application

import (
	"context"
	"strings"
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

func (uc *ClubUseCases) PublishNews(ctx context.Context, clubID, title, content, imageURL string, notify bool) (*domain.News, error) {
	news := &domain.News{
		ClubID:    clubID,
		Title:     title,
		Content:   content,
		ImageURL:  imageURL,
		Published: true, // Auto publish for now
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.newsRepo.CreateNews(ctx, news); err != nil {
		return nil, err
	}

	if notify && uc.notifier != nil {
		// Broadcast Notification Async with restricted concurrency
		go func() {
			bgCtx := context.Background()
			emails, err := uc.clubRepo.GetMemberEmails(bgCtx, clubID)
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
						Title:       subject,
						Body:        body,
						ActionURL:   "https://" + clubID + ".club-pulse.com/news", // Example deep link
					})
				}(email)
			}
			wg.Wait()
		}()
	}

	return news, nil
}

func (uc *ClubUseCases) GetPublicNews(ctx context.Context, slug string) ([]domain.News, error) {
	club, err := uc.clubRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return uc.newsRepo.GetPublicNewsByClub(ctx, club.ID, 10, 0)
}

// ... (Sponsor methods below)

// --- Club Management (Super Admin) ---

func (uc *ClubUseCases) CreateClub(ctx context.Context, name, slug, domainStr, logoURL, primaryColor, secondaryColor, contactEmail, contactPhone, themeConfig, settings string) (*domain.Club, error) {
	if slug == "" {
		// Simple slug generation: lower case, replace spaces with hyphens, remove non-alphanumeric
		// In production, checking for uniqueness loop might be needed or let DB fail
		slug = strings.ToLower(name)
		slug = strings.ReplaceAll(slug, " ", "-")
		// Remove special chars could be done with regex, keeping it simple for now or relying on DB constraint
	}

	club := &domain.Club{
		ID:             uuid.New().String(),
		Name:           name,
		Slug:           slug,
		Domain:         domainStr,
		LogoURL:        logoURL,
		PrimaryColor:   primaryColor,
		SecondaryColor: secondaryColor,
		ContactEmail:   contactEmail,
		ContactPhone:   contactPhone,
		ThemeConfig:    themeConfig,
		Status:         domain.ClubStatusActive,
		Settings:       settings,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := uc.clubRepo.Create(ctx, club); err != nil {
		return nil, err
	}
	return club, nil
}

func (uc *ClubUseCases) GetClub(ctx context.Context, id string) (*domain.Club, error) {
	return uc.clubRepo.GetByID(ctx, id)
}

func (uc *ClubUseCases) GetClubBySlug(ctx context.Context, slug string) (*domain.Club, error) {
	return uc.clubRepo.GetBySlug(ctx, slug)
}

func (uc *ClubUseCases) ListClubs(ctx context.Context, limit, offset int) ([]domain.Club, error) {
	return uc.clubRepo.List(ctx, limit, offset)
}

func (uc *ClubUseCases) UpdateClub(ctx context.Context, id string, name, domainStr, logoURL, primaryColor, secondaryColor, contactEmail, contactPhone, themeConfig, settings string, status domain.ClubStatus) (*domain.Club, error) {
	club, err := uc.clubRepo.GetByID(ctx, id)
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
	if logoURL != "" {
		club.LogoURL = logoURL
	}
	if primaryColor != "" {
		club.PrimaryColor = primaryColor
	}
	if secondaryColor != "" {
		club.SecondaryColor = secondaryColor
	}
	if contactEmail != "" {
		club.ContactEmail = contactEmail
	}
	if contactPhone != "" {
		club.ContactPhone = contactPhone
	}
	if themeConfig != "" {
		club.ThemeConfig = themeConfig
	}
	if settings != "" {
		club.Settings = settings
	}
	if status != "" {
		club.Status = status
	}
	club.UpdatedAt = time.Now()

	if err := uc.clubRepo.Update(ctx, club); err != nil {
		return nil, err
	}
	return club, nil
}

func (uc *ClubUseCases) DeleteClub(ctx context.Context, id string) error {
	return uc.clubRepo.Delete(ctx, id)
}

// --- Sponsor Management ---

func (uc *ClubUseCases) RegisterSponsor(ctx context.Context, clubID, name, contactInfo, logoURL string) (*domain.Sponsor, error) {
	sponsor := &domain.Sponsor{
		ID:          uuid.New(),
		ClubID:      clubID,
		Name:        name,
		ContactInfo: contactInfo,
		LogoURL:     logoURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.sponsorRepo.CreateSponsor(ctx, sponsor); err != nil {
		return nil, err
	}
	return sponsor, nil
}

func (uc *ClubUseCases) CreateAdPlacement(ctx context.Context, sponsorID string, locationType domain.LocationType, detail string, endDate time.Time, amount float64) (*domain.AdPlacement, error) {
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

	if err := uc.sponsorRepo.CreateAdPlacement(ctx, ad); err != nil {
		return nil, err
	}
	return ad, nil
}

func (uc *ClubUseCases) GetActiveAds(ctx context.Context, clubID string) ([]domain.AdPlacement, error) {
	return uc.sponsorRepo.GetActiveAds(ctx, clubID)
}
