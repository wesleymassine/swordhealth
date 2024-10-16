package domain

import "context"

type TaskRepository interface {
    Save(task Task) error
    List() ([]Task, error)
    ListForUser(userID int64) ([]Task, error)
    UpdateStatus(taskID int64, status string) error
    GetTaskByID(taskID int64) (Task, error)
    GetUserName(ctx context.Context, id int64) (*string, error)
}