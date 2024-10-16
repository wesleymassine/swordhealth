package domain

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`    // Don't expose in JSON responses
	Role     string `json:"role"` // super_admin, manager, technical
}

// TODO: Refactoring response
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`    // Don't expose in JSON responses
	Role     string `json:"role"` // super_admin, manager, technical
}
