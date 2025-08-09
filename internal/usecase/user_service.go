package usecase

import (
	"errors"
	"strings"

	"cleanarch/internal/domain"
)

// UserService implements application-specific use cases around the User aggregate.
type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(name, email string) (*domain.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}
	return s.repo.Create(&domain.User{Name: name, Email: email})
}

func (s *UserService) GetUser(id int64) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) ListUsers() ([]*domain.User, error) {
	return s.repo.List()
}

func (s *UserService) UpdateUser(id int64, name, email string) (*domain.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}
	return s.repo.Update(&domain.User{ID: id, Name: name, Email: email})
}

func (s *UserService) DeleteUser(id int64) error {
	return s.repo.Delete(id)
}
