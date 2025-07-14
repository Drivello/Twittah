package kafka

// KafkaEventRequest is a generic container for Kafka events.
type KafkaEventRequest[T any] struct {
	EventType string `json:"event_type"`
	Payload   T      `json:"payload,omitempty"`
}

// KafkaUserCreatePayload represents the payload for creating a user.
type KafkaUserCreatePayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// KafkaFollowPayload represents the payload for follow/unfollow events.
type KafkaFollowPayload struct {
	FollowerID int64 `json:"follower_id"`
	FolloweeID int64 `json:"followee_id"`
}

// KafkaUserCreatedPayload represents the payload for a created user.
type KafkaUserCreatedPayload struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

// KafkaTweetCreatePayload represents the payload for creating a tweet.
type KafkaTweetCreatePayload struct {
	AuthorID int64  `json:"author_id"`
	Content  string `json:"content"`
}

// KafkaTweetDeletePayload represents the payload for deleting a tweet.
type KafkaTweetDeletePayload struct {
	TweetID int64 `json:"tweet_id"`
}

// KafkaEventError encapsulates errors related to Kafka events.
type KafkaEventError struct {
	EventType string `json:"event_type"`
	Error     error  `json:"error"`
}
