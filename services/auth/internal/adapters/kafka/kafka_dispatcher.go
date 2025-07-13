package kafka

import (
	"context"
	"encoding/json"

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
		common.Logger().Error("[KafkaEventDispatcher] Invalid event message. Skipping.", zap.Error(err))
		return nil
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
		return nil

	default:
		common.Logger().Error("[KafkaEventDispatcher] Unknown event type. Skipping.", zap.String("event_type", req.EventType))
		return nil
	}
}
