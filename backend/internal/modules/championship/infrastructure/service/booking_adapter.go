package service

import (
	"time"

	"github.com/google/uuid"
	bookingApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	bookingDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
)

type ChampionshipBookingAdapter struct {
	bookingUC *bookingApp.BookingUseCases
}

func NewChampionshipBookingAdapter(bookingUC *bookingApp.BookingUseCases) *ChampionshipBookingAdapter {
	return &ChampionshipBookingAdapter{bookingUC: bookingUC}
}

func (a *ChampionshipBookingAdapter) CreateSystemBooking(clubID, courtID string, startTime, endTime time.Time, notes string) (*uuid.UUID, error) {
	// We use the Booking UseCase directly.

	// Placeholder System User ID (should be a real UUID string)
	systemUserID := "00000000-0000-0000-0000-000000000000"

	dto := bookingApp.CreateBookingDTO{
		UserID:       systemUserID,
		FacilityID:   courtID,
		StartTime:    startTime,
		EndTime:      endTime,
		GuestDetails: []bookingDomain.GuestDetail{}, // or nil
	}

	// CreateBooking(clubID string, dto CreateBookingDTO)
	booking, err := a.bookingUC.CreateBooking(clubID, dto)
	if err != nil {
		return nil, err
	}

	return &booking.ID, nil
}
