package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
}
