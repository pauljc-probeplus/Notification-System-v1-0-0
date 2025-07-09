package model

import (
	"time"
)

// Notification represents a notification payload
// swagger:model
type Notification struct {
	NotificationID   string    `bson:"notification_id" json:"notification_id" validate:"required,notification_id_format"`  // custom ID
	Type             string    `bson:"type" json:"type"`                         // transactional, promotional, etc.
	Channels         []string  `bson:"channels" json:"channels"`                // email, sms, etc.
	UserId           string    `bson:"user_id" json:"user_id"  validate:"required,user_id_format"`
	Message          string    `bson:"message" json:"message"`
	SendAt           string    `bson:"send_at" json:"send_at"`
	Priority         string    `bson:"priority" json:"priority"`
	MaximumRetries   string    `bson:"maximum_retries" json:"maximum_retries"`
	// CreatedDate is the time the notification was created
	// example: 2025-07-08T14:00:00Z
	CreatedDate      time.Time `json:"created_date" bson:"created_date"`
	CreatedByName    string    `json:"created_by_name" bson:"created_by_name"`
	CreatedById      string    `json:"created_by_id" bson:"created_by_id"`
	ModifiedDate     time.Time `json:"modified_date" bson:"modified_date"`
	ModifiedByName   string    `json:"modified_by_name" bson:"modified_by_name"`
	ModifiedById     string    `json:"modified_by_id" bson:"modified_by_id"`
}