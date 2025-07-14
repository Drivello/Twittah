package adapters_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters"
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/dto"
	"github.com/stretchr/testify/assert"
)

func newUserFollowersTestServer(status int, body interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status) // no error to check
		if body != nil {
			b, _ := json.Marshal(body)
			_, _ = w.Write(b)
		}
	}))
}

func TestGetFollowers_Success(t *testing.T) {
	// arrange
	resp := dto.GetFollowersResponse{
		Followers: []dto.FollowerUserDTO{
			{ID: 1, Username: "alice"},
			{ID: 2, Username: "bob"},
		},
	}
	ts := newUserFollowersTestServer(200, resp)
	defer shutdown(ts)
	adapter := adapters.NewUserServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetFollowers(context.Background(), 999)

	// assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, "alice", result[0].Username)
	assert.Equal(t, int64(2), result[1].ID)
	assert.Equal(t, "bob", result[1].Username)
}

func TestGetFollowers_NonOKStatus(t *testing.T) {
	// arrange
	ts := newUserFollowersTestServer(404, nil)
	defer shutdown(ts)
	adapter := adapters.NewUserServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetFollowers(context.Background(), 123)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetFollowers_DecodeError(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200) // no error to check
		_, _ = io.WriteString(w, "{not-json")
	}))
	defer shutdown(ts)
	adapter := adapters.NewUserServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetFollowers(context.Background(), 123)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetFollowers_RequestError(t *testing.T) {
	// arrange
	adapter := adapters.NewUserServiceHTTPAdapter("http://127.0.0.1:0")

	// act
	result, err := adapter.GetFollowers(context.Background(), 123)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetFollowers_InvalidUserID(t *testing.T) {
	// arrange
	adapter := adapters.NewUserServiceHTTPAdapter("http://localhost")

	// act
	result, err := adapter.GetFollowers(context.Background(), -7)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
