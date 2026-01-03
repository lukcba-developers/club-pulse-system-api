package repository

import (
	"errors"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/domain"
	"gorm.io/gorm"
)

type PostgresClubRepository struct {
	db *gorm.DB
}

func NewPostgresClubRepository(db *gorm.DB) *PostgresClubRepository {
	_ = db.AutoMigrate(&domain.Club{})
	return &PostgresClubRepository{db: db}
}

func (r *PostgresClubRepository) Create(club *domain.Club) error {
	if club.CreatedAt.IsZero() {
		club.CreatedAt = time.Now()
	}
	if club.UpdatedAt.IsZero() {
		club.UpdatedAt = time.Now()
	}
	return r.db.Create(club).Error
}

func (r *PostgresClubRepository) GetByID(id string) (*domain.Club, error) {
	var club domain.Club
	if err := r.db.First(&club, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &club, nil
}

func (r *PostgresClubRepository) List(limit, offset int) ([]domain.Club, error) {
	var clubs []domain.Club
	if err := r.db.Limit(limit).Offset(offset).Find(&clubs).Error; err != nil {
		return nil, err
	}
	return clubs, nil
}

func (r *PostgresClubRepository) Update(club *domain.Club) error {
	club.UpdatedAt = time.Now()
	return r.db.Save(club).Error
}

func (r *PostgresClubRepository) Delete(id string) error {
	return r.db.Delete(&domain.Club{}, "id = ?", id).Error
}
