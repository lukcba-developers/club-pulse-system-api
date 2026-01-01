package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Financial fields
	DateOfBirth       *time.Time             `json:"date_of_birth,omitempty"`
	SportsPreferences map[string]interface{} `json:"sports_preferences,omitempty"`
}

// CalculateCategory returns the user's category based on birth year (e.g., "2012")
func (u *User) CalculateCategory() string {
	if u.DateOfBirth == nil {
		return "Files" // Default category if unknown
	}
	return u.DateOfBirth.Format("2006")
}

type UserRepository interface {
	GetByID(id string) (*User, error)
	// Update updates the non-auth fields of the user
	Update(user *User) error
	Delete(id string) error
	List(limit, offset int, filters map[string]interface{}) ([]User, error)
}
