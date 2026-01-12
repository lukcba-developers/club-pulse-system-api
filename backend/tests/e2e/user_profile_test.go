package e2e_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	userHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/http"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
)

func TestUserProfileCategory(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Ensure clean state
	_ = db.Migrator().DropTable(&domain.UserStats{}, &domain.Wallet{}, &repository.UserModel{})

	// Migrate schema for test (Split to avoid GORM duplicate constraint issue on CREATE TABLE)
	_ = db.AutoMigrate(&domain.UserStats{}, &domain.Wallet{})
	_ = db.AutoMigrate(&repository.UserModel{})

	// 2. Execute

	repo := repository.NewPostgresUserRepository(db)
	useCase := application.NewUserUseCases(repo, nil) // nil FamilyGroupRepo for this test
	handler := userHttp.NewUserHandler(useCase)

	r := gin.Default()

	// Mock Auth Middleware just sets userID in context
	r.Use(func(c *gin.Context) {
		// Create a test user in DB first
		birthDate := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
		user := &domain.User{
			ID:          "test-user-category-id",
			ClubID:      "test-club-1",
			Name:        "Kid Player",
			Email:       "kid@test.com",
			DateOfBirth: &birthDate,
		}
		// Clean up before creating just in case
		// Clean up before creating just in case (Hard Delete to allow re-insert of same ID)
		db.Unscoped().Where("id = ?", user.ID).Delete(&repository.UserModel{})

		if err := repo.Update(context.Background(), user); err != nil {
			// Handle error or just ignore for test setup of mock
			_ = err
		}
		// Actually relying on GORM to Upsert or create helper

		// Let's manually create via GORM for test setup
		if err := db.Create(&repository.UserModel{
			ID:          user.ID,
			ClubID:      user.ClubID,
			Name:        user.Name,
			Email:       user.Email,
			DateOfBirth: user.DateOfBirth,
		}).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		c.Set("userID", user.ID)
		c.Set("clubID", "test-club-1")
		c.Next()
	})

	r.GET("/users/me", handler.GetProfile)

	// 2. Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/me", nil)
	r.ServeHTTP(w, req)

	// 3. Verify
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Debug
	if _, exists := response["category"]; !exists {
		t.Logf("Response Body: %s", w.Body.String())
	}

	category, ok := response["category"].(string)
	assert.True(t, ok, "Category field should be present")
	assert.Equal(t, "2015", category, "Category should be 2015")
}
