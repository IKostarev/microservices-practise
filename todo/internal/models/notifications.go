package models

const (
	TodoEventTypeEmailVerification = "todo_verify_email"
)

type TodoMailItem struct {
	TodoEventType string   `json:"todo_event_type"`
	Receivers     []string `json:"receivers"`
	Link          string   `json:"link"`
}
