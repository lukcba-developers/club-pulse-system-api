package domain

import (
	"time"

	"github.com/google/uuid"
)

// ConsentType represents the types of consent that can be recorded
type ConsentType string

const (
	ConsentTypeTerms      ConsentType = "TERMS"
	ConsentTypePrivacy    ConsentType = "PRIVACY"
	ConsentTypeMarketing  ConsentType = "MARKETING"
	ConsentTypeHealthData ConsentType = "HEALTH_DATA"
)

// ConsentRecord represents a user's consent for data processing
// Required by GDPR Article 7 - Conditions for consent
type ConsentRecord struct {
	ID           uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID       string      `json:"club_id" gorm:"index;not null"`
	UserID       string      `json:"user_id" gorm:"index;not null"`
	ConsentType  ConsentType `json:"consent_type" gorm:"not null"`
	Version      string      `json:"version" gorm:"not null"` // e.g., "2026-01"
	Accepted     bool        `json:"accepted" gorm:"default:true"`
	AcceptedAt   time.Time   `json:"accepted_at" gorm:"not null"`
	IPAddress    string      `json:"ip_address,omitempty"`
	UserAgent    string      `json:"user_agent,omitempty"`
	ParentUserID *string     `json:"parent_user_id,omitempty"` // For minors
	RevokedAt    *time.Time  `json:"revoked_at,omitempty"`
	CreatedAt    time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (ConsentRecord) TableName() string {
	return "consent_records"
}

// IsActive returns true if the consent is currently active (not revoked)
func (c *ConsentRecord) IsActive() bool {
	return c.Accepted && c.RevokedAt == nil
}

// HealthDataAccessLog tracks access to special category data (GDPR Article 9)
type HealthDataAccessLog struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID            string     `json:"club_id" gorm:"not null"`
	AccessedUserID    string     `json:"accessed_user_id" gorm:"index;not null"`
	AccessingUserID   string     `json:"accessing_user_id" gorm:"index;not null"`
	AccessingUserRole string     `json:"accessing_user_role" gorm:"not null"`
	DocumentID        *uuid.UUID `json:"document_id,omitempty" gorm:"type:uuid"`
	Action            string     `json:"action" gorm:"not null"` // VIEW, DOWNLOAD, VALIDATE, DELETE
	IPAddress         string     `json:"ip_address,omitempty"`
	UserAgent         string     `json:"user_agent,omitempty"`
	AccessedAt        time.Time  `json:"accessed_at" gorm:"not null;default:now()"`
}

// TableName specifies the table name for GORM
func (HealthDataAccessLog) TableName() string {
	return "health_data_access_log"
}

// GDPRErasureRequest tracks right to erasure requests (GDPR Article 17)
type GDPRErasureRequest struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID      string     `json:"club_id" gorm:"not null"`
	UserID      string     `json:"user_id" gorm:"index;not null"`
	RequestedAt time.Time  `json:"requested_at" gorm:"not null;default:now()"`
	ExecutedAt  *time.Time `json:"executed_at,omitempty"`
	ExecutedBy  *string    `json:"executed_by,omitempty"`
	Status      string     `json:"status" gorm:"default:'PENDING'"` // PENDING, COMPLETED, FAILED
	Notes       string     `json:"notes,omitempty" gorm:"type:text"`
}

// TableName specifies the table name for GORM
func (GDPRErasureRequest) TableName() string {
	return "gdpr_erasure_requests"
}

// GDPRExportData represents the data package for right to portability (GDPR Article 20)
type GDPRExportData struct {
	ExportedAt  time.Time              `json:"exported_at"`
	UserProfile map[string]interface{} `json:"user_profile"`
	Documents   []DocumentExport       `json:"documents,omitempty"`
	Consents    []ConsentExport        `json:"consents,omitempty"`
	Bookings    []BookingExport        `json:"bookings,omitempty"`
	Payments    []PaymentExport        `json:"payments,omitempty"`
	FamilyGroup *FamilyGroupExport     `json:"family_group,omitempty"`
}

// Export sub-types for clean JSON structure
type DocumentExport struct {
	Type           string     `json:"type"`
	Status         string     `json:"status"`
	UploadedAt     time.Time  `json:"uploaded_at"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
}

type ConsentExport struct {
	Type       string    `json:"type"`
	Version    string    `json:"version"`
	AcceptedAt time.Time `json:"accepted_at"`
	Active     bool      `json:"active"`
}

type BookingExport struct {
	ID         string    `json:"id"`
	FacilityID string    `json:"facility_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
}

type PaymentExport struct {
	ID        string    `json:"id"`
	Amount    string    `json:"amount"`
	Status    string    `json:"status"`
	Method    string    `json:"method"`
	CreatedAt time.Time `json:"created_at"`
}

type FamilyGroupExport struct {
	Name      string   `json:"name"`
	IsHead    bool     `json:"is_head"`
	MemberIDs []string `json:"member_ids,omitempty"`
}

// ConsentRepository defines the operations for consent management
type ConsentRepository interface {
	Create(record *ConsentRecord) error
	GetByUserID(clubID, userID string) ([]ConsentRecord, error)
	GetActiveByType(clubID, userID string, consentType ConsentType) (*ConsentRecord, error)
	Revoke(clubID, userID string, consentType ConsentType) error
	LogHealthDataAccess(log *HealthDataAccessLog) error
	CreateErasureRequest(req *GDPRErasureRequest) error
	UpdateErasureRequest(req *GDPRErasureRequest) error
}
