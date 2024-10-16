package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/wesleymassine/swordhealth/user-management/domain"
	"github.com/wesleymassine/swordhealth/user-management/usecase"
)

// Mock repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUser(ctx context.Context, id int) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Test Suite
type UserUsecaseTestSuite struct {
	suite.Suite
	MockRepo *MockUserRepository
	Usecase  *usecase.UserUsecase
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.MockRepo = new(MockUserRepository)
	suite.Usecase = usecase.NewUserUsecase(suite.MockRepo)
}

func (suite *UserUsecaseTestSuite) TestCreateUser() {
	user := &domain.User{Username: "testuser", Email: "test@example.com"}
	suite.MockRepo.On("Create", mock.Anything, user).Return(nil)

	err := suite.Usecase.CreateUser(context.Background(), user)
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestGetUser() {
	user := &domain.User{ID: 1, Username: "testuser"}
	suite.MockRepo.On("GetUser", mock.Anything, 1).Return(user, nil)

	result, err := suite.Usecase.GetUser(context.Background(), 1)
	suite.NoError(err)
	suite.Equal(user, result)
}

func (suite *UserUsecaseTestSuite) TestUpdateUser() {
	user := &domain.User{ID: 1, Username: "updateduser"}
	suite.MockRepo.On("UpdateUser", mock.Anything, user).Return(nil)

	err := suite.Usecase.UpdateUser(context.Background(), user)
	suite.NoError(err)
}

func (suite *UserUsecaseTestSuite) TestDeleteUser() {
	suite.MockRepo.On("DeleteUser", mock.Anything, 1).Return(nil)

	err := suite.Usecase.DeleteUser(context.Background(), 1)
	suite.NoError(err)
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}
