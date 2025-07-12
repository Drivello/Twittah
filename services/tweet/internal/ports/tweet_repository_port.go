package ports

import "github.com/Drivello/Twittah/services/tweet/internal/domain"

type TweetRepository interface {
	Save(tweet *domain.Tweet) error
	Delete(tweetID int64) error
	FindAllByUserIDs(userIDs []int64) ([]*domain.Tweet, error)
}
