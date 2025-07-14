package usecase_test

import (
	"context"
	"testing"

	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestUnfollowUser_Success(t *testing.T) {
	// arrange
	uc := usecase.NewUnfollowUseCase(&mockFollowProducer{})
	// act
	err := uc.Execute(context.Background(), 1, 2)
	// assert
	assert.NoError(t, err)
}

func TestUnfollowUser_ProducerError_ReturnsError(t *testing.T) {
	// arrange
	uc := usecase.NewUnfollowUseCase(&errorFollowProducer{})
	// act
	err := uc.Execute(context.Background(), 1, 2)
	// assert
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}
