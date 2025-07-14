package route

import (
	"github.com/gofiber/fiber/v2"
	"notification-system/internal/notification/handler"
    "notification-system/internal/notification/service"
    "notification-system/internal/notification/repository"
	upref_repo "notification-system/internal/userpreference/repository"
	sched_repo "notification-system/internal/scheduler/repository"
	upref_svc "notification-system/internal/userpreference/service"
	sched_svc "notification-system/internal/scheduler/service"
	"go.mongodb.org/mongo-driver/mongo"
)

// @title Notification System Routes
// @description All routes for notification APIs
func RegisterRoutes(app fiber.Router, db *mongo.Database) {
	
	//Routes registered before coding scheduler
	/*repo := repository.NewNotificationRepository(db)
	svc := service.NewNotificationService(repo)
	handler := handler.NewNotificationHandler(svc)*/

	//handler.H=handlerInstance

	// notificationGroup := app.Group("/notifications")
	// notificationGroup.Post("/", handler.CreateNotification)

	repo := repository.NewNotificationRepository(db)
	prefRepo := upref_repo.NewUserPreferenceRepository(db)
	schedulerRepo := sched_repo.NewSchedulerEntryRepository(db)

	svc := service.NewNotificationService(
		repo,
		upref_svc.NewUserPreferenceService(prefRepo),
		sched_svc.NewSchedulerEntryService(schedulerRepo),
	)

	handler := handler.NewNotificationHandler(svc)


	app.Post("/notifications", handler.CreateNotification)
}
