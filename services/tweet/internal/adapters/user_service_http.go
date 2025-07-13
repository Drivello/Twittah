package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"go.uber.org/zap"
)

type UserServiceHTTPAdapter struct {
	BaseURL string
}

type GetFollowingResponse struct {
	Following []*domain.User `json:"following"`
}

func NewUserServiceHTTPAdapter(baseURL string) *UserServiceHTTPAdapter {
	common.Logger().Debug("[TweetService] Creating UserServiceHTTPAdapter with url: ", zap.String("url", baseURL))
	return &UserServiceHTTPAdapter{BaseURL: baseURL}
}

func (a *UserServiceHTTPAdapter) GetFollowing(ctx context.Context, userID int64) ([]*domain.User, error) {
	common.Logger().Debug("[TweetService] UserServiceHTTPAdapter.GetFollowing called", zap.Int64("user_id", userID), zap.String("url", a.BaseURL+"/following/"+strconv.FormatInt(userID, 10)))
	req, err := http.NewRequestWithContext(ctx, "GET", a.BaseURL+"/following/"+strconv.FormatInt(userID, 10), nil)
	if err != nil {
		common.Logger().Debug("[TweetService] UserServiceHTTPAdapter.NewRequest error", zap.Error(err))
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		common.Logger().Debug("[TweetService] UserServiceHTTPAdapter.Do error", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		common.Logger().Debug("[TweetService] UserServiceHTTPAdapter non-OK status", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("failed to get following: %s", resp.Status)
	}

	var followingDTO GetFollowingResponse
	if err := json.NewDecoder(resp.Body).Decode(&followingDTO); err != nil {
		common.Logger().Debug("[TweetService] UserServiceHTTPAdapter decode error", zap.Error(err))
		return nil, err
	}

	users := make([]*domain.User, len(followingDTO.Following))
	for i, f := range followingDTO.Following {
		users[i] = &domain.User{ID: f.ID, Username: f.Username}
	}
	common.Logger().Debug("[TweetService] UserServiceHTTPAdapter.GetFollowing success", zap.Int64("user_id", userID), zap.Any("following", users))
	return users, nil
}
