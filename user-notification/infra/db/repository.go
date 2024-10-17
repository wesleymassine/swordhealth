package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/wesleymassine/swordhealth/user-notification/domain"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) domain.NotificationRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) UpsertNotification(ctx context.Context, msg domain.Notification) error {
	// Check if the task exists in the tasks table
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", msg.TaskID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if task exists: %v", err)
		return err
	}

	if !exists {
		return fmt.Errorf("Task with ID %d does not exist", msg.TaskID)
	}

	// Proceed to insert or update the notification
	query := `INSERT INTO notifications 
        (task_id, notification_body, notification_status, sent_at) 
        VALUES (?, ?, ?, ?) 
        ON DUPLICATE KEY UPDATE 
            notification_body = VALUES(notification_body), 
            notification_status = VALUES(notification_status), 
            sent_at = VALUES(sent_at)`

	_, err = r.db.ExecContext(ctx, query, msg.TaskID, msg.NotificationBody, msg.NotificationStatus, msg.SentAt)

	if err != nil {
		log.Printf("Error saving notification: %v", err)
	}

	return err
}

// func (r *MySQLRepository) GetManagerByTaskID(taskID int64) (*domain.User, error) {
// 	query := `
//         SELECT u.name, u.email, u.role 
//         FROM users u
//         LEFT JOIN tasks t ON u.id = t.assigned_to
//         WHERE u.id = ? limit 1;
//     `
// 	var user domain.User
// 	err := r.db.QueryRow(query, taskID).Scan(&user.Username, &user.Email, &user.Role)

// 	if err != nil {
// 		log.Printf("Error fetching manager email: %v", err)
// 		return nil, err
// 	}

// 	return &user, nil
// }

func (r *MySQLRepository) ListLatestNotifications(limit int) (domain.Notifications, error) {
	var list domain.Notifications

	query := `
        SELECT task_id, sent_at, notification_status, notification_body 
        FROM notifications 
        ORDER BY sent_at DESC 
        LIMIT ?
    `

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg domain.Notification
		var sentAt []byte // Use []byte to scan the sent_at field as a byte slice

		err := rows.Scan(&msg.TaskID, &sentAt, &msg.NotificationStatus, &msg.NotificationBody)
		if err != nil {
			return list, err
		}

		// Manually parse the sentAt field into a time.Time
		if len(sentAt) > 0 {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", string(sentAt)) // Assuming DATETIME format
			if err != nil {
				log.Printf("Error parsing sent_at field: %v", err)
				msg.SentAt = time.Time{} // Handle invalid date format by setting to zero value
			} else {
				msg.SentAt = parsedTime
			}
		} else {
			msg.SentAt = time.Time{} // If sent_at is NULL, use the zero value
		}

		list.Notification = append(list.Notification, msg)
	}

	return list, nil
}
