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
	// Auto-Migrate the new entities
	_ = db.AutoMigrate(&domain.UserStats{}, &domain.Wallet{})
	return &PostgresUserRepository{db: db}
}

// UserModel mirrors the database schema defined in Auth module.
// In a larger system, we might use a shared kernel, or duplicate/map.
// We duplicate here to keep modules decoupled in code, even if coupling in DB.
type UserModel struct {
	ID                string `gorm:"primaryKey"`
	Name              string
	Email             string
	Password          string `gorm:"not null"` // Added to support creation via User module, though mainly Auth managed.
	Role              string
	DateOfBirth       *time.Time             `gorm:"type:date"`
	SportsPreferences map[string]interface{} `gorm:"serializer:json"`
	ParentID          *string                `gorm:"index"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	// Join fields for GORM Preloading (mapping back to domain entities)
	Stats  *domain.UserStats `gorm:"foreignKey:UserID;references:ID"`
	Wallet *domain.Wallet    `gorm:"foreignKey:UserID;references:ID"`
	ClubID string            `gorm:"index;not null"`
}

func (UserModel) TableName() string {
	return "users"
}

func (r *PostgresUserRepository) GetByID(clubID, id string) (*domain.User, error) {
	var model UserModel
	// Preload Stats and Wallet
	result := r.db.Preload("Stats").Preload("Wallet").Where("id = ? AND club_id = ?", id, clubID).First(&model)
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
		ParentID:          model.ParentID,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		Stats:             model.Stats,
		Wallet:            model.Wallet,
		ClubID:            model.ClubID,
	}, nil
}

func (r *PostgresUserRepository) Update(user *domain.User) error {
	// We only update fields allowed by this module (e.g. Name, DoB, Preferences).
	updates := map[string]interface{}{
		"name":       user.Name,
		"updated_at": user.UpdatedAt,
		// ClubID is typically immutable or handled via admin, but if needed:
		// "club_id": user.ClubID,
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

func (r *PostgresUserRepository) Delete(clubID, id string) error {
	return r.db.Delete(&UserModel{}, "id = ? AND club_id = ?", id, clubID).Error
}

func (r *PostgresUserRepository) List(clubID string, limit, offset int, filters map[string]interface{}) ([]domain.User, error) {
	var models []UserModel
	query := r.db.Model(&UserModel{}).Where("club_id = ?", clubID).Limit(limit).Offset(offset)

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
			ID:                model.ID,
			Name:              model.Name,
			Email:             model.Email,
			Role:              model.Role,
			DateOfBirth:       model.DateOfBirth,
			SportsPreferences: model.SportsPreferences,
			ParentID:          model.ParentID,
			CreatedAt:         model.CreatedAt,
			UpdatedAt:         model.UpdatedAt,
			ClubID:            model.ClubID,
		}
	}
	return users, nil
}

func (r *PostgresUserRepository) FindChildren(clubID, parentID string) ([]domain.User, error) {
	var models []UserModel
	result := r.db.Where("parent_id = ? AND club_id = ?", parentID, clubID).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	users := make([]domain.User, len(models))
	for i, m := range models {
		users[i] = domain.User{
			ID:                m.ID,
			Name:              m.Name,
			Email:             m.Email,
			Role:              m.Role,
			DateOfBirth:       m.DateOfBirth,
			ParentID:          m.ParentID,
			SportsPreferences: m.SportsPreferences,
			CreatedAt:         m.CreatedAt,
			UpdatedAt:         m.UpdatedAt,
			ClubID:            m.ClubID,
		}
	}
	return users, nil
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
	model := UserModel{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		Role:              user.Role,
		DateOfBirth:       user.DateOfBirth,
		SportsPreferences: user.SportsPreferences,
		ParentID:          user.ParentID,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
		ClubID:            user.ClubID,
	}

	// Note: We are relying on the DB/GORM to ignore or default fields not present here (like password)
	// But since this is creating a "User", the Auth module model requires Password.
	// For "Children" managed by parents, they might not have login credentials initially,
	// OR we generate a placeholder.
	// However, GORM might fail constraint "not null" on password if defined in migration.
	// Let's assume we handle creation gracefully, potentially setting a default hash if needed by DB.
	// In the actual system, Auth logic handles creation.
	// If we use User module to create, we might be bypassing Auth constraints.
	// BUT, for Phase 7/8, we are pragmatic.
	// Let's assume we set a dummy password hash if empty? Or the DB allows null?
	// Checking Auth migration: Password string `gorm:"not null"`.
	// So we MUST provide a password.
	// We'll set a placeholder in the UseCase, here we just save what is given.
	if model.CreatedAt.IsZero() {
		model.CreatedAt = time.Now()
	}
	if model.UpdatedAt.IsZero() {
		model.UpdatedAt = time.Now()
	}
	model.Password = "$2a$10$PlaceholderHashForChildAcc" // Dummy hash to satisfy constraint if not provided in struct (which isn't)
	// Wait, UserModel in this file does NOT have Password field.
	// If we create here, GORM uses THIS struct model.
	// If Schema has Password column NOT NULL, and we insert without it, Postgres will Error.
	// We must add Password to local UserModel or handle it.
	// Let's add Password to UserModel in this file to support creation.
	// Update: also need to set ClubID
	model.ClubID = user.ClubID

	return r.db.Create(&model).Error
}
