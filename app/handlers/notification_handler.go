package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

type discordNotifier interface {
	Send(content string) error
}

type NotificationHandler struct {
	discordNotifier discordNotifier
}

type sendDiscordNotificationRequest struct {
	Message string `json:"message"`
}

func NewNotificationHandler(discordNotifier discordNotifier) *NotificationHandler {
	return &NotificationHandler{
		discordNotifier: discordNotifier,
	}
}

func (h *NotificationHandler) SendDiscordNotification(c fiber.Ctx) error {
	var req sendDiscordNotificationRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "message is required",
		})
	}

	if h.discordNotifier == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"message": "discord notifier is unavailable",
		})
	}

	if err := h.discordNotifier.Send(req.Message); err != nil {
		statusCode := fiber.StatusInternalServerError
		if strings.Contains(err.Error(), "DISCORD_WEBHOOK_URL is not configured") {
			statusCode = fiber.StatusServiceUnavailable
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"message": "failed to send discord notification",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "discord notification sent",
	})
}
