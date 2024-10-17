package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wesleymassine/swordhealth/task-management/domain"
	"github.com/wesleymassine/swordhealth/task-management/usecase"
)

// Mocking TaskRepository
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) UserExists(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockTaskRepository) Save(ctx context.Context, task domain.Task) (int64, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTaskRepository) UpdateTaskStatus(ctx context.Context, taskID int64, userID int64, performedAt time.Time, status string) error {
	args := m.Called(ctx, taskID, userID, performedAt, status)
	return args.Error(0)
}

func (m *MockTaskRepository) GetTaskByID(ctx context.Context, taskID int64) (domain.Task, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) List(ctx context.Context) ([]domain.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) ListForUser(ctx context.Context, userID int64) ([]domain.Task, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetUserByAssignedTask(ctx context.Context, taskID int64) (*domain.User, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).(*domain.User), args.Error(1)
}

// Mocking notification.PublishToTopicExchange
var mockPublishToTopicExchange = func(routingKey string, task domain.Task) error {
	return nil
}

func TestCreateTask_UserNotFound(t *testing.T) {
	// Initialize mocks
	mockRepo := new(MockTaskRepository)
	taskUsecase := usecase.NewTaskUseCase(mockRepo)
	// notification.PublishToTopicExchange = mockPublishToTopicExchange

	// Define inputs and expected outputs
	ctx := context.Background()
	task := domain.Task{
		AssignedTo: 1,
		Title:      "Test Task",
	}
	mockRepo.On("UserExists", ctx, int64(1)).Return(errors.New("user not found"))

	// Call the CreateTask method
	createdTask, err := taskUsecase.CreateTask(ctx, task)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, createdTask)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTaskStatus_Success(t *testing.T) {
	// Initialize mocks
	mockRepo := new(MockTaskRepository)
	taskUsecase := usecase.NewTaskUseCase(mockRepo)

	// Define inputs and expected outputs
	ctx := context.Background()
	taskID := int64(1)
	userID := int64(2)
	status := "completed"
	userRole := "manager"
	task := domain.Task{
		ID:        taskID,
		AssignedTo: userID,
		Status:    status,
	}

	mockRepo.On("UserExists", ctx, userID).Return(nil)
	mockRepo.On("UpdateTaskStatus", ctx, taskID, userID, mock.AnythingOfType("time.Time"), status).Return(nil)
	mockRepo.On("GetTaskByID", ctx, taskID).Return(task, nil)

	// Call the UpdateTaskStatus method
	err := taskUsecase.UpdateTaskStatus(ctx, taskID, userID, userRole, status)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestListTasks_Manager(t *testing.T) {
	// Initialize mocks
	mockRepo := new(MockTaskRepository)
	taskUsecase := usecase.NewTaskUseCase(mockRepo)

	// Define inputs and expected outputs
	ctx := context.Background()
	userID := int64(1)
	userRole := "manager"
	tasks := []domain.Task{
		{ID: 1, Title: "Task 1"},
		{ID: 2, Title: "Task 2"},
	}
	mockRepo.On("UserExists", ctx, userID).Return(nil)
	mockRepo.On("List", ctx).Return(tasks, nil)

	// Call the ListTasks method
	result, err := taskUsecase.ListTasks(ctx, userRole, userID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByAssignedTask_Success(t *testing.T) {
	// Initialize mocks
	mockRepo := new(MockTaskRepository)
	taskUsecase := usecase.NewTaskUseCase(mockRepo)

	// Define inputs and expected outputs
	ctx := context.Background()
	taskID := int64(1)
	user := &domain.User{
		ID:   1,
		Username: "Manager",
	}

	mockRepo.On("GetUserByAssignedTask", ctx, taskID).Return(user, nil)

	// Call the GetUserByAssignedTask method
	result, err := taskUsecase.GetUserByAssignedTask(ctx, taskID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	mockRepo.AssertExpectations(t)
}
