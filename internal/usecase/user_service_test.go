package usecase

import (
	"errors"
	"testing"
	"time"

	"cleanarch/internal/domain"
)

// MockUserRepository implements domain.UserRepository for testing
type MockUserRepository struct {
	users  map[int64]*domain.User
	nextID int64
	fail   bool // for testing error scenarios
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[int64]*domain.User),
		nextID: 1,
	}
}

func (m *MockUserRepository) SetFail(fail bool) {
	m.fail = fail
}

func (m *MockUserRepository) Create(user *domain.User) (*domain.User, error) {
	if m.fail {
		return nil, errors.New("repository error")
	}
	now := time.Now().UTC()
	created := &domain.User{
		ID:        m.nextID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.users[m.nextID] = created
	m.nextID++
	return created, nil
}

func (m *MockUserRepository) GetByID(id int64) (*domain.User, error) {
	if m.fail {
		return nil, errors.New("repository error")
	}
	user, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) List() ([]*domain.User, error) {
	if m.fail {
		return nil, errors.New("repository error")
	}
	result := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		result = append(result, user)
	}
	return result, nil
}

func (m *MockUserRepository) Update(user *domain.User) (*domain.User, error) {
	if m.fail {
		return nil, errors.New("repository error")
	}
	existing, ok := m.users[user.ID]
	if !ok {
		return nil, errors.New("user not found")
	}
	existing.Name = user.Name
	existing.Email = user.Email
	existing.UpdatedAt = time.Now().UTC()
	return existing, nil
}

func (m *MockUserRepository) Delete(id int64) error {
	if m.fail {
		return errors.New("repository error")
	}
	if _, ok := m.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

func TestUserService_CreateUser(t *testing.T) {
	t.Run("Create user with valid data", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		user, err := service.CreateUser("John Doe", "john@example.com")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user.Name != "John Doe" {
			t.Errorf("expected name 'John Doe', got %s", user.Name)
		}
		if user.Email != "john@example.com" {
			t.Errorf("expected email 'john@example.com', got %s", user.Email)
		}
		if user.ID == 0 {
			t.Error("expected ID to be set")
		}
	})

	t.Run("Create user with empty name", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		_, err := service.CreateUser("", "john@example.com")
		if err == nil {
			t.Error("expected error for empty name")
		}
		expectedMsg := "name and email are required"
		if err.Error() != expectedMsg {
			t.Errorf("expected error '%s', got '%s'", expectedMsg, err.Error())
		}
	})

	t.Run("Create user with empty email", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		_, err := service.CreateUser("John Doe", "")
		if err == nil {
			t.Error("expected error for empty email")
		}
	})

	t.Run("Create user with whitespace-only name and email", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		_, err := service.CreateUser("   ", "   ")
		if err == nil {
			t.Error("expected error for whitespace-only name and email")
		}
	})

	t.Run("Repository error handling", func(t *testing.T) {
		repo := NewMockUserRepository()
		repo.SetFail(true)
		service := NewUserService(repo)

		_, err := service.CreateUser("John Doe", "john@example.com")
		if err == nil {
			t.Error("expected error from repository")
		}
	})
}

func TestUserService_GetUser(t *testing.T) {
	t.Run("Get existing user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		// First create a user
		created, _ := service.CreateUser("John Doe", "john@example.com")

		// Then get it
		user, err := service.GetUser(created.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user.Name != "John Doe" {
			t.Errorf("expected name 'John Doe', got %s", user.Name)
		}
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		_, err := service.GetUser(999)
		if err == nil {
			t.Error("expected error for non-existent user")
		}
	})
}

func TestUserService_ListUsers(t *testing.T) {
	t.Run("List users", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		// Create some users
		_, _ = service.CreateUser("John Doe", "john@example.com")
		_, _ = service.CreateUser("Jane Doe", "jane@example.com")

		users, err := service.ListUsers()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(users) != 2 {
			t.Errorf("expected 2 users, got %d", len(users))
		}
	})

	t.Run("List empty users", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		users, err := service.ListUsers()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(users) != 0 {
			t.Errorf("expected 0 users, got %d", len(users))
		}
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	t.Run("Update existing user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		// First create a user
		created, _ := service.CreateUser("John Doe", "john@example.com")

		// Then update it
		updated, err := service.UpdateUser(created.ID, "Jane Doe", "jane@example.com")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if updated.Name != "Jane Doe" {
			t.Errorf("expected name 'Jane Doe', got %s", updated.Name)
		}
		if updated.Email != "jane@example.com" {
			t.Errorf("expected email 'jane@example.com', got %s", updated.Email)
		}
	})

	t.Run("Update with empty name", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		_, err := service.UpdateUser(1, "", "john@example.com")
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("Update with empty email", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		_, err := service.UpdateUser(1, "John Doe", "")
		if err == nil {
			t.Error("expected error for empty email")
		}
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Run("Delete existing user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		// First create a user
		created, _ := service.CreateUser("John Doe", "john@example.com")

		// Then delete it
		err := service.DeleteUser(created.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify it's deleted
		_, err = service.GetUser(created.ID)
		if err == nil {
			t.Error("expected error for deleted user")
		}
	})

	t.Run("Delete non-existent user", func(t *testing.T) {
		repo := NewMockUserRepository()
		service := NewUserService(repo)

		err := service.DeleteUser(999)
		if err == nil {
			t.Error("expected error for non-existent user")
		}
	})
}
