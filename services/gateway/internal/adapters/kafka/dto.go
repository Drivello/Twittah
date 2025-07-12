package kafka

// KafkaEventRequest es un contenedor genérico para eventos Kafka
type KafkaEventRequest[T any] struct {
	EventType string `json:"event_type"`
	Payload   T      `json:"payload,omitempty"`
}

// KafkaUserCreatePayload representa el payload para crear un usuario
type KafkaUserCreatePayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// KafkaFollowPayload representa el payload para follow/unfollow
type KafkaFollowPayload struct {
	FollowerID int64 `json:"follower_id"`
	FolloweeID int64 `json:"followee_id"`
}

// KafkaUserCreatedPayload representa el payload de usuario creado
type KafkaUserCreatedPayload struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

// KafkaTweetCreatePayload representa el payload para crear un tweet
type KafkaTweetCreatePayload struct {
	TweetID  int64  `json:"tweet_id"`
	AuthorID int64  `json:"author_id"`
	Content  string `json:"content"`
}

// KafkaEventError encapsula errores relacionados a eventos Kafka
type KafkaEventError struct {
	EventType string `json:"event_type"`
	Error     error  `json:"error"`
}
