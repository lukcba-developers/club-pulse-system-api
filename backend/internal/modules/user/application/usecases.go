package application

import (
	"errors"
	"time"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/domain"
)

type UserUseCases struct {
	repo domain.UserRepository
}

func NewUserUseCases(repo domain.UserRepository) *UserUseCases {
	return &UserUseCases{
		repo: repo,
	}
}

func (uc *UserUseCases) GetProfile(userID string) (*domain.User, error) {
	if userID == "" {
		return nil, errors.New("invalid user ID")
	}
	return uc.repo.GetByID(userID)
}

type UpdateProfileDTO struct {
	Name string `json:"name"`
}

func (uc *UserUseCases) UpdateProfile(userID string, dto UpdateProfileDTO) (*domain.User, error) {
	user, err := uc.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update fields
	if dto.Name != "" {
		user.Name = dto.Name
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	if err := uc.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCases) DeleteUser(id string, requesterID string) error {
	if id == requesterID {
		return errors.New("cannot delete yourself")
	}
	return uc.repo.Delete(id)
}

func (uc *UserUseCases) ListUsers(limit, offset int, search string) ([]domain.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	filters := make(map[string]interface{})
	if search != "" {
		filters["search"] = search
	}

	return uc.repo.List(limit, offset, filters)
}
