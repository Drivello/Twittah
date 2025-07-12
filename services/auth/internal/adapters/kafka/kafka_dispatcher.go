package kafka

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
	"go.uber.org/zap"
)

type KafkaEventDispatcher struct {
	registerUserUseCasesPort ports.RegisterUserUseCasesPort
}

func NewKafkaEventDispatcher(authUC ports.RegisterUserUseCasesPort) *KafkaEventDispatcher {
	return &KafkaEventDispatcher{
		registerUserUseCasesPort: authUC,
	}
}

func (d *KafkaEventDispatcher) Dispatch(ctx context.Context, data []byte) *KafkaEventError {

	var req KafkaEventRequest
	if err := json.Unmarshal(data, &req); err != nil {
		common.Logger().Error("[KafkaEventDispatcher] Invalid event message, skipping", zap.Error(err))
		return &KafkaEventError{
			EventType: "invalid.event",
			Error:     err,
		}
	}

	switch req.EventType {
	case "users.create":

		var userReq KafkaUserCreateRequest
		if err := json.Unmarshal(data, &userReq); err != nil {
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}
		err := HandleUserCreate(ctx, d.registerUserUseCasesPort, userReq.Payload)
		if err != nil {
			return &KafkaEventError{
				EventType: req.EventType,
				Error:     err,
			}
		}

	default:
		common.Logger().Error("[KafkaEventDispatcher] Unknown event type", zap.String("event_type", req.EventType))
		return &KafkaEventError{
			EventType: "unknown.event",
			Error:     errors.New("unknown event type"),
		}
	}
	return nil

}
