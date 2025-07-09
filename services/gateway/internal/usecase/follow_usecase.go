package usecase

import "github.com/Drivello/Twittah/services/gateway/internal/ports"

// FollowUseCase provides the application logic for follow-related operations.
type FollowUseCase struct {
	Producer ports.FollowEventProducerPort
}

// NewFollowUseCase creates a new FollowUseCase.
func NewFollowUseCase(producer ports.FollowEventProducerPort) *FollowUseCase {
	return &FollowUseCase{Producer: producer}
}

// FollowUser handles follow logic.
func (uc *FollowUseCase) FollowUser(followerID, followeeID string) error {
	return uc.Producer.PublishFollow(followerID, followeeID)
}

// UnfollowUser handles unfollow logic.
func (uc *FollowUseCase) UnfollowUser(followerID, followeeID string) error {
	return uc.Producer.PublishUnfollow(followerID, followeeID)
}
