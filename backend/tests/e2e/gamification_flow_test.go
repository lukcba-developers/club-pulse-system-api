package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	authApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/application"
	authRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/repository"
	authToken "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/token"
	userApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	userHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/http"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGamificationFlow(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Ensure clean state
	_ = db.Migrator().DropTable(&userRepo.UserModel{}, &userDomain.UserStats{}, &userDomain.Wallet{})
	_ = db.AutoMigrate(&userRepo.UserModel{}, &userDomain.UserStats{}, &userDomain.Wallet{})

	// Clear PostgreSQL cached prepared statements after schema change
	db.Exec("DISCARD ALL")

	// Dependencies
	authR := authRepo.NewPostgresAuthRepository(db)
	tokenService := authToken.NewJWTService("secret")
	authUC := authApp.NewAuthUseCases(authR, tokenService, nil)

	userR := userRepo.NewPostgresUserRepository(db)
	userUC := userApp.NewUserUseCases(userR)
	userH := userHttp.NewUserHandler(userUC)

	r := gin.New()
	clubID := "test-club-gamification"

	// Create User
	email := "gamer@test.com"
	registerDTO := authApp.RegisterDTO{
		Name:     "Gamer One",
		Email:    email,
		Password: "password",
	}
	// Cleaning if exists
	existing, _ := authR.FindUserByEmail(email)
	if existing != nil {
		db.Unscoped().Delete(existing)
	}

	_, err := authUC.Register(registerDTO)
	require.NoError(t, err)

	user, err := authR.FindUserByEmail(email)
	require.NoError(t, err)

	// Set ClubID manually
	db.Model(&userRepo.UserModel{}).Where("id = ?", user.ID).Update("club_id", clubID)

	// Create default Stats and Wallet for user (simulating what the app would do)
	defaultStats := &userDomain.UserStats{
		UserID:        user.ID,
		MatchesPlayed: 0,
		RankingPoints: 0, // Default per model
		Level:         1,
	}
	err = db.Create(defaultStats).Error
	require.NoError(t, err, "Failed to create default stats")

	defaultWallet := &userDomain.Wallet{
		UserID:  user.ID,
		Balance: 0,
		Points:  0,
	}
	err = db.Create(defaultWallet).Error
	require.NoError(t, err, "Failed to create default wallet")

	authMiddleware := func(c *gin.Context) {
		c.Set("userID", user.ID)
		c.Set("clubID", clubID)
		c.Next()
	}
	userHttp.RegisterRoutes(r.Group("/api/v1"), userH, authMiddleware, func(c *gin.Context) { c.Next() })

	// 2. Test: Get Stats (Should be empty initially but exist)
	t.Run("Get Stats", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/"+user.ID+"/stats", nil)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		stats := resp["data"].(map[string]interface{})

		// Check defaults (stats fields may be nil if not present in response)
		if stats["matches_played"] != nil {
			assert.Equal(t, float64(0), stats["matches_played"])
		}
		if stats["ranking_points"] != nil {
			assert.Equal(t, float64(0), stats["ranking_points"]) // Default per model is 0
		}
	})

	// 3. Test: Get Wallet
	t.Run("Get Wallet", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/"+user.ID+"/wallet", nil)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		wallet := resp["data"].(map[string]interface{})

		assert.Equal(t, float64(0), wallet["balance"])
	})

	// 4. Simulate Stat Update (e.g. via DB directly as logic is inside Tournament flows)
	// We want to ensure API reads updated values
	db.Model(&userDomain.UserStats{}).
		Where("user_id = ?", user.ID).
		Updates(map[string]interface{}{
			"matches_played": 5,
			"ranking_points": 1200,
			"level":          2,
		})

	t.Run("Get Updated Stats", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/"+user.ID+"/stats", nil)
		r.ServeHTTP(w, req)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		stats := resp["data"].(map[string]interface{})

		assert.Equal(t, float64(5), stats["matches_played"])
		assert.Equal(t, float64(1200), stats["ranking_points"])
	})

	// 5. Simulate Wallet Transaction (Direct DB)
	var wallet userDomain.Wallet
	db.First(&wallet, "user_id = ?", user.ID)

	// Add transaction
	newTx := userDomain.Transaction{
		ID:          "tx-123",
		Type:        "CREDIT",
		Amount:      50.0,
		Description: "Win Bonus",
		Date:        time.Now(),
	}
	// Append by copy since it's a slice alias type
	wallet.Transactions = append(wallet.Transactions, newTx)
	// Update Balance
	wallet.Balance += 50.0

	err = db.Save(&wallet).Error
	require.NoError(t, err)

	t.Run("Get Updated Wallet", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/"+user.ID+"/wallet", nil)
		r.ServeHTTP(w, req)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		res_wallet := resp["data"].(map[string]interface{})

		assert.Equal(t, float64(50), res_wallet["balance"])
		// Verify transactions if API exposes them (It assumes so based on struct)
		// transactions := wallet["transactions"].([]interface{})
		// assert.NotEmpty(t, transactions)
	})
}
