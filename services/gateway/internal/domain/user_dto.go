package domain

type UserEvent struct {
	EventType  string `json:"event_type"`
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}
