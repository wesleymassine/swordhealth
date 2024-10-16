package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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
	// TODO: middleware
	api.Get("/notifications", h.listLatestNotifications)
}

func (h *HTTPHandler) listLatestNotifications(c *fiber.Ctx) error {
	limit := 10 // TODO: query params LIMIT OFFSET
	notifications, err := h.notification.ListLatestNotifications(limit)

	fmt.Println("DEBUG", err)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch notifications")
	}

	return c.JSON(notifications)
}
