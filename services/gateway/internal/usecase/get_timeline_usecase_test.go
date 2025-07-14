package usecase_test

import (
	"context"
	"testing"

	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestGetTimeline_Success(t *testing.T) {
	// arrange
	uc := usecase.NewGetTimelineUseCase(&mockTweetService{})
	// act
	tweets, err := uc.Execute(context.Background(), 1)
	// assert
	assert.NoError(t, err)
	assert.Len(t, tweets, 2)
	assert.Equal(t, int64(1), tweets[0].ID)
	assert.Equal(t, "First tweet", tweets[0].Content)
}

func TestGetTimeline_Error(t *testing.T) {
	// arrange
	uc := usecase.NewGetTimelineUseCase(&errorTweetService{})
	// act
	tweets, err := uc.Execute(context.Background(), 1)
	// assert
	assert.Error(t, err)
	assert.Nil(t, tweets)
	assert.EqualError(t, err, "service unavailable")
}
