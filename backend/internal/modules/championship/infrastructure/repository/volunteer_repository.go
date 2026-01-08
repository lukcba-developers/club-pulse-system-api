package repository

import (
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
func (r *PostgresVolunteerRepository) Create(assignment *domain.VolunteerAssignment) error {
	return r.db.Create(assignment).Error
}

// GetByMatchID obtiene todas las asignaciones de un partido
func (r *PostgresVolunteerRepository) GetByMatchID(clubID string, matchID uuid.UUID) ([]domain.VolunteerAssignment, error) {
	var assignments []domain.VolunteerAssignment
	err := r.db.Where("club_id = ? AND match_id = ?", clubID, matchID).
		Order("role ASC, assigned_at ASC").
		Find(&assignments).Error
	return assignments, err
}

// GetByUserID obtiene todas las asignaciones de un usuario
func (r *PostgresVolunteerRepository) GetByUserID(clubID, userID string) ([]domain.VolunteerAssignment, error) {
	var assignments []domain.VolunteerAssignment
	err := r.db.Where("club_id = ? AND user_id = ?", clubID, userID).
		Order("assigned_at DESC").
		Find(&assignments).Error
	return assignments, err
}

// Delete elimina una asignación de voluntario
func (r *PostgresVolunteerRepository) Delete(clubID string, id uuid.UUID) error {
	return r.db.Where("club_id = ? AND id = ?", clubID, id).
		Delete(&domain.VolunteerAssignment{}).Error
}

// GetByRoleAndMatch obtiene voluntarios de un rol específico en un partido
func (r *PostgresVolunteerRepository) GetByRoleAndMatch(clubID string, matchID uuid.UUID, role domain.VolunteerRole) ([]domain.VolunteerAssignment, error) {
	var assignments []domain.VolunteerAssignment
	err := r.db.Where("club_id = ? AND match_id = ? AND role = ?", clubID, matchID, role).
		Order("assigned_at ASC").
		Find(&assignments).Error
	return assignments, err
}

// Update actualiza una asignación existente
func (r *PostgresVolunteerRepository) Update(assignment *domain.VolunteerAssignment) error {
	return r.db.Save(assignment).Error
}
