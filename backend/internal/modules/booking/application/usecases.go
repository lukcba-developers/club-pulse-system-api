package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	bookingDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	facilityDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

// DTOs
type CreateBookingDTO struct {
	UserID       string                      `json:"user_id" binding:"required"`
	FacilityID   string                      `json:"facility_id" binding:"required"`
	StartTime    time.Time                   `json:"start_time" binding:"required"`
	EndTime      time.Time                   `json:"end_time" binding:"required"`
	GuestDetails []bookingDomain.GuestDetail `json:"guest_details"`
}

type CreateRecurringRuleDTO struct {
	FacilityID string                       `json:"facility_id" binding:"required"`
	Type       bookingDomain.RecurrenceType `json:"type" binding:"required"`
	DayOfWeek  int                          `json:"day_of_week" binding:"gte=0,lte=6"`
	StartTime  time.Time                    `json:"start_time" binding:"required"`
	EndTime    time.Time                    `json:"end_time" binding:"required"`
	StartDate  string                       `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate    string                       `json:"end_date" binding:"required"`   // YYYY-MM-DD
}

// BookingUseCases handles core booking logic.
// Refactored to follow SOLID principles:
// - Logic separated into private methods (SRP).
// - Depend on interfaces (DIP).
type BookingUseCases struct {
	repo          bookingDomain.BookingRepository
	recurringRepo bookingDomain.RecurringRepository
	facilityRepo  facilityDomain.FacilityRepository
	userRepo      userDomain.UserRepository
	notifier      service.NotificationSender
}

func NewBookingUseCases(
	repo bookingDomain.BookingRepository,
	recurringRepo bookingDomain.RecurringRepository,
	facilityRepo facilityDomain.FacilityRepository,
	userRepo userDomain.UserRepository,
	notifier service.NotificationSender,
) *BookingUseCases {
	return &BookingUseCases{
		repo:          repo,
		recurringRepo: recurringRepo,
		facilityRepo:  facilityRepo,
		userRepo:      userRepo,
		notifier:      notifier,
	}
}

// CreateBooking orchestrates the booking creation flow: Validate -> Conflict Check -> Persist -> Notify.
func (uc *BookingUseCases) CreateBooking(clubID string, dto CreateBookingDTO) (*bookingDomain.Booking, error) {
	// 1. Parsing & Basic Validation
	userID, facilityID, err := parseBookingIDs(dto)
	if err != nil {
		return nil, err
	}

	if dto.StartTime.After(dto.EndTime) {
		return nil, errors.New("start time must be before end time")
	}

	// 2. Business Rule Validation (Facility Status & Conflicts)
	facility, err := uc.validateBookingRules(clubID, dto.FacilityID, facilityID, dto.StartTime, dto.EndTime)
	if err != nil {
		return nil, err
	}

	// 2.1. Validate User Medical Certificate
	if err := uc.validateUserHealth(clubID, userID.String()); err != nil {
		return nil, err
	}

	// 2.2 Calculate Price (Using already fetched facility)
	dtoDuration := dto.EndTime.Sub(dto.StartTime).Hours()
	basePrice := dtoDuration * facility.HourlyRate
	guestPrice := float64(len(dto.GuestDetails)) * facility.GuestFee
	totalPrice := basePrice + guestPrice

	// 3. Entity Construction
	// Determine initial status based on whether payment is required
	initialStatus := bookingDomain.BookingStatusConfirmed
	var paymentExpiry *time.Time
	if totalPrice > 0 {
		initialStatus = bookingDomain.BookingStatusPendingPayment
		// SECURITY FIX (VUL-001): Set payment expiry to 15 minutes
		expiry := time.Now().Add(15 * time.Minute)
		paymentExpiry = &expiry
	}

	booking := &bookingDomain.Booking{
		ID:            uuid.New(),
		UserID:        userID,
		FacilityID:    facilityID,
		ClubID:        clubID,
		StartTime:     dto.StartTime,
		EndTime:       dto.EndTime,
		TotalPrice:    totalPrice,
		Status:        initialStatus,
		GuestDetails:  dto.GuestDetails,
		PaymentExpiry: paymentExpiry,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 4. Persistence
	if err := uc.repo.Create(booking); err != nil {
		return nil, err
	}

	// 5. Side Effects (Notifications) - Only send confirmation if no payment required
	if initialStatus == bookingDomain.BookingStatusConfirmed {
		uc.notifyAsync(userID.String(), booking.ID.String())
	}

	return booking, nil
}

// ListBookings retrieves bookings with optional filtering.
func (uc *BookingUseCases) ListBookings(clubID string, userID string) ([]bookingDomain.Booking, error) {
	filter := make(map[string]interface{})
	if userID != "" {
		if uid, err := uuid.Parse(userID); err == nil {
			filter["user_id"] = uid
		}
	}
	return uc.repo.List(clubID, filter)
}

// ListClubBookings retrieves all bookings for a club, typically for Admin dashboard.
// Supports filtering by facility and date range.
func (uc *BookingUseCases) ListClubBookings(clubID string, facilityID string, from, to *time.Time) ([]bookingDomain.Booking, error) {
	filter := make(map[string]interface{})

	if facilityID != "" {
		if fid, err := uuid.Parse(facilityID); err == nil {
			filter["facility_id"] = fid
		}
	}

	return uc.repo.ListAll(clubID, filter, from, to)
}

// CancelBooking handles cancellation with authorization check.
func (uc *BookingUseCases) CancelBooking(clubID, bookingID, requestingUserID string) error {
	bID, err := uuid.Parse(bookingID)
	if err != nil {
		return errors.New("invalid booking id")
	}

	booking, err := uc.repo.GetByID(clubID, bID)
	if err != nil {
		return err
	}
	if booking == nil {
		return errors.New("booking not found")
	}

	// Authorization Check
	if requestingUserID != "" && booking.UserID.String() != requestingUserID {
		return errors.New("unauthorized to cancel this booking")
	}

	booking.Status = bookingDomain.BookingStatusCancelled
	booking.UpdatedAt = time.Now()

	if err := uc.repo.Update(booking); err != nil {
		return err
	}

	// Waitlist Logic
	ctx := context.Background()
	next, err := uc.repo.GetNextInLine(ctx, clubID, booking.FacilityID, booking.StartTime)
	if err == nil && next != nil {
		_ = uc.notifier.Send(ctx, service.Notification{
			RecipientID: next.UserID.String(),
			Type:        service.NotificationTypeEmail,
			Subject:     "Slot Available!",
			Message:     "Good news! A slot has opened up for your waitlisted time: " + booking.StartTime.String(),
		})
	}

	return nil
}

// GetAvailability calculates available slots based on business hours and existing bookings.
func (uc *BookingUseCases) GetAvailability(clubID, facilityID string, date time.Time) ([]map[string]interface{}, error) {
	facUUID, err := uuid.Parse(facilityID)
	if err != nil {
		return nil, errors.New("invalid facility id")
	}

	// 1. Fetch dependencies (Facility & Existing Bookings)
	// Suggestion: Use errgroup here for parallel fetching in High Performance scenarios.
	facility, err := uc.facilityRepo.GetByID(clubID, facilityID)
	if err != nil {
		return nil, err
	}
	if facility == nil {
		return nil, errors.New("facility not found")
	}

	bookings, err := uc.repo.ListByFacilityAndDate(clubID, facUUID, date)
	if err != nil {
		return nil, err
	}

	// 1.5. OPTIMIZATION: Fetch Maintenance Tasks Upfront (Avoid N+1)
	allMaintenance, err := uc.facilityRepo.ListMaintenanceByFacility(facilityID)
	if err != nil {
		return nil, err
	}

	// Filter maintenance for this day (Simple In-Memory Filter)
	// In a real high-scale system, we'd add 'ListMaintenanceByDateRange' to the repository.
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var dailyMaintenance []facilityDomain.MaintenanceTask
	for _, m := range allMaintenance {
		// active status check
		if m.Status == facilityDomain.MaintenanceStatusScheduled || m.Status == facilityDomain.MaintenanceStatusInProgress {
			// overlap check
			if m.StartTime.Before(dayEnd) && m.EndTime.After(dayStart) {
				dailyMaintenance = append(dailyMaintenance, *m)
			}
		}
	}

	// 2. Calculate Slots
	startHour := facility.OpeningHour
	endHour := facility.ClosingHour

	// Safety Defaults
	if startHour == 0 && endHour == 0 {
		startHour = 8
		endHour = 23
	}

	var slots []map[string]interface{}

	for h := startHour; h < endHour; h++ {
		slotStart := time.Date(date.Year(), date.Month(), date.Day(), h, 0, 0, 0, date.Location())
		slotEnd := slotStart.Add(1 * time.Hour)

		// Pass pre-fetched maintenance to helper
		status := uc.determineSlotStatusInMemory(slotStart, slotEnd, bookings, dailyMaintenance)

		slots = append(slots, map[string]interface{}{
			"start_time": slotStart.Format("15:04"),
			"end_time":   slotEnd.Format("15:04"),
			"available":  status == "available",
			"status":     status,
		})
	}

	return slots, nil
}

// CreateRecurringRule creates a pattern for future bookings.
func (uc *BookingUseCases) CreateRecurringRule(clubID string, dto CreateRecurringRuleDTO) (*bookingDomain.RecurringRule, error) {
	facID, err := uuid.Parse(dto.FacilityID)
	if err != nil {
		return nil, errors.New("invalid facility id")
	}

	startD, err := time.Parse("2006-01-02", dto.StartDate)
	if err != nil {
		return nil, errors.New("invalid start date format (YYYY-MM-DD)")
	}
	endD, err := time.Parse("2006-01-02", dto.EndDate)
	if err != nil {
		return nil, errors.New("invalid end date format (YYYY-MM-DD)")
	}

	rule := &bookingDomain.RecurringRule{
		ID:         uuid.New(),
		FacilityID: facID,
		ClubID:     clubID,
		Type:       dto.Type,
		DayOfWeek:  dto.DayOfWeek,
		StartTime:  dto.StartTime,
		EndTime:    dto.EndTime,
		StartDate:  startD,
		EndDate:    endD,
	}

	if err := uc.recurringRepo.Create(context.Background(), rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// GenerateBookingsFromRules looks ahead and materializes recurring bookings.
// Refactored to separate logic from loop complexity.
func (uc *BookingUseCases) GenerateBookingsFromRules(ctx context.Context, clubID string, weeks int) error {
	rules, err := uc.recurringRepo.GetAllActive(ctx, clubID)
	if err != nil {
		return err
	}

	horizon := time.Now().AddDate(0, 0, weeks*7)
	generatedCount := 0

	for _, rule := range rules {
		bookings := uc.calculateRecurringBookings(rule, horizon)
		for _, bk := range bookings {
			// Check conflict before creation (Double-check safety)
			conflict, _ := uc.repo.HasTimeConflict(clubID, bk.FacilityID, bk.StartTime, bk.EndTime)
			if !conflict {
				if err := uc.repo.Create(&bk); err == nil {
					generatedCount++
				}
			}
		}
	}
	return nil
}

// --- Private Helpers (The "Clean Code" Section) ---

func parseBookingIDs(dto CreateBookingDTO) (uuid.UUID, uuid.UUID, error) {
	usrID, err := uuid.Parse(dto.UserID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid user id")
	}
	facID, err := uuid.Parse(dto.FacilityID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid facility id")
	}
	return usrID, facID, nil
}

func (uc *BookingUseCases) validateBookingRules(clubID, facilityIDStr string, facilityID uuid.UUID, start, end time.Time) (*facilityDomain.Facility, error) {
	// 1. Check Facility Existence & Status
	facility, err := uc.facilityRepo.GetByID(clubID, facilityIDStr)
	if err != nil {
		return nil, err
	}
	if facility == nil {
		return nil, errors.New("facility not found")
	}
	if facility.Status != facilityDomain.FacilityStatusActive {
		return nil, errors.New("facility is not active (current status: " + string(facility.Status) + ")")
	}

	// 2. Check Existing Bookings
	conflict, err := uc.repo.HasTimeConflict(clubID, facilityID, start, end)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, errors.New("booking time conflict: facility is already booked for this requested time")
	}

	// 3. Check Maintenance Schedules
	maintConflict, err := uc.facilityRepo.HasConflict(clubID, facilityIDStr, start, end)
	if err != nil {
		return nil, err
	}
	if maintConflict {
		return nil, errors.New("booking time conflict: facility is scheduled for maintenance during this time")
	}

	return facility, nil
}

func (uc *BookingUseCases) validateUserHealth(clubID, userID string) error {
	user, err := uc.userRepo.GetByID(clubID, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if user.MedicalCertStatus == nil || *user.MedicalCertStatus != userDomain.MedicalCertStatusValid {
		return errors.New("medical certificate expired or invalid")
	}

	if user.MedicalCertExpiry != nil && user.MedicalCertExpiry.Before(time.Now()) {
		return errors.New("medical certificate expired")
	}

	return nil
}

func (uc *BookingUseCases) notifyAsync(userID, bookingID string) {
	go func() {
		err := uc.notifier.Send(context.Background(), service.Notification{
			RecipientID: userID,
			Type:        service.NotificationTypeEmail,
			Subject:     "Booking Confirmed",
			Message:     "Booking Confirmed: " + bookingID,
		})
		if err != nil {
			// Logger should be injected, but ignoring for now as per previous lint strategy
			_ = err
		}
	}()
}

func (uc *BookingUseCases) determineSlotStatusInMemory(start, end time.Time, bookings []bookingDomain.Booking, maintenance []facilityDomain.MaintenanceTask) string {
	// 1. Check Overlap with Bookings
	for _, b := range bookings {
		if b.StartTime.Before(end) && b.EndTime.After(start) {
			return "booked"
		}
	}

	// 2. Check Overlap with Maintenance (In-Memory)
	for _, m := range maintenance {
		if m.StartTime.Before(end) && m.EndTime.After(start) {
			return "maintenance"
		}
	}

	return "available"
}

func (uc *BookingUseCases) calculateRecurringBookings(rule bookingDomain.RecurringRule, horizon time.Time) []bookingDomain.Booking {
	var bookings []bookingDomain.Booking
	now := time.Now()

	// Determine iteration range
	start := rule.StartDate
	if start.Before(now) {
		start = now
	}
	end := rule.EndDate
	if end.After(horizon) {
		end = horizon
	}

	// System user ID placeholder (extracted from original code)
	systemUser := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	current := start
	for !current.After(end) {
		if int(current.Weekday()) == rule.DayOfWeek {
			// Combine Date(current) + Time(rule.StartTime)
			y, m, d := current.Date()

			h, min, s := rule.StartTime.Clock()
			bookingStart := time.Date(y, m, d, h, min, s, 0, rule.StartTime.Location())

			hEnd, minEnd, sEnd := rule.EndTime.Clock()
			bookingEnd := time.Date(y, m, d, hEnd, minEnd, sEnd, 0, rule.EndTime.Location())

			bookings = append(bookings, bookingDomain.Booking{
				ID:         uuid.New(),
				UserID:     systemUser,
				ClubID:     rule.ClubID,
				FacilityID: rule.FacilityID,
				StartTime:  bookingStart,
				EndTime:    bookingEnd,
				Status:     bookingDomain.BookingStatusConfirmed,
				CreatedAt:  now,
				UpdatedAt:  now,
			})
		}
		current = current.AddDate(0, 0, 1)
	}

	return bookings
}

type JoinWaitlistDTO struct {
	UserID     string    `json:"user_id" binding:"required"`
	ResourceID string    `json:"resource_id" binding:"required"`
	TargetDate time.Time `json:"target_date" binding:"required"`
}

func (uc *BookingUseCases) JoinWaitlist(clubID string, dto JoinWaitlistDTO) (*bookingDomain.Waitlist, error) {
	uid, err := uuid.Parse(dto.UserID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	rid, err := uuid.Parse(dto.ResourceID)
	if err != nil {
		return nil, errors.New("invalid resource id")
	}

	entry := &bookingDomain.Waitlist{
		ID:         uuid.New(),
		ClubID:     clubID,
		UserID:     uid,
		ResourceID: rid,
		TargetDate: dto.TargetDate,
		Status:     "PENDING",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.repo.AddToWaitlist(context.Background(), entry); err != nil {
		return nil, err
	}
	return entry, nil
}
