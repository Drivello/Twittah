package usecase

import (
	"context"
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

type TweetUsecase struct {
	Publisher ports.TweetPublisherPort
}

func NewTweetUsecase(publisher ports.TweetPublisherPort) *TweetUsecase {
	return &TweetUsecase{Publisher: publisher}
}

type TweetInput struct {
	AuthorID string
	Content  string
}

func (uc *TweetUsecase) PublishTweet(ctx context.Context, input TweetInput) error {
	event := kafka.TweetPublishedEventDTO{
		ID:        "",
		AuthorID:  input.AuthorID,
		Content:   input.Content,
		CreatedAt: "",
	}
	return uc.Publisher.PublishTweet(ctx, event)
}
