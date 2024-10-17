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

func (r *MySQLTaskRepository) Save(ctx context.Context, task domain.Task) (int64, error) {
	query := "INSERT INTO tasks (title, description, status, assigned_to) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(query, task.Title, task.Description, task.Status, task.AssignedTo)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *MySQLTaskRepository) List(ctx context.Context) ([]domain.Task, error) {
	query := "SELECT id, title, description, status, assigned_to, performed_by, performed_at, created_at FROM tasks"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []domain.Task{}
	for rows.Next() {
		var task domain.Task
		var performedBy sql.NullInt64     // Use sql.NullInt64 to handle NULL values in performed_by
		var createdAt, performedAt []byte // Use sql.NullTime to handle NULL values in performed_at

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.AssignedTo, &performedBy, &performedAt, &createdAt)
		if err != nil {
			return nil, err
		}

		// If performed_by is valid, assign the value to the task
		if performedBy.Valid {
			task.PerformedBy = performedBy.Int64
		} else {
			task.PerformedBy = 0 // or use another default value or pointer if preferred
		}

		if len(performedAt) > 0 {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", string(performedAt))
			if err != nil {
				return nil, fmt.Errorf("error parsing performed_at: %v", err)
			}
			task.PerformedAt = parsedTime
		} else {
			task.PerformedAt = time.Time{} // or handle as required (e.g., using a nil pointer)
		}

		if len(createdAt) > 0 {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
			if err != nil {
				return nil, fmt.Errorf("error parsing performed_at: %v", err)
			}
			task.CreatedAt = parsedTime
		} else {
			task.CreatedAt = time.Time{} // or handle as required (e.g., using a nil pointer)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *MySQLTaskRepository) ListForUser(ctx context.Context, userID int64) ([]domain.Task, error) {
	query := `SELECT id, title, description, assigned_to, status, performed_at, created_at FROM tasks WHERE performed_by = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		var createdAt, performedAt []byte

		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.AssignedTo, &task.Status, &performedAt, &createdAt); err != nil {
			return nil, err
		}

		// Convert performedAt from []byte to time.Time
		if len(performedAt) > 0 {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", string(performedAt)) // Adjust the format based on your MySQL `DATETIME` or `TIMESTAMP` format
			if err != nil {
				return nil, fmt.Errorf("error parsing performed_at: %v", err)
			}
			task.PerformedAt = parsedTime
		} else {
			task.PerformedAt = time.Time{} // Zero value if NULL
		}

		if len(createdAt) > 0 {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
			if err != nil {
				return nil, fmt.Errorf("error parsing performed_at: %v", err)
			}
			task.CreatedAt = parsedTime
		} else {
			task.CreatedAt = time.Time{} // Zero value if NULL
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *MySQLTaskRepository) UpdateTaskStatus(ctx context.Context, taskID int64, performedBy int64, performedAt time.Time, status string) error {
	query := `
		UPDATE tasks 
		SET performed_by = ?, performed_at = ?, status = ? 
		WHERE id = ?
	`

	// Execute the update query
	result, err := r.db.ExecContext(ctx, query, performedBy, performedAt, status, taskID)
	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no task found with id %d", taskID)
	}

	return nil
}

func (r *MySQLTaskRepository) GetTaskByID(ctx context.Context, taskID int64) (domain.Task, error) {
	var task domain.Task
	var performedBy sql.NullInt64 // Handle NULL values for performed_by
	var performedAt []byte        // Store performed_at as []byte to handle byte slices
	var createdAt []byte          // Store performed_at as []byte to handle byte slices

	query := `
		SELECT id, title, description, assigned_to, status, performed_by, performed_at, created_at
		FROM tasks 
		WHERE id = ?
	`

	// Query the database with the provided taskID
	err := r.db.QueryRowContext(ctx, query, taskID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.AssignedTo,
		&task.Status,
		&performedBy,
		&performedAt,
		&createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no task is found with the provided ID, return an appropriate error
			return task, fmt.Errorf("no task found with id %d", taskID)
		}
		return task, err
	}

	// Convert sql.NullInt64 to a regular int64, handling NULL values
	if performedBy.Valid {
		task.PerformedBy = performedBy.Int64
	} else {
		task.PerformedBy = 0 // Or another default value if needed
	}

	// Convert performedAt from []byte to time.Time
	if len(performedAt) > 0 {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(performedAt)) // Adjust the format based on your MySQL `DATETIME` or `TIMESTAMP` format
		if err != nil {
			return task, fmt.Errorf("error parsing performed_at: %v", err)
		}
		task.PerformedAt = parsedTime
	} else {
		task.PerformedAt = time.Time{} // Zero value if NULL
	}

	if len(createdAt) > 0 {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return task, fmt.Errorf("error parsing performed_at: %v", err)
		}
		task.CreatedAt = parsedTime
	} else {
		task.CreatedAt = time.Time{} // Zero value if NULL
	}

	return task, nil
}

func (r *MySQLTaskRepository) UserExists(ctx context.Context, userID int64) error {
	// SQL query to check if the user exists by userID
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`

	var exists bool

	// Query the database to check if the user exists
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	// Check if the user does not exist
	if !exists {
		return fmt.Errorf("user with ID %d does not exist", userID)
	}

	// If the user exists, return nil (no error)
	return nil
}

func (r *MySQLTaskRepository) GetUserByAssignedTask(ctx context.Context, taskID int64) (*domain.User, error) {
	query := `
        SELECT u.id,u.name, u.email, u.role 
        FROM users u
        LEFT JOIN tasks t ON u.id = t.assigned_to
        WHERE t.id = ?;
    `
	var user domain.User
	err := r.db.QueryRow(query, taskID).Scan(&user.ID, &user.Username, &user.Email, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
