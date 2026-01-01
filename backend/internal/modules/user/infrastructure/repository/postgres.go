package repository

import (
	"errors"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// UserModel mirrors the database schema defined in Auth module.
// In a larger system, we might use a shared kernel, or duplicate/map.
// We duplicate here to keep modules decoupled in code, even if coupling in DB.
type UserModel struct {
	ID                string `gorm:"primaryKey"`
	Name              string
	Email             string
	Role              string
	DateOfBirth       *time.Time             `gorm:"type:date"`
	SportsPreferences map[string]interface{} `gorm:"serializer:json"` // Requires GORM JSON serializer
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string {
	return "users"
}

func (r *PostgresUserRepository) GetByID(id string) (*domain.User, error) {
	var model UserModel
	result := r.db.Where("id = ?", id).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Or specific domain error
		}
		return nil, result.Error
	}

	return &domain.User{
		ID:                model.ID,
		Name:              model.Name,
		Email:             model.Email,
		Role:              model.Role,
		DateOfBirth:       model.DateOfBirth,
		SportsPreferences: model.SportsPreferences,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
	}, nil
}

func (r *PostgresUserRepository) Update(user *domain.User) error {
	// We only update fields allowed by this module (e.g. Name, DoB, Preferences).
	updates := map[string]interface{}{
		"name":       user.Name,
		"updated_at": user.UpdatedAt,
	}
	if user.DateOfBirth != nil {
		updates["date_of_birth"] = user.DateOfBirth
	}
	if user.SportsPreferences != nil {
		updates["sports_preferences"] = user.SportsPreferences
	}

	result := r.db.Model(&UserModel{ID: user.ID}).Updates(updates)

	return result.Error
}

func (r *PostgresUserRepository) Delete(id string) error {
	return r.db.Delete(&UserModel{}, "id = ?", id).Error
}

func (r *PostgresUserRepository) List(limit, offset int, filters map[string]interface{}) ([]domain.User, error) {
	var models []UserModel
	query := r.db.Model(&UserModel{}).Limit(limit).Offset(offset)

	if search, ok := filters["search"].(string); ok && search != "" {
		// PostgreSQL ILIKE
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if category, ok := filters["category"].(string); ok && category != "" {
		// Filter by Year of DateOfBirth
		// SQLite/Postgres syntax might differ slightly, but standard SQL is EXTRACT(YEAR FROM ...)
		// GORM: datatypes might be sensitive.
		query = query.Where("EXTRACT(YEAR FROM date_of_birth) = ?", category)
	}

	result := query.Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	users := make([]domain.User, len(models))
	for i, model := range models {
		users[i] = domain.User{
			ID:        model.ID,
			Name:      model.Name,
			Email:     model.Email,
			Role:      model.Role,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		}
	}
	return users, nil
}
