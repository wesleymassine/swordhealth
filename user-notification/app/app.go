package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/wesleymassine/swordhealth/user-notification/api"
	"github.com/wesleymassine/swordhealth/user-notification/config"
	"github.com/wesleymassine/swordhealth/user-notification/infra/db"
	"github.com/wesleymassine/swordhealth/user-notification/infra/mq"
	userservice "github.com/wesleymassine/swordhealth/user-notification/infra/service-external/user-service"
	"github.com/wesleymassine/swordhealth/user-notification/usecase"
	"go.uber.org/fx"
)

func Run() {
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal catching for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fx.New(
		fx.Provide(
			fiber.New,
			config.NewRabbitMQConnection,
			config.NewRabbitMQChannel,
			config.NewMySQLConnection,
			mq.NewRabbitMQConsumer,
			db.NewMySQLRepository,
			usecase.NewNotificationService,
			api.NewHTTPHandler,
		),
		userservice.Module,

		fx.Invoke(
			func(service *usecase.NotificationService, app *fiber.App, handler *api.HTTPHandler) {
				handler.SetupRoutes(app)

				// Start the notification service
				go service.Start(appCtx)

				// Listen for system signals in a separate goroutine
				go func() {
					sig := <-sigs
					log.Printf("Received signal: %v", sig)

					// Signal cancellation and allow time for cleanup
					cancel()

					service.Shutdown()
					log.Println("Exiting...")
					os.Exit(0)
				}()

				// Start the HTTP server (non-blocking)
				app.Listen(":8083")
			},
		),
	).Run()
}
