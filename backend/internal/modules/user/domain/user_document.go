package domain

import (
	"time"

	"github.com/google/uuid"
)

// DocumentType representa los tipos de documentos que pueden subir los usuarios
type DocumentType string

const (
	DocumentTypeDNIFront     DocumentType = "DNI_FRONT"
	DocumentTypeDNIBack      DocumentType = "DNI_BACK"
	DocumentTypeEMMACMedical DocumentType = "EMMAC_MEDICAL"
	DocumentTypeLeagueForm   DocumentType = "LEAGUE_FORM"
	DocumentTypeInsurance    DocumentType = "INSURANCE"
)

// DocumentStatus representa el estado de validación de un documento
type DocumentStatus string

const (
	DocumentStatusPending  DocumentStatus = "PENDING"
	DocumentStatusValid    DocumentStatus = "VALID"
	DocumentStatusRejected DocumentStatus = "REJECTED"
	DocumentStatusExpired  DocumentStatus = "EXPIRED"
)

// UserDocument representa un documento subido por un usuario (DNI, apto médico, etc.)
type UserDocument struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID         string         `json:"club_id" gorm:"index;not null"`
	UserID         string         `json:"user_id" gorm:"index;not null"`
	Type           DocumentType   `json:"type" gorm:"not null"`
	FileURL        string         `json:"file_url" gorm:"not null"`
	Status         DocumentStatus `json:"status" gorm:"default:'PENDING'"`
	ExpirationDate *time.Time     `json:"expiration_date,omitempty"`
	RejectionNotes string         `json:"rejection_notes,omitempty"`
	UploadedAt     time.Time      `json:"uploaded_at" gorm:"autoCreateTime"`
	ValidatedAt    *time.Time     `json:"validated_at,omitempty"`
	ValidatedBy    *string        `json:"validated_by,omitempty"` // Admin UserID
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (UserDocument) TableName() string {
	return "user_documents"
}

// IsExpired verifica si el documento ha vencido
func (d *UserDocument) IsExpired() bool {
	if d.ExpirationDate == nil {
		return false
	}
	return time.Now().After(*d.ExpirationDate)
}

// DaysUntilExpiration retorna los días hasta el vencimiento
// Retorna -1 si el documento no tiene fecha de vencimiento
func (d *UserDocument) DaysUntilExpiration() int {
	if d.ExpirationDate == nil {
		return -1
	}
	duration := time.Until(*d.ExpirationDate)
	return int(duration.Hours() / 24)
}

// IsValid verifica si el documento es válido y no ha vencido
func (d *UserDocument) IsValid() bool {
	return d.Status == DocumentStatusValid && !d.IsExpired()
}

// CanBeValidated verifica si el documento puede ser validado
func (d *UserDocument) CanBeValidated() bool {
	// No se puede validar un documento vencido
	if d.IsExpired() {
		return false
	}
	// Solo se pueden validar documentos pendientes
	return d.Status == DocumentStatusPending
}

// UserDocumentRepository define las operaciones de persistencia para documentos de usuario
type UserDocumentRepository interface {
	Create(doc *UserDocument) error
	GetByID(clubID string, id uuid.UUID) (*UserDocument, error)
	GetByUserID(clubID, userID string) ([]UserDocument, error)
	GetByUserAndType(clubID, userID string, docType DocumentType) (*UserDocument, error)
	Update(doc *UserDocument) error
	Delete(clubID string, id uuid.UUID) error

	// Para el Cron Job de vencimientos
	GetExpiringDocuments(clubID string, daysUntilExpiration int) ([]UserDocument, error)
	GetExpiredDocuments(clubID string) ([]UserDocument, error)

	// Operaciones masivas
	GetAllByType(clubID string, docType DocumentType) ([]UserDocument, error)
	GetPendingValidation(clubID string) ([]UserDocument, error)
}
