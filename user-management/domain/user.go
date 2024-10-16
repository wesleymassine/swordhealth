package domain

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`    // Don't expose in JSON responses
	Role      string    `json:"role"` // super_admin, manager, technical
	CreatedAt time.Time `json:"created_at"`
}

// TODO: Refactoring response
type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`    // Don't expose in JSON responses
	Role     string `json:"role"` // super_admin, manager, technical
}
