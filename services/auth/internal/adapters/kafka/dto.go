package kafka

type KafkaEventRequest struct {
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload,omitempty"`
}

type KafkaEventError struct {
	EventType string `json:"event_type"`
	Error     error  `json:"error"`
}

type KafkaUserCreateRequest struct {
	EventType string                 `json:"event_type"`
	Payload   KafkaUserCreatePayload `json:"payload,omitempty"`
}
type KafkaUserCreatePayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Created is for user service
type KafkaUserCreatedRequest struct {
	EventType string                  `json:"event_type"`
	Payload   KafkaUserCreatedPayload `json:"payload,omitempty"`
}
type KafkaUserCreatedPayload struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

// NewUserCreatedPayload construye el payload tipado para el evento users.created
func NewUserCreatedPayload(userID int64, username string) KafkaUserCreatedPayload {
	return KafkaUserCreatedPayload{
		UserID:   userID,
		Username: username,
	}
}
