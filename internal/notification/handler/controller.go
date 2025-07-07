package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pauljc-probeplus/Notification-System-v1-0-0/internal/notification"
	"github.com/pauljc-probeplus/Notification-System-v1-0-0/internal/notification/service"
)

type NotificationHandler struct {
	svc service.NotificationService
}

func NewNotificationHandler(svc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	var req notification.Notification
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	err := h.svc.CreateNotification(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create notification"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "notification created"})
}
