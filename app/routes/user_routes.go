package routes

import (
	"github.com/KaningNoppasin/embedded-system-lab-backend/app/handlers"
	ws "github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(app *fiber.App, userHandler *handlers.UserHandler) {
	users := app.Group("/users")
	upgrader := ws.FastHTTPUpgrader{}

	users.Get("/ws", func(c fiber.Ctx) error {
		if !ws.FastHTTPIsWebSocketUpgrade(c.RequestCtx()) {
			return c.SendStatus(fiber.StatusUpgradeRequired)
		}

		return upgrader.Upgrade(c.RequestCtx(), userHandler.RFIDWebSocket)
	})
	users.Post("/rfid-action", userHandler.SendRFIDToWebSocket)
	users.Get("/", userHandler.GetAllUsers)
	users.Post("/", userHandler.CreateUser)
	users.Delete("/:rfid", userHandler.DeleteUserByRFID)
	users.Get("/:rfid/amount", userHandler.GetAmountByRFID)
	users.Patch("/:rfid/amount", userHandler.UpdateAmountByRFID)
}
