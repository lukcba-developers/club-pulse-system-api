package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	"gorm.io/gorm"
)

type PostgresClubRepository struct {
	db *gorm.DB
}

func NewPostgresClubRepository(db *gorm.DB) *PostgresClubRepository {
	return &PostgresClubRepository{db: db}
}

func (r *PostgresClubRepository) CreateSponsor(ctx context.Context, sponsor *domain.Sponsor) error {
	return r.db.WithContext(ctx).Create(sponsor).Error
}

func (r *PostgresClubRepository) CreateAdPlacement(ctx context.Context, ad *domain.AdPlacement) error {
	return r.db.WithContext(ctx).Create(ad).Error
}

func (r *PostgresClubRepository) GetActiveAds(ctx context.Context, clubID string) ([]domain.AdPlacement, error) {
	var ads []domain.AdPlacement
	// Join with Sponsor to filter by ClubID if needed, assuming Sponsor has ClubID
	now := time.Now()
	err := r.db.WithContext(ctx).Preload("Sponsor").
		Joins("JOIN sponsors ON sponsors.id = ad_placements.sponsor_id").
		Where("sponsors.club_id = ?", clubID).
		Where("ad_placements.is_active = ?", true).
		Where("ad_placements.contract_end >= ?", now).
		Find(&ads).Error
	return ads, err
}

// --- Club CRUD Implementation (Restoring lost functionality) ---

func (r *PostgresClubRepository) Create(ctx context.Context, club *domain.Club) error {
	return r.db.WithContext(ctx).Create(club).Error
}

func (r *PostgresClubRepository) GetByID(ctx context.Context, id string) (*domain.Club, error) {
	var club domain.Club
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&club).Error; err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *PostgresClubRepository) GetBySlug(ctx context.Context, slug string) (*domain.Club, error) {
	var club domain.Club
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&club).Error; err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *PostgresClubRepository) GetMemberEmails(ctx context.Context, clubID string) ([]string, error) {
	var emails []string
	// Assuming 'users' table has 'email' and 'club_id'.
	// Note: 'users' table name is implicit, but better to be explicit if model name differs.
	// userRepo.UserModel uses table "users".
	err := r.db.WithContext(ctx).Table("users").
		Where("club_id = ?", clubID).
		Pluck("email", &emails).Error
	return emails, err
}

func (r *PostgresClubRepository) List(ctx context.Context, limit, offset int) ([]domain.Club, error) {
	var clubs []domain.Club
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&clubs).Error; err != nil {
		return nil, err
	}
	return clubs, nil
}

func (r *PostgresClubRepository) Update(ctx context.Context, club *domain.Club) error {
	return r.db.WithContext(ctx).Save(club).Error
}

func (r *PostgresClubRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Club{}, "id = ?", id).Error
}

// --- News Repository ---

func (r *PostgresClubRepository) CreateNews(ctx context.Context, news *domain.News) error {
	return r.db.WithContext(ctx).Create(news).Error
}

func (r *PostgresClubRepository) GetNewsByClub(ctx context.Context, clubID string, limit, offset int) ([]domain.News, error) {
	var news []domain.News
	err := r.db.WithContext(ctx).Where("club_id = ?", clubID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&news).Error
	return news, err
}

func (r *PostgresClubRepository) GetPublicNewsByClub(ctx context.Context, clubID string, limit, offset int) ([]domain.News, error) {
	var news []domain.News
	err := r.db.WithContext(ctx).Where("club_id = ? AND published = ?", clubID, true).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&news).Error
	return news, err
}

func (r *PostgresClubRepository) GetNewsByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.News, error) {
	var news domain.News
	if err := r.db.WithContext(ctx).Where("id = ? AND club_id = ?", id, clubID).First(&news).Error; err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *PostgresClubRepository) UpdateNews(ctx context.Context, clubID string, news *domain.News) error {
	return r.db.WithContext(ctx).Where("id = ? AND club_id = ?", news.ID, clubID).Updates(news).Error
}

func (r *PostgresClubRepository) DeleteNews(ctx context.Context, clubID string, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.News{}, "id = ? AND club_id = ?", id, clubID).Error
}
