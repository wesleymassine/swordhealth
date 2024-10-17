package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wesleymassine/swordhealth/task-management/api/middleware"
)

type HTTPHandler struct {
	taskHandler *TaskHandler
}

func NewHTTPHandler(taskHandler *TaskHandler) *HTTPHandler {
	return &HTTPHandler{taskHandler: taskHandler}
}

func (h *HTTPHandler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/tasks/healthcheck", h.taskHandler.HealthCheck)

	// Protect these routes with authentication middleware
	api.Post("/tasks", middleware.AuthRequired, h.taskHandler.CreateTaskHandler)
	api.Get("/tasks", middleware.AuthRequired, h.taskHandler.ListTasksHandler)
	api.Patch("/tasks/:id/status", middleware.AuthRequired, h.taskHandler.UpdateTaskStatusHandler)
	// Protect these routes with x-api-key middleware integration of external services
	api.Get("/tasks/:task_id", middleware.ApiKeyMiddleware, h.taskHandler.GetUserAssignedTaskHandler)

	RegisterSwagger(app)
}
