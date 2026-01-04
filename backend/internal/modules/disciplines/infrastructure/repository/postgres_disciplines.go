package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/domain"
	"gorm.io/gorm"
)

type PostgresDisciplineRepository struct {
	db *gorm.DB
}

func NewPostgresDisciplineRepository(db *gorm.DB) *PostgresDisciplineRepository {
	return &PostgresDisciplineRepository{db: db}
}

func (r *PostgresDisciplineRepository) CreateDiscipline(discipline *domain.Discipline) error {
	return r.db.Create(discipline).Error
}

func (r *PostgresDisciplineRepository) ListDisciplines(clubID string) ([]domain.Discipline, error) {
	var disciplines []domain.Discipline
	err := r.db.Where("is_active = ? AND club_id = ?", true, clubID).Find(&disciplines).Error
	return disciplines, err
}

func (r *PostgresDisciplineRepository) GetDisciplineByID(clubID string, id uuid.UUID) (*domain.Discipline, error) {
	var discipline domain.Discipline
	if err := r.db.First(&discipline, "id = ? AND club_id = ?", id, clubID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &discipline, nil
}

func (r *PostgresDisciplineRepository) CreateGroup(group *domain.TrainingGroup) error {
	return r.db.Create(group).Error
}

func (r *PostgresDisciplineRepository) ListGroups(clubID string, filter map[string]interface{}) ([]domain.TrainingGroup, error) {
	var groups []domain.TrainingGroup
	query := r.db.Preload("Discipline").Where("club_id = ?", clubID)
	if dID, ok := filter["discipline_id"]; ok {
		query = query.Where("discipline_id = ?", dID)
	}
	if category, ok := filter["category"]; ok {
		query = query.Where("category = ?", category)
	}
	if coachID, ok := filter["coach_id"]; ok {
		query = query.Where("coach_id = ?", coachID)
	}
	if year, ok := filter["category_year"]; ok {
		query = query.Where("category_year = ?", year)
	}
	err := query.Find(&groups).Error
	return groups, err
}

func (r *PostgresDisciplineRepository) GetGroupByID(clubID string, id uuid.UUID) (*domain.TrainingGroup, error) {
	var group domain.TrainingGroup
	if err := r.db.Preload("Discipline").First(&group, "id = ? AND club_id = ?", id, clubID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}
