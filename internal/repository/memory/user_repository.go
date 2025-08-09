package memory

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"cleanarch/internal/domain"
)

// InMemoryUserRepository is a threadsafe in-memory implementation of UserRepository.
type InMemoryUserRepository struct {
	mu        sync.RWMutex
	autoIncID int64
	users     map[int64]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[int64]*domain.User),
	}
}

func (r *InMemoryUserRepository) Create(user *domain.User) (*domain.User, error) {
	if user == nil {
		return nil, errors.New("nil user")
	}
	id := atomic.AddInt64(&r.autoIncID, 1)
	now := time.Now().UTC()

	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *user
	copy.ID = id
	copy.CreatedAt = now
	copy.UpdatedAt = now
	r.users[id] = &copy
	return &copy, nil
}

func (r *InMemoryUserRepository) GetByID(id int64) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	copy := *u
	return &copy, nil
}

func (r *InMemoryUserRepository) List() ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.User, 0, len(r.users))
	for _, u := range r.users {
		copy := *u
		result = append(result, &copy)
	}
	return result, nil
}

func (r *InMemoryUserRepository) Update(user *domain.User) (*domain.User, error) {
	if user == nil {
		return nil, errors.New("nil user")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.users[user.ID]
	if !ok {
		return nil, errors.New("user not found")
	}
	existing.Name = user.Name
	existing.Email = user.Email
	existing.UpdatedAt = time.Now().UTC()
	copy := *existing
	return &copy, nil
}

func (r *InMemoryUserRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}
