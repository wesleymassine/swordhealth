package usecase

import (
	"testing"

	"github.com/stretchr/testify/mock"
	// "github.com/wesleymassine/swordhealth/user-notification/domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetManagerByTaskID(taskID int64) (string, error) {
	args := m.Called(taskID)
	return args.String(0), args.Error(1)
}

type MockRabbitMQConsumer struct {
	mock.Mock
}

func (m *MockRabbitMQConsumer) StartConsuming(queueName string) {
	m.Called(queueName)
}

func TestNotificationService_Notify(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockRepo.On("GetManagerByTaskID", int64(1)).Return("manager@example.com", nil)

	// mockConsumer := new(MockRabbitMQConsumer)

	// notificationService := NewNotificationService(mockConsumer, mockRepo)
	// task := domain.Task{
	//     ID:         1,
	//     Title:      "Fix server",
	//     PerformedBy: "John Doe",
	// }

	// // Act
	// notificationService.Notify(task, true, false)

	// Assert
	mockRepo.AssertCalled(t, "GetManagerByTaskID", int64(1))
	// Add more assertions to check that email sending is triggered, etc.
}
