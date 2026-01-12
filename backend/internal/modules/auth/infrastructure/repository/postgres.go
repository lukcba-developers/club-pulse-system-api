package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/domain"
	"gorm.io/gorm"
)

type PostgresAuthRepository struct {
	db *gorm.DB
}

func NewPostgresAuthRepository(db *gorm.DB) *PostgresAuthRepository {
	// AutoMigrate should ideally be done in a separate migration step,
	// but for MVP/Development it's acceptable here or in main.
	// We'll trust the main setup or do safe migration here.
	_ = db.AutoMigrate(&UserModel{}, &RefreshTokenModel{}, &AuthenticationLogModel{})
	return &PostgresAuthRepository{db: db}
}

// UserModel is the Infrastructure representation of the User
type UserModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Email       string `gorm:"uniqueIndex;not null"`
	Password    string `gorm:"column:password;not null"`
	Role        string `gorm:"default:'USER'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DateOfBirth *time.Time
	ParentID    *string `gorm:"index"`
	ClubID      string  `gorm:"index"`
	GoogleID    string  `gorm:"index"`
	AvatarURL   string
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// GDPR Compliance Fields
	TermsAcceptedAt      *time.Time
	PrivacyPolicyVersion string
	DataRetentionUntil   *time.Time
}

// RefreshTokenModel represents the refresh token in DB
type RefreshTokenModel struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"not null;index"`
	DeviceID  string    `gorm:"index"`
	Token     string    `gorm:"not null;unique;index"`
	ExpiresAt time.Time `gorm:"not null"`
	IsRevoked bool      `gorm:"default:false"`
	RevokedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (RefreshTokenModel) TableName() string {
	return "refresh_tokens"

}

// AuthenticationLogModel represents auth logs in DB
type AuthenticationLogModel struct {
	ID            string `gorm:"primaryKey"`
	UserID        string `gorm:"index"`
	Type          string `gorm:"index"` // LOGIN, LOGOUT
	IPAddress     string
	UserAgent     string
	Success       bool
	FailureReason string
	CreatedAt     time.Time
}

func (AuthenticationLogModel) TableName() string {
	return "auth_logs"
}

// TableName overrides the table name
func (UserModel) TableName() string {
	return "users"
}

func (r *PostgresAuthRepository) SaveUser(ctx context.Context, user *domain.User) error {
	userModel := UserModel{
		ID:                   user.ID,
		Name:                 user.Name,
		Email:                user.Email,
		Password:             user.Password,
		Role:                 user.Role,
		CreatedAt:            user.CreatedAt,
		UpdatedAt:            user.UpdatedAt,
		DateOfBirth:          user.DateOfBirth,
		ParentID:             user.ParentID,
		ClubID:               user.ClubID,
		GoogleID:             user.GoogleID,
		AvatarURL:            user.AvatarURL,
		TermsAcceptedAt:      user.TermsAcceptedAt,
		PrivacyPolicyVersion: user.PrivacyPolicyVersion,
		DataRetentionUntil:   user.DataRetentionUntil,
	}

	result := r.db.WithContext(ctx).Create(&userModel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostgresAuthRepository) FindUserByEmail(ctx context.Context, email, clubID string) (*domain.User, error) {
	var userModel UserModel
	// Ensure we only find users within the requesting club context
	result := r.db.WithContext(ctx).Where("email = ? AND club_id = ?", email, clubID).First(&userModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &domain.User{
		ID:                   userModel.ID,
		Name:                 userModel.Name,
		Email:                userModel.Email,
		Password:             userModel.Password,
		Role:                 userModel.Role,
		CreatedAt:            userModel.CreatedAt,
		UpdatedAt:            userModel.UpdatedAt,
		DateOfBirth:          userModel.DateOfBirth,
		ParentID:             userModel.ParentID,
		ClubID:               userModel.ClubID,
		GoogleID:             userModel.GoogleID,
		AvatarURL:            userModel.AvatarURL,
		TermsAcceptedAt:      userModel.TermsAcceptedAt,
		PrivacyPolicyVersion: userModel.PrivacyPolicyVersion,
		DataRetentionUntil:   userModel.DataRetentionUntil,
	}, nil
}

func (r *PostgresAuthRepository) FindUserByID(ctx context.Context, id, clubID string) (*domain.User, error) {
	var userModel UserModel
	// Ensure we only find users within the requesting club context
	result := r.db.WithContext(ctx).Where("id = ? AND club_id = ?", id, clubID).First(&userModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &domain.User{
		ID:                   userModel.ID,
		Name:                 userModel.Name,
		Email:                userModel.Email,
		Password:             userModel.Password,
		Role:                 userModel.Role,
		CreatedAt:            userModel.CreatedAt,
		UpdatedAt:            userModel.UpdatedAt,
		DateOfBirth:          userModel.DateOfBirth,
		ParentID:             userModel.ParentID,
		ClubID:               userModel.ClubID,
		TermsAcceptedAt:      userModel.TermsAcceptedAt,
		PrivacyPolicyVersion: userModel.PrivacyPolicyVersion,
		DataRetentionUntil:   userModel.DataRetentionUntil,
	}, nil
}

func (r *PostgresAuthRepository) SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	tokenModel := RefreshTokenModel{
		ID:        token.ID,
		UserID:    token.UserID,
		DeviceID:  token.DeviceID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		IsRevoked: token.IsRevoked,
		RevokedAt: token.RevokedAt,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}

	result := r.db.WithContext(ctx).Create(&tokenModel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostgresAuthRepository) GetRefreshToken(ctx context.Context, token, clubID string) (*domain.RefreshToken, error) {
	var tokenModel RefreshTokenModel
	// SECURITY: Validate token belongs to a user within the requesting club
	result := r.db.WithContext(ctx).
		Joins("JOIN users ON users.id = refresh_tokens.user_id").
		Where("refresh_tokens.token = ? AND users.club_id = ?", token, clubID).
		First(&tokenModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &domain.RefreshToken{
		ID:        tokenModel.ID,
		UserID:    tokenModel.UserID,
		DeviceID:  tokenModel.DeviceID,
		Token:     tokenModel.Token,
		ExpiresAt: tokenModel.ExpiresAt,
		IsRevoked: tokenModel.IsRevoked,
		RevokedAt: tokenModel.RevokedAt,
		CreatedAt: tokenModel.CreatedAt,
		UpdatedAt: tokenModel.UpdatedAt,
	}, nil
}

func (r *PostgresAuthRepository) RevokeRefreshToken(ctx context.Context, tokenID, userID string) error {
	now := time.Now()
	// Using map to update multiple fields and avoid zero-value issues with boolean
	updates := map[string]interface{}{
		"is_revoked": true,
		"revoked_at": now,
		"updated_at": now,
	}
	// Security: Ensure the token belongs to the user
	result := r.db.WithContext(ctx).Model(&RefreshTokenModel{}).
		Where("id = ? AND user_id = ?", tokenID, userID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("token not found or access denied")
	}
	return nil
}

func (r *PostgresAuthRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"is_revoked": true,
		"revoked_at": now,
		"updated_at": now,
	}
	// Revoke all non-revoked tokens for the user
	return r.db.WithContext(ctx).Model(&RefreshTokenModel{}).Where("user_id = ? AND is_revoked = ?", userID, false).Updates(updates).Error
}

func (r *PostgresAuthRepository) ListUserSessions(ctx context.Context, userID string) ([]domain.RefreshToken, error) {
	var models []RefreshTokenModel
	// Only return active (non-revoked) sessions for MVP view
	// In production, might want 'all' sessions with status.
	// Let's return all but order by created_at desc
	result := r.db.WithContext(ctx).Where("user_id = ? AND is_revoked = ?", userID, false).Order("created_at desc").Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	tokens := make([]domain.RefreshToken, len(models))
	for i, m := range models {
		tokens[i] = domain.RefreshToken{
			ID:        m.ID,
			UserID:    m.UserID,
			DeviceID:  m.DeviceID,
			Token:     m.Token,
			ExpiresAt: m.ExpiresAt,
			IsRevoked: m.IsRevoked,
			RevokedAt: m.RevokedAt,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	return tokens, nil
}

func (r *PostgresAuthRepository) LogAuthentication(ctx context.Context, log *domain.AuthenticationLog) error {
	model := AuthenticationLogModel{
		ID:            log.ID,
		UserID:        log.UserID,
		Type:          log.Type,
		IPAddress:     log.IPAddress,
		UserAgent:     log.UserAgent,
		Success:       log.Success,
		FailureReason: log.FailureReason,
		CreatedAt:     log.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}
