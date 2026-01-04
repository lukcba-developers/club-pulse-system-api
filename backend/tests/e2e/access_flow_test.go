package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/application"
	accessApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/application"
	accessHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/infrastructure/http"
	accessRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/infrastructure/repository"

	authApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/application"
	authRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/repository"
	authToken "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/token"

	membershipDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	membershipRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/repository"

	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"

	accessDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/domain"

	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessFlow(t *testing.T) {
	// 1. Setup Environment
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Ensure clean state for ALL related tables
	_ = db.Migrator().DropTable(&membershipDomain.MembershipTier{}, &membershipDomain.Membership{}, &userRepo.UserModel{}, &accessDomain.AccessLog{})

	// AutoMigrate test dependencies
	// Note: Auth Repo typically migrates 'users' but might miss 'club_id' if using an older model definition.
	// We force migration with UserRepo's model which has it.
	_ = db.AutoMigrate(&userRepo.UserModel{})
	// MembershipTier and Membership are migrated by NewPostgresMembershipRepository

	// 2. Setup Dependencies
	// Auth
	authR := authRepo.NewPostgresAuthRepository(db)
	tokenService := authToken.NewJWTService("secret")
	authUC := authApp.NewAuthUseCases(authR, tokenService, nil) // Google Auth not needed for this test

	// User
	userR := userRepo.NewPostgresUserRepository(db)
	// userUC := userApp.NewUserUseCases(userR)

	// Membership (this does AutoMigrate internally)
	memR := membershipRepo.NewPostgresMembershipRepository(db)
	// memUC := membershipApp.NewMembershipUseCases(memR)

	// Migrate access_logs AFTER membership repo (which does its own migrations)
	_ = db.AutoMigrate(&accessDomain.AccessLog{})

	// Access
	accessR := accessRepo.NewPostgresAccessRepository(db)
	accessUC := accessApp.NewAccessUseCases(accessR, userR, memR)
	accessH := accessHttp.NewAccessHandler(accessUC)

	// Router
	r := gin.New()
	authMiddleware := func(c *gin.Context) {
		// Mock Auth for simplicity or use real one
		c.Set("userID", "ignored_for_public_access_scan_usually")
		c.Set("clubID", "test-club-1")
		// Access Entry might be authenticated by an Admin/Device or Public with Token.
		// In our plan: POST /access/entry checks the ID in body. Middleware protects the Endpoint itself?
		// Let's assume the endpoint is valid for a "Gate" device token.
		c.Next()
	}
	accessHttp.RegisterRoutes(r.Group("/api/v1"), accessH, authMiddleware, func(c *gin.Context) { c.Next() })

	// 3. Create Test Data
	// Create User
	email := "access_test_" + uuid.New().String() + "@example.com"
	_, err := authUC.Register(authApp.RegisterDTO{
		Name:     "Access User",
		Email:    email,
		Password: "password",
	})
	require.NoError(t, err)

	user, err := authR.FindUserByEmail(email)
	require.NoError(t, err)
	require.NotNil(t, user, "User should not be nil")
	userID := user.ID
	// Manually set ClubID for User
	db.Model(&userRepo.UserModel{}).Where("id = ?", userID).Update("club_id", "test-club-1")

	// 4. Test Case 1: Denied (No Membership)
	reqBody := application.EntryRequest{
		UserID:    userID,
		Direction: "IN",
	}
	w := httptest.NewRecorder()
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/access/entry", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	// In handler we implemented: if DENIED -> 403 Forbidden.
	// Let's verify what we wrote: Match the Handler implementation.
	// if log.Status == "DENIED" { statusCode = http.StatusForbidden }
	// So expect 403.
	// Actually, let's check the Response.

	if w.Code == http.StatusForbidden {
		// Good
		var resp gin.H
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, "DENIED", data["status"])
	} else {
		// Just in case we changed mind
		assert.Equal(t, http.StatusForbidden, w.Code)
	}

	// 5. Test Case 2: Granted (Active Membership)
	// Create Tier
	tier := &membershipDomain.MembershipTier{
		ID:         uuid.New(),
		ClubID:     "test-club-1",
		Name:       "Access Tier",
		MonthlyFee: decimal.NewFromFloat(100),
	}
	db.Create(tier)

	// Create Membership
	membership := &membershipDomain.Membership{
		ID:                 uuid.New(),
		ClubID:             "test-club-1",
		UserID:             uuid.MustParse(userID),
		MembershipTierID:   tier.ID,
		MembershipTier:     *tier,
		Status:             membershipDomain.MembershipStatusActive,
		StartDate:          time.Now(),
		NextBillingDate:    time.Now().AddDate(0, 1, 0),
		OutstandingBalance: decimal.NewFromFloat(0),
	}
	_ = memR.Create(context.Background(), membership)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/access/entry", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp2 gin.H
	_ = json.Unmarshal(w2.Body.Bytes(), &resp2)
	data2 := resp2["data"].(map[string]interface{})
	assert.Equal(t, "GRANTED", data2["status"])

	// 6. Test Case 3: Denied (Debt)
	membership.OutstandingBalance = decimal.NewFromFloat(50)
	db.Save(membership) // Update directly to simulate debt

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/api/v1/access/entry", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusForbidden, w3.Code)
	var resp3 gin.H
	_ = json.Unmarshal(w3.Body.Bytes(), &resp3)
	data3 := resp3["data"].(map[string]interface{})
	assert.Equal(t, "DENIED", data3["status"])
	assert.Contains(t, data3["reason"], "debt")
}
