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
	"github.com/wesleymassine/swordhealth/user-management/domain"
	"github.com/wesleymassine/swordhealth/user-management/usecase"
)

// Mock use case
type MockUserUsecase struct {
	mock.Mock
	userUsecase usecase.UserUsecase
}

func (m *MockUserUsecase) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) GetUser(ctx context.Context, id int) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Test Suite
func TestCreateUser(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)
	// cases := usecase.NewUserUsecase(mockUsecase.userUsecase.UserRepo)
	handler := api.NewUserHandler(&mockUsecase.userUsecase)
	handler.RegisterRoutes(app)
	user := domain.User{Username: "testuser", Email: "test@example.com"}
	mockUsecase.On("CreateUser", mock.Anything, &user).Return(nil)

	jsonUser, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

// func TestGetUser(t *testing.T) {
//     app := fiber.New()
//     mockUsecase := new(MockUserUsecase)
//     handler := api.NewUserHandler(mockUsecase)
//     handler.RegisterRoutes(app)

//     user := domain.User{ID: 1, Username: "testuser"}
//     mockUsecase.On("GetUser", mock.Anything, 1).Return(&user, nil)

//     req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
//     resp, _ := app.Test(req)

//     assert.Equal(t, http.StatusOK, resp.StatusCode)
// }

// func TestUpdateUser(t *testing.T) {
//     app := fiber.New()
//     mockUsecase := new(MockUserUsecase)
//     handler := api.NewUserHandler(mockUsecase)
//     handler.RegisterRoutes(app)

//     user := domain.User{ID: 1, Username: "updateduser"}
//     mockUsecase.On("UpdateUser", mock.Anything, &user).Return(nil)

//     jsonUser, _ := json.Marshal(user)
//     req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(jsonUser))
//     req.Header.Set("Content-Type", "application/json")
//     resp, _ := app.Test(req)

//     assert.Equal(t, http.StatusOK, resp.StatusCode)
// }

// func TestDeleteUser(t *testing.T) {
//     app := fiber.New()
//     mockUsecase := new(MockUserUsecase)
//     handler := api.NewUserHandler(mockUsecase)
//     handler.RegisterRoutes(app)

//     mockUsecase.On("DeleteUser", mock.Anything, 1).Return(nil)

//     req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
//     resp, _ := app.Test(req)

//     assert.Equal(t, http.StatusNoContent, resp.StatusCode)
// }
