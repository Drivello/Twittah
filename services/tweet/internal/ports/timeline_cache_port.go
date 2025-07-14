package ports

import (
	"context"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
)

// TimelineCachePort define el contrato para el cache de timelines.
type TimelineCachePort interface {
	GetTimeline(ctx context.Context, userID string) ([]domain.Tweet, error)
	SetTimeline(ctx context.Context, userID string, tweets []domain.Tweet) error
	InvalidateTimeline(ctx context.Context, userID string) error
}
