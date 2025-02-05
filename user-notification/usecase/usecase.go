package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wesleymassine/swordhealth/user-notification/domain"
)

type NotificationService struct {
	rabbitMQConsumer  domain.NotificationEvent
	notificationRepo  domain.NotificationRepository
	userServiceClient domain.UserServiceClient
	taskChannel       chan domain.Task
	wg                sync.WaitGroup // WaitGroup to control goroutines
}

func NewNotificationService(consumer domain.NotificationEvent, userRepository domain.NotificationRepository,
	userClient domain.UserServiceClient) *NotificationService {
	return &NotificationService{
		rabbitMQConsumer:  consumer,
		notificationRepo:  userRepository,
		userServiceClient: userClient,
		taskChannel:       make(chan domain.Task),
	}
}

func (n *NotificationService) Start(ctx context.Context) {
	n.wg.Add(2) // Increment WaitGroup counter for the consumer goroutine

	// Start consuming from task_status_create_queue
	go n.rabbitMQConsumer.StartConsuming(ctx, "task_status_create_queue", n.taskChannel, &n.wg)

	// Start consuming from task_status_update_queue
	go n.rabbitMQConsumer.StartConsuming(ctx, "task_status_update_queue", n.taskChannel, &n.wg)

	n.wg.Add(1) // Start another goroutine to process tasks
	go n.handleTasks(ctx)
}

// handleTasks listens for tasks on the channel and processes them with Notify
func (n *NotificationService) handleTasks(ctx context.Context) {
	defer n.wg.Done() // Decrement the WaitGroup counter when done

	for {
		select {
		case <-ctx.Done(): // Graceful shutdown signal from context
			log.Println("Shutting down task handler...")
			return

		case task, ok := <-n.taskChannel: // Read from the task channel
			if !ok {
				log.Println("Task channel closed, stopping task handler...")
				return
			}

			n.Notify(ctx, task)
		}
	}
}

func (n *NotificationService) Notify(ctx context.Context, task domain.Task) {
	var msg string

	switch task.Event {
	case "task.status.create":
		taskCreatedAt := task.CreatedAt.Format("January 02, 2006 15:04:05")
		msg = fmt.Sprintf("Task %s created on date %v success\n", task.Title, taskCreatedAt)

	case "task.status.update":
		taskPerformedAt := task.PerformedAt.Format("January 02, 2006 15:04:05")
		msg = fmt.Sprintf("The tech %v performed the task %s on date %v\n", task.PerformedBy, task.Title, taskPerformedAt)
	}

	notification := domain.Notification{
		TaskID:             task.ID,
		NotificationBody:   msg,
		NotificationStatus: string(domain.StatusPending),
		SentAt:             time.Now().UTC().Local(),
		ByEmail:            true, // TODO: flag via environment variable
		ByPush:             false,
	}

	manager, err := n.userServiceClient.GetUserByTaskID(ctx, task.ID)

	if err != nil {
		log.Printf("Error fetching manager email: %v", err)
		return
	}

	if err := n.notificationRepo.UpsertNotification(ctx, notification); err != nil {
		log.Printf("Error notification status pending: %v", err)
		return
	}

	if notification.ByEmail {
		if err = n.sendEmailNotification(manager.Email, msg); err != nil {
			log.Printf("Error notify by email: %v", err)
			return
		}
	}

	if notification.ByPush {
		if err = n.notifyPush(manager.Username, msg); err != nil {
			log.Printf("Error notify by push: %v", err)
			return
		}
	}

	notification.NotificationStatus = string(domain.StatusSent)

	if err := n.notificationRepo.UpsertNotification(ctx, notification); err != nil {
		log.Printf("Error notification status sent: %v", err)
		return
	}
}

func (n *NotificationService) sendEmailNotification(userEmail, msg string) error {
	// TODO: Email sending logic here
	log.Printf("sending email to %s...", userEmail)
	log.Printf("notification: %s", msg)
	log.Printf("status %s...", domain.StatusSent)

	return nil
}

func (n *NotificationService) notifyPush(userName, msg string) error {
	// TODO: Email sending logic here
	log.Printf("sending push to manager %s...", userName)
	log.Printf("notification: %s", msg)
	log.Printf("status %s...", domain.StatusSent)

	return nil
}

// Shutdown waits for all goroutines to finish before exiting
func (n *NotificationService) Shutdown() {
	log.Println("Waiting for all goroutines to complete...")
	n.wg.Wait()
	log.Println("Notification service gracefully shut down.")
}

func (n *NotificationService) ListLatestNotifications(limit int) (domain.Notifications, error) {
	return n.notificationRepo.ListLatestNotifications(limit)
}
