package usecase

import (
	"fmt"
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

func (u *TaskUseCase) CreateTask(task domain.Task) error {
	// Check if the user (performed_by) exists
	// userExists, err := u.repo.Exists(task.PerformedBy)
	// if err != nil || !userExists {
	// 	return fmt.Errorf("User with ID %d does not exist", task.PerformedBy)
	// }

	task.ID = 19

	task.PerformedAt = time.Now().Local()

	err := u.repo.Save(task)

	if err != nil {
		log.Errorf("Failed to save task %v", err)
		return err
	}

	// Publish the message to the RabbitMQ exchange using the task.status.create routing key
	if err = notification.PublishToTopicExchange("task.status.create", task); err != nil {
		log.Errorf("Failed to publish task creation to RabbitMQ %v", err)
		return err
	}

	return nil
}

func (u *TaskUseCase) UpdateTaskStatus(taskID int64, status string, userID int64, userRole string) error {
	// userExists, err := u.repo.Exists(task.PerformedBy)
	// if err != nil || !userExists {
	// 	return fmt.Errorf("User with ID %d does not exist", task.PerformedBy)
	// }

	// Fetch task from the repository
	task, err := u.repo.GetTaskByID(taskID)

	if err != nil {
		log.Errorf("Failed to fetch task %v", err)
		return err
	}

	fmt.Println("taskID", taskID)

	// Update task status
	task.Status = status
	err = u.repo.UpdateStatus(taskID, status)

	if err != nil {
		log.Errorf("Failed to update task %v", err)
		return err
	}

	// Notify manager if the user is not a manager or admin
	if userRole == "manager" {
		err = notification.PublishToTopicExchange("task.status.update", task)
		if err != nil {
			log.Errorf("Failed to publish task update to RabbitMQ %v", err)
			return err
		}
	}

	return nil
}

func (u *TaskUseCase) ListTasks(userRole string, userID int64) ([]domain.Task, error) {
	// userExists, err := u.repo.Exists(task.PerformedBy)
	// if err != nil || !userExists {
	// 	return fmt.Errorf("User with ID %d does not exist", task.PerformedBy)
	// }

	if userRole == "manager" {
		return u.repo.List() // Return all tasks
	}

	return u.repo.ListForUser(userID) // Return only tasks assigned to the user
}
