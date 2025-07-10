package kafka

type UserEventRequest struct {
	EventType string                 `json:"event_type"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}
