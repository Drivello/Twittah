package usecase_test

import (
	"context"
	"testing"

	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestDeleteTweet_Success(t *testing.T) {
	// arrange
	uc := usecase.NewDeleteTweetUseCase(&mockTweetProducer{})
	// act
	err := uc.Execute(context.Background(), 1)
	// assert
	assert.NoError(t, err)
}

func TestDeleteTweet_ProducerError_ReturnsError(t *testing.T) {
	// arrange
	uc := usecase.NewDeleteTweetUseCase(&errorTweetProducer{})
	// act
	err := uc.Execute(context.Background(), 1)
	// assert
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}
