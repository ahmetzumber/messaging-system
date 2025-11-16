package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"messaging-system/app/cache"
	"messaging-system/app/client"
	"messaging-system/app/handler"
	"messaging-system/app/processor"
	"messaging-system/app/repository"
	"messaging-system/app/service"
	"messaging-system/config"
	_ "messaging-system/docs"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Messaging System API
// @version 1.0
// @description API for managing message processor and retrieving sent messages
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	appConfig, err := config.NewConfig(".config", os.Getenv("APP_ENV"))
	if err != nil {
		log.Fatal(err)
	}
	appConfig.Print()
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mongoRepo, err := repository.New(ctx, appConfig.Mongo)
	if err != nil {
		log.Fatal(err)
	}

	webhookClient := client.NewClient(appConfig.Client, logger)
	redis := cache.NewRedis(appConfig.Redis)
	messageService := service.NewMessageService(mongoRepo)
	messageProcessor := processor.NewMessageProcessor(messageService, webhookClient, redis, logger)
	messageHandler := handler.NewMessageHandler(messageProcessor)

	server := fiber.New()
	server.Use(
		cors.New(cors.ConfigDefault),
	)
	server.Get("/swagger/*", fiberSwagger.WrapHandler)

	messageHandler.RegisterRoutes(server)
	go log.Fatal(
		server.Listen(fmt.Sprintf(":%d", appConfig.Server.Port)),
	)
}
