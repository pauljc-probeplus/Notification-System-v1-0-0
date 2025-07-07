package notification

import (
	"github.com/gofiber/fiber/v2"
	"notification-system/internal/notification/handler"
    "notification-system/internal/notification/service"
    "notification-system/internal/notification/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(app fiber.Router, db *mongo.Database) {
	repo := repository.NewNotificationRepository(db)
	svc := service.NewNotificationService(repo)
	handlerInstance := handler.NewNotificationHandler(svc)

	handler.H=handlerInstance

	notificationGroup := app.Group("/notifications")
	notificationGroup.Post("/", handler.CreateNotification)
}
