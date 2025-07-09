// Package domain contains core business entities for the User Service.
package domain

// FollowEvent represents a follow/unfollow event in the domain layer.
type FollowEvent struct {
	EventType   string `json:"event_type"` // Type of event (e.g., follow_created, follow_deleted)
	FollowerID  string `json:"follower_id"` // ID of the user who follows
	FolloweeID  string `json:"followee_id"` // ID of the user being followed
}
