package repository

import (
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

func (r *PostgresBookingRepository) Create(booking *domain.Booking) error {
	return r.db.Create(booking).Error
}

func (r *PostgresBookingRepository) GetByID(id uuid.UUID) (*domain.Booking, error) {
	var booking domain.Booking
	if err := r.db.First(&booking, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *PostgresBookingRepository) List(filter map[string]interface{}) ([]domain.Booking, error) {
	var bookings []domain.Booking
	query := r.db.Model(&domain.Booking{})

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

func (r *PostgresBookingRepository) Update(booking *domain.Booking) error {
	return r.db.Save(booking).Error
}

func (r *PostgresBookingRepository) HasTimeConflict(facilityID uuid.UUID, start, end time.Time) (bool, error) {
	var count int64
	// Check for any confirmed booking that overlaps with [start, end)
	// Overlap condition: (ExistingStart < NewEnd) AND (ExistingEnd > NewStart)
	err := r.db.Model(&domain.Booking{}).
		Where("facility_id = ?", facilityID).
		Where("status = ?", domain.BookingStatusConfirmed).
		Where("start_time < ?", end).
		Where("end_time > ?", start).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
