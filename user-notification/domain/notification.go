package domain

import "time"

type status string

var (
	StatusSent    status = "sent"
	StatusPending status = "pending"
	StatusFailed  status = "failed"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssignedTo  int64     `json:"assigned_to"`
	PerformedBy int64     `json:"performed_by"`
	PerformedAt time.Time `json:"performed_at"`
	CreatedAt   time.Time `json:"created_at"`
	Event       string    `json:"-"`
}

type Notifications struct {
	Notification []Notification `json:"itens"`
}

type Notification struct {
	TaskID             int64     `json:"task_id"`
	SentAt             time.Time `json:"sent_at"`
	NotificationBody   string    `json:"message"`
	NotificationStatus string    `json:"status"`
	ByEmail            bool      `json:"-"`
	ByPush             bool      `json:"-"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
