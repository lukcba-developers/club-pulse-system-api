package application_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	bookingDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	facilityDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	paymentDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

// --- Mocks ---

type MockBookingRepo struct {
	mock.Mock
}

func (m *MockBookingRepo) Create(ctx context.Context, booking *bookingDomain.Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func (m *MockBookingRepo) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*bookingDomain.Booking, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingDomain.Booking), args.Error(1)
}

func (m *MockBookingRepo) List(ctx context.Context, clubID string, filter map[string]interface{}) ([]bookingDomain.Booking, error) {
	args := m.Called(ctx, clubID, filter)
	var res []bookingDomain.Booking
	if args.Get(0) != nil {
		res = args.Get(0).([]bookingDomain.Booking)
	}
	return res, args.Error(1)
}

func (m *MockBookingRepo) Update(ctx context.Context, booking *bookingDomain.Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func (m *MockBookingRepo) HasTimeConflict(ctx context.Context, clubID string, facilityID uuid.UUID, start, end time.Time) (bool, error) {
	args := m.Called(ctx, clubID, facilityID, start, end)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookingRepo) ListByFacilityAndDate(ctx context.Context, clubID string, facilityID uuid.UUID, date time.Time) ([]bookingDomain.Booking, error) {
	args := m.Called(ctx, clubID, facilityID, date)
	return args.Get(0).([]bookingDomain.Booking), args.Error(1)
}

func (m *MockBookingRepo) ListAll(ctx context.Context, clubID string, filter map[string]interface{}, from, to *time.Time) ([]bookingDomain.Booking, error) {
	args := m.Called(ctx, clubID, filter, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]bookingDomain.Booking), args.Error(1)
}

func (m *MockBookingRepo) AddToWaitlist(ctx context.Context, entry *bookingDomain.Waitlist) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *MockBookingRepo) GetNextInLine(ctx context.Context, clubID string, resourceID uuid.UUID, date time.Time) (*bookingDomain.Waitlist, error) {
	args := m.Called(ctx, clubID, resourceID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingDomain.Waitlist), args.Error(1)
}

type MockFacilityRepo struct {
	mock.Mock
}

func (m *MockFacilityRepo) Create(ctx context.Context, facility *facilityDomain.Facility) error {
	args := m.Called(ctx, facility)
	return args.Error(0)
}

func (m *MockFacilityRepo) GetByID(ctx context.Context, clubID, id string) (*facilityDomain.Facility, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*facilityDomain.Facility), args.Error(1)
}

func (m *MockFacilityRepo) List(ctx context.Context, clubID string, limit, offset int) ([]*facilityDomain.Facility, error) {
	args := m.Called(ctx, clubID, limit, offset)
	return args.Get(0).([]*facilityDomain.Facility), args.Error(1)
}

func (m *MockFacilityRepo) Update(ctx context.Context, facility *facilityDomain.Facility) error {
	args := m.Called(ctx, facility)
	return args.Error(0)
}

func (m *MockFacilityRepo) HasConflict(ctx context.Context, clubID, facilityID string, startTime, endTime time.Time) (bool, error) {
	args := m.Called(ctx, clubID, facilityID, startTime, endTime)
	return args.Bool(0), args.Error(1)
}

func (m *MockFacilityRepo) ListMaintenanceByFacility(ctx context.Context, facilityID string) ([]*facilityDomain.MaintenanceTask, error) {
	args := m.Called(ctx, facilityID)
	return args.Get(0).([]*facilityDomain.MaintenanceTask), args.Error(1)
}

func (m *MockFacilityRepo) SemanticSearch(ctx context.Context, clubID string, embedding []float32, limit int) ([]*facilityDomain.FacilityWithSimilarity, error) {
	args := m.Called(ctx, clubID, embedding, limit)
	return args.Get(0).([]*facilityDomain.FacilityWithSimilarity), args.Error(1)
}

func (m *MockFacilityRepo) UpdateEmbedding(ctx context.Context, facilityID string, embedding []float32) error {
	args := m.Called(ctx, facilityID, embedding)
	return args.Error(0)
}

func (m *MockFacilityRepo) CreateEquipment(ctx context.Context, equipment *facilityDomain.Equipment) error {
	args := m.Called(ctx, equipment)
	return args.Error(0)
}

func (m *MockFacilityRepo) GetEquipmentByID(ctx context.Context, id string) (*facilityDomain.Equipment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*facilityDomain.Equipment), args.Error(1)
}

func (m *MockFacilityRepo) ListEquipmentByFacility(ctx context.Context, facilityID string) ([]*facilityDomain.Equipment, error) {
	args := m.Called(ctx, facilityID)
	return args.Get(0).([]*facilityDomain.Equipment), args.Error(1)
}

func (m *MockFacilityRepo) UpdateEquipment(ctx context.Context, equipment *facilityDomain.Equipment) error {
	args := m.Called(ctx, equipment)
	return args.Error(0)
}

func (m *MockFacilityRepo) LoanEquipmentAtomic(ctx context.Context, loan *facilityDomain.EquipmentLoan, equipmentID string) error {
	args := m.Called(ctx, loan, equipmentID)
	return args.Error(0)
}

type MockNotificationSender struct {
	mock.Mock
}

func (m *MockNotificationSender) Send(ctx context.Context, n service.Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

type MockRefundService struct {
	mock.Mock
}

func (m *MockRefundService) Refund(ctx context.Context, clubID string, referenceID uuid.UUID, referenceType string) error {
	args := m.Called(ctx, clubID, referenceID, referenceType)
	return args.Error(0)
}

type MockRecurringRepo struct {
	mock.Mock
}

func (m *MockRecurringRepo) Create(ctx context.Context, rule *bookingDomain.RecurringRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockRecurringRepo) GetByFacility(ctx context.Context, clubID string, facilityID uuid.UUID) ([]bookingDomain.RecurringRule, error) {
	args := m.Called(ctx, clubID, facilityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]bookingDomain.RecurringRule), args.Error(1)
}

func (m *MockRecurringRepo) GetAllActive(ctx context.Context, clubID string) ([]bookingDomain.RecurringRule, error) {
	args := m.Called(ctx, clubID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]bookingDomain.RecurringRule), args.Error(1)
}

func (m *MockRecurringRepo) Delete(ctx context.Context, clubID string, id uuid.UUID) error {
	args := m.Called(ctx, clubID, id)
	return args.Error(0)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, clubID, id string) (*userDomain.User, error) {
	args := m.Called(ctx, clubID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDomain.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *userDomain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, clubID, id string) error {
	args := m.Called(ctx, clubID, id)
	return args.Error(0)
}

func (m *MockUserRepo) List(ctx context.Context, clubID string, limit, offset int, filters map[string]interface{}) ([]userDomain.User, error) {
	args := m.Called(ctx, clubID, limit, offset, filters)
	return args.Get(0).([]userDomain.User), args.Error(1)
}

func (m *MockUserRepo) FindChildren(ctx context.Context, clubID, parentID string) ([]userDomain.User, error) {
	args := m.Called(ctx, clubID, parentID)
	return args.Get(0).([]userDomain.User), args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, user *userDomain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) CreateIncident(ctx context.Context, incident *userDomain.IncidentLog) error {
	args := m.Called(ctx, incident)
	return args.Error(0)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDomain.User), args.Error(1)
}

func (m *MockUserRepo) ListByIDs(ctx context.Context, clubID string, ids []string) ([]userDomain.User, error) {
	args := m.Called(ctx, clubID, ids)
	return args.Get(0).([]userDomain.User), args.Error(1)
}

func (m *MockUserRepo) AnonymizeForGDPR(ctx context.Context, clubID, id string) error {
	args := m.Called(ctx, clubID, id)
	return args.Error(0)
}

// --- Tests ---

func TestCreateBooking(t *testing.T) {
	userID := uuid.New().String()
	facilityID := uuid.New().String()
	startTime := time.Now().Add(1 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)

	tests := []struct {
		name          string
		dto           application.CreateBookingDTO
		setupMocks    func(mbr *MockBookingRepo, mfr *MockFacilityRepo, mns *MockNotificationSender, mur *MockUserRepo)
		expectedError string
		checkResult   func(t *testing.T, booking *bookingDomain.Booking)
	}{
		{
			name: "Success: Booking created successfully",
			dto: application.CreateBookingDTO{
				UserID:     userID,
				FacilityID: facilityID,
				StartTime:  startTime,
				EndTime:    endTime,
			},
			setupMocks: func(mbr *MockBookingRepo, mfr *MockFacilityRepo, mns *MockNotificationSender, mur *MockUserRepo) {
				// 0. User Health -> Active
				status := userDomain.MedicalCertStatusValid
				now := time.Now().Add(24 * time.Hour)
				mur.On("GetByID", mock.Anything, "test-club", userID).Return(&userDomain.User{
					ID:                userID,
					MedicalCertStatus: &status,
					MedicalCertExpiry: &now,
				}, nil).Once()

				// 1. Get Facility -> Active
				mfr.On("GetByID", mock.Anything, "test-club", facilityID).Return(&facilityDomain.Facility{
					ID:     facilityID,
					Status: facilityDomain.FacilityStatusActive,
				}, nil).Once()

				// 2. Check Conflict -> False
				mbr.On("HasTimeConflict", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()

				// 3. Maintenance Conflict -> False
				mfr.On("HasConflict", mock.Anything, "test-club", facilityID, startTime, endTime).Return(false, nil).Once()

				// 4. Create -> Success
				mbr.On("Create", mock.Anything, mock.AnythingOfType("*domain.Booking")).Return(nil).Once()

				// 5. Notification -> Called (Async)
				mns.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			expectedError: "",
			checkResult: func(t *testing.T, booking *bookingDomain.Booking) {
				assert.NotNil(t, booking)
				assert.Equal(t, bookingDomain.BookingStatusConfirmed, booking.Status)
			},
		},
		{
			name: "Fail: Facility in maintenance",
			dto: application.CreateBookingDTO{
				UserID:     userID,
				FacilityID: facilityID,
				StartTime:  startTime,
				EndTime:    endTime,
			},
			setupMocks: func(mbr *MockBookingRepo, mfr *MockFacilityRepo, mns *MockNotificationSender, mur *MockUserRepo) {
				mfr.On("GetByID", mock.Anything, "test-club", facilityID).Return(&facilityDomain.Facility{
					ID:     facilityID,
					Status: facilityDomain.FacilityStatusActive,
				}, nil).Once()

				// 2. Check Conflict -> False
				mbr.On("HasTimeConflict", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()

				// 3. Maintenance Conflict -> True
				mfr.On("HasConflict", mock.Anything, "test-club", facilityID, startTime, endTime).Return(true, nil).Once()
			},
			expectedError: "scheduled for maintenance",
			checkResult: func(t *testing.T, booking *bookingDomain.Booking) {
				assert.Nil(t, booking)
			},
		},
		{
			name: "Fail: Time conflict",
			dto: application.CreateBookingDTO{
				UserID:     userID,
				FacilityID: facilityID,
				StartTime:  startTime,
				EndTime:    endTime,
			},
			setupMocks: func(mbr *MockBookingRepo, mfr *MockFacilityRepo, mns *MockNotificationSender, mur *MockUserRepo) {
				mfr.On("GetByID", mock.Anything, "test-club", facilityID).Return(&facilityDomain.Facility{
					ID:     facilityID,
					Status: facilityDomain.FacilityStatusActive,
				}, nil).Once()

				mbr.On("HasTimeConflict", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
			},
			expectedError: "booking time conflict",
			checkResult: func(t *testing.T, booking *bookingDomain.Booking) {
				assert.Nil(t, booking)
			},
		},
		{
			name: "Fail: Start time after end time",
			dto: application.CreateBookingDTO{
				UserID:     userID,
				FacilityID: facilityID,
				StartTime:  endTime, // Invalid
				EndTime:    startTime,
			},
			setupMocks: func(mbr *MockBookingRepo, mfr *MockFacilityRepo, mns *MockNotificationSender, mur *MockUserRepo) {
				// No mocks needed
			},
			expectedError: "start time must be before end time",
			checkResult: func(t *testing.T, booking *bookingDomain.Booking) {
				assert.Nil(t, booking)
			},
		},
		{
			name: "Fail: Facility not found",
			dto: application.CreateBookingDTO{
				UserID:     userID,
				FacilityID: facilityID,
				StartTime:  startTime,
				EndTime:    endTime,
			},
			setupMocks: func(mbr *MockBookingRepo, mfr *MockFacilityRepo, mns *MockNotificationSender, mur *MockUserRepo) {
				mfr.On("GetByID", mock.Anything, "test-club", facilityID).Return(nil, nil).Once()
			},
			expectedError: "facility not found",
			checkResult: func(t *testing.T, booking *bookingDomain.Booking) {
				assert.Nil(t, booking)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockBookingRepo := new(MockBookingRepo)
			mockRecurringRepo := new(MockRecurringRepo)
			mockFacilityRepo := new(MockFacilityRepo)
			mockUserRepo := new(MockUserRepo)
			mockNotificationSender := new(MockNotificationSender)
			mockRefundService := new(MockRefundService)
			useCase := application.NewBookingUseCases(mockBookingRepo, mockRecurringRepo, mockFacilityRepo, mockUserRepo, mockNotificationSender, mockRefundService)

			if tc.setupMocks != nil {
				tc.setupMocks(mockBookingRepo, mockFacilityRepo, mockNotificationSender, mockUserRepo)
			}

			// Execution
			booking, err := useCase.CreateBooking(context.Background(), "test-club", tc.dto)

			// Assertions
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if tc.checkResult != nil {
				tc.checkResult(t, booking)
			}

			mockBookingRepo.AssertExpectations(t)
			mockFacilityRepo.AssertExpectations(t)
		})
	}
}

func TestCancelBooking(t *testing.T) {
	clubID := "test-club"
	userID := uuid.New().String()
	bookingID := uuid.New()

	tests := []struct {
		name          string
		setupMocks    func(mbr *MockBookingRepo, mrs *MockRefundService)
		expectedError string
	}{
		{
			name: "Success: User cancels own booking",
			setupMocks: func(mbr *MockBookingRepo, mrs *MockRefundService) {
				mbr.On("GetByID", mock.Anything, clubID, bookingID).Return(&bookingDomain.Booking{
					ID:        bookingID,
					UserID:    uuid.MustParse(userID),
					Status:    bookingDomain.BookingStatusConfirmed,
					StartTime: time.Now().Add(24 * time.Hour),
				}, nil).Once()
				mbr.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
				mbr.On("GetNextInLine", mock.Anything, clubID, mock.Anything, mock.Anything).Return(nil, nil).Once()
				mrs.On("Refund", mock.Anything, clubID, bookingID, "BOOKING").Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name: "Fail: Not owner",
			setupMocks: func(mbr *MockBookingRepo, mrs *MockRefundService) {
				mbr.On("GetByID", mock.Anything, clubID, bookingID).Return(&bookingDomain.Booking{
					ID:     bookingID,
					UserID: uuid.New(), // Different user
				}, nil).Once()
			},
			expectedError: "unauthorized to cancel",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mbr := new(MockBookingRepo)
			mrs := new(MockRefundService)
			uc := application.NewBookingUseCases(mbr, nil, nil, nil, nil, mrs)
			tc.setupMocks(mbr, mrs)

			err := uc.CancelBooking(context.Background(), clubID, bookingID.String(), userID)
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mbr.AssertExpectations(t)
		})
	}
}

func TestGetAvailability(t *testing.T) {
	clubID := "test-club"
	facilityID := uuid.New().String()
	date := time.Now().AddDate(0, 0, 1)

	mbr := new(MockBookingRepo)
	mfr := new(MockFacilityRepo)
	uc := application.NewBookingUseCases(mbr, nil, mfr, nil, nil, nil)

	t.Run("Success: Returns slots", func(t *testing.T) {
		mfr.On("GetByID", mock.Anything, clubID, facilityID).Return(&facilityDomain.Facility{
			ID:          facilityID,
			Status:      facilityDomain.FacilityStatusActive,
			OpeningHour: 8,
			ClosingHour: 22,
		}, nil).Once()
		mfr.On("ListMaintenanceByFacility", mock.Anything, facilityID).Return([]*facilityDomain.MaintenanceTask{}, nil).Once()
		mbr.On("ListByFacilityAndDate", mock.Anything, clubID, uuid.MustParse(facilityID), mock.Anything).Return([]bookingDomain.Booking{}, nil).Once()

		slots, err := uc.GetAvailability(context.Background(), clubID, facilityID, date)
		assert.NoError(t, err)
		assert.NotEmpty(t, slots)
	})
}

func TestCreateRecurringRule(t *testing.T) {
	clubID := "test-club"
	facilityID := uuid.New().String()

	mrr := new(MockRecurringRepo)
	uc := application.NewBookingUseCases(nil, mrr, nil, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		mrr.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		dto := application.CreateRecurringRuleDTO{
			FacilityID: facilityID,
			Type:       bookingDomain.RecurrenceTypeFixed,
			DayOfWeek:  1,
			StartTime:  time.Now(),
			EndTime:    time.Now().Add(1 * time.Hour),
			StartDate:  "2024-01-01",
			EndDate:    "2024-12-31",
		}
		res, err := uc.CreateRecurringRule(context.Background(), clubID, dto)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestJoinWaitlist(t *testing.T) {
	clubID := "test-club"

	mbr := new(MockBookingRepo)
	uc := application.NewBookingUseCases(mbr, nil, nil, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		mbr.On("AddToWaitlist", mock.Anything, mock.Anything).Return(nil).Once()
		dto := application.JoinWaitlistDTO{
			UserID:     uuid.New().String(),
			ResourceID: uuid.New().String(),
			TargetDate: time.Now(),
		}
		res, err := uc.JoinWaitlist(context.Background(), clubID, dto)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestListClubBookings(t *testing.T) {
	clubID := "test-club"
	mbr := new(MockBookingRepo)
	uc := application.NewBookingUseCases(mbr, nil, nil, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		mbr.On("ListAll", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return([]bookingDomain.Booking{}, nil).Once()
		res, err := uc.ListClubBookings(context.Background(), clubID, "", nil, nil)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestOnPaymentStatusChanged(t *testing.T) {
	clubID := "test-club"
	userID := uuid.New().String()
	bookingID := uuid.New()
	mbr := new(MockBookingRepo)
	mur := new(MockUserRepo)
	mns := new(MockNotificationSender)
	mns.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

	uc := application.NewBookingUseCases(mbr, nil, nil, mur, mns, nil)

	t.Run("Payment Completed", func(t *testing.T) {
		status := bookingDomain.BookingStatusPendingPayment
		mbr.On("GetByID", mock.Anything, clubID, bookingID).Return(&bookingDomain.Booking{
			ID: bookingID, Status: status, UserID: uuid.New(),
		}, nil).Once()
		mbr.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		mur.On("GetByID", mock.Anything, clubID, mock.Anything).Return(&userDomain.User{ID: "user-1"}, nil).Maybe()
		mbr.On("ListAll", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return([]bookingDomain.Booking{}, nil).Maybe()
		mur.On("Update", mock.Anything, mock.Anything).Return(nil).Maybe()

		err := uc.OnPaymentStatusChanged(context.Background(), clubID, bookingID, paymentDomain.PaymentStatusCompleted)
		assert.NoError(t, err)
		// Wait a bit for async notification
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Payment Completed - Streak Increment", func(t *testing.T) {
		yesterday := time.Now().AddDate(0, 0, -1)
		// Mock GetByID for the user to return a user with existing stats
		mur.On("GetByID", mock.Anything, clubID, userID).Return(&userDomain.User{
			ID: userID,
			Stats: &userDomain.UserStats{
				UserID:           userID,
				LastActivityDate: &yesterday,
				CurrentStreak:    5,
				LongestStreak:    5,
			},
		}, nil).Once()

		// Mock Update for the user to check if streak is incremented
		mur.On("Update", mock.Anything, mock.MatchedBy(func(u *userDomain.User) bool {
			return u.Stats != nil && u.Stats.CurrentStreak == 6
		})).Return(nil).Once()

		// Mock GetByID for the booking
		mbr.On("GetByID", mock.Anything, clubID, bookingID).Return(&bookingDomain.Booking{
			ID: bookingID, Status: bookingDomain.BookingStatusPendingPayment, UserID: uuid.MustParse(userID),
		}, nil).Once()
		// Mock Update for the booking
		mbr.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		// Mock ListAll for notifications (can be empty)
		mbr.On("ListAll", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return([]bookingDomain.Booking{}, nil).Maybe()

		err := uc.OnPaymentStatusChanged(context.Background(), clubID, bookingID, paymentDomain.PaymentStatusCompleted)
		assert.NoError(t, err)
		time.Sleep(10 * time.Millisecond) // Wait for async notification
	})
}

func TestGenerateBookingsFromRules(t *testing.T) {
	clubID := "test-club"
	mrr := new(MockRecurringRepo)
	mbr := new(MockBookingRepo)
	uc := application.NewBookingUseCases(mbr, mrr, nil, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		startDate := time.Now()
		endDate := startDate.AddDate(0, 0, 30)
		mrr.On("GetAllActive", mock.Anything, clubID).Return([]bookingDomain.RecurringRule{
			{
				ID: uuid.New(), ClubID: clubID, FacilityID: uuid.New(),
				Type: bookingDomain.RecurrenceTypeFixed, DayOfWeek: int(startDate.Weekday()),
				StartTime: startDate, EndTime: startDate.Add(1 * time.Hour),
				StartDate: startDate, EndDate: endDate,
			},
		}, nil).Once()
		mbr.On("HasTimeConflict", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
		mbr.On("Create", mock.Anything, mock.Anything).Return(nil)

		err := uc.GenerateBookingsFromRules(context.Background(), clubID, 1)
		assert.NoError(t, err)
	})
}

func TestListBookings(t *testing.T) {
	clubID := "test-club"
	userID := uuid.New().String()
	mbr := new(MockBookingRepo)
	uc := application.NewBookingUseCases(mbr, nil, nil, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		mbr.On("List", mock.Anything, clubID, mock.Anything).Return([]bookingDomain.Booking{}, nil).Once()
		res, err := uc.ListBookings(context.Background(), clubID, userID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestGetAvailabilityDetailed(t *testing.T) {
	clubID := "test-club"
	facilityID := uuid.New().String()
	date := time.Now().AddDate(0, 0, 1)
	mfr := new(MockFacilityRepo)
	mbr := new(MockBookingRepo)
	uc := application.NewBookingUseCases(mbr, nil, mfr, nil, nil, nil)

	t.Run("With Bookings and Maintenance", func(t *testing.T) {
		mfr.On("GetByID", mock.Anything, clubID, facilityID).Return(&facilityDomain.Facility{
			ID: facilityID, Status: facilityDomain.FacilityStatusActive,
			OpeningHour: 8, ClosingHour: 10,
		}, nil).Once()

		mfr.On("ListMaintenanceByFacility", mock.Anything, facilityID).Return([]*facilityDomain.MaintenanceTask{
			{StartTime: time.Now().AddDate(0, 0, 1)},
		}, nil).Once()

		mbr.On("ListByFacilityAndDate", mock.Anything, clubID, uuid.MustParse(facilityID), mock.Anything).Return([]bookingDomain.Booking{
			{StartTime: time.Now().AddDate(0, 0, 1), EndTime: time.Now().AddDate(0, 0, 1).Add(1 * time.Hour)},
		}, nil).Once()

		slots, err := uc.GetAvailability(context.Background(), clubID, facilityID, date)
		assert.NoError(t, err)
		assert.NotEmpty(t, slots)
	})
}

func TestCreateBooking_MedicalFail(t *testing.T) {
	clubID := "test-club"
	userID := uuid.New().String()
	facilityID := uuid.New()

	mur := new(MockUserRepo)
	mn := new(MockNotificationSender)
	mn.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()
	mfr := new(MockFacilityRepo)
	mbr := new(MockBookingRepo)
	uc := application.NewBookingUseCases(mbr, nil, mfr, mur, mn, nil)

	t.Run("CreateBooking_ZeroPrice", func(t *testing.T) {
		mr := new(MockBookingRepo)
		fr := new(MockFacilityRepo)
		ur := new(MockUserRepo)
		mn := new(MockNotificationSender)
		mn.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()
		uCases := application.NewBookingUseCases(mr, nil, fr, ur, mn, nil)

		facilityID := uuid.New()
		uID := uuid.New()
		now := time.Now().Add(1 * time.Hour)

		medicalStatus := userDomain.MedicalCertStatusValid
		ur.On("GetByID", mock.Anything, clubID, uID.String()).Return(&userDomain.User{
			ID: uID.String(), MedicalCertStatus: &medicalStatus,
		}, nil).Once()
		fr.On("GetByID", mock.Anything, clubID, facilityID.String()).Return(&facilityDomain.Facility{
			ID:         facilityID.String(),
			Status:     facilityDomain.FacilityStatusActive,
			HourlyRate: 0,
			GuestFee:   0,
		}, nil).Once()
		mr.On("HasTimeConflict", mock.Anything, clubID, facilityID, mock.Anything, mock.Anything).Return(false, nil).Once()
		fr.On("HasConflict", mock.Anything, clubID, facilityID.String(), mock.Anything, mock.Anything).Return(false, nil).Once()
		mr.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		dto := application.CreateBookingDTO{
			UserID: uID.String(), FacilityID: facilityID.String(),
			StartTime: now, EndTime: now.Add(1 * time.Hour),
		}
		res, err := uCases.CreateBooking(context.Background(), clubID, dto)
		assert.NoError(t, err)
		assert.Equal(t, bookingDomain.BookingStatusConfirmed, res.Status)
	})

	t.Run("Medical Certificate Expired", func(t *testing.T) {
		status := userDomain.MedicalCertStatusExpired
		mur.On("GetByID", mock.Anything, clubID, userID).Return(&userDomain.User{
			ID: userID, MedicalCertStatus: &status,
		}, nil).Once()
		mfr.On("GetByID", mock.Anything, clubID, facilityID.String()).Return(&facilityDomain.Facility{
			ID: facilityID.String(), Status: facilityDomain.FacilityStatusActive,
		}, nil).Once()

		// Reach medical validation by ensuring previous checks pass
		mbr.On("HasTimeConflict", mock.Anything, clubID, facilityID, mock.Anything, mock.Anything).Return(false, nil).Once()
		mfr.On("HasConflict", mock.Anything, clubID, facilityID.String(), mock.Anything, mock.Anything).Return(false, nil).Once()

		_, err := uc.CreateBooking(context.Background(), clubID, application.CreateBookingDTO{
			UserID:     userID,
			FacilityID: facilityID.String(),
			StartTime:  time.Now(),
			EndTime:    time.Now().Add(1 * time.Hour),
		})
		assert.Error(t, err)
	})
}

func TestGetAvailabilityExtra(t *testing.T) {
	clubID := "test-club"
	facilityID := uuid.New()
	mbr := new(MockBookingRepo)
	mfr := new(MockFacilityRepo)
	uc := application.NewBookingUseCases(mbr, nil, mfr, nil, nil, nil)

	t.Run("Complex Availability - Multiple Bookings", func(t *testing.T) {
		date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		mfr.On("GetByID", mock.Anything, clubID, facilityID.String()).Return(&facilityDomain.Facility{
			ID: facilityID.String(), Status: facilityDomain.FacilityStatusActive, OpeningHour: 8, ClosingHour: 10,
		}, nil).Once()
		mfr.On("ListMaintenanceByFacility", mock.Anything, facilityID.String()).Return([]*facilityDomain.MaintenanceTask{}, nil).Once()

		startTime := date.Add(8 * time.Hour)
		mbr.On("ListByFacilityAndDate", mock.Anything, clubID, facilityID, mock.MatchedBy(func(d time.Time) bool {
			return d.Format("2006-01-02") == date.Format("2006-01-02")
		})).Return([]bookingDomain.Booking{
			{ID: uuid.New(), StartTime: startTime, EndTime: startTime.Add(30 * time.Minute), Status: bookingDomain.BookingStatusConfirmed},
		}, nil).Once()

		res, err := uc.GetAvailability(context.Background(), clubID, facilityID.String(), date)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		found := false
		for _, s := range res {
			if s["start_time"] == "08:00" {
				assert.Equal(t, "booked", s["status"])
				found = true
			}
		}
		assert.True(t, found)
	})
}

func TestRecurringRuleEdgeCases(t *testing.T) {
	clubID := "test-club"
	mrr := new(MockRecurringRepo)
	uc := application.NewBookingUseCases(nil, mrr, nil, nil, nil, nil)

	t.Run("RecurringRuleEdgeCases/Create_Recurring_Rule_-_Invalid_Dates", func(t *testing.T) {
		dto := application.CreateRecurringRuleDTO{
			FacilityID: uuid.New().String(),
			StartDate:  time.Now().Add(24 * time.Hour).Format("2006-01-02"),
			EndDate:    time.Now().Format("2006-01-02"), // End before start
			StartTime:  time.Now(),
			EndTime:    time.Now().Add(1 * time.Hour),
		}
		_, err := uc.CreateRecurringRule(context.Background(), clubID, dto)
		assert.Error(t, err)
	})

	t.Run("UpdateUserStreak_NoUpdate", func(t *testing.T) {
		mr := new(MockBookingRepo)
		ur := new(MockUserRepo)
		mn := new(MockNotificationSender)
		uCases := application.NewBookingUseCases(mr, nil, nil, ur, mn, nil)
		bookingID := uuid.New()
		uID := uuid.New()

		// Mock GetByID for the initial booking fetch
		mr.On("GetByID", mock.Anything, clubID, bookingID).Return(&bookingDomain.Booking{
			ID:     bookingID,
			UserID: uID,
			Status: bookingDomain.BookingStatusConfirmed,
		}, nil).Once()

		// Mock List to return existing bookings this month
		mr.On("List", mock.Anything, clubID, mock.Anything).Return([]bookingDomain.Booking{{ID: uuid.New()}}, nil).Once()

		// Mock Notification
		mn.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

		// Mock User fetch for XP (async)
		ur.On("GetByID", mock.Anything, clubID, uID.String()).Return(&userDomain.User{ID: uID.String()}, nil).Maybe()

		// Mock User Update for XP (async)
		ur.On("Update", mock.Anything, mock.Anything).Return(nil).Maybe()

		// Mock Repo Update (final step of OnPaymentStatusChanged)
		mr.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		// Mock ListAll for isFirstBookingOfMonth (called by awardBookingXP goroutine)
		mr.On("ListAll", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return([]bookingDomain.Booking{}, nil).Maybe()

		err := uCases.OnPaymentStatusChanged(context.Background(), clubID, bookingID, "COMPLETED")
		assert.NoError(t, err)

		// Give a tiny bit of time for async goroutines to avoid panic if they run after test ends
		time.Sleep(20 * time.Millisecond)
	})

	t.Run("CreateBooking_MedicalCertInvalid", func(t *testing.T) {
		mr := new(MockBookingRepo)
		fr := new(MockFacilityRepo)
		ur := new(MockUserRepo)
		uCases := application.NewBookingUseCases(mr, nil, fr, ur, nil, nil)
		uID := uuid.New().String()

		fr.On("GetByID", mock.Anything, clubID, mock.Anything).Return(&facilityDomain.Facility{
			ID:     uuid.New().String(),
			Status: facilityDomain.FacilityStatusActive,
		}, nil).Once()
		fr.On("ListMaintenanceByFacility", mock.Anything, mock.Anything).Return([]*facilityDomain.MaintenanceTask{}, nil).Once()
		fr.On("HasConflict", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()
		mr.On("HasTimeConflict", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()

		status := userDomain.MedicalCertStatusExpired
		ur.On("GetByID", mock.Anything, clubID, uID).Return(&userDomain.User{
			ID:                uID,
			MedicalCertStatus: &status,
		}, nil).Once()

		dto := application.CreateBookingDTO{
			UserID:     uID,
			FacilityID: uuid.New().String(),
			StartTime:  time.Now().Add(1 * time.Hour),
			EndTime:    time.Now().Add(2 * time.Hour),
		}
		_, err := uCases.CreateBooking(context.Background(), clubID, dto)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "medical certificate expired or invalid")
	})

	t.Run("CreateBooking_MedicalCertExpired", func(t *testing.T) {
		mr := new(MockBookingRepo)
		fr := new(MockFacilityRepo)
		ur := new(MockUserRepo)
		uCases := application.NewBookingUseCases(mr, nil, fr, ur, nil, nil)
		uID := uuid.New().String()

		fr.On("GetByID", mock.Anything, clubID, mock.Anything).Return(&facilityDomain.Facility{
			ID:     uuid.New().String(),
			Status: facilityDomain.FacilityStatusActive,
		}, nil).Once()
		fr.On("ListMaintenanceByFacility", mock.Anything, mock.Anything).Return([]*facilityDomain.MaintenanceTask{}, nil).Once()
		fr.On("HasConflict", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()
		mr.On("HasTimeConflict", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()

		status := userDomain.MedicalCertStatusValid
		expiry := time.Now().Add(-1 * time.Hour)
		ur.On("GetByID", mock.Anything, clubID, uID).Return(&userDomain.User{
			ID:                uID,
			MedicalCertStatus: &status,
			MedicalCertExpiry: &expiry,
		}, nil).Once()

		dto := application.CreateBookingDTO{
			UserID:     uID,
			FacilityID: uuid.New().String(),
			StartTime:  time.Now().Add(1 * time.Hour),
			EndTime:    time.Now().Add(2 * time.Hour),
		}
		_, err := uCases.CreateBooking(context.Background(), clubID, dto)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "medical certificate expired")
	})

	t.Run("GetAvailability_DetailedStatus", func(t *testing.T) {
		mr := new(MockBookingRepo)
		fr := new(MockFacilityRepo)
		uCases := application.NewBookingUseCases(mr, nil, fr, nil, nil, nil)
		fID := uuid.New().String()

		fr.On("GetByID", mock.Anything, clubID, fID).Return(&facilityDomain.Facility{
			ID:          fID,
			OpeningHour: 8,
			ClosingHour: 22,
		}, nil).Once()
		fr.On("ListMaintenanceByFacility", mock.Anything, fID).Return([]*facilityDomain.MaintenanceTask{}, nil).Once()

		now := time.Now().Truncate(24 * time.Hour).Add(10 * time.Hour) // 10:00 AM
		mr.On("ListByFacilityAndDate", mock.Anything, clubID, uuid.MustParse(fID), mock.Anything).Return([]bookingDomain.Booking{
			{StartTime: now, EndTime: now.Add(1 * time.Hour), Status: bookingDomain.BookingStatusConfirmed},
		}, nil).Once()

		res, err := uCases.GetAvailability(context.Background(), clubID, fID, now)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("CreateRecurringRule_InvalidFacilityID", func(t *testing.T) {
		uCases := application.NewBookingUseCases(nil, nil, nil, nil, nil, nil)
		dto := application.CreateRecurringRuleDTO{
			FacilityID: "invalid-uuid",
			Type:       bookingDomain.RecurrenceTypeFixed,
			DayOfWeek:  1,
			StartTime:  time.Now(),
			EndTime:    time.Now().Add(1 * time.Hour),
			StartDate:  "2025-01-01",
			EndDate:    "2025-01-02",
		}
		_, err := uCases.CreateRecurringRule(context.Background(), clubID, dto)
		assert.Error(t, err)
	})

	t.Run("OnPaymentStatusChanged_FirstBookingError", func(t *testing.T) {
		mr := new(MockBookingRepo)
		ur := new(MockUserRepo)
		mn := new(MockNotificationSender)
		uCases := application.NewBookingUseCases(mr, nil, nil, ur, mn, nil)
		bookingID := uuid.New()
		uID := uuid.New()

		mr.On("GetByID", mock.Anything, clubID, bookingID).Return(&bookingDomain.Booking{
			ID: bookingID, UserID: uID, Status: bookingDomain.BookingStatusConfirmed,
		}, nil).Once()
		mr.On("ListAll", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("list error")).Once()
		mr.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		mn.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()
		ur.On("GetByID", mock.Anything, clubID, uID.String()).Return(&userDomain.User{ID: uID.String(), Stats: &userDomain.UserStats{}}, nil).Maybe()
		ur.On("Update", mock.Anything, mock.Anything).Return(nil).Maybe()

		err := uCases.OnPaymentStatusChanged(context.Background(), clubID, bookingID, "COMPLETED")
		assert.NoError(t, err)            // Error in firstBookingOfMonth is logged but doesn't fail the whole process
		time.Sleep(20 * time.Millisecond) // Give time for async notify
	})

	t.Run("CreateRecurringRule_InvalidDate", func(t *testing.T) {
		uCases := application.NewBookingUseCases(nil, nil, nil, nil, nil, nil)
		dto := application.CreateRecurringRuleDTO{
			FacilityID: uuid.New().String(),
			Type:       bookingDomain.RecurrenceTypeFixed,
			DayOfWeek:  1,
			StartTime:  time.Now(),
			EndTime:    time.Now().Add(1 * time.Hour),
			StartDate:  "invalid-date",
			EndDate:    "2025-01-01",
		}
		_, err := uCases.CreateRecurringRule(context.Background(), clubID, dto)
		assert.Error(t, err)
	})
}

func TestListClubBookings_Extra(t *testing.T) {
	clubID := "test-club"
	t.Run("GetAvailability_RepoError", func(t *testing.T) {
		fr := new(MockFacilityRepo)
		uCases := application.NewBookingUseCases(nil, nil, fr, nil, nil, nil)
		fID := uuid.New().String()
		fr.On("GetByID", mock.Anything, clubID, fID).Return(nil, fmt.Errorf("db error")).Once()

		_, err := uCases.GetAvailability(context.Background(), clubID, fID, time.Now())
		assert.Error(t, err)
	})

	t.Run("OnPaymentStatusChanged_NonCompleted", func(t *testing.T) {
		mr := new(MockBookingRepo)
		uCases := application.NewBookingUseCases(mr, nil, nil, nil, nil, nil)
		bookingID := uuid.New()
		mr.On("GetByID", mock.Anything, clubID, bookingID).Return(&bookingDomain.Booking{ID: bookingID}, nil).Once()
		mr.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		err := uCases.OnPaymentStatusChanged(context.Background(), clubID, bookingID, "PENDING")
		assert.NoError(t, err)
	})

	t.Run("GenerateBookingsFromRules_RepoError", func(t *testing.T) {
		rr := new(MockRecurringRepo)
		uCases := application.NewBookingUseCases(nil, rr, nil, nil, nil, nil)
		rr.On("GetAllActive", mock.Anything, clubID).Return(nil, fmt.Errorf("db error")).Once()

		err := uCases.GenerateBookingsFromRules(context.Background(), clubID, 4)
		assert.Error(t, err)
	})

	t.Run("ListBookings_RepoError", func(t *testing.T) {
		mbr := new(MockBookingRepo)
		uc := application.NewBookingUseCases(mbr, nil, nil, nil, nil, nil)
		uID := uuid.New().String()
		mbr.On("List", mock.Anything, clubID, mock.Anything).Return(nil, fmt.Errorf("db error")).Once()

		_, err := uc.ListBookings(context.Background(), clubID, uID)
		assert.Error(t, err)
	})

	t.Run("CreateRecurringRule_FacilityError", func(t *testing.T) {
		fr := new(MockFacilityRepo)
		uCases := application.NewBookingUseCases(nil, nil, fr, nil, nil, nil)
		fID := uuid.New().String()
		fr.On("GetByID", mock.Anything, clubID, fID).Return(nil, fmt.Errorf("db error")).Once()

		_, err := uCases.CreateRecurringRule(context.Background(), clubID, application.CreateRecurringRuleDTO{
			FacilityID: fID, StartDate: "2024-01-01", EndDate: "2024-02-01",
		})
		assert.Error(t, err)
	})

	t.Run("CreateRecurringRule_CreateError", func(t *testing.T) {
		fr := new(MockFacilityRepo)
		rr := new(MockRecurringRepo)
		uCases := application.NewBookingUseCases(nil, rr, fr, nil, nil, nil)
		fID := uuid.New().String()
		fr.On("GetByID", mock.Anything, clubID, fID).Return(&facilityDomain.Facility{ID: fID}, nil).Once()
		rr.On("Create", mock.Anything, mock.Anything).Return(fmt.Errorf("db error")).Once()

		_, err := uCases.CreateRecurringRule(context.Background(), clubID, application.CreateRecurringRuleDTO{
			FacilityID: fID, StartDate: "2024-01-01", EndDate: "2024-02-01",
		})
		assert.Error(t, err)
	})

	t.Run("OnPaymentStatusChanged_RepoError", func(t *testing.T) {
		mr := new(MockBookingRepo)
		uCases := application.NewBookingUseCases(mr, nil, nil, nil, nil, nil)
		bookingID := uuid.New()
		mr.On("GetByID", mock.Anything, clubID, bookingID).Return(nil, fmt.Errorf("db error")).Once()

		err := uCases.OnPaymentStatusChanged(context.Background(), clubID, bookingID, "COMPLETED")
		assert.Error(t, err)
	})

	t.Run("ListClubBookings_RepoError", func(t *testing.T) {
		mr := new(MockBookingRepo)
		uCases := application.NewBookingUseCases(mr, nil, nil, nil, nil, nil)
		mr.On("ListAll", mock.Anything, clubID, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error")).Once()

		_, err := uCases.ListClubBookings(context.Background(), clubID, "", nil, nil)
		assert.Error(t, err)
	})

	t.Run("CreateBooking_WithPrice", func(t *testing.T) {
		mr := new(MockBookingRepo)
		fr := new(MockFacilityRepo)
		ur := new(MockUserRepo)
		mn := new(MockNotificationSender)
		mn.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

		status := userDomain.MedicalCertStatusValid
		uCases := application.NewBookingUseCases(mr, nil, fr, ur, mn, nil)

		fID := uuid.New()
		uID := uuid.New()

		ur.On("GetByID", mock.Anything, clubID, uID.String()).Return(&userDomain.User{ID: uID.String(), MedicalCertStatus: &status}, nil).Once()
		fr.On("GetByID", mock.Anything, clubID, fID.String()).Return(&facilityDomain.Facility{
			ID: fID.String(), Status: facilityDomain.FacilityStatusActive, HourlyRate: 50.0,
		}, nil).Once()
		mr.On("HasTimeConflict", mock.Anything, clubID, fID, mock.Anything, mock.Anything).Return(false, nil).Once()
		fr.On("HasConflict", mock.Anything, clubID, fID.String(), mock.Anything, mock.Anything).Return(false, nil).Once()
		mr.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		_, err := uCases.CreateBooking(context.Background(), clubID, application.CreateBookingDTO{
			UserID: uID.String(), FacilityID: fID.String(),
			StartTime: time.Now().Add(1 * time.Hour), EndTime: time.Now().Add(2 * time.Hour),
		})
		assert.NoError(t, err)
	})

	t.Run("CreateBooking_MaintenanceError", func(t *testing.T) {
		mr := new(MockBookingRepo)
		fr := new(MockFacilityRepo)
		ur := new(MockUserRepo)
		mn := new(MockNotificationSender)
		mn.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

		status := userDomain.MedicalCertStatusValid
		uCases := application.NewBookingUseCases(mr, nil, fr, ur, mn, nil)

		fID := uuid.New()
		uID := uuid.New()

		ur.On("GetByID", mock.Anything, clubID, uID.String()).Return(&userDomain.User{ID: uID.String(), MedicalCertStatus: &status}, nil).Once()
		fr.On("GetByID", mock.Anything, clubID, fID.String()).Return(&facilityDomain.Facility{
			ID: fID.String(), Status: facilityDomain.FacilityStatusActive,
		}, nil).Once()
		mr.On("HasTimeConflict", mock.Anything, clubID, fID, mock.Anything, mock.Anything).Return(false, nil).Once()
		fr.On("HasConflict", mock.Anything, clubID, fID.String(), mock.Anything, mock.Anything).Return(false, fmt.Errorf("db error")).Once()

		_, err := uCases.CreateBooking(context.Background(), clubID, application.CreateBookingDTO{
			UserID: uID.String(), FacilityID: fID.String(),
			StartTime: time.Now().Add(1 * time.Hour), EndTime: time.Now().Add(2 * time.Hour),
		})
		assert.Error(t, err)
	})

	t.Run("JoinWaitlist_RepoError", func(t *testing.T) {
		mr := new(MockBookingRepo)
		uc := application.NewBookingUseCases(mr, nil, nil, nil, nil, nil)
		fID := uuid.New().String()
		uID := uuid.New().String()
		mr.On("AddToWaitlist", mock.Anything, mock.Anything).Return(fmt.Errorf("db error")).Once()

		_, err := uc.JoinWaitlist(context.Background(), clubID, application.JoinWaitlistDTO{
			UserID: uID, ResourceID: fID, TargetDate: time.Now(),
		})
		assert.Error(t, err)
	})
}
