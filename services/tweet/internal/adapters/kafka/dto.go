package kafka

// TweetPublishedEventDTO is the DTO for a tweet published event in Kafka
// (adapter-level, not domain)
type TweetPublishedEventDTO struct {
	ID        string `json:"id"`
	AuthorID  string `json:"author_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// FollowCreatedEventDTO is the DTO for a follow created event in Kafka
type FollowCreatedEventDTO struct {
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}

// FollowDeletedEventDTO is the DTO for a follow deleted event in Kafka
type FollowDeletedEventDTO struct {
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}
