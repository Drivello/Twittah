package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/dto"
	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"go.uber.org/zap"
)

type TweetServiceHTTPAdapter struct {
	BaseURL string
}

func NewTweetServiceHTTPAdapter(baseURL string) *TweetServiceHTTPAdapter {
	common.Logger().Debug("[Gateway] Creating TweetServiceHTTPAdapter with url: ", zap.String("url", baseURL))
	return &TweetServiceHTTPAdapter{BaseURL: baseURL}
}

func (a *TweetServiceHTTPAdapter) GetTimeline(ctx context.Context, userID int64) ([]*domain.Tweet, error) {

	common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.GetTimeline called", zap.Int64("user_id", userID), zap.String("url", a.BaseURL+"/timeline/"+strconv.FormatInt(userID, 10)))

	userIDInt, err := common.ValidatePositiveInt(userID)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", a.BaseURL+"/timeline/"+strconv.FormatInt(userIDInt, 10), nil)
	if err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.NewRequest error", zap.Error(err))
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.Do error", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter non-OK status", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("failed to get timeline: %s", resp.Status)
	}

	var timelineDTO dto.GetTimelineResponse
	if err := json.NewDecoder(resp.Body).Decode(&timelineDTO); err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter decode error", zap.Error(err))
		return nil, err
	}

	if err := timelineDTO.Validate(); err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter validate GetTimelineResponse error", zap.Error(err))
		return nil, err
	}

	tweets := make([]*domain.Tweet, len(timelineDTO.Timeline))
	for i, t := range timelineDTO.Timeline {
		tweets[i] = &domain.Tweet{ID: t.ID, Author: t.Author, Content: t.Content, CreatedAt: t.CreatedAt}
	}
	common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.GetTimeline success", zap.Int64("user_id", userID), zap.Any("timeline", tweets))
	return tweets, nil
}

func (a *TweetServiceHTTPAdapter) GetTweetsFromUserID(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
	common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.GetTweetsFromUserID called", zap.Int64("user_id", userID), zap.String("url", a.BaseURL+"/tweets/"+strconv.FormatInt(userID, 10)))

	userIDInt, err := common.ValidatePositiveInt(userID)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", a.BaseURL+"/tweets/"+strconv.FormatInt(userIDInt, 10), nil)
	if err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.NewRequest error", zap.Error(err))
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.Do error", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter non-OK status", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("failed to get tweets: %s", resp.Status)
	}

	var tweetsDTO dto.GetTweetsFromUserIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&tweetsDTO); err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter decode error", zap.Error(err))
		return nil, err
	}
	if err := tweetsDTO.Validate(); err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter validate GetTweetsFromUserIDResponse error", zap.Error(err))
		return nil, err
	}

	tweets := make([]*domain.Tweet, len(tweetsDTO.Tweets))
	for i, t := range tweetsDTO.Tweets {
		tweets[i] = &domain.Tweet{ID: t.ID, Author: t.Author, Content: t.Content, CreatedAt: t.CreatedAt}
	}

	common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.GetTweetsFromUserID success", zap.Int64("user_id", userID), zap.Any("tweets", tweets))
	return tweets, nil
}

func (a *TweetServiceHTTPAdapter) GetTweetsFromMultipleUserIDs(ctx context.Context, userIDs []int64) ([]*domain.Tweet, error) {

	strIDs := make([]string, len(userIDs))
	for i, id := range userIDs {
		strIDs[i] = strconv.FormatInt(id, 10)
	}

	err := common.ValidateIDList(userIDs)
	if err != nil {
		return nil, err
	}

	userIDsQuery := "user_ids=" + strings.Join(strIDs, ",")

	common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.GetTweetsFromMultipleUserIDs called", zap.Any("user_ids", userIDs), zap.String("url", a.BaseURL+"/tweets/"+userIDsQuery))

	req, err := http.NewRequestWithContext(ctx, "GET", a.BaseURL+"/tweets?"+userIDsQuery, nil)
	if err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.NewRequest error", zap.Error(err))
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.Do error", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter non-OK status", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("failed to get tweets: %s", resp.Status)
	}

	var tweetsDTO dto.GetTweetsFromMultipleUserIDsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tweetsDTO); err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter decode error", zap.Error(err))
		return nil, err
	}
	if err := tweetsDTO.Validate(); err != nil {
		common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter validate GetTweetsFromMultipleUserIDsResponse error", zap.Error(err))
		return nil, err
	}

	tweets := make([]*domain.Tweet, len(tweetsDTO.Tweets))
	for i, t := range tweetsDTO.Tweets {
		tweets[i] = &domain.Tweet{ID: t.ID, Author: t.Author, Content: t.Content, CreatedAt: t.CreatedAt}
	}

	common.Logger().Debug("[Gateway] TweetServiceHTTPAdapter.GetTweetsFromMultipleUserIDs success", zap.Any("user_ids", userIDs), zap.Any("tweets", tweets))
	return tweets, nil
}
