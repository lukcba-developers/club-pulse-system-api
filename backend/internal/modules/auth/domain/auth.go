package domain

import (
	"context"
	"time"
)

type User struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	Password          string     `json:"-"` // Hash
	Name              string     `json:"name"`
	Role              string     `json:"role"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DateOfBirth       *time.Time `json:"date_of_birth,omitempty"`
	ParentID          *string    `json:"parent_id,omitempty"`
	ClubID            string     `json:"club_id"`
	GoogleID          string     `json:"google_id,omitempty"`
	AvatarURL         string     `json:"avatar_url,omitempty"`
	MedicalCertStatus *string    `json:"medical_cert_status,omitempty"`
	MedicalCertExpiry *time.Time `json:"medical_cert_expiry,omitempty"`
	FamilyGroupID     *string    `json:"family_group_id,omitempty"`

	// GDPR Compliance Fields
	TermsAcceptedAt      *time.Time `json:"terms_accepted_at,omitempty"`
	PrivacyPolicyVersion string     `json:"privacy_policy_version,omitempty"`
	DataRetentionUntil   *time.Time `json:"data_retention_until,omitempty"`
}

const (
	RoleSuperAdmin   = "SUPER_ADMIN"
	RoleAdmin        = "ADMIN"
	RoleMember       = "MEMBER"
	RoleCoach        = "COACH"
	RoleMedicalStaff = "MEDICAL_STAFF" // GDPR Article 9 - Special category data access
)

type UserClaims struct {
	UserID string
	Role   string
	ClubID string
}

type RefreshToken struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	DeviceID  string     `json:"device_id"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	IsRevoked bool       `json:"is_revoked"`
	RevokedAt *time.Time `json:"revoked_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type AuthenticationLog struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Type          string    `json:"type"` // LOGIN, LOGOUT, REFRESH
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	Success       bool      `json:"success"`
	FailureReason string    `json:"failure_reason,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthRepository interface {
	SaveUser(ctx context.Context, user *User) error
	FindUserByEmail(ctx context.Context, email, clubID string) (*User, error)
	FindUserByID(ctx context.Context, id, clubID string) (*User, error)

	// Refresh Token Methods
	SaveRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshToken(ctx context.Context, token, clubID string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID, userID string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
	ListUserSessions(ctx context.Context, userID string) ([]RefreshToken, error)

	// Auth Logs
	LogAuthentication(ctx context.Context, log *AuthenticationLog) error
}

type TokenService interface {
	GenerateToken(user *User) (*Token, error)
	ValidateToken(token string) (*UserClaims, error)
	GenerateRefreshToken(user *User) (string, error)
	ValidateRefreshToken(token string) (string, error) // Returns UserID or TokenID?
}
