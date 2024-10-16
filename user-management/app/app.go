package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wesleymassine/swordhealth/user-management/api"
	"github.com/wesleymassine/swordhealth/user-management/config"
	"github.com/wesleymassine/swordhealth/user-management/infra/db"
	"github.com/wesleymassine/swordhealth/user-management/usecase"
	"go.uber.org/fx"
)

func Run() {
	app := fiber.New()

	fx.New(
		fx.Provide(
			config.NewMySQLConnection,
			db.NewUserRepository,
			usecase.NewUserUsecase,
			api.NewUserHandler,
		),
		fx.Invoke(
			func(handler *api.UserHandler) {
				handler.RegisterRoutes(app)
				app.Listen(":8081")
			},
		),
	).Run()
}
