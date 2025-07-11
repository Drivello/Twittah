package kafka

type UserEventRequest struct {
	EventType string                 `json:"event_type"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

// UserCreateRequestDTO is the DTO for a user create request in Kafka
type UserCreateRequestDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
