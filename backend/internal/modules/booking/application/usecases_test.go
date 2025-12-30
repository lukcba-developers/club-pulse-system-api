package application_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	bookingDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	facilityDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
)

// --- Mocks ---

type MockBookingRepo struct {
	mock.Mock
}

func (m *MockBookingRepo) Create(booking *bookingDomain.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *MockBookingRepo) GetByID(id uuid.UUID) (*bookingDomain.Booking, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingDomain.Booking), args.Error(1)
}

func (m *MockBookingRepo) List(filter map[string]interface{}) ([]bookingDomain.Booking, error) {
	args := m.Called(filter)
	return args.Get(0).([]bookingDomain.Booking), args.Error(1)
}

func (m *MockBookingRepo) Update(booking *bookingDomain.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *MockBookingRepo) HasTimeConflict(facilityID uuid.UUID, start, end time.Time) (bool, error) {
	args := m.Called(facilityID, start, end)
	return args.Bool(0), args.Error(1)
}

type MockFacilityRepo struct {
	mock.Mock
}

func (m *MockFacilityRepo) Create(facility *facilityDomain.Facility) error {
	args := m.Called(facility)
	return args.Error(0)
}

func (m *MockFacilityRepo) GetByID(id string) (*facilityDomain.Facility, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*facilityDomain.Facility), args.Error(1)
}

func (m *MockFacilityRepo) List(limit, offset int) ([]*facilityDomain.Facility, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*facilityDomain.Facility), args.Error(1)
}

func (m *MockFacilityRepo) Update(facility *facilityDomain.Facility) error {
	args := m.Called(facility)
	return args.Error(0)
}

func (m *MockFacilityRepo) HasConflict(facilityID string, startTime, endTime time.Time) (bool, error) {
	args := m.Called(facilityID, startTime, endTime)
	return args.Bool(0), args.Error(1)
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
		setupMocks    func(mbr *MockBookingRepo, mfr *MockFacilityRepo)
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
			setupMocks: func(mbr *MockBookingRepo, mfr *MockFacilityRepo) {
				// 1. Get Facility -> Active
				mfr.On("GetByID", facilityID).Return(&facilityDomain.Facility{
					ID:     facilityID,
					Status: facilityDomain.FacilityStatusActive,
				}, nil).Once()

				// 2. Check Conflict -> False
				mbr.On("HasTimeConflict", mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()

				// 3. Maintenance Conflict -> False
				mfr.On("HasConflict", mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()

				// 4. Create -> Success
				mbr.On("Create", mock.AnythingOfType("*domain.Booking")).Return(nil).Once()
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
			setupMocks: func(mbr *MockBookingRepo, mfr *MockFacilityRepo) {
				mfr.On("GetByID", facilityID).Return(&facilityDomain.Facility{
					ID:     facilityID,
					Status: facilityDomain.FacilityStatusActive, // Facility is active, but maintenance task exists
				}, nil).Once()

				// 2. Check Conflict -> False (No booking conflict)
				mbr.On("HasTimeConflict", mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()

				// 3. Maintenance Conflict -> True
				mfr.On("HasConflict", mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
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
			setupMocks: func(mbr *MockBookingRepo, mfr *MockFacilityRepo) {
				mfr.On("GetByID", facilityID).Return(&facilityDomain.Facility{
					ID:     facilityID,
					Status: facilityDomain.FacilityStatusActive,
				}, nil).Once()

				mbr.On("HasTimeConflict", mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
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
			setupMocks: func(mbr *MockBookingRepo, mfr *MockFacilityRepo) {
				mfr.On("GetByID", facilityID).Return(&facilityDomain.Facility{
					ID:     facilityID,
					Status: facilityDomain.FacilityStatusActive,
				}, nil).Once()
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
			setupMocks: func(mbr *MockBookingRepo, mfr *MockFacilityRepo) {
				// Return nil directly for facility? Or error?
				// usecases.go checks `if facility == nil` after error check.
				// But typically repo returns error too if not found or nil, nil.
				// Let's assume repo returns nil, nil for not found in this mock setup or a specific error.
				// Looking at usecases logic: `if err != nil { return nil, err }` then `if facility == nil { return nil, not found }`.
				// So we mock returning nil, nil.
				mfr.On("GetByID", facilityID).Return(nil, nil).Once()
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
			mockFacilityRepo := new(MockFacilityRepo)
			useCase := application.NewBookingUseCases(mockBookingRepo, mockFacilityRepo)

			if tc.setupMocks != nil {
				tc.setupMocks(mockBookingRepo, mockFacilityRepo)
			}

			// Execution
			booking, err := useCase.CreateBooking(tc.dto)

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
