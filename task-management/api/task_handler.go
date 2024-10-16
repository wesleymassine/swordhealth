package api

import (
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/wesleymassine/swordhealth/task-management/domain"
	"github.com/wesleymassine/swordhealth/task-management/usecase"
)

type TaskHandler struct {
	usecase *usecase.TaskUseCase
}

func NewTaskHandler(u *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{usecase: u}
}

// CreateTask handles POST /tasks
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var task domain.Task
	if err := c.BodyParser(&task); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// // Get user info from middleware context
	userID := int64(c.Locals("user_id").(float64))
	userRole := c.Locals("role").(string)

	claims := domain.Claims{
		UserID: userID,
		Role:   userRole,
	}

	// Optionally, you can validate user role or perform different actions based on role
	if claims.Role != "super_admin" && claims.Role != "manager" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	// if err := RoleMiddleware(claims.Role); err != nil {
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	// }

	task.AssignedTo = claims.UserID // Assign the task to the user
	if err := h.usecase.CreateTask(task); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task created successfully",
	})
}

// ListTasks handles GET /tasks
func (h *TaskHandler) ListTasks(c *fiber.Ctx) error {
	// Get user info from middleware context
	userID := int64(c.Locals("user_id").(float64))
	userRole := c.Locals("role").(string)

	claims := domain.Claims{
		UserID: userID,
		Role:   userRole,
	}

	tasks, err := h.usecase.ListTasks(claims.Role, claims.UserID)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(tasks)
}

func (h *TaskHandler) UpdateTaskStatus(c *fiber.Ctx) error {
	taskID := c.Params("id")
	status := c.Params("status")

	// // Get user info from middleware context
	userID := int64(c.Locals("user_id").(float64))
	userRole := c.Locals("role").(string)

	claims := domain.Claims{
		UserID: userID,
		Role:   userRole,
	}

	IdTask, _ := strconv.ParseInt(taskID, 10, 64)

	// Call service to update task status and notify manager
	err := h.usecase.UpdateTaskStatus(IdTask, status, claims.UserID, claims.Role)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Task status updated successfully"})
}
