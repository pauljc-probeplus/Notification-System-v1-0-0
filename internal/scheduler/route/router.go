/*package route

import (
    "context"
    "notification-system/internal/scheduler/service"

    "github.com/gofiber/fiber/v2"
)

func RegisterSchedulerRoutes(router *fiber.App, schedulerService service.SchedulerEntryService) {
    scheduler := router.Group("/scheduler")

    scheduler.Get("/ping", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "scheduler running"})
    })

    scheduler.Post("/poll", func(c *fiber.Ctx) error {
        go schedulerService.PollAndDispatch(context.Background())
        return c.JSON(fiber.Map{"message": "Polling started"})
    })
}*/

package route

import (
    "notification-system/internal/scheduler/repository"
    "notification-system/internal/scheduler/service"
	"notification-system/internal/scheduler/controller"

    "github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterSchedulerRoutes(app fiber.Router, db *mongo.Database) {
  
    // Wire repo -> service -> route
    repo := repository.NewSchedulerEntryRepository(db)
    svc := service.NewSchedulerEntryService(repo)
	//handler := controller.PollScheduler(svc)
	
	// Delegate to controller
    app.Get("/ping", controller.Ping)
    app.Post("/poll", controller.PollScheduler(svc))

    //consumer logic
    app.Post("/start-consumer", controller.StartKafkaConsumer)

	
	//func RegisterRoutes(app fiber.Router, db *mongo.Database) {

	
}

