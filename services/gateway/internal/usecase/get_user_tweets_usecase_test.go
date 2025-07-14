package usecase_test

import (
	"context"
	"testing"

	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestGetUserTweets_Success(t *testing.T) {
	// arrange
	uc := usecase.NewGetUserTweetsUseCase(&mockTweetService{})
	// act
	tweets, err := uc.Execute(context.Background(), 1)
	// assert
	assert.NoError(t, err)
	assert.Len(t, tweets, 2)
	assert.Equal(t, int64(1), tweets[0].ID)
	assert.Equal(t, "User tweet 1", tweets[0].Content)
}

func TestGetUserTweets_Error(t *testing.T) {
	// arrange
	uc := usecase.NewGetUserTweetsUseCase(&errorTweetService{})
	// act
	tweets, err := uc.Execute(context.Background(), 1)
	// assert
	assert.Error(t, err)
	assert.Nil(t, tweets)
	assert.EqualError(t, err, "db error")
}
