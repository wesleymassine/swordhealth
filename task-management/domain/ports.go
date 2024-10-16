package domain

import (
	"context"
	"time"
)

type TaskRepository interface {
	Save(ctx context.Context, task Task) (int64, error)
	List(ctx context.Context) ([]Task, error)
	ListForUser(ctx context.Context, userID int64) ([]Task, error)
	UpdateTaskStatus(ctx context.Context, taskID int64, performedBy int64, performedAt time.Time, status string) error
	GetTaskByID(ctx context.Context, taskID int64) (Task, error)
	UserExists(ctx context.Context, userID int64) error
}
