package userpreference

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"notification-system/internal/userpreference/handler"
	"notification-system/internal/userpreference/repository"
	"notification-system/internal/userpreference/service"
)

func RegisterRoutes(app fiber.Router, db *mongo.Database) {
	repo := repository.NewUserPreferenceRepository(db)
	svc := service.NewUserPreferenceService(repo)
	handler := handler.NewUserPreferenceHandler(svc)

	// prefGroup := app.Group("/user-preferences")
	// prefGroup.Post("/", handler.CreateUserPreference)
	// userPrefGroup := app.Group("/user-preferences")
	// userPrefGroup.Put("/:id", handler.UpdateUserPreference)


	app.Post("/user-preferences", handler.CreateUserPreference)
	app.Put("/user-preferences/:id", handler.UpdateUserPreference)
}
