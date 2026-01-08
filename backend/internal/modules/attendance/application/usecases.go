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
	if len(list.Records) == 0 {
		return
	}

	// 1. Collect all UserIDs
	userIDs := make([]string, 0, len(list.Records))
	uuidMap := make(map[string]uuid.UUID) // map string ID to UUID for Membership check

	for _, rec := range list.Records {
		userIDs = append(userIDs, rec.UserID)
		if uid, err := uuid.Parse(rec.UserID); err == nil {
			uuidMap[rec.UserID] = uid
		}
	}

	// 2. Batch Fetch Users
	// Note: We need a map for O(1) assignment
	users, err := uc.userRepo.ListByIDs(clubID, userIDs)
	if err == nil {
		userMap := make(map[string]*userDomain.User)
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}

		for i := range list.Records {
			if u, ok := userMap[list.Records[i].UserID]; ok {
				list.Records[i].User = u
			}
		}
	}

	// 3. Batch Fetch Memberships (for Debt Check)
	// We need to pass a slice of UUIDs
	uuids := make([]uuid.UUID, 0, len(uuidMap))
	for _, uid := range uuidMap {
		uuids = append(uuids, uid)
	}

	memberships, err := uc.membershipRepo.GetByUserIDs(context.Background(), clubID, uuids)
	if err == nil {
		// Map UserID -> []Membership
		memMap := make(map[uuid.UUID][]membershipDomain.Membership)
		for _, m := range memberships {
			memMap[m.UserID] = append(memMap[m.UserID], m)
		}

		for i := range list.Records {
			userIDStr := list.Records[i].UserID
			if uid, ok := uuidMap[userIDStr]; ok {
				if mems, found := memMap[uid]; found {
					for _, m := range mems {
						if m.OutstandingBalance.GreaterThan(decimal.Zero) {
							list.Records[i].HasDebt = true
							break
						}
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
