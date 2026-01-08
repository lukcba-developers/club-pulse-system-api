package domain

import (
	"time"

	"github.com/google/uuid"
)

type FamilyGroup struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClubID     string    `json:"club_id" gorm:"index;not null"`
	Name       string    `json:"name" gorm:"not null"`
	HeadUserID string    `json:"head_user_id" gorm:"not null"`
	Members    []User    `json:"members,omitempty" gorm:"foreignKey:FamilyGroupID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FamilyGroupRepository interface {
	Create(group *FamilyGroup) error
	GetByID(clubID string, id uuid.UUID) (*FamilyGroup, error)
	GetByHeadUserID(clubID, headUserID string) (*FamilyGroup, error)
	GetByMemberID(clubID, userID string) (*FamilyGroup, error)
	AddMember(clubID string, groupID uuid.UUID, userID string) error
	RemoveMember(clubID string, groupID uuid.UUID, userID string) error
}

type User struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Email             string                 `json:"email"`
	Role              string                 `json:"role"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DateOfBirth       *time.Time             `json:"date_of_birth,omitempty"`
	SportsPreferences map[string]interface{} `json:"sports_preferences,omitempty" gorm:"serializer:json"`
	ParentID          *string                `json:"parent_id,omitempty"`
	ClubID            string                 `json:"club_id" gorm:"index;not null"`
	FamilyGroupID     *uuid.UUID             `json:"family_group_id,omitempty" gorm:"type:uuid"`

	// Health
	MedicalCertStatus *MedicalCertStatus `json:"medical_cert_status" gorm:"default:'PENDING'"`
	MedicalCertExpiry *time.Time         `json:"medical_cert_expiry"`

	// Emergency & Security
	EmergencyContactName  string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone string `json:"emergency_contact_phone,omitempty"`
	InsuranceProvider     string `json:"insurance_provider,omitempty"`
	InsuranceNumber       string `json:"insurance_number,omitempty"`

	// GDPR Compliance Fields
	TermsAcceptedAt      *time.Time `json:"terms_accepted_at,omitempty"`
	PrivacyPolicyVersion string     `json:"privacy_policy_version,omitempty"`
	DataRetentionUntil   *time.Time `json:"data_retention_until,omitempty"`

	// Relations (Fetched on demand or preloaded)
	Stats  *UserStats `json:"stats,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Wallet *Wallet    `json:"wallet,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

const (
	RoleSuperAdmin   = "SUPER_ADMIN"
	RoleAdmin        = "ADMIN"
	RoleMember       = "MEMBER"
	RoleCoach        = "COACH"
	RoleMedicalStaff = "MEDICAL_STAFF" // GDPR Article 9 - Special category data access
)

type MedicalCertStatus string

const (
	MedicalCertStatusValid   MedicalCertStatus = "VALID"
	MedicalCertStatusExpired MedicalCertStatus = "EXPIRED"
	MedicalCertStatusPending MedicalCertStatus = "PENDING"
)

// CalculateCategory returns the user's category based on birth year (e.g., "2012")
func (u *User) CalculateCategory() string {
	if u.DateOfBirth == nil {
		return "Files" // Default category if unknown
	}
	return u.DateOfBirth.Format("2006")
}

type UserRepository interface {
	GetByID(clubID, id string) (*User, error)
	// Update updates the non-auth fields of the user
	Update(user *User) error
	Delete(clubID, id string) error
	List(clubID string, limit, offset int, filters map[string]interface{}) ([]User, error)
	ListByIDs(clubID string, ids []string) ([]User, error)
	FindChildren(clubID, parentID string) ([]User, error)
	Create(user *User) error
	CreateIncident(incident *IncidentLog) error
	GetByEmail(email string) (*User, error)
	// GDPR Article 17 - Right to erasure
	AnonymizeForGDPR(clubID, id string) error
}

type IncidentLog struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ClubID        string     `json:"club_id" gorm:"not null;index"`
	InjuredUserID *string    `json:"injured_user_id,omitempty"` // Nullable for visitors
	Description   string     `json:"description" gorm:"not null"`
	Witnesses     string     `json:"witnesses,omitempty"`
	ActionTaken   string     `json:"action_taken,omitempty"`
	ReportedAt    time.Time  `json:"reported_at"`
	CreatedBy     string     `json:"created_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}
