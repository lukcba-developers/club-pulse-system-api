package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"gorm.io/gorm"
)

// UserDocumentRepository implementa el repositorio de documentos de usuario usando PostgreSQL
type UserDocumentRepository struct {
	db *gorm.DB
}

// NewUserDocumentRepository crea una nueva instancia del repositorio
func NewUserDocumentRepository(db *gorm.DB) *UserDocumentRepository {
	return &UserDocumentRepository{db: db}
}

// Create crea un nuevo documento de usuario
func (r *UserDocumentRepository) Create(doc *domain.UserDocument) error {
	return r.db.Create(doc).Error
}

// GetByID obtiene un documento por su ID
func (r *UserDocumentRepository) GetByID(clubID string, id uuid.UUID) (*domain.UserDocument, error) {
	var doc domain.UserDocument
	err := r.db.Where("club_id = ? AND id = ?", clubID, id).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetByUserID obtiene todos los documentos de un usuario
func (r *UserDocumentRepository) GetByUserID(clubID, userID string) ([]domain.UserDocument, error) {
	var docs []domain.UserDocument
	err := r.db.Where("club_id = ? AND user_id = ?", clubID, userID).
		Order("created_at DESC").
		Find(&docs).Error
	return docs, err
}

// GetByUserAndType obtiene un documento específico de un usuario por tipo
func (r *UserDocumentRepository) GetByUserAndType(clubID, userID string, docType domain.DocumentType) (*domain.UserDocument, error) {
	var doc domain.UserDocument
	err := r.db.Where("club_id = ? AND user_id = ? AND type = ?", clubID, userID, docType).
		Order("created_at DESC").
		First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// Update actualiza un documento existente
func (r *UserDocumentRepository) Update(doc *domain.UserDocument) error {
	return r.db.Save(doc).Error
}

// Delete elimina un documento (soft delete si está configurado)
func (r *UserDocumentRepository) Delete(clubID string, id uuid.UUID) error {
	return r.db.Where("club_id = ? AND id = ?", clubID, id).
		Delete(&domain.UserDocument{}).Error
}

// GetExpiringDocuments obtiene documentos que vencen en X días
func (r *UserDocumentRepository) GetExpiringDocuments(clubID string, daysUntilExpiration int) ([]domain.UserDocument, error) {
	var docs []domain.UserDocument

	targetDate := time.Now().AddDate(0, 0, daysUntilExpiration)
	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := r.db.Where("expiration_date >= ? AND expiration_date < ?", startOfDay, endOfDay).
		Where("status = ?", domain.DocumentStatusValid)

	if clubID != "" {
		query = query.Where("club_id = ?", clubID)
	}

	err := query.Find(&docs).Error
	return docs, err
}

// GetExpiredDocuments obtiene documentos que ya han vencido pero aún tienen status VALID
func (r *UserDocumentRepository) GetExpiredDocuments(clubID string) ([]domain.UserDocument, error) {
	var docs []domain.UserDocument

	query := r.db.Where("expiration_date < ?", time.Now()).
		Where("status = ?", domain.DocumentStatusValid)

	if clubID != "" {
		query = query.Where("club_id = ?", clubID)
	}

	err := query.Find(&docs).Error
	return docs, err
}

// GetAllByType obtiene todos los documentos de un tipo específico
func (r *UserDocumentRepository) GetAllByType(clubID string, docType domain.DocumentType) ([]domain.UserDocument, error) {
	var docs []domain.UserDocument
	err := r.db.Where("club_id = ? AND type = ?", clubID, docType).
		Order("created_at DESC").
		Find(&docs).Error
	return docs, err
}

// GetPendingValidation obtiene todos los documentos pendientes de validación
func (r *UserDocumentRepository) GetPendingValidation(clubID string) ([]domain.UserDocument, error) {
	var docs []domain.UserDocument
	err := r.db.Where("club_id = ? AND status = ?", clubID, domain.DocumentStatusPending).
		Order("created_at ASC").
		Find(&docs).Error
	return docs, err
}
