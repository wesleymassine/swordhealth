package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/wesleymassine/swordhealth/user-management/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, ?)"
	_, err := repo.DB.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.Role)
	return err
}

func (repo *UserRepository) GetUser(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User

	query := "SELECT id, name, email, role FROM users WHERE id = ? and deleted_at IS NULL"

	err := repo.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {

		if strings.Contains(err.Error(), domain.ErrNoRowsResult.Error()) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := "SELECT id, name, email, password_hash, role FROM users WHERE email = ? and deleted_at IS NULL"

	err := repo.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	// Fetch current user data
	currentUser, err := repo.GetUser(ctx, user.ID)
	if err != nil {
		return err
	}

	// Set only the fields that are not empty, otherwise use the current data
	if user.Username == "" {
		user.Username = currentUser.Username
	}

	if user.Email == "" {
		user.Email = currentUser.Email
	}

	if user.Role == "" {
		user.Role = currentUser.Role
	}

	query := "UPDATE users SET name = ?, email = ?, role = ?, updated_at = ? WHERE id = ?"
	_, err = repo.DB.ExecContext(ctx, query, user.Username, user.Email, user.Role, time.Now().Local(), user.ID)
	return err
}

func (repo *UserRepository) DeleteUser(ctx context.Context, id int) error {
	// Delete logically, not physically
	query := "UPDATE users SET deleted_at = ? WHERE id = ?"
	_, err := repo.DB.ExecContext(ctx, query, time.Now().Local(), id)
	return err
}
