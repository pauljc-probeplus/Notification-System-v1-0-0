package handler

import (
	"github.com/gofiber/fiber/v2"
	"notification-system/internal/userpreference/model"
	"notification-system/internal/userpreference/service"
)

type UserPreferenceHandler struct {
	svc service.UserPreferenceService
}

func NewUserPreferenceHandler(svc service.UserPreferenceService) *UserPreferenceHandler {
	return &UserPreferenceHandler{svc: svc}
}

// @Summary Create user preference
// @Description Stores user preferences for a given user
// @Tags user-preferences
// @Accept json
// @Produce json
// @Param user_preference body model.UserPreference true "User Preference"
// @Success 201 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /user-preferences [post]
func (h *UserPreferenceHandler) CreateUserPreference(c *fiber.Ctx) error {
	var pref model.UserPreference
	if err := c.BodyParser(&pref); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.svc.CreateUserPreference(c.Context(), &pref); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to store preference"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "preference stored"})
}

// UpdateUserPreference godoc
// @Summary Update user preferences
// @Description Update an existing user preference by ID
// @Tags user-preferences
// @Accept json
// @Produce json
// @Param id path string true "User Preference ID"
// @Param userPreference body model.UserPreference true "User Preference Payload"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /user-preferences/{id} [put]
func (h *UserPreferenceHandler) UpdateUserPreference(c *fiber.Ctx) error {
	//userID := c.Params("user_id")
	var req model.UserPreference

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	err := h.svc.UpdateUserPreference(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "updated successfully"})
}


