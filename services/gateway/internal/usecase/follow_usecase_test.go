package usecase_test

import (
	"context"
	"testing"

	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestFollowUser_Success(t *testing.T) {
	// arrange
	uc := usecase.NewFollowUseCase(&mockFollowProducer{})
	// act
	err := uc.Execute(context.Background(), 1, 2)
	// assert
	assert.NoError(t, err)
}

func TestFollowUser_ProducerError_ReturnsError(t *testing.T) {
	// arrange
	uc := usecase.NewFollowUseCase(&errorFollowProducer{})
	// act
	err := uc.Execute(context.Background(), 1, 2)
	// assert
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}
