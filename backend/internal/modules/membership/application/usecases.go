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
	BillingCycle     domain.BillingCycle `json:"billing_cycle" binding:"required"`
}

type MembershipUseCases struct {
	repo            domain.MembershipRepository
	scholarshipRepo domain.ScholarshipRepository
}

func NewMembershipUseCases(repo domain.MembershipRepository, scholarshipRepo domain.ScholarshipRepository) *MembershipUseCases {
	return &MembershipUseCases{
		repo:            repo,
		scholarshipRepo: scholarshipRepo,
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
	var nextBilling time.Time

	switch req.BillingCycle {
	case domain.BillingCycleMonthly:
		nextBilling = addMonthsRobust(now, 1)
	case domain.BillingCycleQuarterly:
		nextBilling = addMonthsRobust(now, 3)
	case domain.BillingCycleSemiAnnual:
		nextBilling = addMonthsRobust(now, 6)
	case domain.BillingCycleAnnual:
		nextBilling = addMonthsRobust(now, 12)
	default:
		nextBilling = addMonthsRobust(now, 1)
	}

	membership := &domain.Membership{
		UserID:           req.UserID,
		MembershipTierID: req.MembershipTierID,
		MembershipTier:   *tier,
		ClubID:           clubID,
		Status:           domain.MembershipStatusActive, // Auto-activate for MVP
		BillingCycle:     req.BillingCycle,
		StartDate:        startDate,
		NextBillingDate:  nextBilling,
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

// ListAllMemberships returns all memberships for admin view
func (uc *MembershipUseCases) ListAllMemberships(ctx context.Context, clubID string) ([]domain.Membership, error) {
	return uc.repo.ListAll(ctx, clubID)
}

// ProcessMonthlyBilling runs the billing cycle for all active memberships
func (uc *MembershipUseCases) ProcessMonthlyBilling(ctx context.Context, clubID string) (int, error) {
	now := time.Now()
	billable, err := uc.repo.ListBillable(ctx, clubID, now)
	if err != nil {
		return 0, err
	}

	processedCount := 0
	for _, m := range billable {
		// Calculate fee with potential scholarship
		fee := m.MembershipTier.MonthlyFee

		scholarship, err := uc.scholarshipRepo.GetActiveByUserID(m.UserID.String()) // Convert UUID to string if needed
		if err == nil && scholarship != nil {
			fee = scholarship.ApplyDiscount(fee)
		}

		// Calculate new balance
		newBalance := m.OutstandingBalance.Add(fee)

		// Calculate next billing date (next month) - using robust function for end-of-month handling
		nextBilling := addMonthsRobust(m.NextBillingDate, 1)

		if err := uc.repo.UpdateBalance(ctx, clubID, m.ID, newBalance, nextBilling); err == nil {
			processedCount++
		}
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

	if err := uc.scholarshipRepo.Create(scholarship); err != nil {
		return nil, err
	}
	return scholarship, nil
}
