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

	// Protect these routes with authentication middleware
	api.Post("/tasks", middleware.AuthRequired, h.taskHandler.CreateTaskHandler)
	api.Get("/tasks", middleware.AuthRequired, h.taskHandler.ListTasksHandler)
	api.Patch("/tasks/:id/status", middleware.AuthRequired, h.taskHandler.UpdateTaskStatusHandler)
	RegisterSwagger(app)
}
