package ports

type UserEventPublisher interface {
	PublishFollow(followerID, followeeID string) error
	PublishUnfollow(followerID, followeeID string) error
}
