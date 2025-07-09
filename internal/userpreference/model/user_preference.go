package model

type UserPreference struct {
	UserPreferenceID string     `json:"user_preference_id" bson:"user_preference_id"`
	Version          string     `json:"version" bson:"version"`
	UserID           string     `json:"user_id" bson:"user_id"  validate:"required,user_id_format"`
	Preferences      Preference `json:"preferences" bson:"preferences"`

	CreatedDate   string `json:"created_date" bson:"created_date"`
	CreatedByName string `json:"created_by_name" bson:"created_by_name"`
	CreatedById   string `json:"created_by_id" bson:"created_by_id"`
	ModifiedDate  string `json:"modified_date" bson:"modified_date"`
	ModifiedByName string `json:"modified_by_name" bson:"modified_by_name"`
	ModifiedById  string `json:"modified_by_id" bson:"modified_by_id"`
}

type Preference struct {
	Channels       Channel     `json:"channels" bson:"channels"`
	NoDisturb      NoDisturb   `json:"no_disturb_details" bson:"no_disturb_details"`
	DailyLimit     DailyLimit  `json:"daily_limit_details" bson:"daily_limit_details"`
	DeliveryTiming DeliveryTime `json:"delivery_time" bson:"delivery_time"`
}

type Channel struct {
	Transactional []string `json:"transactional" bson:"transactional"`
	Promotional   []string `json:"promotional" bson:"promotional"`
	SystemAlerts  []string `json:"system_alerts" bson:"system_alerts"`
}

type NoDisturb struct {
	Enabled     bool   `json:"enabled" bson:"enabled"`
	StartDateTime string `json:"start_date_time" bson:"start_date_time"`
	EndDateTime   string `json:"end_date_time" bson:"end_date_time"`
	TimeZone      string `json:"time_zone" bson:"time_zone"`
}

type DailyLimit struct {
	PromotionalLimit      string `json:"promotional_limit" bson:"promotional_limit"`
	PromotionalSentToday  string `json:"promotional_sent_today" bson:"promotional_sent_today"`
}

type DeliveryTime struct {
	Enabled    bool   `json:"enabled" bson:"enabled"`
	StartTime  string `json:"start_time" bson:"start_time"`
	EndTime    string `json:"end_time" bson:"end_time"`
	TimeZone   string `json:"time_zone" bson:"time_zone"`
}
