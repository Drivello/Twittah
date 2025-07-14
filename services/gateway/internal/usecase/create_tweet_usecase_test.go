package usecase_test

import (
	"context"
	"testing"

	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestCreateTweet_EmptyContent_ReturnsError(t *testing.T) {
	// arrange
	uc := usecase.NewCreateTweetUseCase(&mockTweetProducer{})
	// act
	err := uc.Execute(context.Background(), 123, "")
	// assert
	assert.Error(t, err)
	assert.Equal(t, "tweet content must not be empty", err.Error())
}

func TestCreateTweet_ValidContent_Success(t *testing.T) {
	// arrange
	producer := &mockTweetProducer{}
	uc := usecase.NewCreateTweetUseCase(producer)
	// act
	err := uc.Execute(context.Background(), 123, "¡Hola mundo!")
	// assert
	assert.NoError(t, err)
}

func TestCreateTweet_ProducerError_ReturnsError(t *testing.T) {
	// arrange
	uc := usecase.NewCreateTweetUseCase(&errorTweetProducer{})
	// act
	err := uc.Execute(context.Background(), 123, "¡Otro tweet!")
	// assert
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}
