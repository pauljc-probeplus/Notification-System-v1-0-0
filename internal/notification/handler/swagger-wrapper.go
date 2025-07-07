package handler

import (
	"github.com/gofiber/fiber/v2"
	//"notification-system/internal/notification/model"
	//"notification-system/internal/notification/service"
)
var H *NotificationHandler // This is your handler instance
// CreateNotification godoc
// @Summary Create a new notification
// @Description Create a notification and store it in MongoDB
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body notification-system/internal/notification/model.Notification true "Notification"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /api/v1/notifications [post]
func CreateNotification(c *fiber.Ctx) error {
	return 	H.CreateNotification(c) // <- call the method here
}