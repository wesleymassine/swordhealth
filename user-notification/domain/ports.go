package domain

import (
	"context"
	"sync"
)

type NotificationEvent interface {
	StartConsuming(ctx context.Context, queueName string, taskChannel chan<- Task, wg *sync.WaitGroup)
}

type NotificationRepository interface {
	UpsertNotification(ctx context.Context, Notification Notification) error
	ListLatestNotifications(limit int) (Notifications, error)
}

type UserServiceClient interface {
	GetUserByTaskID(ctx context.Context, taskID int64) (*User, error)
}
