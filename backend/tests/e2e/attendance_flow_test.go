package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	attendanceApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/domain"
	attendanceHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/infrastructure/http"
	attendanceRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/infrastructure/repository"
	membershipRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/repository"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/stretchr/testify/assert"
)

func TestAttendanceFlow(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	database.InitDB()
	db := database.GetDB()

	// Ensure clean state (Legacy schema conflict prevention)
	_ = db.Migrator().DropTable(&attendanceRepo.AttendanceListModel{}, &attendanceRepo.AttendanceRecordModel{}, &userDomain.UserStats{}, &userDomain.Wallet{}, &userRepo.UserModel{})

	// 2. Setup DB & Migrations
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	err := db.AutoMigrate(&userRepo.UserModel{}, &attendanceRepo.AttendanceListModel{}, &attendanceRepo.AttendanceRecordModel{})
	assert.NoError(t, err)

	// Clear PostgreSQL cached prepared statements after schema change
	db.Exec("DISCARD ALL")

	// Repos
	uRepo := userRepo.NewPostgresUserRepository(db)
	aRepo := attendanceRepo.NewPostgresAttendanceRepository(db)
	mRepo := membershipRepo.NewPostgresMembershipRepository(db)
	aUseCase := attendanceApp.NewAttendanceUseCases(aRepo, uRepo, mRepo)

	aHandler := attendanceHttp.NewAttendanceHandler(aUseCase)

	r := gin.Default()
	testClubID := "test-club-attendance"
	// Mock Auth for Coach
	r.Use(func(c *gin.Context) {
		c.Set("userID", "coach-uuid-123")
		c.Set("userRole", "COACH")
		c.Set("clubID", testClubID)
		c.Next()
	})

	attendanceHttp.RegisterRoutes(r.Group("/api/v1"), aHandler, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	// 2. Seed Data
	// Create a Student user with birth year 2012
	deleteDate := time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC)
	student := &userDomain.User{
		ID:          "student-2012",
		Name:        "Student 2012",
		Email:       "student2012@test.com",
		DateOfBirth: &deleteDate,
	}
	// Clean up in correct order (FK dependency: records -> lists)
	db.Unscoped().Where("user_id = ?", student.ID).Delete(&attendanceRepo.AttendanceRecordModel{})
	db.Unscoped().Where("group_name = ?", "2012").Delete(&attendanceRepo.AttendanceListModel{})
	db.Unscoped().Where("id = ?", student.ID).Delete(&userRepo.UserModel{})

	db.Create(&userRepo.UserModel{
		ID:          student.ID,
		ClubID:      testClubID, // Must match the middleware clubID
		Name:        student.Name,
		Email:       student.Email,
		DateOfBirth: student.DateOfBirth,
	})

	// 3. Get Attendance List (Should Auto-Create)
	// GET /api/v1/attendance/groups/2012?date=2024-01-01
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/attendance/groups/2012?date=2024-01-01", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var list domain.AttendanceList
	err = json.Unmarshal(w.Body.Bytes(), &list)
	assert.NoError(t, err)
	assert.Equal(t, "2012", list.Group)
	assert.NotEmpty(t, list.ID)
	// Should have at least one record (our student)
	assert.NotEmpty(t, list.Records)
	found := false
	for _, rec := range list.Records {
		if rec.UserID == "student-2012" {
			found = true
			assert.Equal(t, domain.StatusAbsent, rec.Status) // Default
		}
	}
	assert.True(t, found, "Student 2012 should be in the list")

	// 4. Mark Attendance
	// POST /api/v1/attendance/:listID/records
	recordDto := map[string]interface{}{
		"user_id": "student-2012",
		"status":  "PRESENT",
		"notes":   "Arrived on time",
	}
	body, _ := json.Marshal(recordDto)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/attendance/"+list.ID.String()+"/records", bytes.NewBuffer(body))
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	// 5. Verify Persistence (use the SAME clubID as middleware)
	updatedList, err := aUseCase.GetOrCreateList(testClubID, "2012", list.Date, "coach-uuid-123")
	assert.NoError(t, err)
	assert.NotNil(t, updatedList)

	var updatedRecord domain.AttendanceRecord
	foundUpdated := false
	for _, rec := range updatedList.Records {
		if rec.UserID == "student-2012" {
			updatedRecord = rec
			foundUpdated = true
		}
	}
	assert.True(t, foundUpdated, "Updated record not found")
	assert.Equal(t, domain.StatusPresent, updatedRecord.Status)
	assert.Equal(t, "Arrived on time", updatedRecord.Notes)
}
