package controller

import (
    "context"
    "notification-system/internal/scheduler/service"

    "github.com/gofiber/fiber/v2"
)

// type SchedulerHandler struct {
// 	svc service.SchedulerEntryService
// }

// func NewSchedulerHandler(svc service.SchedulerEntryService) *SchedulerHandler {
// 	return &SchedulerHandler{svc: svc}
// }



// Ping health check endpoint
func Ping(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"status": "scheduler running"})
}

// PollScheduler triggers background dispatch logic
func PollScheduler(schedulerService service.SchedulerEntryService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        go schedulerService.PollAndDispatch(context.Background())
        return c.JSON(fiber.Map{"message": "Polling started"})
    }
}

//Kafka consumer begins consumption
func StartKafkaConsumer(c *fiber.Ctx) error {
	go service.ConsumeFromKafka("localhost:9092", "lowPriority.sms", service.HandleKafkaNotification)
	return c.JSON(fiber.Map{
		"status":  "Kafka consumer started",
		"message": "Consuming from topic: lowPriority.sms",
	})
}