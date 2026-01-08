package repository

import (
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

func (r *PostgresClubRepository) CreateSponsor(sponsor *domain.Sponsor) error {
	return r.db.Create(sponsor).Error
}

func (r *PostgresClubRepository) CreateAdPlacement(ad *domain.AdPlacement) error {
	return r.db.Create(ad).Error
}

func (r *PostgresClubRepository) GetActiveAds(clubID string) ([]domain.AdPlacement, error) {
	var ads []domain.AdPlacement
	// Join with Sponsor to filter by ClubID if needed, assuming Sponsor has ClubID
	now := time.Now()
	err := r.db.Preload("Sponsor").
		Joins("JOIN sponsors ON sponsors.id = ad_placements.sponsor_id").
		Where("sponsors.club_id = ?", clubID).
		Where("ad_placements.is_active = ?", true).
		Where("ad_placements.contract_end >= ?", now).
		Find(&ads).Error
	return ads, err
}

// --- Club CRUD Implementation (Restoring lost functionality) ---

func (r *PostgresClubRepository) Create(club *domain.Club) error {
	return r.db.Create(club).Error
}

func (r *PostgresClubRepository) GetByID(id string) (*domain.Club, error) {
	var club domain.Club
	if err := r.db.Where("id = ?", id).First(&club).Error; err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *PostgresClubRepository) GetBySlug(slug string) (*domain.Club, error) {
	var club domain.Club
	if err := r.db.Where("slug = ?", slug).First(&club).Error; err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *PostgresClubRepository) GetMemberEmails(clubID string) ([]string, error) {
	var emails []string
	// Assuming 'users' table has 'email' and 'club_id'.
	// Note: 'users' table name is implicit, but better to be explicit if model name differs.
	// userRepo.UserModel uses table "users".
	err := r.db.Table("users").
		Where("club_id = ?", clubID).
		Pluck("email", &emails).Error
	return emails, err
}

func (r *PostgresClubRepository) List(limit, offset int) ([]domain.Club, error) {
	var clubs []domain.Club
	if err := r.db.Limit(limit).Offset(offset).Find(&clubs).Error; err != nil {
		return nil, err
	}
	return clubs, nil
}

func (r *PostgresClubRepository) Update(club *domain.Club) error {
	return r.db.Save(club).Error
}

func (r *PostgresClubRepository) Delete(id string) error {
	return r.db.Delete(&domain.Club{}, "id = ?", id).Error
}

// --- News Repository ---

func (r *PostgresClubRepository) CreateNews(news *domain.News) error {
	return r.db.Create(news).Error
}

func (r *PostgresClubRepository) GetNewsByClub(clubID string, limit, offset int) ([]domain.News, error) {
	var news []domain.News
	err := r.db.Where("club_id = ?", clubID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&news).Error
	return news, err
}

func (r *PostgresClubRepository) GetPublicNewsByClub(clubID string, limit, offset int) ([]domain.News, error) {
	var news []domain.News
	err := r.db.Where("club_id = ? AND published = ?", clubID, true).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&news).Error
	return news, err
}

func (r *PostgresClubRepository) GetNewsByID(id uuid.UUID) (*domain.News, error) {
	var news domain.News
	if err := r.db.Where("id = ?", id).First(&news).Error; err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *PostgresClubRepository) UpdateNews(news *domain.News) error {
	return r.db.Save(news).Error
}

func (r *PostgresClubRepository) DeleteNews(id uuid.UUID) error {
	return r.db.Delete(&domain.News{}, "id = ?", id).Error
}
