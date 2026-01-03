package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
)

type CreateMembershipRequest struct {
	UserID           uuid.UUID           `json:"user_id" binding:"required"`
	MembershipTierID uuid.UUID           `json:"membership_tier_id" binding:"required"`
	BillingCycle     domain.BillingCycle `json:"billing_cycle" binding:"required"`
}

type MembershipUseCases struct {
	repo domain.MembershipRepository
}

func NewMembershipUseCases(repo domain.MembershipRepository) *MembershipUseCases {
	return &MembershipUseCases{
		repo: repo,
	}
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
		nextBilling = now.AddDate(0, 1, 0)
	case domain.BillingCycleQuarterly:
		nextBilling = now.AddDate(0, 3, 0)
	case domain.BillingCycleSemiAnnual:
		nextBilling = now.AddDate(0, 6, 0)
	case domain.BillingCycleAnnual:
		nextBilling = now.AddDate(1, 0, 0)
	default:
		nextBilling = now.AddDate(0, 1, 0)
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

// ProcessMonthlyBilling runs the billing cycle for all active memberships
func (uc *MembershipUseCases) ProcessMonthlyBilling(ctx context.Context, clubID string) (int, error) {
	now := time.Now()
	billable, err := uc.repo.ListBillable(ctx, clubID, now)
	if err != nil {
		return 0, err
	}

	processedCount := 0
	for _, m := range billable {
		// Calculate new balance
		newBalance := m.OutstandingBalance.Add(m.MembershipTier.MonthlyFee)

		// Calculate next billing date (next month)
		nextBilling := m.NextBillingDate.AddDate(0, 1, 0)

		if err := uc.repo.UpdateBalance(ctx, clubID, m.ID, newBalance, nextBilling); err == nil {
			processedCount++
		}
	}

	return processedCount, nil
}
