package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserServiceHTTPAdapter struct {
	BaseURL string
}

func NewUserServiceHTTPAdapter(baseURL string) *UserServiceHTTPAdapter {
	return &UserServiceHTTPAdapter{BaseURL: baseURL}
}

func (a *UserServiceHTTPAdapter) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.BaseURL+"/followers/"+userID, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get followers: %s", resp.Status)
	}

	var followers []string
	if err := json.NewDecoder(resp.Body).Decode(&followers); err != nil {
		return nil, err
	}

	return followers, nil
}
