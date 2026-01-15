package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	bookingDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/domain"
	facilityDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/domain"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	paymentDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/shopspring/decimal"
)

// DTOs
type CreateBookingDTO struct {
	UserID       string                      `json:"user_id"` // Auth override
	FacilityID   string                      `json:"facility_id" binding:"required"`
	StartTime    time.Time                   `json:"start_time" binding:"required"`
	EndTime      time.Time                   `json:"end_time" binding:"required"`
	GuestDetails []bookingDomain.GuestDetail `json:"guest_details"`
}

type CreateRecurringRuleDTO struct {
	FacilityID string                       `json:"facility_id" binding:"required"`
	Type       bookingDomain.RecurrenceType `json:"type" binding:"required"`
	Frequency  string                       `json:"frequency" binding:"required,oneof=WEEKLY MONTHLY"`
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
	refundSvc     bookingDomain.RefundService
}

func NewBookingUseCases(
	repo bookingDomain.BookingRepository,
	recurringRepo bookingDomain.RecurringRepository,
	facilityRepo facilityDomain.FacilityRepository,
	userRepo userDomain.UserRepository,
	notifier service.NotificationSender,
	refundSvc bookingDomain.RefundService,
) *BookingUseCases {
	return &BookingUseCases{
		repo:          repo,
		recurringRepo: recurringRepo,
		facilityRepo:  facilityRepo,
		userRepo:      userRepo,
		notifier:      notifier,
		refundSvc:     refundSvc,
	}
}

// CreateBooking orchestrates the booking creation flow: Validate -> Conflict Check -> Persist -> Notify.
func (uc *BookingUseCases) CreateBooking(ctx context.Context, clubID string, dto CreateBookingDTO) (*bookingDomain.Booking, error) {
	// 1. Parsing & Basic Validation
	userID, facilityID, err := parseBookingIDs(dto)
	if err != nil {
		return nil, err
	}

	if dto.StartTime.After(dto.EndTime) {
		return nil, errors.New("start time must be before end time")
	}

	// Timezone Normalization: Force UTC if offset is present to avoid "Timezone Hell"
	if dto.StartTime.Location() != time.UTC {
		dto.StartTime = dto.StartTime.UTC()
	}
	if dto.EndTime.Location() != time.UTC {
		dto.EndTime = dto.EndTime.UTC()
	}

	if dto.StartTime.Before(time.Now()) {
		return nil, errors.New("cannot book in the past")
	}

	// Validate Guest Details Integrity
	for _, guest := range dto.GuestDetails {
		if guest.Name == "" || guest.DNI == "" {
			return nil, errors.New("guest details must include name and DNI")
		}
	}

	// 2. Business Rule Validation (Facility Status & Conflicts)
	facility, err := uc.validateBookingRules(ctx, clubID, dto.FacilityID, facilityID, dto.StartTime, dto.EndTime)
	if err != nil {
		return nil, err
	}

	// 2.1. Validate User Medical Certificate
	if err := uc.validateUserHealth(ctx, clubID, userID.String()); err != nil {
		return nil, err
	}

	// 2.2 Calculate Price (Using already fetched facility)
	dtoDuration := dto.EndTime.Sub(dto.StartTime).Hours()
	basePrice := decimal.NewFromFloat(facility.HourlyRate).Mul(decimal.NewFromFloat(dtoDuration))
	guestPrice := decimal.NewFromFloat(facility.GuestFee).Mul(decimal.NewFromFloat(float64(len(dto.GuestDetails))))
	totalPrice := basePrice.Add(guestPrice)

	// 3. Entity Construction
	// Determine initial status based on whether payment is required
	initialStatus := bookingDomain.BookingStatusConfirmed
	var paymentExpiry *time.Time
	if totalPrice.GreaterThan(decimal.Zero) {
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
	if err := uc.repo.Create(ctx, booking); err != nil {
		return nil, err
	}

	// 5. Side Effects (Notifications) - Only send confirmation if no payment required
	if initialStatus == bookingDomain.BookingStatusConfirmed {
		uc.notifyAsync(userID.String(), booking.ID.String())
	}

	return booking, nil
}

// ListBookings retrieves bookings with optional filtering.
func (uc *BookingUseCases) ListBookings(ctx context.Context, clubID string, userID string) ([]bookingDomain.Booking, error) {
	filter := make(map[string]interface{})
	if userID != "" {
		if uid, err := uuid.Parse(userID); err == nil {
			filter["user_id"] = uid
		}
	}
	return uc.repo.List(ctx, clubID, filter)
}

// ListClubBookings retrieves all bookings for a club, typically for Admin dashboard.
// Supports filtering by facility and date range.
func (uc *BookingUseCases) ListClubBookings(ctx context.Context, clubID string, facilityID string, from, to *time.Time) ([]bookingDomain.Booking, error) {
	filter := make(map[string]interface{})

	if facilityID != "" {
		if fid, err := uuid.Parse(facilityID); err == nil {
			filter["facility_id"] = fid
		}
	}

	return uc.repo.ListAll(ctx, clubID, filter, from, to)
}

// CancelBooking handles cancellation with authorization check.
func (uc *BookingUseCases) CancelBooking(ctx context.Context, clubID, bookingID, requestingUserID string) error {
	bID, err := uuid.Parse(bookingID)
	if err != nil {
		return errors.New("invalid booking id")
	}

	booking, err := uc.repo.GetByID(ctx, clubID, bID)
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

	if err := uc.repo.Update(ctx, booking); err != nil {
		return err
	}

	// Refund Logic (if reference is paid)
	if uc.refundSvc != nil {
		_ = uc.refundSvc.Refund(ctx, clubID, booking.ID, "BOOKING")
	}

	// Waitlist Logic
	next, err := uc.repo.GetNextInLine(ctx, clubID, booking.FacilityID, booking.StartTime)
	if err == nil && next != nil {
		_ = uc.notifier.Send(ctx, service.Notification{
			RecipientID: next.UserID.String(),
			Type:        service.NotificationTypeEmail,
			Title:       "Slot Available!",
			Body:        "Good news! A slot has opened up for your waitlisted time: " + booking.StartTime.String(),
		})
	}

	return nil
}

// OnPaymentStatusChanged reacts to payment updates to confirm or handle failed bookings.
func (uc *BookingUseCases) OnPaymentStatusChanged(ctx context.Context, clubID string, referenceID uuid.UUID, status paymentDomain.PaymentStatus) error {
	booking, err := uc.repo.GetByID(ctx, clubID, referenceID)
	if err != nil {
		return err
	}
	if booking == nil {
		return nil // Not a booking reference
	}

	if status == paymentDomain.PaymentStatusCompleted {
		booking.Status = bookingDomain.BookingStatusConfirmed
		booking.PaymentExpiry = nil // Clear expiry
		uc.notifyAsync(booking.UserID.String(), booking.ID.String())

		// Gamification: Award XP for completed booking
		go uc.awardBookingXP(clubID, booking.UserID.String())
	}

	booking.UpdatedAt = time.Now()
	return uc.repo.Update(ctx, booking)
}

// awardBookingXP grants XP to a user for completing a booking.
// Runs asynchronously to not block the payment flow.
func (uc *BookingUseCases) awardBookingXP(clubID, userID string) {
	user, err := uc.userRepo.GetByID(context.Background(), clubID, userID)
	if err != nil || user == nil {
		return
	}

	if user.Stats == nil {
		// Initialize stats if missing
		now := time.Now()
		user.Stats = &userDomain.UserStats{
			UserID:        userID,
			MatchesPlayed: 0,
			MatchesWon:    0,
			RankingPoints: 0,
			Level:         1,
			Experience:    0,
			CurrentStreak: 0,
			LongestStreak: 0,
			TotalXP:       0,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
	}

	// Calculate XP with streak multiplier
	baseXP := userDomain.GetXPForAction(userDomain.XPBookingComplete)
	finalXP := userDomain.CalculateXPWithStreak(baseXP, user.Stats.CurrentStreak)

	// Check if first booking of the month for bonus
	if uc.isFirstBookingOfMonth(clubID, userID) {
		bonusXP := userDomain.GetXPForAction(userDomain.XPBookingFirstOfMonth)
		finalXP += userDomain.CalculateXPWithStreak(bonusXP, user.Stats.CurrentStreak)
	}

	user.Stats.Experience += finalXP
	user.Stats.TotalXP += finalXP

	// Update streak (booking counts as activity)
	uc.updateUserStreak(user.Stats)

	// Check for level up (exponential formula: 500 * 1.15^Level)
	for {
		requiredXP := int(500 * pow(1.15, float64(user.Stats.Level)))
		if user.Stats.Experience >= requiredXP {
			user.Stats.Level++
			user.Stats.Experience -= requiredXP
		} else {
			break
		}
	}

	user.Stats.UpdatedAt = time.Now()
	_ = uc.userRepo.Update(context.Background(), user)
}

// isFirstBookingOfMonth checks if this is the user's first booking this month.
func (uc *BookingUseCases) isFirstBookingOfMonth(clubID, userID string) bool {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return false
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	filter := map[string]interface{}{
		"user_id": uid,
	}
	bookings, err := uc.repo.ListAll(context.Background(), clubID, filter, &startOfMonth, &now)
	if err != nil {
		return false
	}

	// If this is the only booking (the current one), it's the first
	return len(bookings) <= 1
}

// updateUserStreak updates streak based on activity today.
func (uc *BookingUseCases) updateUserStreak(stats *userDomain.UserStats) {
	today := time.Now().Truncate(24 * time.Hour)

	if stats.LastActivityDate == nil {
		stats.CurrentStreak = 1
		stats.LongestStreak = 1
		stats.LastActivityDate = &today
		return
	}

	lastActivity := stats.LastActivityDate.Truncate(24 * time.Hour)
	daysSince := int(today.Sub(lastActivity).Hours() / 24)

	switch daysSince {
	case 0:
		return // Same day
	case 1:
		stats.CurrentStreak++
		if stats.CurrentStreak > stats.LongestStreak {
			stats.LongestStreak = stats.CurrentStreak
		}
	default:
		stats.CurrentStreak = 1
	}

	stats.LastActivityDate = &today
}

// pow is a simple power function to avoid importing math in this file.
func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	// Handle fractional exponent approximation
	if exp != float64(int(exp)) {
		// Use simple approximation for 1.15^n
		result *= (1 + (exp-float64(int(exp)))*(base-1))
	}
	return result
}

// GetAvailability calculates available slots based on business hours and existing bookings.
func (uc *BookingUseCases) GetAvailability(ctx context.Context, clubID, facilityID string, date time.Time) ([]map[string]interface{}, error) {
	facUUID, err := uuid.Parse(facilityID)
	if err != nil {
		return nil, errors.New("invalid facility id")
	}

	// 1. Fetch dependencies (Facility & Existing Bookings)
	// Suggestion: Use errgroup here for parallel fetching in High Performance scenarios.
	facility, err := uc.facilityRepo.GetByID(ctx, clubID, facilityID)
	if err != nil {
		return nil, err
	}
	if facility == nil {
		return nil, errors.New("facility not found")
	}

	bookings, err := uc.repo.ListByFacilityAndDate(ctx, clubID, facUUID, date)
	if err != nil {
		return nil, err
	}

	// 1.5. OPTIMIZATION: Fetch Maintenance Tasks Upfront (Avoid N+1)
	allMaintenance, err := uc.facilityRepo.ListMaintenanceByFacility(ctx, facilityID)
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

	// Parse Opening Times
	startH, startM := parseTimeStr(facility.OpeningTime, 8, 0)
	endH, endM := parseTimeStr(facility.ClosingTime, 23, 0)

	var slots []map[string]interface{}

	// Iterate by hour from Opening Time until Closing Time
	// We construct daily dates based on the passed 'date'

	// Start Time for the day
	// Start Time for the day
	loopStart := time.Date(date.Year(), date.Month(), date.Day(), startH, startM, 0, 0, date.Location())

	// End Time for the day
	loopEnd := time.Date(date.Year(), date.Month(), date.Day(), endH, endM, 0, 0, date.Location())

	// Loop in 1-hour increments
	for t := loopStart; t.Before(loopEnd); t = t.Add(1 * time.Hour) {
		slotEnd := t.Add(1 * time.Hour)

		// If slotEnd exceeds ClosingTime, should we include partial slot?
		// Usually booking systems enforce full slots. We'll skip if it exceeds.
		if slotEnd.After(loopEnd) {
			break
		}

		// Pass pre-fetched maintenance to helper
		status := uc.determineSlotStatusInMemory(t, slotEnd, bookings, dailyMaintenance)

		slots = append(slots, map[string]interface{}{
			"start_time": t.Format("15:04"),
			"end_time":   slotEnd.Format("15:04"),
			"available":  status == "available",
			"status":     status,
		})
	}

	return slots, nil
}

// CreateRecurringRule creates a pattern for future bookings.
func (uc *BookingUseCases) CreateRecurringRule(ctx context.Context, clubID string, dto CreateRecurringRuleDTO) (*bookingDomain.RecurringRule, error) {
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

	if endD.Before(startD) {
		return nil, errors.New("end date must be after start date")
	}

	if !dto.EndTime.After(dto.StartTime) {
		return nil, errors.New("end time must be after start time")
	}

	rule := &bookingDomain.RecurringRule{
		ID:         uuid.New(),
		FacilityID: facID,
		ClubID:     clubID,
		Type:       dto.Type,
		Frequency:  dto.Frequency, // Map new field
		DayOfWeek:  dto.DayOfWeek,
		StartTime:  dto.StartTime,
		EndTime:    dto.EndTime,
		StartDate:  startD,
		EndDate:    endD,
	}

	if err := uc.recurringRepo.Create(ctx, rule); err != nil {
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
			conflict, _ := uc.repo.HasTimeConflict(ctx, clubID, bk.FacilityID, bk.StartTime, bk.EndTime)
			if !conflict {
				if err := uc.repo.Create(ctx, &bk); err == nil {
					generatedCount++
				}
			}
		}
	}
	return nil
}

// ListRecurringRules retrieves all active recurring rules for the club.
func (uc *BookingUseCases) ListRecurringRules(ctx context.Context, clubID string) ([]bookingDomain.RecurringRule, error) {
	return uc.recurringRepo.GetAllActive(ctx, clubID)
}

// ExpirePayments finds pending bookings with expired payment timers and marks them as EXPIRED.
// This should be called by a background cron job.
func (uc *BookingUseCases) ExpirePayments(ctx context.Context) error {
	expiredBookings, err := uc.repo.ListExpired(ctx)
	if err != nil {
		return err
	}

	for _, b := range expiredBookings {
		b.Status = bookingDomain.BookingStatusExpired
		b.UpdatedAt = time.Now()
		// We could assume partial failure is acceptable here, or log errors.
		// For now, we try to update all and return last error if any.
		if updateErr := uc.repo.Update(ctx, &b); updateErr != nil {
			err = updateErr
		}
	}
	return err
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

func (uc *BookingUseCases) validateBookingRules(ctx context.Context, clubID, facilityIDStr string, facilityID uuid.UUID, start, end time.Time) (*facilityDomain.Facility, error) {
	// 1. Check Facility Existence & Status
	facility, err := uc.facilityRepo.GetByID(ctx, clubID, facilityIDStr)
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
	conflict, err := uc.repo.HasTimeConflict(ctx, clubID, facilityID, start, end)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, errors.New("booking time conflict: facility is already booked for this requested time")
	}

	// 3. Check Maintenance Schedules
	maintConflict, err := uc.facilityRepo.HasConflict(ctx, clubID, facilityIDStr, start, end)
	if err != nil {
		return nil, err
	}
	if maintConflict {
		return nil, errors.New("booking time conflict: facility is scheduled for maintenance during this time")
	}

	return facility, nil
}

func (uc *BookingUseCases) validateUserHealth(ctx context.Context, clubID, userID string) error {
	user, err := uc.userRepo.GetByID(ctx, clubID, userID)
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
			Title:       "Booking Confirmed",
			Body:        "Booking Confirmed: " + bookingID,
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

func (uc *BookingUseCases) JoinWaitlist(ctx context.Context, clubID string, dto JoinWaitlistDTO) (*bookingDomain.Waitlist, error) {
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

	if err := uc.repo.AddToWaitlist(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func parseTimeStr(timeStr string, defaultH, defaultM int) (int, int) {
	if timeStr == "" {
		return defaultH, defaultM
	}
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return defaultH, defaultM
	}
	return t.Hour(), t.Minute()
}
