package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Future fields: Bio, Phone, AvatarURL, etc.
}

type UserRepository interface {
	GetByID(id string) (*User, error)
	// Update updates the non-auth fields of the user
	Update(user *User) error
	Delete(id string) error
	List(limit, offset int, filters map[string]interface{}) ([]User, error)
}
