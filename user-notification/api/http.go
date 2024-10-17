package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wesleymassine/swordhealth/user-notification/api/middleware"
	"github.com/wesleymassine/swordhealth/user-notification/usecase"
)

type HTTPHandler struct {
	notification *usecase.NotificationService
}

func NewHTTPHandler(notification *usecase.NotificationService) *HTTPHandler {
	return &HTTPHandler{notification: notification}
}

func (h *HTTPHandler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/notifications", middleware.AuthRequired, h.listLatestNotifications)
	// api.Get("/notifications/healthcheck", h.HealthCheck)
}

func (h *HTTPHandler) listLatestNotifications(c *fiber.Ctx) error {
	limit := 10 // TODO: query params LIMIT OFFSET

	notifications, err := h.notification.ListLatestNotifications(limit)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch notifications")
	}

	return c.JSON(notifications)
}
