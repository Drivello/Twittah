package kafka

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/Drivello/Twittah/services/auth/ent"
	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/domain"
	"go.uber.org/zap"
)

func RetryWithTimeoutAndJitter(ctx context.Context, cfg config.RetryConfig, fn func() error) error {
	backoff := cfg.InitialBackoff

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry timed out after %v: %w", cfg.MaxRetryDuration, ctx.Err())
		default:
			err := fn()
			if err == nil {
				return nil
			}

			if !IsRetryableError(err) {
				common.Logger().Error("Non-retryable error detected. Will not retry.", zap.Error(err))
				return err
			}

			common.Logger().Info("Retryable error, retrying.", zap.Error(err), zap.Duration("next_backoff", backoff))

			jitter := time.Duration(rand.Int63n(int64(backoff))) / 2
			sleepDuration := backoff + jitter
			time.Sleep(sleepDuration)

			if backoff < cfg.MaxBackoff {
				backoff *= 2
				if backoff > cfg.MaxBackoff {
					backoff = cfg.MaxBackoff
				}
			}
		}
	}
}

func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	switch err {
	case domain.ErrUserAlreadyExists,
		domain.ErrInvalidPayload:
		return false
	}
	if ent.IsConstraintError(err) {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	if strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "connection refused") ||
		strings.Contains(errMsg, "server is closed") {
		return true
	}
	errMsgNum, err := strconv.Atoi(strings.Split(errMsg, " ")[0])
	if err == nil && errMsgNum > 500 {
		return true
	}
	return false
}
