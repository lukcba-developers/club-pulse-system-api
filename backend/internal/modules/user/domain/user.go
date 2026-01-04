package domain

import (
	"time"

	"github.com/google/uuid"
)

type FamilyGroup struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name       string    `json:"name" gorm:"not null"`
	HeadUserID string    `json:"head_user_id" gorm:"not null"`
	Members    []User    `json:"members,omitempty" gorm:"foreignKey:FamilyGroupID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type User struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Email             string                 `json:"email"`
	Role              string                 `json:"role"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DateOfBirth       *time.Time             `json:"date_of_birth,omitempty"`
	SportsPreferences map[string]interface{} `json:"sports_preferences,omitempty"`
	ParentID          *string                `json:"parent_id,omitempty"`
	ClubID            string                 `json:"club_id" gorm:"index;not null"`
	FamilyGroupID     *uuid.UUID             `json:"family_group_id,omitempty" gorm:"type:uuid"`

	// Health
	MedicalCertStatus *MedicalCertStatus `json:"medical_cert_status" gorm:"default:'PENDING'"`
	MedicalCertExpiry *time.Time         `json:"medical_cert_expiry"`

	// Relations (Fetched on demand or preloaded)
	Stats  *UserStats `json:"stats,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Wallet *Wallet    `json:"wallet,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

const (
	RoleSuperAdmin = "SUPER_ADMIN"
	RoleAdmin      = "ADMIN"
	RoleMember     = "MEMBER"
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
	FindChildren(clubID, parentID string) ([]User, error)
	Create(user *User) error
}
