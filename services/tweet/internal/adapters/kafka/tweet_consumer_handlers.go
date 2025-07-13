package kafka

import (
	"context"
	"errors"

	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"go.uber.org/zap"
)

func HandleTweetCreated(ctx context.Context, uc ports.CreateTweetUsecasePort, payload KafkaCreateTweetPayload) error {
	common.Logger().Debug("[Handler] HandleTweetCreated called", zap.Any("payload", payload))
	if uc == nil {
		return errors.New("TweetCreatedUC is nil")
	}
	return uc.CreateTweet(ctx, &domain.Tweet{
		AuthorID: payload.AuthorID,
		Content:  payload.Content,
	})
}

func HandleTweetDeleted(ctx context.Context, uc ports.DeleteTweetUsecasePort, payload KafkaDeleteTweetPayload) error {
	common.Logger().Debug("[Handler] HandleTweetDeleted called", zap.Any("payload", payload))
	if uc == nil {
		return errors.New("TweetDeletedUC is nil")
	}
	return uc.DeleteTweet(ctx, payload.TweetID)
}
