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

type KafkaCreateTweetRequest struct {
	EventType string                  `json:"event_type"`
	Payload   KafkaCreateTweetPayload `json:"payload,omitempty"`
}

type KafkaCreateTweetPayload struct {
	AuthorID int64  `json:"author_id"`
	Content  string `json:"content"`
}

type KafkaDeleteTweetRequest struct {
	EventType string                  `json:"event_type"`
	Payload   KafkaDeleteTweetPayload `json:"payload,omitempty"`
}

type KafkaDeleteTweetPayload struct {
	TweetID int64 `json:"tweet_id"`
}

type KafkaEventError struct {
	EventType string `json:"event_type"`
	Error     error  `json:"error"`
}
