package routes

import (
	"github.com/KaningNoppasin/embedded-system-lab-backend/app/handlers"
	"github.com/gofiber/fiber/v3"
)

func RegisterNotificationRoutes(app *fiber.App, notificationHandler *handlers.NotificationHandler) {
	notifications := app.Group("/notifications")

	notifications.Post("/discord", notificationHandler.SendDiscordNotification)
}
