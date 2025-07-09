package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	//fiberSwagger "github.com/swaggo/fiber-swagger"  // alias to use in code
	//swaggerFiles "github.com/swaggo/files"          // Swagger UI assets
	_ "notification-system/docs"
	_ "notification-system/internal/notification/handler"

	"github.com/gofiber/swagger"

	"notification-system/internal/config"
	"notification-system/internal/notification/route"

	userpreference "notification-system/internal/userpreference/routes"

	"notification-system/internal/validation"

	"github.com/gofiber/fiber/v2"
)

// @title Notification System API
// @version 1.0
// @description API docs for Notification System
// @contact.name Paul
// @host localhost:4001
// @BasePath /api/v1

func main() {
	// 1. Init MongoDB
	mongoURI := "mongodb://localhost:27017"
	db := config.InitMongoDB(mongoURI)

	// Validation checks
	validation.InitValidator()

	// 2. Init Fiber App
	app := fiber.New()

	// Swagger Setup (new method)
	// swaggerHandler:= swagger.New(swagger.Config{
	// 	URL: "http://localhost:4001/swagger/doc.json",
	// })
	// if err != nil {
	// 	log.Fatalf("Swagger setup failed: %v", err)
	// }
	app.Get("/swagger/*", swagger.New())

	// 3. Register Routes
	api := app.Group("/api/v1")
	
	route.RegisterRoutes(api, db)
	userpreference.RegisterRoutes(api,db)
	app.Get("/ping", Ping)


	// 4. Start Server
	go func() {
		if err := app.Listen(":4001"); err != nil {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()
	log.Println("🚀 Server running at http://localhost:4001")

	// 5. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("🛑 Shutting down...")
	if err := config.MongoClient.Disconnect(context.TODO()); err != nil {
		log.Fatalf("❌ Mongo disconnect failed: %v", err)
	}
	log.Println("✅ MongoDB disconnected")
}

// Ping godoc
// @Summary Ping test
// @Description Just ping
// @Tags test
// @Success 200 {string} string "pong"
// @Router /ping [get]
func Ping(c *fiber.Ctx) error {
	return c.SendString("pong")
}

