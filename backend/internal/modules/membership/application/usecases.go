package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	"github.com/shopspring/decimal"
)

type CreateMembershipRequest struct {
	UserID           uuid.UUID           `json:"user_id" binding:"required"`
	MembershipTierID uuid.UUID           `json:"membership_tier_id" binding:"required"`
	BillingCycle     domain.BillingCycle `json:"billing_cycle"` // Optional, defaults to Monthly
	StartDate        *time.Time          `json:"start_date"`    // Optional, defaults to Now
	AutoRenew        *bool               `json:"auto_renew"`    // Optional, defaults based on tier type
}

type MembershipUseCases struct {
	repo             domain.MembershipRepository
	scholarshipRepo  domain.ScholarshipRepository
	subscriptionRepo domain.SubscriptionRepository
}

func NewMembershipUseCases(repo domain.MembershipRepository, scholarshipRepo domain.ScholarshipRepository, subscriptionRepo domain.SubscriptionRepository) *MembershipUseCases {
	return &MembershipUseCases{
		repo:             repo,
		scholarshipRepo:  scholarshipRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

// addMonthsRobust adds months to a date while handling end-of-month edge cases.
// For example: Jan 31 + 1 month = Feb 28 (or 29 in leap year), not March 3.
func addMonthsRobust(t time.Time, months int) time.Time {
	originalDay := t.Day()
	result := t.AddDate(0, months, 0)

	// If the day changed, we overflowed into the next month
	// (e.g., Jan 31 + 1 month = March 3 instead of Feb 28)
	if result.Day() != originalDay {
		// Go back to the last day of the previous month
		result = result.AddDate(0, 0, -result.Day())
	}
	return result
}

func (uc *MembershipUseCases) ListTiers(ctx context.Context, clubID string) ([]domain.MembershipTier, error) {
	return uc.repo.ListTiers(ctx, clubID)
}

func (uc *MembershipUseCases) CreateMembership(ctx context.Context, clubID string, req CreateMembershipRequest) (*domain.Membership, error) {
	// 1. Get Tier to calculate dates and validate
	tier, err := uc.repo.GetTierByID(ctx, clubID, req.MembershipTierID)
	if err != nil {
		return nil, errors.New("invalid membership tier")
	}

	// 2. Calculate dates
	now := time.Now()
	startDate := now
	if req.StartDate != nil {
		startDate = *req.StartDate
	}

	var nextBilling time.Time
	var endDate *time.Time
	autoRenew := true
	billingCycle := req.BillingCycle
	if billingCycle == "" {
		billingCycle = domain.BillingCycleMonthly
	}

	// Logic for Fixed Duration (Day Pass, Season Pass) vs Recurring
	if tier.DurationDays != nil && *tier.DurationDays > 0 {
		// Fixed Duration
		calculatedEnd := startDate.AddDate(0, 0, *tier.DurationDays)
		endDate = &calculatedEnd
		autoRenew = false
		// For fixed duration, NextBillingDate is irrelevant or could be same as EndDate
		nextBilling = calculatedEnd
	} else {
		// Recurring Subscription
		// autoRenew remains true (default)
		// Calculate next billing
		switch billingCycle {
		case domain.BillingCycleMonthly:
			nextBilling = addMonthsRobust(startDate, 1)
		case domain.BillingCycleQuarterly:
			nextBilling = addMonthsRobust(startDate, 3)
		case domain.BillingCycleSemiAnnual:
			nextBilling = addMonthsRobust(startDate, 6)
		case domain.BillingCycleAnnual:
			nextBilling = addMonthsRobust(startDate, 12)
		default:
			nextBilling = addMonthsRobust(startDate, 1)
		}
	}

	// Allow user to override auto_renew if explicitly provided
	if req.AutoRenew != nil {
		autoRenew = *req.AutoRenew
	}

	membership := &domain.Membership{
		UserID:           req.UserID,
		MembershipTierID: req.MembershipTierID,
		MembershipTier:   *tier,
		ClubID:           clubID,
		Status:           domain.MembershipStatusActive, // Auto-activate for MVP
		BillingCycle:     billingCycle,
		StartDate:        startDate,
		EndDate:          endDate,
		NextBillingDate:  nextBilling,
		AutoRenew:        autoRenew,
	}

	if err := uc.repo.Create(ctx, membership); err != nil {
		return nil, err
	}

	return membership, nil
}

func (uc *MembershipUseCases) GetMembership(ctx context.Context, clubID string, id uuid.UUID) (*domain.Membership, error) {
	return uc.repo.GetByID(ctx, clubID, id)
}

func (uc *MembershipUseCases) ListUserMemberships(ctx context.Context, clubID string, userID uuid.UUID) ([]domain.Membership, error) {
	return uc.repo.GetByUserID(ctx, clubID, userID)
}

type SubscriptionUseCases struct {
	repo domain.SubscriptionRepository
}

func NewSubscriptionUseCases(repo domain.SubscriptionRepository) *SubscriptionUseCases {
	return &SubscriptionUseCases{repo: repo}
}

// SECURITY FIX (VUL-002): Added clubID for tenant isolation
func (uc *MembershipUseCases) ListUserSubscriptions(ctx context.Context, clubID string, userID uuid.UUID) ([]domain.Subscription, error) {
	return uc.subscriptionRepo.GetByUserID(ctx, clubID, userID)
}

// ListAllMemberships returns all memberships for admin view
func (uc *MembershipUseCases) ListAllMemberships(ctx context.Context, clubID string) ([]domain.Membership, error) {
	return uc.repo.ListAll(ctx, clubID)
}

// CancelMembership cancels an active membership
// Returns the cancelled membership or an error if not found/not allowed
func (uc *MembershipUseCases) CancelMembership(ctx context.Context, clubID string, membershipID uuid.UUID, requestingUserID string) (*domain.Membership, error) {
	// 1. Get the membership
	membership, err := uc.repo.GetByID(ctx, clubID, membershipID)
	if err != nil {
		return nil, errors.New("membership not found")
	}
	if membership == nil {
		return nil, errors.New("membership not found")
	}

	// 2. Authorization check: only the owner or admin can cancel
	// Note: Admin check should be done at the handler level
	if membership.UserID.String() != requestingUserID {
		return nil, errors.New("not authorized to cancel this membership")
	}

	// 3. Check if already cancelled
	if membership.Status == domain.MembershipStatusCancelled {
		return nil, errors.New("membership is already cancelled")
	}

	// 4. Update status
	membership.Status = domain.MembershipStatusCancelled
	now := time.Now()
	membership.EndDate = &now

	// 5. Persist the update
	if err := uc.repo.Update(ctx, membership); err != nil {
		return nil, err
	}

	return membership, nil
}

// ProcessMonthlyBilling runs the billing cycle for all active memberships
func (uc *MembershipUseCases) ProcessMonthlyBilling(ctx context.Context, clubID string) (int, error) {
	now := time.Now()
	billable, err := uc.repo.ListBillable(ctx, clubID, now)
	if err != nil {
		return 0, err
	}

	if len(billable) == 0 {
		return 0, nil
	}

	// 1. Batch Fetch Scholarships
	var userIDs []string
	for _, m := range billable {
		userIDs = append(userIDs, m.UserID.String())
	}

	scholarships, err := uc.scholarshipRepo.ListActiveByUserIDs(ctx, userIDs)
	if err != nil {
		return 0, err // Fail entire batch? Or log and proceed? For consistency, fail.
	}

	// 2. Calculate Updates in Memory
	updates := make(map[uuid.UUID]struct {
		Balance     decimal.Decimal
		NextBilling time.Time
	})
	processedCount := 0

	for _, m := range billable {
		// CRITICAL FIX: Skip memberships that don't auto-renew
		if !m.AutoRenew {
			// Fixed-duration membership reached its end. Mark as expired.
			// Note: This is a simplification. In production, use a separate job.
			continue
		}

		// Calculate fee with potential scholarship
		fee := m.MembershipTier.MonthlyFee

		// Memory Lookup
		if s, ok := scholarships[m.UserID.String()]; ok {
			fee = s.ApplyDiscount(fee)
		}

		// Calculate new balance
		newBalance := m.OutstandingBalance.Add(fee)

		// Calculate next billing date (next month) - using robust function for end-of-month handling
		nextBilling := addMonthsRobust(m.NextBillingDate, 1)

		updates[m.ID] = struct {
			Balance     decimal.Decimal
			NextBilling time.Time
		}{
			Balance:     newBalance,
			NextBilling: nextBilling,
		}
		processedCount++
	}

	// 3. Batch Update DB
	if err := uc.repo.UpdateBalancesBatch(ctx, updates); err != nil {
		return 0, err
	}

	return processedCount, nil
}

type AssignScholarshipRequest struct {
	UserID     string          `json:"user_id"`
	Percentage decimal.Decimal `json:"percentage"`
	Reason     string          `json:"reason"`
	ValidUntil *time.Time      `json:"valid_until"`
}

func (uc *MembershipUseCases) AssignScholarship(ctx context.Context, clubID string, req AssignScholarshipRequest, grantorID string) (*domain.Scholarship, error) {
	// Create Scholarship
	scholarship := &domain.Scholarship{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		Percentage: req.Percentage,
		Reason:     req.Reason,
		GrantorID:  grantorID,
		ValidUntil: req.ValidUntil,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.scholarshipRepo.Create(ctx, scholarship); err != nil {
		return nil, err
	}
	return scholarship, nil
}
