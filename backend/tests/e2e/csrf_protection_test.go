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
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCSRFProtection verifies the CSRF middleware behavior
func TestCSRFProtection(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Repos & Services
	authR := authRepo.NewPostgresAuthRepository(db)
	tokenService := authToken.NewJWTService("secret")
	authUC := authApp.NewAuthUseCases(authR, tokenService, nil)
	authH := authHttp.NewAuthHandler(authUC)

	// Create a test user for this test
	testEmail := "csrftest@test.com"
	clubID := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

	// Clean up and create test user
	db.Exec("DELETE FROM users WHERE email = ?", testEmail)
	_, err := authUC.Register(context.Background(), authApp.RegisterDTO{
		Name:                 "CSRF Test",
		Email:                testEmail,
		Password:             "password123",
		AcceptTerms:          true,
		PrivacyPolicyVersion: "2026-01",
	}, clubID)
	if err != nil {
		t.Logf("User may already exist, continuing: %v", err)
	}

	// Router with CSRF middleware
	r := gin.New()

	// Club ID middleware
	r.Use(func(c *gin.Context) {
		clubID := c.GetHeader("X-Club-ID")
		if clubID != "" {
			c.Set("clubID", clubID)
		}
		c.Next()
	})

	// CSRF Middleware
	r.Use(middleware.CSRFMiddleware())

	authHttp.RegisterRoutes(r.Group("/api/v1"), authH, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	// 2. Test: Login should work without CSRF (it's in exclusion list)
	t.Run("Login Without CSRF Token Succeeds", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    testEmail,
			"password": "password123",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Club-ID", clubID)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		// Verify CSRF cookie was set
		cookies := w.Result().Cookies()
		var csrfCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "csrf_token" {
				csrfCookie = c
				break
			}
		}
		require.NotNil(t, csrfCookie, "CSRF cookie should be set after login")
		assert.NotEmpty(t, csrfCookie.Value)
	})

	// 3. Test: Protected endpoint without CSRF token should fail
	t.Run("Protected Endpoint Without CSRF Token Fails", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"refresh_token": "fake-token",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/logout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Club-ID", clubID)
		// NOT setting any CSRF cookie or header
		r.ServeHTTP(w, req)

		// Should fail with 403 Forbidden due to missing CSRF
		require.Equal(t, http.StatusForbidden, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Contains(t, resp["error"], "CSRF")
	})

	// 4. Test: Protected endpoint with valid CSRF token should work
	t.Run("Protected Endpoint With Valid CSRF Token Succeeds", func(t *testing.T) {
		// First, login to get the CSRF token
		loginBody, _ := json.Marshal(map[string]string{
			"email":    testEmail,
			"password": "password123",
		})
		loginW := httptest.NewRecorder()
		loginReq, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
		loginReq.Header.Set("Content-Type", "application/json")
		loginReq.Header.Set("X-Club-ID", clubID)
		r.ServeHTTP(loginW, loginReq)

		require.Equal(t, http.StatusOK, loginW.Code)

		// Get CSRF and refresh tokens from cookies
		cookies := loginW.Result().Cookies()
		var csrfToken, refreshToken string
		for _, c := range cookies {
			if c.Name == "csrf_token" {
				csrfToken = c.Value
			}
			if c.Name == "refresh_token" {
				refreshToken = c.Value
			}
		}
		require.NotEmpty(t, csrfToken)
		require.NotEmpty(t, refreshToken)

		// Now try logout WITH CSRF token
		logoutBody, _ := json.Marshal(map[string]string{
			"refresh_token": refreshToken,
		})
		logoutW := httptest.NewRecorder()
		logoutReq, _ := http.NewRequest("POST", "/api/v1/auth/logout", bytes.NewBuffer(logoutBody))
		logoutReq.Header.Set("Content-Type", "application/json")
		logoutReq.Header.Set("X-Club-ID", clubID)
		logoutReq.Header.Set("X-CSRF-Token", csrfToken)                         // Set CSRF header
		logoutReq.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken}) // Set CSRF cookie

		r.ServeHTTP(logoutW, logoutReq)

		// Should succeed with 204 No Content
		require.Equal(t, http.StatusNoContent, logoutW.Code)
	})

	// 5. Test: Mismatched CSRF tokens should fail
	t.Run("Mismatched CSRF Tokens Fails", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"refresh_token": "fake-token",
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/logout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Club-ID", clubID)
		req.Header.Set("X-CSRF-Token", "wrong-token")
		req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "correct-token"})

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Contains(t, resp["error"], "mismatch")
	})
}
