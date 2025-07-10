package ports

import "context"

type AuthEventProducerPort interface {
	PublishUserCreateRequest(ctx context.Context, eventBytes []byte) error
}
