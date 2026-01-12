package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	"gorm.io/gorm"
)

type PostgresBookingRepository struct {
	db *gorm.DB
}

func NewPostgresBookingRepository(db *gorm.DB) domain.BookingRepository {
	_ = db.AutoMigrate(&domain.Booking{})
	return &PostgresBookingRepository{db: db}
}

func (r *PostgresBookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *PostgresBookingRepository) GetByID(ctx context.Context, clubID string, id uuid.UUID) (*domain.Booking, error) {
	var booking domain.Booking
	if err := r.db.WithContext(ctx).First(&booking, "id = ? AND club_id = ?", id, clubID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *PostgresBookingRepository) List(ctx context.Context, clubID string, filter map[string]interface{}) ([]domain.Booking, error) {
	var bookings []domain.Booking
	query := r.db.WithContext(ctx).Model(&domain.Booking{}).Where("club_id = ?", clubID)

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}
	// Default sort by start time desc
	query = query.Order("start_time desc")

	if err := query.Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *PostgresBookingRepository) Update(ctx context.Context, booking *domain.Booking) error {
	return r.db.WithContext(ctx).Save(booking).Error
}

func (r *PostgresBookingRepository) HasTimeConflict(ctx context.Context, clubID string, facilityID uuid.UUID, start, end time.Time) (bool, error) {
	var count int64
	// Check for any confirmed booking that overlaps with [start, end)
	// Overlap condition: (ExistingStart < NewEnd) AND (ExistingEnd > NewStart)
	err := r.db.WithContext(ctx).Model(&domain.Booking{}).
		Where("club_id = ?", clubID).
		Where("facility_id = ?", facilityID).
		Where("status IN (?)", []domain.BookingStatus{domain.BookingStatusConfirmed, domain.BookingStatusPendingPayment}).
		Where("start_time < ?", end).
		Where("end_time > ?", start).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresBookingRepository) ListByFacilityAndDate(ctx context.Context, clubID string, facilityID uuid.UUID, date time.Time) ([]domain.Booking, error) {
	var bookings []domain.Booking
	// Filter by facility and date range (start of day to end of day)
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// We want bookings that overlap with this day (though usually bookings are contained within a day)
	// Simple overlap check: Start < EndOfDay AND End > StartOfDay
	err := r.db.WithContext(ctx).Model(&domain.Booking{}).
		Where("club_id = ?", clubID).
		Where("facility_id = ?", facilityID).
		Where("status IN (?)", []domain.BookingStatus{domain.BookingStatusConfirmed, domain.BookingStatusPendingPayment}).
		Where("start_time < ?", endOfDay).
		Where("end_time > ?", startOfDay).
		Order("start_time asc").
		Find(&bookings).Error

	return bookings, err
}

func (r *PostgresBookingRepository) ListAll(ctx context.Context, clubID string, filter map[string]interface{}, from, to *time.Time) ([]domain.Booking, error) {
	var bookings []domain.Booking
	query := r.db.WithContext(ctx).Model(&domain.Booking{}).Where("club_id = ?", clubID)

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if from != nil {
		query = query.Where("start_time >= ?", from)
	}
	if to != nil {
		query = query.Where("start_time <= ?", to)
	}

	query = query.Order("start_time desc")

	if err := query.Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *PostgresBookingRepository) AddToWaitlist(ctx context.Context, entry *domain.Waitlist) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *PostgresBookingRepository) GetNextInLine(ctx context.Context, clubID string, resourceID uuid.UUID, date time.Time) (*domain.Waitlist, error) {
	var entry domain.Waitlist
	err := r.db.WithContext(ctx).
		Where("club_id = ?", clubID).
		Where("resource_id = ?", resourceID).
		Where("target_date = ?", date).
		Where("status = ?", "PENDING").
		Order("created_at asc").
		First(&entry).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entry, nil
}
