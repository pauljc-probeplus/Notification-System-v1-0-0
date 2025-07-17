package model

type FailureLog struct {
	LogID          string `json:"log_id" bson:"log_id"` // UUID
	NotificationID string `json:"notification_id" bson:"notification_id"`
	UserID         string `json:"user_id" bson:"user_id"`
	Type           string `json:"type" bson:"type"` // e.g., promotional
	Message        string `json:"message" bson:"message"`
	FailureReason  string `json:"failure_reason" bson:"failure_reason"`
	Timestamp      string `json:"timestamp" bson:"timestamp"` // Use time.Now().Format("2006-01-02T15:04:05")
}
