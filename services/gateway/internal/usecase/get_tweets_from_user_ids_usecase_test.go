package usecase_test

import (
	"context"
	"testing"

	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestGetTweetsFromMultipleUserIDs_Success(t *testing.T) {
	// arrange
	uc := usecase.NewGetTweetsFromMultipleUserIDsUseCase(&mockTweetService{})
	// act
	tweets, err := uc.Execute(context.Background(), []int64{1, 2})
	// assert
	assert.NoError(t, err)
	assert.Len(t, tweets, 2)
	assert.Equal(t, int64(1), tweets[0].ID)
	assert.Equal(t, "User1 tweet", tweets[0].Content)
}

func TestGetTweetsFromMultipleUserIDs_Error(t *testing.T) {
	// arrange
	uc := usecase.NewGetTweetsFromMultipleUserIDsUseCase(&errorTweetService{})
	// act
	tweets, err := uc.Execute(context.Background(), []int64{1, 2})
	// assert
	assert.Error(t, err)
	assert.Nil(t, tweets)
	assert.EqualError(t, err, "some service error")
}
