package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/dto"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"go.uber.org/zap"
)

type UserServiceHTTPAdapter struct {
	BaseURL string
}

func NewUserServiceHTTPAdapter(baseURL string) *UserServiceHTTPAdapter {
	return &UserServiceHTTPAdapter{BaseURL: baseURL}
}

func (a *UserServiceHTTPAdapter) GetFollowers(ctx context.Context, userID string) ([]*domain.User, error) {
	common.Logger().Debug("[Gateway] UserServiceHTTPAdapter.GetFollowers called", zap.String("user_id", userID), zap.String("url", a.BaseURL+"/followers/"+userID))
	req, err := http.NewRequestWithContext(ctx, "GET", a.BaseURL+"/followers/"+userID, nil)
	if err != nil {
		common.Logger().Debug("[Gateway] UserServiceHTTPAdapter.NewRequest error", zap.Error(err))
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		common.Logger().Debug("[Gateway] UserServiceHTTPAdapter.Do error", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		common.Logger().Debug("[Gateway] UserServiceHTTPAdapter non-OK status", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("failed to get followers: %s", resp.Status)
	}

	var followersDTO dto.GetFollowersResponse
	if err := json.NewDecoder(resp.Body).Decode(&followersDTO); err != nil {
		common.Logger().Debug("[Gateway] UserServiceHTTPAdapter decode error", zap.Error(err))
		return nil, err
	}
	// Map DTO to domain.User
	users := make([]*domain.User, len(followersDTO.Followers))
	for i, f := range followersDTO.Followers {
		users[i] = &domain.User{ID: f.ID, Username: f.Username}
	}
	common.Logger().Debug("[Gateway] UserServiceHTTPAdapter.GetFollowers success", zap.String("user_id", userID), zap.Any("followers", users))
	return users, nil
}
