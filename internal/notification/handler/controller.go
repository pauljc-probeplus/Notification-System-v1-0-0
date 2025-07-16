package handler

import (
	"github.com/gofiber/fiber/v2"
	"notification-system/internal/notification/model"
	"notification-system/internal/notification/service"
	"notification-system/internal/validation"
	"time"
	"log"
	"strings"
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
	loc, _ := time.LoadLocation("Asia/Kolkata")
	currentTime := time.Now().In(loc).Format("2006-01-02T15:04:05")
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
	// err := h.svc.CreateNotification(c.Context(), &req)
	// if err != nil {
	// 	log.Println("⚠️ CreateNotification failed:", err.Error())
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create notification"})
	// }

	// return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "notification created"})

	//below is the code befoe error msg handling
	// if err := h.svc.CreateNotification(c.Context(), &req); err != nil {
	// 	log.Println("⚠️ CreateNotification failed:", err.Error())
	
	// 	if strings.Contains(err.Error(), "duplicate") {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "duplicate entry",
	// 		})
	// 	}
	
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "scheduler failed while creating notification",
	// 	})
	// }
	// return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "notification created and entered to scheduler"})

	wasScheduled, err := h.svc.CreateNotification(c.Context(), &req)
	if err != nil {
		log.Println("⚠️ CreateNotification failed:", err.Error())

		if strings.Contains(err.Error(), "duplicate") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "duplicate entry",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "scheduler failed while creating notification",
		})
	}

	msg := "notification created"
	if wasScheduled {
		msg += " and entered to scheduler"
	} else {
		msg += " but not scheduled"
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": msg})



}

// Dummy godoc
// @Summary Dummy check
// @Tags debug
// @Success 200 {string} string "ok"
// @Router /debug [get]
func Dummy(ctx *fiber.Ctx) error {
	return ctx.SendString("ok")
}
