package api

import (
	"github.com/gofiber/fiber/v2"
)

type HTTPHandler struct {
	taskHandler *TaskHandler
}

func NewHTTPHandler(taskHandler *TaskHandler) *HTTPHandler {
	return &HTTPHandler{taskHandler: taskHandler}
}

func (h *HTTPHandler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Post("/login", Login)

	// Protect these routes with authentication middleware
	api.Post("/tasks", AuthRequired, h.taskHandler.CreateTask)
	api.Get("/tasks", AuthRequired, h.taskHandler.ListTasks)
	api.Patch("/tasks/:id/status/:status", AuthRequired, h.taskHandler.UpdateTaskStatus)
	RegisterSwagger(app)
}
