package validation

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()

	// Validate format: user-123
	_ = Validate.RegisterValidation("user_id_format", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString(`^user-\d{3}$`, fl.Field().String())
		return match
	})

	// Validate format: notif-123
	_ = Validate.RegisterValidation("notification_id_format", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString(`^notif-\d{3}$`, fl.Field().String())
		return match
	})
}

func ValidateStruct(s interface{}) error {
	if err := Validate.Struct(s); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
