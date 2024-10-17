package usecase

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/wesleymassine/swordhealth/task-management/domain"
	"github.com/wesleymassine/swordhealth/task-management/infra/notification"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{
		repo: repo,
	}
}

func (u *TaskUseCase) CreateTask(ctx context.Context, task domain.Task) (*domain.Task, error) {
	err := u.repo.UserExists(ctx, task.AssignedTo)

	if err != nil {
		return nil, err
	}

	task.CreatedAt = time.Now().Local()
	task.ID, err = u.repo.Save(ctx, task)

	if err != nil {
		log.Errorf("Failed to save task %v", err)
		return nil, err
	}

	// // Fetch task from the repository
	// taskCreated, err := u.repo.GetTaskByID(ctx, task.ID)

	// if err != nil {
	// 	log.Errorf("Failed to fetch task %v", err)
	// 	return nil, err
	// }

	// Publish the message to the RabbitMQ exchange using the task.status.create routing key
	if err = notification.PublishToTopicExchange("task.status.create", task); err != nil {
		log.Errorf("Failed to publish task creation to RabbitMQ %v", err)
		return nil, err
	}

	return &task, nil
}

func (u *TaskUseCase) UpdateTaskStatus(ctx context.Context, taskID int64, userID int64, userRole string, status string) error {
	err := u.repo.UserExists(ctx, userID)

	if err != nil {
		return err
	}

	performedAt := time.Now()

	// Call the repository method to update the task
	err = u.repo.UpdateTaskStatus(ctx, taskID, userID, performedAt, status)

	if err != nil {
		return err
	}

	// Fetch task from the repository
	task, err := u.repo.GetTaskByID(ctx, taskID)

	if err != nil {
		log.Errorf("Failed to fetch task %v", err)
		return err
	}

	// Notify manager if the user is not a technical or admin
	if userRole == "manager" {
		err = notification.PublishToTopicExchange("task.status.update", task)
		if err != nil {
			log.Errorf("Failed to publish task update to RabbitMQ %v", err)
			return err
		}
	}

	return nil
}

func (u *TaskUseCase) ListTasks(ctx context.Context, userRole string, userID int64) ([]domain.Task, error) {
	err := u.repo.UserExists(ctx, userID)

	if err != nil {
		return nil, err
	}

	if userRole == "manager" {
		return u.repo.List(ctx) // Return all tasks
	}

	return u.repo.ListForUser(ctx, userID) // Return only tasks assigned to the user
}

func (u *TaskUseCase) GetUserByAssignedTask(ctx context.Context, taskID int64) (*domain.User, error) {
	return u.repo.GetUserByAssignedTask(ctx, taskID)
}
