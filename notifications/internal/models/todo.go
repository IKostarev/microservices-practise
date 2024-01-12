package models

type TodoMailItem struct {
	//UserEventType string   `json:"user_event_type"`
	Receivers []string `json:"receivers"`
	//Link          string   `json:"link"`
}
