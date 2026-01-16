package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestBooking struct {
	ID            uuid.UUID            `gorm:"type:text;primary_key"`
	ClubID        string               `gorm:"index;not null"`
	UserID        uuid.UUID            `gorm:"type:text;not null"`
	FacilityID    uuid.UUID            `gorm:"type:text;not null"`
	StartTime     time.Time            `gorm:"not null"`
	EndTime       time.Time            `gorm:"not null"`
	TotalPrice    float64              `gorm:"type:real"`
	Status        domain.BookingStatus `gorm:"type:text"`
	GuestDetails  domain.GuestDetails  `gorm:"type:text"`
	PaymentExpiry *time.Time           `gorm:"index"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (TestBooking) TableName() string { return "bookings" }

type TestWaitlist struct {
	ID         uuid.UUID `gorm:"type:text;primary_key"`
	ClubID     string    `gorm:"index;not null"`
	ResourceID uuid.UUID `gorm:"type:text;not null"`
	UserID     uuid.UUID `gorm:"type:text;not null"`
	TargetDate time.Time `gorm:"not null"`
	Status     string    `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (TestWaitlist) TableName() string { return "waitlists" }

type BookingRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo domain.BookingRepository
}

func (s *BookingRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	s.Require().NoError(err)
	s.db = db

	err = s.db.AutoMigrate(&TestBooking{}, &TestWaitlist{})
	s.Require().NoError(err)

	s.repo = repository.NewPostgresBookingRepository(s.db)
}

func (s *BookingRepositoryTestSuite) TearDownTest() {
	s.db.Exec("DELETE FROM bookings")
	s.db.Exec("DELETE FROM waitlists")
}

func (s *BookingRepositoryTestSuite) TestCreateAndGet() {
	bookingIDs := []uuid.UUID{uuid.New()}
	for _, id := range bookingIDs {
		booking := &domain.Booking{
			ID:         id,
			ClubID:     "club-1",
			UserID:     uuid.New(),
			FacilityID: uuid.New(),
			StartTime:  time.Now(),
			EndTime:    time.Now().Add(1 * time.Hour),
			Status:     domain.BookingStatusConfirmed,
		}

		err := s.repo.Create(context.Background(), booking)
		s.NoError(err)

		result, err := s.repo.GetByID(context.Background(), "club-1", id)
		s.NoError(err)
		s.NotNil(result)
		s.Equal(id, result.ID)
	}
}

func (s *BookingRepositoryTestSuite) TestListAndListAll() {
	clubID := "list-club"
	facilityID := uuid.New()

	for i := 0; i < 3; i++ {
		b := &domain.Booking{
			ID:         uuid.New(),
			ClubID:     clubID,
			FacilityID: facilityID,
			UserID:     uuid.New(),
			StartTime:  time.Now(),
			EndTime:    time.Now().Add(1 * time.Hour),
			Status:     domain.BookingStatusConfirmed,
		}
		err := s.repo.Create(context.Background(), b)
		s.NoError(err)
	}

	list, err := s.repo.List(context.Background(), clubID, map[string]interface{}{"facility_id": facilityID})
	s.NoError(err)
	s.Len(list, 3)

	all, err := s.repo.ListAll(context.Background(), clubID, nil, nil, nil)
	s.NoError(err)
	s.Len(all, 3)
}

func (s *BookingRepositoryTestSuite) TestUpdate() {
	booking := &domain.Booking{
		ID:         uuid.New(),
		ClubID:     "club-1",
		UserID:     uuid.New(),
		FacilityID: uuid.New(),
		Status:     domain.BookingStatusConfirmed,
	}
	err := s.repo.Create(context.Background(), booking)
	s.NoError(err)

	booking.Status = domain.BookingStatusCancelled
	err = s.repo.Update(context.Background(), booking)
	s.NoError(err)

	updated, _ := s.repo.GetByID(context.Background(), "club-1", booking.ID)
	s.Equal(domain.BookingStatusCancelled, updated.Status)
}

func (s *BookingRepositoryTestSuite) TestHasTimeConflict() {
	fID := uuid.New()
	now := time.Now()

	err := s.repo.Create(context.Background(), &domain.Booking{
		ID: uuid.New(), ClubID: "c1", FacilityID: fID,
		StartTime: now, EndTime: now.Add(1 * time.Hour),
		Status: domain.BookingStatusConfirmed,
	})
	s.NoError(err)

	conflict, err := s.repo.HasTimeConflict(context.Background(), "c1", fID, now.Add(30*time.Minute), now.Add(90*time.Minute))
	s.NoError(err)
	s.True(conflict, "Should conflict with overlapping end")

	// Case 2: Exact match
	conflict, err = s.repo.HasTimeConflict(context.Background(), "c1", fID, now, now.Add(1*time.Hour))
	s.NoError(err)
	s.True(conflict, "Should conflict with exact match")

	// Case 3: Enclosed
	conflict, err = s.repo.HasTimeConflict(context.Background(), "c1", fID, now.Add(10*time.Minute), now.Add(50*time.Minute))
	s.NoError(err)
	s.True(conflict, "Should conflict with enclosed interval")

	// Case 4: Enclosing
	conflict, err = s.repo.HasTimeConflict(context.Background(), "c1", fID, now.Add(-10*time.Minute), now.Add(70*time.Minute))
	s.NoError(err)
	s.True(conflict, "Should conflict with enclosing interval")

	// Case 5: Adjacent (End touches Start) - Should NOT conflict usually
	conflict, err = s.repo.HasTimeConflict(context.Background(), "c1", fID, now.Add(-1*time.Hour), now)
	s.NoError(err)
	s.False(conflict, "Should not conflict if end touches start")

	// Case 6: Adjacent (Start touches End)
	conflict, err = s.repo.HasTimeConflict(context.Background(), "c1", fID, now.Add(1*time.Hour), now.Add(2*time.Hour))
	s.NoError(err)
	s.False(conflict, "Should not conflict if start touches end")
}

func (s *BookingRepositoryTestSuite) TestGetNotFound() {
	res, err := s.repo.GetByID(context.Background(), "club-1", uuid.New())
	s.NoError(err)
	s.Nil(res)
}

func (s *BookingRepositoryTestSuite) TestListWithFilters() {
	clubID := "filter-club"
	fID := uuid.New()
	uID := uuid.New()

	err := s.repo.Create(context.Background(), &domain.Booking{
		ID: uuid.New(), ClubID: clubID, FacilityID: fID, UserID: uID,
		Status: domain.BookingStatusConfirmed, StartTime: time.Now(), EndTime: time.Now().Add(1 * time.Hour),
	})
	s.NoError(err)

	// Filter by both
	list, err := s.repo.List(context.Background(), clubID, map[string]interface{}{
		"facility_id": fID,
		"user_id":     uID,
	})
	s.NoError(err)
	s.Len(list, 1)
}

func (s *BookingRepositoryTestSuite) TestListAllWithDates() {
	clubID := "date-club"
	now := time.Now().Truncate(time.Hour)

	err := s.repo.Create(context.Background(), &domain.Booking{
		ID: uuid.New(), ClubID: clubID, UserID: uuid.New(), FacilityID: uuid.New(),
		StartTime: now.Add(-48 * time.Hour), EndTime: now.Add(-47 * time.Hour),
		Status: domain.BookingStatusConfirmed,
	})
	s.NoError(err)
	err = s.repo.Create(context.Background(), &domain.Booking{
		ID: uuid.New(), ClubID: clubID, UserID: uuid.New(), FacilityID: uuid.New(),
		StartTime: now.Add(24 * time.Hour), EndTime: now.Add(25 * time.Hour),
		Status: domain.BookingStatusConfirmed,
	})
	s.NoError(err)

	from := now.Add(-1 * time.Hour)
	to := now.Add(48 * time.Hour)

	list, err := s.repo.ListAll(context.Background(), clubID, nil, &from, &to)
	s.NoError(err)
	s.Len(list, 1)

	// Filter by facility_id
	fID := uuid.New()
	err = s.repo.Create(context.Background(), &domain.Booking{
		ID: uuid.New(), ClubID: clubID, UserID: uuid.New(), FacilityID: fID,
		Status: domain.BookingStatusConfirmed, StartTime: now, EndTime: now.Add(1 * time.Hour),
	})
	s.NoError(err)
	list, err = s.repo.ListAll(context.Background(), clubID, map[string]interface{}{"facility_id": fID}, nil, nil)
	s.NoError(err)
	s.Len(list, 1)
}

func (s *BookingRepositoryTestSuite) TestWaitlistMultiple() {
	resourceID := uuid.New()
	date := time.Now()
	clubID := "c1"

	for i := 0; i < 3; i++ {
		entry := &domain.Waitlist{
			ID:         uuid.New(),
			ClubID:     clubID,
			ResourceID: resourceID,
			UserID:     uuid.New(),
			TargetDate: date,
			Status:     "PENDING",
			CreatedAt:  time.Now().Add(time.Duration(i) * time.Minute),
		}
		err := s.repo.AddToWaitlist(context.Background(), entry)
		s.NoError(err)
	}

	next, err := s.repo.GetNextInLine(context.Background(), clubID, resourceID, date)
	s.NoError(err)
	s.NotNil(next)
	// Should be the first one added (i=0)
}

func (s *BookingRepositoryTestSuite) TestWaitlist() {
	resourceID := uuid.New()
	date := time.Now()

	entry := &domain.Waitlist{
		ID:         uuid.New(),
		ClubID:     "c1",
		ResourceID: resourceID,
		UserID:     uuid.New(),
		TargetDate: date,
		Status:     "PENDING",
		CreatedAt:  time.Now(),
	}

	err := s.repo.AddToWaitlist(context.Background(), entry)
	s.NoError(err)

	next, err := s.repo.GetNextInLine(context.Background(), "c1", resourceID, date)
	s.NoError(err)
	s.NotNil(next)
	s.Equal(entry.UserID, next.UserID)
}
func (s *BookingRepositoryTestSuite) TestListByFacilityAndDate() {
	clubID := "list-fac-club"
	fID := uuid.New()
	now := time.Now().UTC().Truncate(24 * time.Hour)

	err := s.repo.Create(context.Background(), &domain.Booking{
		ID: uuid.New(), ClubID: clubID, FacilityID: fID, UserID: uuid.New(),
		StartTime: now.Add(10 * time.Hour), EndTime: now.Add(11 * time.Hour),
		Status: domain.BookingStatusConfirmed,
	})
	s.NoError(err)

	list, err := s.repo.ListByFacilityAndDate(context.Background(), clubID, fID, now)
	s.NoError(err)
	s.Len(list, 1)
}

func (s *BookingRepositoryTestSuite) TestWaitlistEmpty() {
	next, err := s.repo.GetNextInLine(context.Background(), "c1", uuid.New(), time.Now())
	s.NoError(err)
	s.Nil(next)
}

func TestBookingRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BookingRepositoryTestSuite))
}
