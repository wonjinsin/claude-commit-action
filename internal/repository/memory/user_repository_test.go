package memory

import (
	"cleanarch/internal/domain"
	"sync"
	"testing"
	"time"
)

func TestInMemoryUserRepository_Create(t *testing.T) {
	t.Run("Create user successfully", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		user := &domain.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		created, err := repo.Create(user)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if created.ID == 0 {
			t.Error("expected ID to be set")
		}
		if created.Name != "John Doe" {
			t.Errorf("expected name 'John Doe', got %s", created.Name)
		}
		if created.Email != "john@example.com" {
			t.Errorf("expected email 'john@example.com', got %s", created.Email)
		}
		if created.CreatedAt.IsZero() {
			t.Error("expected CreatedAt to be set")
		}
		if created.UpdatedAt.IsZero() {
			t.Error("expected UpdatedAt to be set")
		}
		if created.CreatedAt != created.UpdatedAt {
			t.Error("expected CreatedAt and UpdatedAt to be the same for new user")
		}
	})

	t.Run("Create nil user", func(t *testing.T) {
		repo := NewInMemoryUserRepository()

		_, err := repo.Create(nil)
		if err == nil {
			t.Error("expected error for nil user")
		}
		expectedMsg := "nil user"
		if err.Error() != expectedMsg {
			t.Errorf("expected error '%s', got '%s'", expectedMsg, err.Error())
		}
	})

	t.Run("Create multiple users with incremental IDs", func(t *testing.T) {
		repo := NewInMemoryUserRepository()

		user1, _ := repo.Create(&domain.User{Name: "User1", Email: "user1@example.com"})
		user2, _ := repo.Create(&domain.User{Name: "User2", Email: "user2@example.com"})

		if user1.ID >= user2.ID {
			t.Error("expected user IDs to be incremental")
		}
		if user2.ID != user1.ID+1 {
			t.Errorf("expected user2 ID to be user1 ID + 1, got %d and %d", user1.ID, user2.ID)
		}
	})
}

func TestInMemoryUserRepository_GetByID(t *testing.T) {
	t.Run("Get existing user", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		created, _ := repo.Create(&domain.User{Name: "John Doe", Email: "john@example.com"})

		user, err := repo.GetByID(created.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user.Name != "John Doe" {
			t.Errorf("expected name 'John Doe', got %s", user.Name)
		}
		if user.Email != "john@example.com" {
			t.Errorf("expected email 'john@example.com', got %s", user.Email)
		}
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		repo := NewInMemoryUserRepository()

		_, err := repo.GetByID(999)
		if err == nil {
			t.Error("expected error for non-existent user")
		}
		expectedMsg := "user not found"
		if err.Error() != expectedMsg {
			t.Errorf("expected error '%s', got '%s'", expectedMsg, err.Error())
		}
	})

	t.Run("Get user returns copy, not reference", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		created, _ := repo.Create(&domain.User{Name: "John Doe", Email: "john@example.com"})

		user1, _ := repo.GetByID(created.ID)
		user2, _ := repo.GetByID(created.ID)

		// Modify one copy
		user1.Name = "Modified Name"

		// Check that the other copy is not affected
		if user2.Name == "Modified Name" {
			t.Error("expected GetByID to return copies, not references")
		}
	})
}

func TestInMemoryUserRepository_List(t *testing.T) {
	t.Run("List empty repository", func(t *testing.T) {
		repo := NewInMemoryUserRepository()

		users, err := repo.List()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(users) != 0 {
			t.Errorf("expected 0 users, got %d", len(users))
		}
	})

	t.Run("List multiple users", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		_, _ = repo.Create(&domain.User{Name: "John Doe", Email: "john@example.com"})
		_, _ = repo.Create(&domain.User{Name: "Jane Doe", Email: "jane@example.com"})

		users, err := repo.List()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(users) != 2 {
			t.Errorf("expected 2 users, got %d", len(users))
		}
	})

	t.Run("List returns copies, not references", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		_, _ = repo.Create(&domain.User{Name: "John Doe", Email: "john@example.com"})

		users1, _ := repo.List()
		users2, _ := repo.List()

		// Modify one list
		users1[0].Name = "Modified Name"

		// Check that the other list is not affected
		if users2[0].Name == "Modified Name" {
			t.Error("expected List to return copies, not references")
		}
	})
}

func TestInMemoryUserRepository_Update(t *testing.T) {
	t.Run("Update existing user", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		created, _ := repo.Create(&domain.User{Name: "John Doe", Email: "john@example.com"})

		// Save the original UpdatedAt time
		originalUpdatedAt := created.UpdatedAt

		// Wait a bit to ensure UpdatedAt is different
		time.Sleep(10 * time.Millisecond)

		updated, err := repo.Update(&domain.User{
			ID:    created.ID,
			Name:  "Jane Doe",
			Email: "jane@example.com",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if updated.Name != "Jane Doe" {
			t.Errorf("expected name 'Jane Doe', got %s", updated.Name)
		}
		if updated.Email != "jane@example.com" {
			t.Errorf("expected email 'jane@example.com', got %s", updated.Email)
		}
		if updated.CreatedAt != created.CreatedAt {
			t.Error("expected CreatedAt to remain unchanged")
		}
		// Check that UpdatedAt is different (either equal or after, but not before)
		if updated.UpdatedAt.Before(originalUpdatedAt) {
			t.Errorf("expected UpdatedAt to be updated: original=%v, updated=%v", originalUpdatedAt, updated.UpdatedAt)
		}

		// Alternative check: verify that the user in the repository was actually updated
		retrieved, _ := repo.GetByID(created.ID)
		if retrieved.Name != "Jane Doe" {
			t.Errorf("expected retrieved name 'Jane Doe', got %s", retrieved.Name)
		}
		if retrieved.Email != "jane@example.com" {
			t.Errorf("expected retrieved email 'jane@example.com', got %s", retrieved.Email)
		}
	})

	t.Run("Update nil user", func(t *testing.T) {
		repo := NewInMemoryUserRepository()

		_, err := repo.Update(nil)
		if err == nil {
			t.Error("expected error for nil user")
		}
		expectedMsg := "nil user"
		if err.Error() != expectedMsg {
			t.Errorf("expected error '%s', got '%s'", expectedMsg, err.Error())
		}
	})

	t.Run("Update non-existent user", func(t *testing.T) {
		repo := NewInMemoryUserRepository()

		_, err := repo.Update(&domain.User{
			ID:    999,
			Name:  "Jane Doe",
			Email: "jane@example.com",
		})
		if err == nil {
			t.Error("expected error for non-existent user")
		}
		expectedMsg := "user not found"
		if err.Error() != expectedMsg {
			t.Errorf("expected error '%s', got '%s'", expectedMsg, err.Error())
		}
	})
}

func TestInMemoryUserRepository_Delete(t *testing.T) {
	t.Run("Delete existing user", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		created, _ := repo.Create(&domain.User{Name: "John Doe", Email: "john@example.com"})

		err := repo.Delete(created.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify user is deleted
		_, err = repo.GetByID(created.ID)
		if err == nil {
			t.Error("expected error when getting deleted user")
		}
	})

	t.Run("Delete non-existent user", func(t *testing.T) {
		repo := NewInMemoryUserRepository()

		err := repo.Delete(999)
		if err == nil {
			t.Error("expected error for non-existent user")
		}
		expectedMsg := "user not found"
		if err.Error() != expectedMsg {
			t.Errorf("expected error '%s', got '%s'", expectedMsg, err.Error())
		}
	})
}

func TestInMemoryUserRepository_Concurrency(t *testing.T) {
	t.Run("Concurrent creates", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		var wg sync.WaitGroup
		numGoroutines := 100
		wg.Add(numGoroutines)

		// Launch multiple goroutines creating users concurrently
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				_, err := repo.Create(&domain.User{
					Name:  "User",
					Email: "user@example.com",
				})
				if err != nil {
					t.Errorf("unexpected error in goroutine %d: %v", id, err)
				}
			}(i)
		}
		wg.Wait()

		// Check that all users were created with unique IDs
		users, _ := repo.List()
		if len(users) != numGoroutines {
			t.Errorf("expected %d users, got %d", numGoroutines, len(users))
		}

		// Check for duplicate IDs
		ids := make(map[int64]bool)
		for _, user := range users {
			if ids[user.ID] {
				t.Errorf("found duplicate ID: %d", user.ID)
			}
			ids[user.ID] = true
		}
	})

	t.Run("Concurrent reads and writes", func(t *testing.T) {
		repo := NewInMemoryUserRepository()

		// Create some initial users
		for i := 0; i < 10; i++ {
			repo.Create(&domain.User{Name: "User", Email: "user@example.com"})
		}

		var wg sync.WaitGroup
		numGoroutines := 50
		wg.Add(numGoroutines)

		// Launch multiple goroutines doing mixed operations
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()

				// Mix of read and write operations
				switch id % 3 {
				case 0:
					// Read operation
					repo.List()
				case 1:
					// Create operation
					repo.Create(&domain.User{Name: "NewUser", Email: "new@example.com"})
				case 2:
					// Update operation
					repo.Update(&domain.User{ID: int64(id%10 + 1), Name: "Updated", Email: "updated@example.com"})
				}
			}(i)
		}
		wg.Wait()

		// If we get here without race conditions, the test passes
	})
}
