package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Save(ctx context.Context, tweet *domain.Tweet) error {
	args := m.Called(ctx, tweet)
	return args.Error(0)
}
func (m *mockRepo) GetRecentTweets(ctx context.Context, userIDs []string, limit int) ([]domain.Tweet, error) {
	return nil, nil
}

type mockCache struct{ mock.Mock }

func (m *mockCache) AddTweetToTimelines(ctx context.Context, tweet domain.Tweet, followerIDs []string) error {
	args := m.Called(ctx, tweet, followerIDs)
	return args.Error(0)
}
func (m *mockCache) GetTimeline(ctx context.Context, userID string, limit int) ([]domain.Tweet, error) {
	return nil, nil
}
func (m *mockCache) RemoveTweetsFromTimeline(ctx context.Context, userID, authorID string) error {
	return nil
}

func TestPublishTweet_Execute_Success(t *testing.T) {
	repo := new(mockRepo)
	cache := new(mockCache)
	uc := NewPublishTweet(repo, cache)

	tweet := &domain.Tweet{ID: "t1", AuthorID: "u1", Content: "Hola mundo!"}
	followerIDs := []string{"f1", "f2"}

	repo.On("Save", mock.Anything, tweet).Return(nil)
	cache.On("AddTweetToTimelines", mock.Anything, *tweet, followerIDs).Return(nil)

	err := uc.Execute(context.Background(), tweet, followerIDs)
	assert.NoError(t, err)
}

func TestPublishTweet_Execute_SaveFail(t *testing.T) {
	repo := new(mockRepo)
	cache := new(mockCache)
	uc := NewPublishTweet(repo, cache)

	tweet := &domain.Tweet{ID: "t1", AuthorID: "u1", Content: "Hola mundo!"}
	followerIDs := []string{"f1"}

	repo.On("Save", mock.Anything, tweet).Return(errors.New("db fail"))

	err := uc.Execute(context.Background(), tweet, followerIDs)
	assert.Error(t, err)
}
