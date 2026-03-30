package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KaningNoppasin/embedded-system-lab-backend/app/config"
	"github.com/KaningNoppasin/embedded-system-lab-backend/app/database"
	"github.com/KaningNoppasin/embedded-system-lab-backend/app/handlers"
	"github.com/KaningNoppasin/embedded-system-lab-backend/app/mqtt"
	"github.com/KaningNoppasin/embedded-system-lab-backend/app/repositories"
	"github.com/KaningNoppasin/embedded-system-lab-backend/app/routes"
	"github.com/KaningNoppasin/embedded-system-lab-backend/app/services"
	"github.com/KaningNoppasin/embedded-system-lab-backend/app/timeseries"
	"github.com/gofiber/fiber/v3"
)

func main() {
	if err := config.LoadEnvFile(".env"); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	mongoClient, userCollection, transactionCollection, err := database.ConnectMongo()
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("failed to disconnect mongo: %v", err)
		}
	}()

	app := fiber.New()
	userRepository, err := repositories.NewUserRepository(userCollection)
	if err != nil {
		log.Fatalf("failed to setup user repository: %v", err)
	}
	transactionRepository, err := repositories.NewTransactionRepository(transactionCollection)
	if err != nil {
		log.Fatalf("failed to setup transaction repository: %v", err)
	}
	rfidWebSocketHub := services.NewRFIDWebSocketHub()
	discordNotifier := services.NewDiscordNotifier()
	lineNotifier := services.NewLineNotifier()
	userHandler := handlers.NewUserHandler(userRepository, rfidWebSocketHub)
	notificationHandler := handlers.NewNotificationHandler(discordNotifier, lineNotifier)
	mqttPublisher, err := mqtt.NewPublisher()
	if err != nil {
		log.Fatalf("failed to connect mqtt publisher: %v", err)
	}
	defer mqttPublisher.Close()
	transactionHandler := handlers.NewTransactionHandler(transactionRepository, userRepository, mqttPublisher)

	routes.RegisterUserRoutes(app, userHandler)
	routes.RegisterNotificationRoutes(app, notificationHandler)
	routes.RegisterTransactionRoutes(app, transactionHandler)

	influxWriter, err := timeseries.NewInfluxWriter()
	if err != nil {
		log.Fatalf("failed to connect influxdb: %v", err)
	}
	defer influxWriter.Close()

	temperatureSubscriber, err := mqtt.NewTemperatureSubscriber(influxWriter, rfidWebSocketHub)
	if err != nil {
		log.Fatalf("failed to connect mqtt broker: %v", err)
	}
	defer temperatureSubscriber.Close()

	listenAddr := ":" + getEnv("APP_PORT", "8080")

	go func() {
		if err := app.Listen(listenAddr); err != nil {
			log.Fatalf("fiber server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("failed to shutdown app: %v", err)
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
