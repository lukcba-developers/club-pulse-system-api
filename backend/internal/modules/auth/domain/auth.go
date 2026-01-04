package domain

import (
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
}

const (
	RoleSuperAdmin = "SUPER_ADMIN"
	RoleAdmin      = "ADMIN"
	RoleMember     = "MEMBER"
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
	SaveUser(user *User) error
	FindUserByEmail(email string) (*User, error)
	FindUserByID(id string) (*User, error)

	// Refresh Token Methods
	SaveRefreshToken(token *RefreshToken) error
	GetRefreshToken(token string) (*RefreshToken, error)
	RevokeRefreshToken(tokenID string) error
	RevokeAllUserTokens(userID string) error
	ListUserSessions(userID string) ([]RefreshToken, error)

	// Auth Logs
	LogAuthentication(log *AuthenticationLog) error
}

type TokenService interface {
	GenerateToken(user *User) (*Token, error)
	ValidateToken(token string) (*UserClaims, error)
	GenerateRefreshToken(user *User) (string, error)
	ValidateRefreshToken(token string) (string, error) // Returns UserID or TokenID?
}
