package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wesleymassine/swordhealth/user-management/api/middleware"
	"github.com/wesleymassine/swordhealth/user-management/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: uc}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Public routes
	api.Post("/users/login", h.Login)

	auth := app.Group("/api/v1")

	// Protected routes
	auth.Use(middleware.JWTMiddleware())

	auth.Post("/users/register", h.CreateUser)
	auth.Get("/users/profile/:id", h.GetUser)
	auth.Put("/users/update/:id", h.UpdateUser)
	auth.Delete("/users/delete/:id", h.DeleteUser)
}