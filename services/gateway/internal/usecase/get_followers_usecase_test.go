package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/stretchr/testify/assert"
)

// Mock successful UserService
type mockUserService struct{}

func (m *mockUserService) GetFollowers(ctx context.Context, userID int64) ([]*domain.User, error) {
	return []*domain.User{
		{ID: 2, Username: "follower1"},
		{ID: 3, Username: "follower2"},
	}, nil
}

// Mock UserService that returns error
type errorUserService struct{}

func (m *errorUserService) GetFollowers(ctx context.Context, userID int64) ([]*domain.User, error) {
	return nil, errors.New("db error")
}

func TestGetFollowers_Success(t *testing.T) {
	// arrange
	uc := usecase.NewGetFollowersUseCase(&mockUserService{})
	// act
	followers, err := uc.Execute(context.Background(), 1)
	// assert
	assert.NoError(t, err)
	assert.Len(t, followers, 2)
	assert.Equal(t, int64(2), followers[0].ID)
	assert.Equal(t, "follower1", followers[0].Username)
}

func TestGetFollowers_Error(t *testing.T) {
	// arrange
	uc := usecase.NewGetFollowersUseCase(&errorUserService{})
	// act
	followers, err := uc.Execute(context.Background(), 1)
	// assert
	assert.Error(t, err)
	assert.Nil(t, followers)
	assert.EqualError(t, err, "db error")
}
