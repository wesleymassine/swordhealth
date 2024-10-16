package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/wesleymassine/swordhealth/user-management/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := "INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, ?)"

	// Execute the insert query
	result, err := repo.DB.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.Role)
	if err != nil {
		return nil, fmt.Errorf("error inserting user: %v", err)
	}

	// Get the last inserted ID (user's ID)
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}

	// Optionally, query the database to get the full user details (if needed)
	createdUser, err := repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching created user: %v", err)
	}

	return createdUser, nil
}

func (repo *UserRepository) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	query := "SELECT id, name, email, role, created_at FROM users WHERE id = ? AND deleted_at IS NULL"
	var user domain.User
	var createdAt []byte // Store created_at as a byte slice ([]byte)

	err := repo.DB.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &createdAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, err
	}

	if len(createdAt) > 0 {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %v", err)
		}
		user.CreatedAt = parsedTime
	} else {
		user.CreatedAt = time.Time{} // Set to zero value if created_at is NULL
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
	currentUser, err := repo.GetUserByID(ctx, user.ID)

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

func (repo *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	// Delete logically, not physically
	query := "UPDATE users SET deleted_at = ? WHERE id = ?"
	_, err := repo.DB.ExecContext(ctx, query, time.Now().Local(), id)

	return err
}
