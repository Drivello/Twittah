package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/dto"
	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"go.uber.org/zap"
)

type UserServiceHTTPAdapter struct {
	BaseURL string
}

func NewUserServiceHTTPAdapter(baseURL string) *UserServiceHTTPAdapter {
	common.Logger().Debug("[Gateway] Creating UserServiceHTTPAdapter with url: ", zap.String("url", baseURL))
	return &UserServiceHTTPAdapter{BaseURL: baseURL}
}

func (a *UserServiceHTTPAdapter) GetFollowers(ctx context.Context, userID int64) ([]*domain.User, error) {

	_, err := common.ValidatePositiveInt(userID)
	if err != nil {
		return nil, err
	}

	common.Logger().Debug("[Gateway] UserServiceHTTPAdapter.GetFollowers called", zap.Int64("user_id", userID), zap.String("url", a.BaseURL+"/followers/"+strconv.FormatInt(userID, 10)))
	req, err := http.NewRequestWithContext(ctx, "GET", a.BaseURL+"/followers/"+strconv.FormatInt(userID, 10), nil)
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

	if err := followersDTO.Validate(); err != nil {
		common.Logger().Debug("[Gateway] UserServiceHTTPAdapter validate GetFollowersResponse error", zap.Error(err))
		return nil, err
	}

	users := make([]*domain.User, len(followersDTO.Followers))
	for i, f := range followersDTO.Followers {
		users[i] = &domain.User{ID: f.ID, Username: f.Username}
	}

	common.Logger().Debug("[Gateway] UserServiceHTTPAdapter.GetFollowers success", zap.Int64("user_id", userID), zap.Any("followers", users))
	return users, nil
}
