package model

import(
	//"time"
)

type SchedulerEntry struct {
	SchedulerEntryID string `bson:"scheduler_entry_id" json:"scheduler_entry_id"` 
	NotificationID string `bson:"notification_id" json:"notification_id"` 
	NotificationType string `bson:"notification_type" json:"notification_type"` 
	Topic string `bson:"topic" json:"topic"` 
	UserId         string `bson:"user_id" json:"user_id"`
	Message string   `bson:"message" json:"message"` 
	SendAt string   `bson:"send_at" json:"send_at"`
	//Priority string   `bson:"priority" json:"priority"` 
	MaximumRetries string   `bson:"maximum_retries" json:"maximum_retries"`
	Attempts string   `bson:"attempts" json:"attempts"` 
	Status         string    `json:"status"` // e.g., pending, sent, failed
	CreatedDate string `json:"created_date" bson:"created_date"`
	CreatedByName string `json:"created_by_name" bson:"created_by_name"`
	CreatedById string `json:"created_by_id" bson:"created_by_id"`
	ModifiedDate string `json:"modified_date" bson:"modified_date"`
	ModifiedByName string `json:"modified_by_name" bson:"modified_by_name"`
	ModifiedById string `json:"modified_by_id" bson:"modified_by_id"`
}
