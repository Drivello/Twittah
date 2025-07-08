package kafka

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/Drivello/Twittah/services/user/internal/domain"
	"github.com/Drivello/Twittah/services/user/internal/usecase"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type UserEventConsumer struct {
	group      sarama.ConsumerGroup
	topic      string
	followUC   *usecase.FollowUserUseCase
	unfollowUC *usecase.UnfollowUserUseCase
}

func NewUserEventConsumer(brokers []string, groupID, topic string, followUC *usecase.FollowUserUseCase, unfollowUC *usecase.UnfollowUserUseCase) (*UserEventConsumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}
	return &UserEventConsumer{
		group:      group,
		topic:      topic,
		followUC:   followUC,
		unfollowUC: unfollowUC,
	}, nil
}

type FollowEvent struct {
	EventType  string `json:"event_type"`
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}

func (c *UserEventConsumer) Start(ctx context.Context) error {
	consumer := &userEventHandler{followUC: c.followUC, unfollowUC: c.unfollowUC}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm
		cancel()
	}()
	for {
		if err := c.group.Consume(ctx, []string{c.topic}, consumer); err != nil {
			zap.S().Errorw("Kafka consume error", "error", err)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type userEventHandler struct {
	followUC   *usecase.FollowUserUseCase
	unfollowUC *usecase.UnfollowUserUseCase
}

func (h *userEventHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *userEventHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *userEventHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event domain.FollowEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			zap.S().Errorw("Failed to unmarshal follow event", "error", err)
			sess.MarkMessage(msg, "")
			continue
		}
		switch event.EventType {
		case "follow_created":
			if err := h.followUC.Execute(context.Background(), event.FollowerID, event.FolloweeID); err != nil {
				zap.S().Errorw("Failed to process follow_created", "error", err)
			}
		case "follow_deleted":
			if err := h.unfollowUC.Execute(context.Background(), event.FollowerID, event.FolloweeID); err != nil {
				zap.S().Errorw("Failed to process follow_deleted", "error", err)
			}
		default:
			zap.S().Warnw("Unknown event type", "type", event.EventType)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
