package api

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Health represents the structure for the health check response
type Health struct {
	HealthCheck string    `json:"health_check" example:"alive"`
	Time        time.Time `json:"time" example:"2023-04-17T20:07:50Z"`
}

// HealthCheck performs a health check for the API
//
//	@Router			/api/v1/healthcheck [get]
//	@Summary		Health-Check
//	@Description	Check API Alive
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	Health
func (h *HTTPHandler) HealthCheck(ctx *fiber.Ctx) error {
	// Use shared utility to get current time
	now := time.Now().Local()

	// Create the health check response
	response := Health{
		HealthCheck: "alive",
		Time:        now,
	}

	// Return the health check response as JSON
	return ctx.Status(http.StatusOK).JSON(response)
}
