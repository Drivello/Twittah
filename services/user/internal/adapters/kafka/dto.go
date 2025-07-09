package kafka

// FollowEventDTO is the DTO for a follow event in Kafka
type FollowEventDTO struct {
	EventType  string `json:"event_type"`
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}
