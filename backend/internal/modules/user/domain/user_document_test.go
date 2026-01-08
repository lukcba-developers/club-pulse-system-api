package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserDocument_IsExpired(t *testing.T) {
	tests := []struct {
		name           string
		expirationDate *time.Time
		expected       bool
	}{
		{
			name:           "Sin fecha de vencimiento",
			expirationDate: nil,
			expected:       false,
		},
		{
			name:           "Fecha futura - No vencido",
			expirationDate: ptr(time.Now().AddDate(0, 1, 0)),
			expected:       false,
		},
		{
			name:           "Fecha pasada - Vencido",
			expirationDate: ptr(time.Now().AddDate(0, -1, 0)),
			expected:       true,
		},
		{
			name:           "Mismo día pero hora futura - No vencido",
			expirationDate: ptr(time.Now().Add(1 * time.Hour)),
			expected:       false,
		},
		{
			name:           "Mismo día pero hora pasada - Vencido",
			expirationDate: ptr(time.Now().Add(-1 * time.Hour)),
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &domain.UserDocument{
				ID:             uuid.New(),
				ExpirationDate: tt.expirationDate,
			}

			result := doc.IsExpired()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserDocument_DaysUntilExpiration(t *testing.T) {
	tests := []struct {
		name           string
		expirationDate *time.Time
		expectedMin    int
		expectedMax    int
	}{
		{
			name:           "Sin fecha de vencimiento",
			expirationDate: nil,
			expectedMin:    -1,
			expectedMax:    -1,
		},
		{
			name:           "30 días en el futuro",
			expirationDate: ptr(time.Now().AddDate(0, 0, 30)),
			expectedMin:    29,
			expectedMax:    30,
		},
		{
			name:           "7 días en el futuro",
			expirationDate: ptr(time.Now().AddDate(0, 0, 7)),
			expectedMin:    6,
			expectedMax:    7,
		},
		{
			name:           "Ya vencido",
			expirationDate: ptr(time.Now().AddDate(0, 0, -5)),
			expectedMin:    -6,
			expectedMax:    -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &domain.UserDocument{
				ID:             uuid.New(),
				ExpirationDate: tt.expirationDate,
			}

			result := doc.DaysUntilExpiration()
			assert.GreaterOrEqual(t, result, tt.expectedMin)
			assert.LessOrEqual(t, result, tt.expectedMax)
		})
	}
}

func TestUserDocument_IsValid(t *testing.T) {
	futureDate := time.Now().AddDate(0, 6, 0)
	pastDate := time.Now().AddDate(0, -1, 0)

	tests := []struct {
		name     string
		doc      *domain.UserDocument
		expected bool
	}{
		{
			name: "Documento válido y no vencido",
			doc: &domain.UserDocument{
				Status:         domain.DocumentStatusValid,
				ExpirationDate: &futureDate,
			},
			expected: true,
		},
		{
			name: "Documento válido sin fecha de vencimiento",
			doc: &domain.UserDocument{
				Status:         domain.DocumentStatusValid,
				ExpirationDate: nil,
			},
			expected: true,
		},
		{
			name: "Documento válido pero vencido",
			doc: &domain.UserDocument{
				Status:         domain.DocumentStatusValid,
				ExpirationDate: &pastDate,
			},
			expected: false,
		},
		{
			name: "Documento pendiente",
			doc: &domain.UserDocument{
				Status:         domain.DocumentStatusPending,
				ExpirationDate: &futureDate,
			},
			expected: false,
		},
		{
			name: "Documento rechazado",
			doc: &domain.UserDocument{
				Status:         domain.DocumentStatusRejected,
				ExpirationDate: &futureDate,
			},
			expected: false,
		},
		{
			name: "Documento expirado",
			doc: &domain.UserDocument{
				Status:         domain.DocumentStatusExpired,
				ExpirationDate: &pastDate,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.doc.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserDocument_CanBeValidated(t *testing.T) {
	futureDate := time.Now().AddDate(0, 6, 0)
	pastDate := time.Now().AddDate(0, -1, 0)

	tests := []struct {
		name     string
		doc      *domain.UserDocument
		expected bool
	}{
		{
			name: "Documento pendiente no vencido - Puede validarse",
			doc: &domain.UserDocument{
				Status:         domain.DocumentStatusPending,
				ExpirationDate: &futureDate,
			},
			expected: true,
		},
		{
			name: "Documento pendiente sin vencimiento - Puede validarse",
			doc: &domain.UserDocument{
				Status:         domain.DocumentStatusPending,
				ExpirationDate: nil,
			},
			expected: true,
		},
		{
			name: "Documento pendiente pero vencido - No puede validarse",
			doc: &domain.UserDocument{
				Status:         domain.DocumentStatusPending,
				ExpirationDate: &pastDate,
			},
			expected: false,
		},
		{
			name: "Documento ya válido - No puede validarse",
			doc: &domain.UserDocument{
				Status:         domain.DocumentStatusValid,
				ExpirationDate: &futureDate,
			},
			expected: false,
		},
		{
			name: "Documento rechazado - No puede validarse",
			doc: &domain.UserDocument{
				Status:         domain.DocumentStatusRejected,
				ExpirationDate: &futureDate,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.doc.CanBeValidated()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDocumentTypes(t *testing.T) {
	// Verificar que los tipos de documento están definidos correctamente
	assert.Equal(t, domain.DocumentType("DNI_FRONT"), domain.DocumentTypeDNIFront)
	assert.Equal(t, domain.DocumentType("DNI_BACK"), domain.DocumentTypeDNIBack)
	assert.Equal(t, domain.DocumentType("EMMAC_MEDICAL"), domain.DocumentTypeEMMACMedical)
	assert.Equal(t, domain.DocumentType("LEAGUE_FORM"), domain.DocumentTypeLeagueForm)
	assert.Equal(t, domain.DocumentType("INSURANCE"), domain.DocumentTypeInsurance)
}

func TestDocumentStatuses(t *testing.T) {
	// Verificar que los estados están definidos correctamente
	assert.Equal(t, domain.DocumentStatus("PENDING"), domain.DocumentStatusPending)
	assert.Equal(t, domain.DocumentStatus("VALID"), domain.DocumentStatusValid)
	assert.Equal(t, domain.DocumentStatus("REJECTED"), domain.DocumentStatusRejected)
	assert.Equal(t, domain.DocumentStatus("EXPIRED"), domain.DocumentStatusExpired)
}

// Helper para crear punteros a time.Time
func ptr(t time.Time) *time.Time {
	return &t
}
