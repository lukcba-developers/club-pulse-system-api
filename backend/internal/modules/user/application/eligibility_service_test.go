package application_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserDocumentRepository es un mock del repositorio de documentos
type MockUserDocumentRepository struct {
	mock.Mock
}

func (m *MockUserDocumentRepository) Create(doc *domain.UserDocument) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockUserDocumentRepository) GetByID(clubID string, id uuid.UUID) (*domain.UserDocument, error) {
	args := m.Called(clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserDocument), args.Error(1)
}

func (m *MockUserDocumentRepository) GetByUserID(clubID, userID string) ([]domain.UserDocument, error) {
	args := m.Called(clubID, userID)
	return args.Get(0).([]domain.UserDocument), args.Error(1)
}

func (m *MockUserDocumentRepository) GetByUserAndType(clubID, userID string, docType domain.DocumentType) (*domain.UserDocument, error) {
	args := m.Called(clubID, userID, docType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserDocument), args.Error(1)
}

func (m *MockUserDocumentRepository) Update(doc *domain.UserDocument) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockUserDocumentRepository) Delete(clubID string, id uuid.UUID) error {
	args := m.Called(clubID, id)
	return args.Error(0)
}

func (m *MockUserDocumentRepository) GetExpiringDocuments(clubID string, daysUntilExpiration int) ([]domain.UserDocument, error) {
	args := m.Called(clubID, daysUntilExpiration)
	return args.Get(0).([]domain.UserDocument), args.Error(1)
}

func (m *MockUserDocumentRepository) GetExpiredDocuments(clubID string) ([]domain.UserDocument, error) {
	args := m.Called(clubID)
	return args.Get(0).([]domain.UserDocument), args.Error(1)
}

func (m *MockUserDocumentRepository) GetAllByType(clubID string, docType domain.DocumentType) ([]domain.UserDocument, error) {
	args := m.Called(clubID, docType)
	return args.Get(0).([]domain.UserDocument), args.Error(1)
}

func (m *MockUserDocumentRepository) GetPendingValidation(clubID string) ([]domain.UserDocument, error) {
	args := m.Called(clubID)
	return args.Get(0).([]domain.UserDocument), args.Error(1)
}

// Tests para EligibilityService

func TestCheckEligibility_EligibleUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserDocumentRepository)
	service := application.NewEligibilityService(mockRepo)

	clubID := "club-123"
	userID := "user-456"

	// Usuario con DNI y EMMAC válidos
	futureDate := time.Now().AddDate(0, 6, 0) // 6 meses en el futuro
	docs := []domain.UserDocument{
		{
			ID:     uuid.New(),
			Type:   domain.DocumentTypeDNIFront,
			Status: domain.DocumentStatusValid,
		},
		{
			ID:             uuid.New(),
			Type:           domain.DocumentTypeEMMACMedical,
			Status:         domain.DocumentStatusValid,
			ExpirationDate: &futureDate,
		},
	}

	mockRepo.On("GetByUserID", clubID, userID).Return(docs, nil)

	// Act
	result, err := service.CheckEligibility(clubID, userID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result.IsEligible)
	assert.True(t, result.HasDNI)
	assert.True(t, result.HasEMMAC)
	assert.Empty(t, result.Issues)

	mockRepo.AssertExpectations(t)
}

func TestCheckEligibility_MissingDNI(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserDocumentRepository)
	service := application.NewEligibilityService(mockRepo)

	clubID := "club-123"
	userID := "user-456"

	// Usuario solo con EMMAC válido (sin DNI)
	futureDate := time.Now().AddDate(0, 6, 0)
	docs := []domain.UserDocument{
		{
			ID:             uuid.New(),
			Type:           domain.DocumentTypeEMMACMedical,
			Status:         domain.DocumentStatusValid,
			ExpirationDate: &futureDate,
		},
	}

	mockRepo.On("GetByUserID", clubID, userID).Return(docs, nil)

	// Act
	result, err := service.CheckEligibility(clubID, userID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, result.IsEligible)
	assert.False(t, result.HasDNI)
	assert.True(t, result.HasEMMAC)
	assert.Contains(t, result.Issues, "DNI faltante o no validado")

	mockRepo.AssertExpectations(t)
}

func TestCheckEligibility_ExpiredEMMAC(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserDocumentRepository)
	service := application.NewEligibilityService(mockRepo)

	clubID := "club-123"
	userID := "user-456"

	// Usuario con DNI válido pero EMMAC vencido
	pastDate := time.Now().AddDate(0, -1, 0) // 1 mes atrás
	docs := []domain.UserDocument{
		{
			ID:     uuid.New(),
			Type:   domain.DocumentTypeDNIFront,
			Status: domain.DocumentStatusValid,
		},
		{
			ID:             uuid.New(),
			Type:           domain.DocumentTypeEMMACMedical,
			Status:         domain.DocumentStatusValid,
			ExpirationDate: &pastDate,
		},
	}

	mockRepo.On("GetByUserID", clubID, userID).Return(docs, nil)

	// Act
	result, err := service.CheckEligibility(clubID, userID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, result.IsEligible)
	assert.True(t, result.HasDNI)
	assert.False(t, result.HasEMMAC)
	assert.Contains(t, result.Issues, "Apto físico vencido")

	mockRepo.AssertExpectations(t)
}

func TestCheckEligibility_PendingValidation(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserDocumentRepository)
	service := application.NewEligibilityService(mockRepo)

	clubID := "club-123"
	userID := "user-456"

	// Usuario con DNI válido pero EMMAC pendiente
	futureDate := time.Now().AddDate(0, 6, 0)
	docs := []domain.UserDocument{
		{
			ID:     uuid.New(),
			Type:   domain.DocumentTypeDNIFront,
			Status: domain.DocumentStatusValid,
		},
		{
			ID:             uuid.New(),
			Type:           domain.DocumentTypeEMMACMedical,
			Status:         domain.DocumentStatusPending,
			ExpirationDate: &futureDate,
		},
	}

	mockRepo.On("GetByUserID", clubID, userID).Return(docs, nil)

	// Act
	result, err := service.CheckEligibility(clubID, userID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, result.IsEligible)
	assert.True(t, result.HasDNI)
	assert.False(t, result.HasEMMAC)
	assert.Contains(t, result.Issues, "Apto físico pendiente de validación")

	mockRepo.AssertExpectations(t)
}

func TestCheckEligibility_NoDocuments(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserDocumentRepository)
	service := application.NewEligibilityService(mockRepo)

	clubID := "club-123"
	userID := "user-456"

	// Usuario sin documentos
	docs := []domain.UserDocument{}

	mockRepo.On("GetByUserID", clubID, userID).Return(docs, nil)

	// Act
	result, err := service.CheckEligibility(clubID, userID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, result.IsEligible)
	assert.False(t, result.HasDNI)
	assert.False(t, result.HasEMMAC)
	assert.Len(t, result.Issues, 2) // DNI faltante + Apto físico faltante

	mockRepo.AssertExpectations(t)
}

func TestValidateDocument_Approve(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserDocumentRepository)
	service := application.NewEligibilityService(mockRepo)

	clubID := "club-123"
	docID := uuid.New()
	validatorID := "admin-789"

	doc := &domain.UserDocument{
		ID:     docID,
		ClubID: clubID,
		Status: domain.DocumentStatusPending,
	}

	mockRepo.On("GetByID", clubID, docID).Return(doc, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.UserDocument")).Return(nil)

	// Act
	err := service.ValidateDocument(clubID, docID.String(), validatorID, true, "")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, domain.DocumentStatusValid, doc.Status)
	assert.NotNil(t, doc.ValidatedAt)
	assert.Equal(t, &validatorID, doc.ValidatedBy)

	mockRepo.AssertExpectations(t)
}

func TestValidateDocument_Reject(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserDocumentRepository)
	service := application.NewEligibilityService(mockRepo)

	clubID := "club-123"
	docID := uuid.New()
	validatorID := "admin-789"
	rejectionNotes := "Imagen borrosa, por favor subir nuevamente"

	doc := &domain.UserDocument{
		ID:     docID,
		ClubID: clubID,
		Status: domain.DocumentStatusPending,
	}

	mockRepo.On("GetByID", clubID, docID).Return(doc, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.UserDocument")).Return(nil)

	// Act
	err := service.ValidateDocument(clubID, docID.String(), validatorID, false, rejectionNotes)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, domain.DocumentStatusRejected, doc.Status)
	assert.Equal(t, rejectionNotes, doc.RejectionNotes)

	mockRepo.AssertExpectations(t)
}

func TestValidateDocument_CannotValidateExpired(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserDocumentRepository)
	service := application.NewEligibilityService(mockRepo)

	clubID := "club-123"
	docID := uuid.New()
	validatorID := "admin-789"

	pastDate := time.Now().AddDate(0, -1, 0) // Ya vencido
	doc := &domain.UserDocument{
		ID:             docID,
		ClubID:         clubID,
		Status:         domain.DocumentStatusPending,
		ExpirationDate: &pastDate,
	}

	mockRepo.On("GetByID", clubID, docID).Return(doc, nil)

	// Act
	err := service.ValidateDocument(clubID, docID.String(), validatorID, true, "")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no puede ser validado")

	mockRepo.AssertExpectations(t)
}
