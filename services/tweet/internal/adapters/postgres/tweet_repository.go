package postgres

import (
	"context"
	"time"

	"github.com/Drivello/Twittah/services/tweet/ent"
	tweetEnt "github.com/Drivello/Twittah/services/tweet/ent/tweet"
	userEnt "github.com/Drivello/Twittah/services/tweet/ent/user"
	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"go.uber.org/zap"
)

type TweetRepository struct {
	Client *ent.Client
}

func NewTweetRepository(client *ent.Client) *TweetRepository {
	return &TweetRepository{Client: client}
}

func (r *TweetRepository) Save(t *domain.Tweet) error {
	common.Logger().Debug("[TweetRepository] Save called", zap.Any("tweet", t))
	_, err := r.Client.Tweet.Create().
		SetAuthorID(t.AuthorID).
		SetContent(t.Content).
		SetCreatedAt(time.Now()).
		Save(context.Background())
	return err
}

func (r *TweetRepository) Delete(tweetID int64) error {
	common.Logger().Debug("[TweetRepository] Delete called", zap.Int64("tweet_id", tweetID))
	return r.Client.Tweet.DeleteOneID(tweetID).Exec(context.Background())
}

func (r *TweetRepository) FindByID(tweetID int64) (*domain.Tweet, error) {
	common.Logger().Debug("[TweetRepository] FindByID called", zap.Int64("tweet_id", tweetID))
	tweet, err := r.Client.Tweet.
		Query().
		Where(tweetEnt.ID(tweetID)).
		WithAuthor().
		Only(context.Background())
	if err != nil {
		return nil, err
	}
	return &domain.Tweet{
		ID:        tweet.ID,
		AuthorID:  tweet.Edges.Author.ID,
		Content:   tweet.Content,
		CreatedAt: tweet.CreatedAt.Unix(),
	}, nil
}

func (r *TweetRepository) FindAllByUserID(userID int64) ([]*domain.Tweet, error) {
	common.Logger().Debug("[TweetRepository] FindAllByUserID called", zap.Int64("user_id", userID))
	tweets, err := r.Client.Tweet.
		Query().
		Where(tweetEnt.HasAuthorWith(userEnt.ID(userID))).
		WithAuthor().
		All(context.Background())
	if err != nil {
		return nil, err
	}

	var domainTweets []*domain.Tweet
	for _, tweet := range tweets {
		domainTweets = append(domainTweets, &domain.Tweet{
			ID:        tweet.ID,
			AuthorID:  tweet.Edges.Author.ID,
			Content:   tweet.Content,
			CreatedAt: tweet.CreatedAt.Unix(),
		})
	}

	return domainTweets, nil
}

func (r *TweetRepository) FindAllByMultipleUserIDs(userIDs []int64) ([]*domain.Tweet, error) {
	common.Logger().Debug("[TweetRepository] FindAllByMultipleUserIDs called", zap.Any("user_ids", userIDs))
	tweets, err := r.Client.Tweet.
		Query().
		Where(tweetEnt.HasAuthorWith(userEnt.IDIn(userIDs...))).
		WithAuthor().
		All(context.Background())
	if err != nil {
		return nil, err
	}

	var domainTweets []*domain.Tweet
	for _, tweet := range tweets {
		domainTweets = append(domainTweets, &domain.Tweet{
			ID:        tweet.ID,
			AuthorID:  tweet.Edges.Author.ID,
			Content:   tweet.Content,
			CreatedAt: tweet.CreatedAt.Unix(),
		})
	}

	common.Logger().Debug("[TweetRepository] FindAllByMultipleUserIDs success", zap.Any("user_ids", userIDs), zap.Any("tweets", domainTweets))

	return domainTweets, nil
}
