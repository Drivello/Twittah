package usecase_test

import (
	"context"
	"errors"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Mock successful producer
type mockTweetProducer struct{}

func (m *mockTweetProducer) PublishEvent(event kafka.KafkaEventRequest[kafka.TweetPayload]) error {
	return nil
}

// Mock producer that returns error
type errorTweetProducer struct{}

func (m *errorTweetProducer) PublishEvent(event kafka.KafkaEventRequest[kafka.TweetPayload]) error {
	return assert.AnError
}

// Mock successful producer
type mockFollowProducer struct{}

func (m *mockFollowProducer) PublishEvent(event kafka.KafkaEventRequest[kafka.UserPayload]) error {
	return nil
}

// Mock producer that returns error
type errorFollowProducer struct{}

func (m *errorFollowProducer) PublishEvent(event kafka.KafkaEventRequest[kafka.UserPayload]) error {
	return assert.AnError
}

// Mock successful producer
type mockAuthProducer struct{}

func (m *mockAuthProducer) PublishEvent(event kafka.KafkaEventRequest[kafka.AuthPayload]) error {
	return nil
}

// Mock producer that returns error
type errorAuthProducer struct{}

func (m *errorAuthProducer) PublishEvent(event kafka.KafkaEventRequest[kafka.AuthPayload]) error {
	return assert.AnError
}

// Mock successful TweetService
type mockTweetService struct{}

func (m *mockTweetService) GetTweetsFromUserID(context.Context, int64) ([]*domain.Tweet, error) {
	return []*domain.Tweet{
		{ID: 1, Content: "User tweet 1"},
		{ID: 2, Content: "User tweet 2"},
	}, nil
}

func (m *mockTweetService) GetTweetsFromMultipleUserIDs(context.Context, []int64) ([]*domain.Tweet, error) {
	return []*domain.Tweet{
		{ID: 1, Content: "User1 tweet"},
		{ID: 2, Content: "User2 tweet"},
	}, nil
}

func (m *mockTweetService) GetTimeline(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
	return []*domain.Tweet{
		{ID: 1, Content: "First tweet"},
		{ID: 2, Content: "Second tweet"},
	}, nil
}

// Mock TweetService that returns error
type errorTweetService struct{}

func (m *errorTweetService) GetTweetsFromUserID(context.Context, int64) ([]*domain.Tweet, error) {
	return nil, errors.New("db error")
}

func (m *errorTweetService) GetTweetsFromMultipleUserIDs(context.Context, []int64) ([]*domain.Tweet, error) {
	return nil, errors.New("some service error")
}

func (m *errorTweetService) GetTimeline(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
	return nil, errors.New("service unavailable")
}
