package kafka

import (
	"context"
	"errors"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"go.uber.org/zap"
)

func HandleTweetCreated(ctx context.Context, uc ports.CreateTweetUsecasePort, payload KafkaCreateTweetPayload) error {
	common.Logger().Debug("[Handler] HandleTweetCreated called", zap.Any("payload", payload))
	if uc == nil {
		return errors.New("TweetCreatedUC is nil")
	}
	return uc.Execute(ctx, &domain.Tweet{
		AuthorID: payload.AuthorID,
		Content:  payload.Content,
	})
}

func HandleTweetDeleted(ctx context.Context, uc ports.DeleteTweetUsecasePort, payload KafkaDeleteTweetPayload) error {
	common.Logger().Debug("[Handler] HandleTweetDeleted called", zap.Any("payload", payload))
	if uc == nil {
		return errors.New("TweetDeletedUC is nil")
	}
	return uc.Execute(ctx, payload.TweetID)
}
