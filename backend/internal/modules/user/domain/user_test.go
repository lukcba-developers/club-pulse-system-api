package domain_test

import (
	"testing"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/stretchr/testify/assert"
)

func TestCalculateCategory(t *testing.T) {
	// Setup
	birthDate := time.Date(2012, 5, 15, 0, 0, 0, 0, time.UTC)
	user := domain.User{
		DateOfBirth: &birthDate,
	}

	// Execution
	category := user.CalculateCategory()

	// Validation
	assert.Equal(t, "2012", category, "Category should be the birth year")
}

func TestCalculateCategory_NoDateOfBirth(t *testing.T) {
	// Setup
	user := domain.User{
		DateOfBirth: nil,
	}

	// Execution
	category := user.CalculateCategory()

	// Validation
	assert.Equal(t, "Files", category, "Default category should be 'Files'")
}
