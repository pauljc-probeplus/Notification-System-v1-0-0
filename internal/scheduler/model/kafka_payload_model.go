package model

type KafkaPayload struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}
