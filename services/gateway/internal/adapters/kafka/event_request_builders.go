package kafka

import (
	"fmt"
)

// Payload constrains
type AuthPayloadTypeConstraint interface {
	KafkaUserCreatePayload
}

type UserPayloadTypeConstraint interface {
	KafkaUserCreatedPayload | KafkaFollowPayload
}

type TweetPayloadTypeConstraint interface {
	KafkaTweetCreatePayload
}

// Payload interfaces
type AuthPayload interface {
	isAuthPayload()
}

func (KafkaUserCreatePayload) isAuthPayload() {}

type UserPayload interface {
	isUserPayload()
}

func (KafkaUserCreatedPayload) isUserPayload() {}
func (KafkaFollowPayload) isUserPayload()      {}

type TweetPayload interface {
	isTweetPayload()
}

func (KafkaTweetCreatePayload) isTweetPayload() {}

func AuthEventRequestBuilder(eventType string, payload AuthPayload) (KafkaEventRequest[AuthPayload], error) {
	switch eventType {
	case "users.create":
		return KafkaEventRequest[AuthPayload]{
			EventType: eventType,
			Payload:   payload,
		}, nil
	default:
		return KafkaEventRequest[AuthPayload]{}, fmt.Errorf("unsupported auth event type: %s", eventType)
	}
}

func UserEventRequestBuilder(eventType string, payload UserPayload) (KafkaEventRequest[UserPayload], error) {
	switch eventType {
	case "users.create", "users.created":
		if _, ok := payload.(KafkaUserCreatedPayload); !ok {
			return KafkaEventRequest[UserPayload]{}, fmt.Errorf("payload must be KafkaUserCreatedPayload for %s", eventType)
		}
	case "users.follow", "users.unfollow":
		if _, ok := payload.(KafkaFollowPayload); !ok {
			return KafkaEventRequest[UserPayload]{}, fmt.Errorf("payload must be KafkaFollowPayload for %s", eventType)
		}
	default:
		return KafkaEventRequest[UserPayload]{}, fmt.Errorf("unsupported user event type: %s", eventType)
	}
	return KafkaEventRequest[UserPayload]{
		EventType: eventType,
		Payload:   payload,
	}, nil
}

func TweetEventRequestBuilder(eventType string, payload TweetPayload) (KafkaEventRequest[TweetPayload], error) {
	switch eventType {
	case "tweets.create":
		return KafkaEventRequest[TweetPayload]{
			EventType: eventType,
			Payload:   payload,
		}, nil
	default:
		return KafkaEventRequest[TweetPayload]{}, fmt.Errorf("unsupported tweet event type: %s", eventType)
	}
}
