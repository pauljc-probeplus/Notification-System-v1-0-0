package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/pauljc-probeplus/Notification-System-v1-0-0/internal/config"
	"github.com/pauljc-probeplus/Notification-System-v1-0-0/internal/notification"
)

func main() {
	// 1. Init MongoDB
	mongoURI := "mongodb://localhost:27017"
	db := config.InitMongoDB(mongoURI)

	// 2. Init Fiber App
	app := fiber.New()

	// 3. Register Routes
	api := app.Group("/api/v1")
	notification.RegisterRoutes(api, db)

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
	if err := config.MongoClient.Disconnect(nil); err != nil {
		log.Fatalf("❌ Mongo disconnect failed: %v", err)
	}
	log.Println("✅ MongoDB disconnected")
}
