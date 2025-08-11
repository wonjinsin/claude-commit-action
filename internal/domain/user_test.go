package domain

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	t.Run("Create user with valid data", func(t *testing.T) {
		now := time.Now().UTC()
		user := &User{
			ID:        1,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: now,
			UpdatedAt: now,
		}

		if user.ID != 1 {
			t.Errorf("expected ID to be 1, got %d", user.ID)
		}
		if user.Name != "John Doe" {
			t.Errorf("expected Name to be 'John Doe', got %s", user.Name)
		}
		if user.Email != "john@example.com" {
			t.Errorf("expected Email to be 'john@example.com', got %s", user.Email)
		}
		if user.CreatedAt != now {
			t.Errorf("expected CreatedAt to be %v, got %v", now, user.CreatedAt)
		}
		if user.UpdatedAt != now {
			t.Errorf("expected UpdatedAt to be %v, got %v", now, user.UpdatedAt)
		}
	})

	t.Run("User struct fields are properly accessible", func(t *testing.T) {
		user := &User{}

		// Test that we can set and get all fields
		user.ID = 123
		user.Name = "Test User"
		user.Email = "test@example.com"
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		if user.ID == 0 {
			t.Error("ID field should be accessible")
		}
		if user.Name == "" {
			t.Error("Name field should be accessible")
		}
		if user.Email == "" {
			t.Error("Email field should be accessible")
		}
		if user.CreatedAt.IsZero() {
			t.Error("CreatedAt field should be accessible")
		}
		if user.UpdatedAt.IsZero() {
			t.Error("UpdatedAt field should be accessible")
		}
	})
}
