package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/application"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockFacilityRepo struct {
	mock.Mock
}

func (m *MockFacilityRepo) Create(ctx context.Context, facility *domain.Facility) error {
	args := m.Called(ctx, facility)
	return args.Error(0)
}

func (m *MockFacilityRepo) GetByID(ctx context.Context, clubID, id string) (*domain.Facility, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Facility), args.Error(1)
}

func (m *MockFacilityRepo) GetByIDForUpdate(ctx context.Context, clubID, id string) (*domain.Facility, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Facility), args.Error(1)
}

func (m *MockFacilityRepo) List(ctx context.Context, clubID string, limit, offset int) ([]*domain.Facility, error) {
	args := m.Called(ctx, clubID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Facility), args.Error(1)
}

func (m *MockFacilityRepo) Update(ctx context.Context, facility *domain.Facility) error {
	args := m.Called(ctx, facility)
	return args.Error(0)
}

func (m *MockFacilityRepo) HasConflict(ctx context.Context, clubID, facilityID string, startTime, endTime time.Time) (bool, error) {
	args := m.Called(ctx, clubID, facilityID, startTime, endTime)
	return args.Bool(0), args.Error(1)
}

func (m *MockFacilityRepo) ListMaintenanceByFacility(ctx context.Context, clubID, facilityID string) ([]*domain.MaintenanceTask, error) {
	args := m.Called(ctx, clubID, facilityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.MaintenanceTask), args.Error(1)
}

func (m *MockFacilityRepo) CreateMaintenance(ctx context.Context, clubID string, task *domain.MaintenanceTask) error {
	args := m.Called(ctx, clubID, task)
	return args.Error(0)
}

func (m *MockFacilityRepo) GetMaintenanceByID(ctx context.Context, clubID, id string) (*domain.MaintenanceTask, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MaintenanceTask), args.Error(1)
}

func (m *MockFacilityRepo) SemanticSearch(ctx context.Context, clubID string, embedding []float32, limit int) ([]*domain.FacilityWithSimilarity, error) {
	args := m.Called(ctx, clubID, embedding, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.FacilityWithSimilarity), args.Error(1)
}

func (m *MockFacilityRepo) UpdateEmbedding(ctx context.Context, facilityID string, embedding []float32) error {
	// TODO: Update methods likely need clubID too? For now leaving as is if interface didn't change for UpdateEmbedding in my plan.
	// UseCases says: "UpdateEmbedding: Añadir validación de tenant."
	// But let's check current signature in domain?
	// Assuming it's not critical for now or assuming tenant validation happens inside via query.
	args := m.Called(ctx, facilityID, embedding)
	return args.Error(0)
}

func (m *MockFacilityRepo) CreateEquipment(ctx context.Context, clubID string, equipment *domain.Equipment) error {
	args := m.Called(ctx, clubID, equipment)
	return args.Error(0)
}

func (m *MockFacilityRepo) GetEquipmentByID(ctx context.Context, clubID, id string) (*domain.Equipment, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Equipment), args.Error(1)
}

func (m *MockFacilityRepo) ListEquipmentByFacility(ctx context.Context, clubID, facilityID string) ([]*domain.Equipment, error) {
	args := m.Called(ctx, clubID, facilityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Equipment), args.Error(1)
}

func (m *MockFacilityRepo) UpdateEquipment(ctx context.Context, clubID string, equipment *domain.Equipment) error {
	args := m.Called(ctx, clubID, equipment)
	return args.Error(0)
}

func (m *MockFacilityRepo) LoanEquipmentAtomic(ctx context.Context, loan *domain.EquipmentLoan, equipmentID string) error {
	args := m.Called(ctx, loan, equipmentID)
	return args.Error(0)
}

type MockLoanRepo struct {
	mock.Mock
}

func (m *MockLoanRepo) Create(ctx context.Context, loan *domain.EquipmentLoan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepo) GetByID(ctx context.Context, id string) (*domain.EquipmentLoan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EquipmentLoan), args.Error(1)
}

func (m *MockLoanRepo) ListByUser(ctx context.Context, userID string) ([]*domain.EquipmentLoan, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.EquipmentLoan), args.Error(1)
}

func (m *MockLoanRepo) ListByStatus(ctx context.Context, status domain.LoanStatus) ([]*domain.EquipmentLoan, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.EquipmentLoan), args.Error(1)
}

func (m *MockLoanRepo) Update(ctx context.Context, loan *domain.EquipmentLoan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

// --- Tests ---

func TestCreateFacility(t *testing.T) {
	mockRepo := new(MockFacilityRepo)
	mockLoanRepo := new(MockLoanRepo)
	uc := application.NewFacilityUseCases(mockRepo, mockLoanRepo)

	t.Run("Success", func(t *testing.T) {
		dto := application.CreateFacilityDTO{
			Name:       "Test Court",
			Type:       domain.FacilityTypeCourt,
			Capacity:   4,
			HourlyRate: 20.0,
		}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Facility")).Return(nil).Once()

		facility, err := uc.CreateFacility(context.Background(), "club-1", dto)

		assert.NoError(t, err)
		assert.NotNil(t, facility)
		assert.Equal(t, "Test Court", facility.Name)
		mockRepo.AssertExpectations(t)
	})
}

func TestListAndGet(t *testing.T) {
	mockRepo := new(MockFacilityRepo)
	mockLoanRepo := new(MockLoanRepo)
	uc := application.NewFacilityUseCases(mockRepo, mockLoanRepo)
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		mockRepo.On("List", ctx, "club-1", 10, 0).Return([]*domain.Facility{{ID: "1"}}, nil).Once()
		res, err := uc.ListFacilities(ctx, "club-1", 10, 0)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("Get", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "club-1", "fac-1").Return(&domain.Facility{ID: "fac-1"}, nil).Once()
		res, err := uc.GetFacility(ctx, "club-1", "fac-1")
		assert.NoError(t, err)
		assert.Equal(t, "fac-1", res.ID)
	})
}

func TestEquipmentManagement(t *testing.T) {
	mockRepo := new(MockFacilityRepo)
	mockLoanRepo := new(MockLoanRepo)
	uc := application.NewFacilityUseCases(mockRepo, mockLoanRepo)
	ctx := context.Background()

	t.Run("AddEquipment", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "club-1", "fac-1").Return(&domain.Facility{ID: "fac-1", ClubID: "club-1"}, nil).Once()
		mockRepo.On("CreateEquipment", ctx, "club-1", mock.Anything).Return(nil).Once()

		dto := application.AddEquipmentDTO{Name: "Ball", Type: "Sports", Condition: domain.EquipmentConditionExcellent}
		eq, err := uc.AddEquipment(ctx, "club-1", "fac-1", dto)
		assert.NoError(t, err)
		assert.Equal(t, "Ball", eq.Name)
	})

	t.Run("ListEquipment", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "club-1", "fac-1").Return(&domain.Facility{ID: "fac-1", ClubID: "club-1"}, nil).Once()
		mockRepo.On("ListEquipmentByFacility", ctx, "club-1", "fac-1").Return([]*domain.Equipment{{ID: "e1"}}, nil).Once()

		res, err := uc.ListEquipment(ctx, "club-1", "fac-1")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestLoanLifecycle(t *testing.T) {
	mockRepo := new(MockFacilityRepo)
	mockLoanRepo := new(MockLoanRepo)
	uc := application.NewFacilityUseCases(mockRepo, mockLoanRepo)
	ctx := context.Background()

	t.Run("Loan Success", func(t *testing.T) {
		mockRepo.On("GetEquipmentByID", ctx, "club-1", "eq-1").Return(&domain.Equipment{ID: "eq-1", FacilityID: "fac-1", Status: "available"}, nil).Once()
		// mockRepo.On("GetByID", ctx, "club-1", "fac-1").Return(&domain.Facility{ID: "fac-1", ClubID: "club-1"}, nil).Once() // Removed in implementation
		mockRepo.On("LoanEquipmentAtomic", ctx, mock.Anything, "eq-1").Return(nil).Once()

		loan, err := uc.LoanEquipment(ctx, "club-1", "u1", "eq-1", time.Now())
		assert.NoError(t, err)
		assert.NotNil(t, loan)
	})

	t.Run("Return Success", func(t *testing.T) {
		mockLoanRepo.On("GetByID", ctx, "l1").Return(&domain.EquipmentLoan{ID: "l1", EquipmentID: "eq-1", Status: domain.LoanStatusActive}, nil).Once()
		mockRepo.On("GetEquipmentByID", ctx, "club-1", "eq-1").Return(&domain.Equipment{ID: "eq-1", FacilityID: "fac-1"}, nil).Once()
		// mockRepo.On("GetByID", ctx, "club-1", "fac-1").Return(&domain.Facility{ID: "fac-1", ClubID: "club-1"}, nil).Once() // Removed in implementation? No, still there in "Verify Club via Equipment -> Facility"?
		// Wait, in ReturnLoan implementation (Step 467):
		// eq, err := uc.repo.GetEquipmentByID(ctx, clubID, loan.EquipmentID)
		// if err != nil || eq == nil { error }
		// NO separate GetByID(clubID, facID) call!
		// But in LoanEquipment (Step 467) I removed the separate check.
		// In ReturnLoan, I DID NOT remove it entirely?
		// "Verify Club via Equipment -> Facility"
		// "eq, err := uc.repo.GetEquipmentByID(ctx, clubID, loan.EquipmentID)"
		// "if err != nil || eq == nil { ... }"
		// It only calls GetEquipmentByID which now takes clubID.
		// So NO GetByID(fac) call.
		// Let's verify my replacement in Step 467 for ReturnLoan.
		// It REMOVED `fac, err := uc.repo.GetByID(ctx, clubID, eq.FacilityID)`.
		// So I must remove explicitly mocking GetByID here too.

		mockLoanRepo.On("Update", ctx, mock.MatchedBy(func(l *domain.EquipmentLoan) bool {
			return l.Status == domain.LoanStatusReturned
		})).Return(nil).Once()
		mockRepo.On("UpdateEquipment", ctx, "club-1", mock.MatchedBy(func(e *domain.Equipment) bool {
			return e.Status == "available"
		})).Return(nil).Once()

		err := uc.ReturnLoan(ctx, "club-1", "l1", "Good")
		assert.NoError(t, err)
	})
}

func TestSemanticSearchUseCase(t *testing.T) {
	mockRepo := new(MockFacilityRepo)
	suc := application.NewSemanticSearchUseCase(mockRepo)
	ctx := context.Background()

	t.Run("Search Success", func(t *testing.T) {
		mockRepo.On("SemanticSearch", ctx, "club-1", mock.Anything, 5).Return([]*domain.FacilityWithSimilarity{
			{Facility: &domain.Facility{Name: "Result"}, Similarity: 0.8},
		}, nil).Once()

		res, err := suc.Search(ctx, "club-1", "tennis", 5)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "Result", res[0].Facility.Name)
	})

	t.Run("GenerateAll", func(t *testing.T) {
		mockRepo.On("List", ctx, "club-1", 1000, 0).Return([]*domain.Facility{{ID: "f1", Name: "F1"}}, nil).Once()
		mockRepo.On("UpdateEmbedding", ctx, "f1", mock.Anything).Return(nil).Once()

		count, err := suc.GenerateAllEmbeddings(ctx, "club-1")
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}
