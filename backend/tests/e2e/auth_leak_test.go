package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	authApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/application"
	authHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/http"
	authRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/repository"
	authToken "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/token"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthDataLeak(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Clear DB
	db.Exec("DELETE FROM users")

	// Repos & Services
	authR := authRepo.NewPostgresAuthRepository(db)
	tokenService := authToken.NewJWTService("secret")
	authUC := authApp.NewAuthUseCases(authR, tokenService, nil)
	authH := authHttp.NewAuthHandler(authUC)

	// Router
	r := gin.New()

	// Middleware to simulate Club Context from Header
	r.Use(func(c *gin.Context) {
		clubID := c.GetHeader("X-Club-ID")
		if clubID != "" {
			c.Set("clubID", clubID)
		}
		c.Next()
	})

	authHttp.RegisterRoutes(r.Group("/api/v1"), authH, func(c *gin.Context) { c.Next() })

	// 2. Create User in CLUB A
	email := "leaker@test.com"
	clubA := "club-A"
	clubB := "club-B"

	// Register in Club A (Direct UseCase call)
	_, err := authUC.Register(context.Background(), authApp.RegisterDTO{
		Name:     "Leaker User",
		Email:    email,
		Password: "password",
	}, clubA)
	require.NoError(t, err)

	user, _ := authR.FindUserByEmail(context.Background(), email, clubA)
	require.NotNil(t, user)
	assert.Equal(t, clubA, user.ClubID)

	// 3. Attempt Login in CLUB A (Should Success)
	t.Run("Login in Correct Club", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": "password",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("X-Club-ID", clubA)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	// 4. Attempt Login in CLUB B (Should Fail - This proves the fix)
	t.Run("Login in Wrong Club", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": "password",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("X-Club-ID", clubB) // Different Club!
		r.ServeHTTP(w, req)

		// Expect 401 Unauthorized because user shouldn't be found in this club context
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
