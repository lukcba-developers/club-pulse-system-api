package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuthServiceImpl struct {
	config *oauth2.Config
}

func NewGoogleAuthService() *GoogleAuthServiceImpl {
	return &GoogleAuthServiceImpl{
		config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (s *GoogleAuthServiceImpl) GetUserInfo(ctx context.Context, code string) (*domain.GoogleUserInfo, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := s.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api returned status: %d", resp.StatusCode)
	}

	var userInfo domain.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}
