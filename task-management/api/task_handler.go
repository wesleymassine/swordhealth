package api

import (
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3/log"
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
func (h *TaskHandler) CreateTaskHandler(ctx *fiber.Ctx) error {
	var task domain.Task
	if err := ctx.BodyParser(&task); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// // Get user info from middleware context
	userID := int64(ctx.Locals("user_id").(float64))
	userRole := ctx.Locals("role").(string)

	// Optionally, you can validate user role or perform different actions based on role
	if userRole != "super_admin" && userRole != "manager" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	// if err := RoleMiddleware(claims.Role); err != nil {
	// 	return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	// }

	task.AssignedTo = userID // Assign the task to the user

	taskResponse, err := h.usecase.CreateTask(ctx.Context(), task)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	log.Info("Task created successfully")

	return ctx.Status(fiber.StatusCreated).JSON(taskResponse)
}

// ListTasks handles GET /tasks
func (h *TaskHandler) ListTasksHandler(ctx *fiber.Ctx) error {
	// Get user info from middleware context
	userID := int64(ctx.Locals("user_id").(float64))
	userRole := ctx.Locals("role").(string)

	tasks, err := h.usecase.ListTasks(ctx.Context(), userRole, userID)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(tasks)
}

// UpdateTaskStatusHandler handles PATCH requests to update the task status
func (c *TaskHandler) UpdateTaskStatusHandler(ctx *fiber.Ctx) error {
	// Parse task ID from the URL path
	taskIDParam := ctx.Params("id")

	taskID, err := strconv.ParseInt(taskIDParam, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid task ID"})
	}

	// Get the user ID and role from the middleware (assuming user is authenticated)
	userID := int64(ctx.Locals("user_id").(float64))
	userRole := ctx.Locals("role").(string)

	// TODO: Parse request body for new status
	type UpdateTaskRequest struct {
		Status string `json:"status"`
	}

	var req UpdateTaskRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Call the service to update the task status
	err = c.usecase.UpdateTaskStatus(ctx.Context(), taskID, userID, userRole, req.Status)

	if err != nil {
		if err.Error() == "unauthorized" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "task updated successfully"})
}
