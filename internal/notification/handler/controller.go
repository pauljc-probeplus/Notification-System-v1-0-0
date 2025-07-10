package handler

import (
	"github.com/gofiber/fiber/v2"
	"notification-system/internal/notification/model"
	"notification-system/internal/notification/service"
	"notification-system/internal/validation"
	"time"
)

type NotificationHandler struct {
	svc service.NotificationService
}

func NewNotificationHandler(svc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}


// Test godoc
// @Summary Test
// @Description Dummy route
// @Tags test
// @Success 200 {string} string "test"
// @Router /test [get]
func (c *NotificationHandler) Test(ctx *fiber.Ctx) error {
	return ctx.SendString("test")
}



// CreateNotification godoc
// @Summary Create a new notification
// @Description Create a notification and store it in MongoDB
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body model.Notification true "Notification payload"
// @Success 201 {object} model.Notification
// @Router  /notifications [post]
func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	var req model.Notification
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	// Apply default values
	currentTime := time.Now().UTC().Format(time.RFC3339)
	req.CreatedDate = currentTime
	req.ModifiedDate = currentTime
	req.CreatedByName = req.UserId
	req.CreatedByID = req.UserId
	req.ModifiedByName = req.UserId
	req.ModifiedByID = req.UserId
	
	// 🔍 Perform validation
	if err := validation.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save to DB
	err := h.svc.CreateNotification(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create notification"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "notification created"})
}

// Dummy godoc
// @Summary Dummy check
// @Tags debug
// @Success 200 {string} string "ok"
// @Router /debug [get]
func Dummy(ctx *fiber.Ctx) error {
	return ctx.SendString("ok")
}
