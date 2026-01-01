package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/domain"
	membershipDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

type AccessUseCases struct {
	accessRepo     domain.AccessRepository
	userRepo       userDomain.UserRepository
	membershipRepo membershipDomain.MembershipRepository
}

func NewAccessUseCases(
	accessRepo domain.AccessRepository,
	userRepo userDomain.UserRepository,
	membershipRepo membershipDomain.MembershipRepository,
) *AccessUseCases {
	return &AccessUseCases{
		accessRepo:     accessRepo,
		userRepo:       userRepo,
		membershipRepo: membershipRepo,
	}
}

type EntryRequest struct {
	UserID     string     `json:"user_id"`
	FacilityID *uuid.UUID `json:"facility_id"`
	Direction  string     `json:"direction"` // IN, OUT
}

func (uc *AccessUseCases) RequestEntry(ctx context.Context, req EntryRequest) (*domain.AccessLog, error) {
	// 1. Validate Input
	if req.UserID == "" {
		return nil, errors.New("user_id required")
	}

	// 2. Validate User Exists
	user, err := uc.userRepo.GetByID(req.UserID)
	if err != nil || user == nil {
		return uc.logAccess(ctx, req, domain.AccessStatusDenied, "User not found")
	}

	// 3. For Entry (IN), validate Membership
	if req.Direction == "IN" || req.Direction == "" {
		// Convert String ID to UUID for Membership Repo
		userUUID, err := uuid.Parse(req.UserID)
		if err != nil {
			return uc.logAccess(ctx, req, domain.AccessStatusDenied, "Invalid User ID format")
		}

		memberships, err := uc.membershipRepo.GetByUserID(ctx, userUUID)
		if err != nil {
			return uc.logAccess(ctx, req, domain.AccessStatusDenied, "Error fetching memberships")
		}

		hasActive := false
		hasDebt := false

		for _, m := range memberships {
			if m.Status == membershipDomain.MembershipStatusActive {
				hasActive = true
				if m.OutstandingBalance.IsPositive() {
					hasDebt = true
				}
			}
		}

		if !hasActive {
			return uc.logAccess(ctx, req, domain.AccessStatusDenied, "No active membership")
		}
		if hasDebt {
			return uc.logAccess(ctx, req, domain.AccessStatusDenied, "Outstanding debt")
		}
	}

	// 4. Grant Access
	return uc.logAccess(ctx, req, domain.AccessStatusGranted, "Access Granted")
}

func (uc *AccessUseCases) logAccess(ctx context.Context, req EntryRequest, status domain.AccessStatus, reason string) (*domain.AccessLog, error) {
	dir := domain.AccessDirectionIn
	if req.Direction == "OUT" {
		dir = domain.AccessDirectionOut
	}

	log := &domain.AccessLog{
		ID:         uuid.New(),
		UserID:     req.UserID,
		FacilityID: req.FacilityID,
		Direction:  dir,
		Status:     status,
		Reason:     reason,
		Timestamp:  time.Now(),
		CreatedAt:  time.Now(),
	}

	uc.accessRepo.Create(ctx, log)
	return log, nil
}
