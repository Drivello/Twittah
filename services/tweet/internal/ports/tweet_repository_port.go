package ports

import "github.com/Drivello/Twittah/services/tweet/internal/domain"

type TweetRepository interface {
	Save(tweet *domain.Tweet) error
	Delete(tweetID int64) error
	FindAllByUserID(userID int64) ([]*domain.Tweet, error)
	FindAllByMultipleUserIDs(userIDs []int64) ([]*domain.Tweet, error)
}
