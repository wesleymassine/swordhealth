package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/wesleymassine/swordhealth/task-management/domain"
)

type MySQLTaskRepository struct {
	db *sql.DB
}

func NewMySQLTaskRepository(db *sql.DB) domain.TaskRepository {
	return &MySQLTaskRepository{db: db}
}

func (r *MySQLTaskRepository) Save(task domain.Task) error {
	query := "INSERT INTO tasks (title, description, status, assigned_to) VALUES (?, ?, ?, ?)"
	_, err := r.db.Exec(query, task.Title, task.Description, task.Status, task.AssignedTo)
	return err
}

func (r *MySQLTaskRepository) List() ([]domain.Task, error) {
	query := "SELECT id, title, description, status, assigned_to, performed_by, performed_at FROM tasks"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []domain.Task{}
	for rows.Next() {
		var task domain.Task
		var performedBy sql.NullInt64   // Use sql.NullInt64 to handle NULL values in performed_by
		var performedAt sql.NullTime    // Use sql.NullTime to handle NULL values in performed_at

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.AssignedTo, &performedBy, &performedAt)
		if err != nil {
			return nil, err
		}

		// If performed_by is valid, assign the value to the task
		if performedBy.Valid {
			task.PerformedBy = performedBy.Int64
		} else {
			task.PerformedBy = 0  // or use another default value or pointer if preferred
		}

		// If performed_at is valid, assign the time
		if performedAt.Valid {
			task.PerformedAt = performedAt.Time
		} else {
			task.PerformedAt = time.Time{}  // or handle as required (e.g., using a nil pointer)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *MySQLTaskRepository) ListForUser(userID int64) ([]domain.Task, error) {
	query := `SELECT id, title, description, assigned_to, status, created_at FROM tasks WHERE assigned_to = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		var createdAt sql.NullTime // Use sql.NullTime to handle NULL values or incorrect data types

		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.AssignedTo, &task.Status, &createdAt); err != nil {
			return nil, err
		}

		// Check if the value is valid before assigning it
		if createdAt.Valid {
			task.CreatedAt = createdAt.Time
		} else {
			task.CreatedAt = time.Time{} // Assign zero value if created_at is NULL
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *MySQLTaskRepository) UpdateStatus(taskID int64, status string) error {
	fmt.Println("UpdateStatus", taskID, status)
	query := "UPDATE tasks SET status = ? WHERE id = ?"
	_, err := r.db.Exec(query, status, taskID)
	return err
}

func (r *MySQLTaskRepository) GetTaskByID(taskID int64) (domain.Task, error) {
	return domain.Task{
		ID:          taskID,
		Title:       "Task XPTO",
		Description: "Tarefa de refatoracao",
		Status:      "",
		AssignedTo:  0,
		PerformedBy: 0,
		PerformedAt: time.Time{},
		CreatedAt:   time.Time{},
	}, nil
}

func (r *MySQLTaskRepository) GetUserName(ctx context.Context, id int64) (*string, error) {
	var userName string

	query := "SELECT name FROM users WHERE id = ? and deleted_at IS NULL"

	err := r.db.QueryRowContext(ctx, query, id).Scan(&userName)
	if err != nil {

		return nil, err
	}

	return &userName, nil
}
