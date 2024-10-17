package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wesleymassine/swordhealth/user-management/api"
	"github.com/wesleymassine/swordhealth/user-management/api/security"
	"github.com/wesleymassine/swordhealth/user-management/domain"
	"github.com/wesleymassine/swordhealth/user-management/usecase"
)

// Mock Usecase
type MockAuthUserUsecase struct {
	mock.Mock
	userUsecase usecase.UserUsecase
}

func (m *MockAuthUserUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestLogin(t *testing.T) {
	app := fiber.New()

	mockUsecase := new(MockAuthUserUsecase)
	authHandler := api.NewUserHandler(&mockUsecase.userUsecase)

	app.Post("/login", authHandler.Login)

	// Mock user and hashed password
	hashedPassword, _ := security.HashPassword("password")
	mockUser := &domain.User{Email: "test@example.com", Password: hashedPassword, ID: 1, Role: "admin"}

	mockUsecase.On("GetUserByEmail", mock.Anything, "test@example.com").Return(mockUser, nil)

	// Create login request
	loginRequest := map[string]string{
		"email":    "test@example.com",
		"password": "password",
	}
	jsonReq, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
