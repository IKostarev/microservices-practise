package models

const (
	UserEventTypeEmailVerification = "user_verify_email"
)

type UserMailItem struct {
	UserEventType string   `json:"user_event_type"`
	Receivers     []string `json:"receivers"`
	Link          string   `json:"link"`
}
