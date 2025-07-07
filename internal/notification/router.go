package notification

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pauljc-probeplus/Notification-System-v1-0-0/notification/handler"
	"github.com/pauljc-probeplus/Notification-System-v1-0-0/internal/notification/repository"
	"github.com/pauljc-probeplus/Notification-System-v1-0-0/internal/notification/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(app fiber.Router, db *mongo.Database) {
	repo := repository.NewNotificationRepository(db)
	svc := service.NewNotificationService(repo)
	handler := handler.NewNotificationHandler(svc)

	notificationGroup := app.Group("/notifications")
	notificationGroup.Post("/", handler.CreateNotification)
}
