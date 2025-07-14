package kafka

type KafkaEventRequest struct {
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload,omitempty"`
}

type KafkaUserCreatedRequest struct {
	EventType string                  `json:"event_type"`
	Payload   KafkaUserCreatedPayload `json:"payload,omitempty"`
}
type KafkaUserCreatedPayload struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}
type KafkaFollowRequest struct {
	EventType string             `json:"event_type"`
	Payload   KafkaFollowPayload `json:"payload,omitempty"`
}

type KafkaFollowPayload struct {
	FollowerID int64 `json:"follower_id"`
	FolloweeID int64 `json:"followee_id"`
}

type KafkaEventError struct {
	EventType string `json:"event_type"`
	Error     error  `json:"error"`
}
