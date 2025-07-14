// Package ports defines interfaces (ports) for driving and driven adapters in the hexagonal architecture.
package ports

// UserEventPublisher defines the contract for publishing user-related events (outbound port).
type UserEventPublisher interface {
	// PublishFollow publishes a follow event.
	PublishFollow(followerID, followeeID string) error
	// PublishUnfollow publishes an unfollow event.
	PublishUnfollow(followerID, followeeID string) error
}
