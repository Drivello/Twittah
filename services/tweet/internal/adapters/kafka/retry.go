package kafka

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Drivello/Twittah/services/tweet/config"
	"go.uber.org/zap"
)

// retryWithTimeoutAndJitter retries a function using exponential backoff and jitter until timeout.
func retryWithTimeoutAndJitter(ctx context.Context, cfg config.RetryConfig, fn func() error, logger *zap.Logger) error {
	backoff := cfg.InitialBackoff
	rand.Seed(time.Now().UnixNano()) // for jitter
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry timed out after %v: %w", cfg.MaxRetryDuration, ctx.Err())
		default:
			err := fn()
			if err == nil {
				return nil // success
			}

			logger.Warn("Retryable error, will retry", zap.Error(err), zap.Duration("next_backoff", backoff))

			// Sleep with jitter (randomize between 50% and 150% of backoff)
			jitter := time.Duration(rand.Int63n(int64(backoff))) / 2
			sleepDuration := backoff + jitter
			time.Sleep(sleepDuration)

			// Exponential backoff up to max
			if backoff < cfg.MaxBackoff {
				backoff *= 2
				if backoff > cfg.MaxBackoff {
					backoff = cfg.MaxBackoff
				}
			}

			if time.Since(start) > cfg.MaxRetryDuration {
				return fmt.Errorf("retry timed out after %v", cfg.MaxRetryDuration)
			}
		}
	}
}
