package domain

import "time"

// User represents the core domain entity.
// In a real system, avoid exposing persistence-specific concerns here.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository defines the persistence port for the User aggregate.
type UserRepository interface {
	Create(user *User) (*User, error)
	GetByID(id int64) (*User, error)
	List() ([]*User, error)
	Update(user *User) (*User, error)
	Delete(id int64) error
}
