package app

import (
	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/wesleymassine/swordhealth/task-management/api"
	"github.com/wesleymassine/swordhealth/task-management/config"
	"github.com/wesleymassine/swordhealth/task-management/infra/db"
	"github.com/wesleymassine/swordhealth/task-management/infra/notification"
	"github.com/wesleymassine/swordhealth/task-management/usecase"

	"go.uber.org/fx"
)

func Run() {
	fx.New(
		fx.Provide(fiber.New),
		fx.Provide(config.NewMySQLConnection),
		fx.Provide(db.NewMySQLTaskRepository),
		fx.Provide(notification.SetupRabbitMQConnection),
		fx.Provide(usecase.NewTaskUseCase),
		fx.Provide(api.NewTaskHandler),
		fx.Provide(api.NewHTTPHandler),
		fx.Invoke(
			func(app *fiber.App, handler *api.HTTPHandler) {
				handler.SetupRoutes(app)
				app.Listen(":8082") // TODO PORTS
			},
		),
	).Run()
}
