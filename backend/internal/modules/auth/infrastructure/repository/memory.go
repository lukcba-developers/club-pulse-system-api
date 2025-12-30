package repository

import (
	"sync"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/domain"
)

type InMemoryAuthRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

func NewInMemoryAuthRepository() *InMemoryAuthRepository {
	return &InMemoryAuthRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *InMemoryAuthRepository) SaveUser(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.Email] = user
	return nil
}

func (r *InMemoryAuthRepository) FindUserByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if user, exists := r.users[email]; exists {
		return user, nil
	}
	return nil, nil // Not found is not an error here, just nil
}
