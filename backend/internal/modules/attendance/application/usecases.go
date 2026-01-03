package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/domain"
	membershipDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/domain"
	userDomain "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
	"github.com/shopspring/decimal"
)

type AttendanceUseCases struct {
	repo           domain.AttendanceRepository
	userRepo       userDomain.UserRepository
	membershipRepo membershipDomain.MembershipRepository
}

func NewAttendanceUseCases(repo domain.AttendanceRepository, userRepo userDomain.UserRepository, membershipRepo membershipDomain.MembershipRepository) *AttendanceUseCases {
	return &AttendanceUseCases{
		repo:           repo,
		userRepo:       userRepo,
		membershipRepo: membershipRepo,
	}
}

// GetOrCreateList returns the attendance list for a group and date.
// If it implies "Coach View", if it doesn't exist, we might create it empty or just return 404.
// Ideally, for "View Student List", we might want to return potential students even if list doesn't exist?
// For MVP: We return the list. if not exists, we can create it or return empty structure.
// Let's implement: Get List. If not found, return a new transient list (not persisted until saved) or persist it immediately?
// Let's persist immediately for simplicity of "Starting a class".
// GetOrCreateList returns the attendance list for a group and date.
func (uc *AttendanceUseCases) GetOrCreateList(clubID, group string, date time.Time, coachID string) (*domain.AttendanceList, error) {
	list, err := uc.repo.GetListByGroupAndDate(clubID, group, date)
	if err != nil {
		return nil, err
	}
	if list != nil {
		return list, nil
	}

	// 2. Create new if not exists
	newList := &domain.AttendanceList{
		ID:        uuid.New(),
		ClubID:    clubID,
		Date:      date,
		Group:     group,
		CoachID:   coachID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 3. Auto-populate?
	// Phase 4 goal: "Auto-assignment".
	// We should fetch users of this 'category/group' and add them as "ABSENT" or "PENDING" records?
	// For MVP, lets just create the header. The frontend can fetch users separately or we populate here.
	// Populating here is better for "Digital Attendance List".

	users, err := uc.userRepo.List(clubID, 100, 0, map[string]interface{}{"category": group}) // Assuming group matches category year
	if err == nil && len(users) > 0 {
		records := make([]domain.AttendanceRecord, len(users))
		for i, u := range users {
			records[i] = domain.AttendanceRecord{
				ID:               uuid.New(),
				AttendanceListID: newList.ID,
				UserID:           u.ID,
				Status:           domain.StatusAbsent, // Default
			}
		}
		newList.Records = records
	}

	if err := uc.repo.CreateList(newList); err != nil {
		return nil, err
	}

	// Create records if any
	for _, rec := range newList.Records {
		if err := uc.repo.UpsertRecord(&rec); err != nil {
			return nil, err
		}
	}

	return newList, nil
}

func (uc *AttendanceUseCases) GetOrCreateListByTrainingGroup(clubID string, groupID uuid.UUID, groupName string, category string, date time.Time, coachID string) (*domain.AttendanceList, error) {
	list, err := uc.repo.GetListByTrainingGroupAndDate(clubID, groupID, date)
	if err != nil {
		return nil, err
	}
	if list != nil {
		uc.populateRecords(clubID, list)
		return list, nil
	}

	newList := &domain.AttendanceList{
		ID:              uuid.New(),
		ClubID:          clubID,
		Date:            date,
		Group:           groupName,
		TrainingGroupID: &groupID,
		CoachID:         coachID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// UserRepo call
	users, err := uc.userRepo.List(clubID, 100, 0, map[string]interface{}{"category": category})
	if err == nil && len(users) > 0 {
		records := make([]domain.AttendanceRecord, len(users))
		for i, u := range users {
			records[i] = domain.AttendanceRecord{
				ID:               uuid.New(),
				AttendanceListID: newList.ID,
				UserID:           u.ID,
				Status:           domain.StatusAbsent,
			}
		}
		newList.Records = records
	}

	if err := uc.repo.CreateList(newList); err != nil {
		return nil, err
	}

	for _, rec := range newList.Records {
		if err := uc.repo.UpsertRecord(&rec); err != nil {
			return nil, err
		}
	}

	uc.populateRecords(clubID, newList)
	return newList, nil
}

func (uc *AttendanceUseCases) populateRecords(clubID string, list *domain.AttendanceList) {
	for i := range list.Records {
		rec := &list.Records[i]
		u, err := uc.userRepo.GetByID(clubID, rec.UserID)
		if err == nil {
			rec.User = u
		}

		// Check Debt
		uid, err := uuid.Parse(rec.UserID)
		if err == nil {
			memberships, err := uc.membershipRepo.GetByUserID(context.Background(), clubID, uid)
			if err == nil {
				for _, m := range memberships {
					if m.OutstandingBalance.GreaterThan(decimal.Zero) {
						rec.HasDebt = true
						break
					}
				}
			}
		}
	}
}

type MarkAttendanceDTO struct {
	UserID string                  `json:"user_id"`
	Status domain.AttendanceStatus `json:"status"`
	Notes  string                  `json:"notes"`
}

func (uc *AttendanceUseCases) MarkAttendance(clubID string, listID uuid.UUID, dto MarkAttendanceDTO) error {
	list, err := uc.repo.GetListByID(clubID, listID)
	if err != nil {
		return err
	}
	if list == nil {
		return errors.New("list not found")
	}

	record := &domain.AttendanceRecord{
		ID: uuid.New(), // Might need to check if exists to keep ID? Repo Upsert handles it?
		// Logic: If user already in list, update. If not, insert.
		AttendanceListID: listID,
		UserID:           dto.UserID,
		Status:           dto.Status,
		Notes:            dto.Notes,
	}

	// For Upsert to work on ID, we need the ID.
	// The repo implementation uses ID as PK.
	// We need to find the specific record ID if it exists.
	// The repo `GetListByID` preloads records. We can search there.
	for _, r := range list.Records {
		if r.UserID == dto.UserID {
			record.ID = r.ID
			break
		}
	}
	// If record.ID is new, Upsert (GORM Save) will Insert.

	return uc.repo.UpsertRecord(record)
}
