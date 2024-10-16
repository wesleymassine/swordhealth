package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3" // Use SQLite for testing
	"github.com/stretchr/testify/assert"
	"github.com/wesleymassine/swordhealth/user-management/domain"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// Create users table
	createTable := `
    CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        email TEXT,
        password TEXT,
        role TEXT
    );`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestCreateUser(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	repo := NewUserRepository(db)

	user := &domain.User{Username: "testuser", Email: "test@example.com", Password: "password", Role: "user"}
	err = repo.Create(context.Background(), user)
	assert.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	repo := NewUserRepository(db)

	// Insert a test user
	db.Exec("INSERT INTO users (name, email, password_hash, role) VALUES ('testuser', 'test@example.com', 'password', 'user')")

	user, err := repo.GetUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUpdateUser(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	repo := NewUserRepository(db)

	// Insert a test user
	db.Exec("INSERT INTO users (name, email, password_hash, role) VALUES ('testuser', 'test@example.com', 'password', 'user')")

	user := &domain.User{ID: 1, Username: "updateduser", Email: "updated@example.com", Role: "user"}
	err = repo.UpdateUser(context.Background(), user)
	assert.NoError(t, err)

	// Check if the user was updated
	updatedUser, err := repo.GetUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updatedUser.Username)
	assert.Equal(t, "updated@example.com", updatedUser.Email)
}

func TestDeleteUser(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	repo := NewUserRepository(db)

	// Insert a test user
	db.Exec("INSERT INTO users (name, email, password_hash, role) VALUES ('testuser', 'test@example.com', 'password', 'user')")

	err = repo.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)

	// Try to fetch the deleted user
	user, err := repo.GetUser(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, user)
}
