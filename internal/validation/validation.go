package validation

import (
	"fmt"
	"regexp"
	"strings"

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

	// Register custom validation for notification type
	Validate.RegisterValidation("notificationType", func(fl validator.FieldLevel) bool {
		val := strings.ToLower(fl.Field().String())
		switch val {
		case "transactional", "promotional", "system_alert":
			return true
		}
		return false
	})

	// ✅ Custom validator for channels
	Validate.RegisterValidation("channels", func(fl validator.FieldLevel) bool {
		channels,ok := fl.Field().Interface().([]string)
		if !ok || len(channels) == 0 {
			// Reject nil or empty slice
			return false
		}
	
		allowed := map[string]bool{
			"email": true,
			"sms":   true,
		}
		for _, ch := range channels {
			if !allowed[strings.ToLower(ch)] {
				return false
			}
		}
		return true
	})

	//Custom validor for user preference id
	Validate.RegisterValidation("user_preference_id_format", func(fl validator.FieldLevel) bool {
		id := fl.Field().String()
		match, _ := regexp.MatchString(`^upref-\d{3}$`, strings.ToLower(id))
		return match
	})
}

func ValidateStruct(s interface{}) error {
	if err := Validate.Struct(s); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
