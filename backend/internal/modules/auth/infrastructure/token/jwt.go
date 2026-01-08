package token

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/domain"
)

type JWTService struct {
	secretKey []byte
	issuer    string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
		issuer:    "club-pulse-api",
	}
}

func (s *JWTService) GenerateToken(user *domain.User) (*domain.Token, error) {
	expiration := time.Now().Add(15 * time.Minute) // Access Token 15m expiry (Security Best Practice)

	claims := jwt.MapClaims{
		"sub":     user.ID,
		"name":    user.Name,
		"role":    user.Role,
		"club_id": user.ClubID,
		"iss":     s.issuer,
		"exp":     expiration.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.Token{
		AccessToken:  signedToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15m in seconds
	}, nil
}

func (s *JWTService) GenerateRefreshToken(user *domain.User) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *JWTService) ValidateRefreshToken(token string) (string, error) {
	// For opaque tokens, we don't validate in token service (it's a DB lookup)
	// Return the token itself as ID / or mock validation
	return token, nil
}

func (s *JWTService) ValidateToken(tokenString string) (*domain.UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, okSub := claims["sub"].(string)
		role, _ := claims["role"].(string) // Optional or default to USER
		clubID, _ := claims["club_id"].(string)

		if okSub {
			return &domain.UserClaims{
				UserID: sub,
				Role:   role,
				ClubID: clubID,
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid token claims")
}
