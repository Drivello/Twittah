package ports

// UserEventProducerPort defines the contract for producing user events.
type UserEventProducerPort interface {
	PublishUserCreateRequest(username, email, password string) error
	PublishUserCreated(id, username string) error
	PublishUserDeleted(id string) error
}

// FollowEventProducerPort defines the contract for producing follow events.
type FollowEventProducerPort interface {
	PublishFollow(followerID, followeeID string) error
	PublishUnfollow(followerID, followeeID string) error
}
