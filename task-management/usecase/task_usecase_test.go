package usecase

// import (
//     "testing"
//     "github.com/stretchr/testify/assert"
//     "github.com/swordhealth-root/internal/task-management/domain"
//     "github.com/swordhealth-root/internal/task-management/usecase"
//     "github.com/swordhealth-root/mocks"
// )

// func TestCreateTask(t *testing.T) {
//     repo := new(mocks.TaskRepository)
//     producer := new(mocks.RabbitMQProducer)
//     useCase := usecase.NewTaskUseCase(repo, producer)

//     task := domain.Task{
//         Title:       "Task 1",
//         Description: "Test Task",
//         Status:      "Pending",
//         PerformedBy: "Technician 1",
//     }

//     repo.On("Save", task).Return(nil)
//     producer.On("PublishNotification", mock.Anything).Return(nil)

//     err := useCase.CreateTask(task)
//     assert.Nil(t, err)
// }
