package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UserQueryClient implementa ports.UserQueryPort para consultar followers vía HTTP
// al microservicio de usuarios.
type UserQueryClient struct {
	BaseURL string
}

func NewUserQueryClient(baseURL string) *UserQueryClient {
	return &UserQueryClient{BaseURL: baseURL}
}

func (c *UserQueryClient) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	url := fmt.Sprintf("%s/followers/%s", c.BaseURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var dto struct {
		Followers []string `json:"followers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, err
	}
	return dto.Followers, nil
}
