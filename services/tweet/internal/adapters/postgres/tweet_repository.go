package postgres

import (
	"context"
	"time"

	"github.com/Drivello/Twittah/services/tweet/ent"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
)

type TweetRepository struct {
	Client *ent.Client
}

func NewTweetRepository(client *ent.Client) *TweetRepository {
	return &TweetRepository{Client: client}
}

func (r *TweetRepository) Save(t *domain.Tweet) error {
	_, err := r.Client.Tweet.Create().
		SetContent(t.Content).
		SetCreatedAt(time.Now()).
		Save(context.Background())
	return err
}

func (r *TweetRepository) Delete(tweetID int64) error {
	return r.Client.Tweet.DeleteOneID(tweetID).Exec(context.Background())
}
