package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterSwagger(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}
