package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/championship/domain"
	"gorm.io/gorm"
)

// PostgresVolunteerRepository implementa el repositorio de voluntarios usando PostgreSQL
type PostgresVolunteerRepository struct {
	db *gorm.DB
}

// NewPostgresVolunteerRepository crea una nueva instancia del repositorio
func NewPostgresVolunteerRepository(db *gorm.DB) *PostgresVolunteerRepository {
	return &PostgresVolunteerRepository{db: db}
}

// Create crea una nueva asignación de voluntario
func (r *PostgresVolunteerRepository) Create(ctx context.Context, assignment *domain.VolunteerAssignment) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

// GetByMatchID obtiene todas las asignaciones de un partido
func (r *PostgresVolunteerRepository) GetByMatchID(ctx context.Context, clubID string, matchID uuid.UUID) ([]domain.VolunteerAssignment, error) {
	var assignments []domain.VolunteerAssignment
	err := r.db.WithContext(ctx).Where("club_id = ? AND match_id = ?", clubID, matchID).
		Order("role ASC, assigned_at ASC").
		Find(&assignments).Error
	return assignments, err
}

// GetByUserID obtiene todas las asignaciones de un usuario
func (r *PostgresVolunteerRepository) GetByUserID(ctx context.Context, clubID, userID string) ([]domain.VolunteerAssignment, error) {
	var assignments []domain.VolunteerAssignment
	err := r.db.WithContext(ctx).Where("club_id = ? AND user_id = ?", clubID, userID).
		Order("assigned_at DESC").
		Find(&assignments).Error
	return assignments, err
}

// Delete elimina una asignación de voluntario
func (r *PostgresVolunteerRepository) Delete(ctx context.Context, clubID string, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("club_id = ? AND id = ?", clubID, id).
		Delete(&domain.VolunteerAssignment{}).Error
}

// GetByRoleAndMatch obtiene voluntarios de un rol específico en un partido
func (r *PostgresVolunteerRepository) GetByRoleAndMatch(ctx context.Context, clubID string, matchID uuid.UUID, role domain.VolunteerRole) ([]domain.VolunteerAssignment, error) {
	var assignments []domain.VolunteerAssignment
	err := r.db.WithContext(ctx).Where("club_id = ? AND match_id = ? AND role = ?", clubID, matchID, role).
		Order("assigned_at ASC").
		Find(&assignments).Error
	return assignments, err
}

// Update actualiza una asignación existente
func (r *PostgresVolunteerRepository) Update(ctx context.Context, assignment *domain.VolunteerAssignment) error {
	return r.db.WithContext(ctx).Save(assignment).Error
}
