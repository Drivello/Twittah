package adapters_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters"
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/dto"
	"github.com/stretchr/testify/assert"
)

// ------ HELPERS ------

func newTestServer(status int, body interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status) // no error to check
		if body != nil {
			b, _ := json.Marshal(body)
			_, _ = w.Write(b)
		}
	}))
}

func shutdown(ts *httptest.Server) {
	if ts != nil {
		ts.Close()
	}
}

// ------ TESTS ------

func TestGetTimeline_Success(t *testing.T) {
	// arrange
	resp := dto.GetTimelineResponse{
		Timeline: []dto.TweetDTO{
			{ID: 1, Author: 2, Content: "hi", CreatedAt: time.Now()},
		},
	}
	ts := newTestServer(200, resp)
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTimeline(context.Background(), 7)

	// assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), result[0].ID)
}

func TestGetTimeline_BadStatus(t *testing.T) {
	// arrange
	ts := newTestServer(404, nil)
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTimeline(context.Background(), 1)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTimeline_DecodeError(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("{invalid-json"))
	}))
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTimeline(context.Background(), 2)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTimeline_ValidateError(t *testing.T) {
	// arrange
	badResp := dto.GetTimelineResponse{
		Timeline: []dto.TweetDTO{
			{ID: 0, Author: 2, Content: "x", CreatedAt: time.Now()},
		},
	}
	ts := newTestServer(200, badResp)
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTimeline(context.Background(), 3)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTimeline_RequestError(t *testing.T) {
	// arrange
	adapter := adapters.NewTweetServiceHTTPAdapter("http://127.0.0.1:0") // invalid port, will fail to connect

	// act
	result, err := adapter.GetTimeline(context.Background(), 123)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTimeline_InvalidID(t *testing.T) {
	// arrange
	adapter := adapters.NewTweetServiceHTTPAdapter("http://127.0.0.1:1234")

	// act
	result, err := adapter.GetTimeline(context.Background(), 0) // not positive

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTweetsFromUserID_Success(t *testing.T) {
	// arrange
	resp := dto.GetTweetsFromUserIDResponse{
		Tweets: []dto.TweetDTO{
			{ID: 44, Author: 5, Content: "hello", CreatedAt: time.Now()},
		},
	}
	ts := newTestServer(200, resp)
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTweetsFromUserID(context.Background(), 12)

	// assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(44), result[0].ID)
}

func TestGetTweetsFromUserID_BadStatus(t *testing.T) {
	// arrange
	ts := newTestServer(403, nil)
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTweetsFromUserID(context.Background(), 1)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTweetsFromUserID_DecodeError(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = io.WriteString(w, "not-json")
	}))
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTweetsFromUserID(context.Background(), 5)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTweetsFromUserID_ValidateError(t *testing.T) {
	// arrange
	badResp := dto.GetTweetsFromUserIDResponse{
		Tweets: []dto.TweetDTO{
			{ID: -1, Author: 2, Content: "xx", CreatedAt: time.Now()},
		},
	}
	ts := newTestServer(200, badResp)
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTweetsFromUserID(context.Background(), 4)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTweetsFromUserID_RequestError(t *testing.T) {
	// arrange
	adapter := adapters.NewTweetServiceHTTPAdapter("http://127.0.0.1:0")

	// act
	result, err := adapter.GetTweetsFromUserID(context.Background(), 123)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTweetsFromUserID_InvalidID(t *testing.T) {
	// arrange
	adapter := adapters.NewTweetServiceHTTPAdapter("http://localhost:1111")

	// act
	result, err := adapter.GetTweetsFromUserID(context.Background(), -99)

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTweetsFromMultipleUserIDs_Success(t *testing.T) {
	// arrange
	resp := dto.GetTweetsFromMultipleUserIDsResponse{
		Tweets: []dto.TweetDTO{
			{ID: 99, Author: 1, Content: "multi", CreatedAt: time.Now()},
		},
	}
	ts := newTestServer(200, resp)
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTweetsFromMultipleUserIDs(context.Background(), []int64{111})

	// assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(99), result[0].ID)
}

func TestGetTweetsFromMultipleUserIDs_BadStatus(t *testing.T) {
	// arrange
	ts := newTestServer(418, nil)
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTweetsFromMultipleUserIDs(context.Background(), []int64{1, 2})

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTweetsFromMultipleUserIDs_DecodeError(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("broken"))
	}))
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTweetsFromMultipleUserIDs(context.Background(), []int64{1, 2, 3})

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTweetsFromMultipleUserIDs_ValidateError(t *testing.T) {
	// arrange
	badResp := dto.GetTweetsFromMultipleUserIDsResponse{
		Tweets: []dto.TweetDTO{
			{ID: 0, Author: 2, Content: "multi", CreatedAt: time.Now()},
		},
	}
	ts := newTestServer(200, badResp)
	defer shutdown(ts)
	adapter := adapters.NewTweetServiceHTTPAdapter(ts.URL)

	// act
	result, err := adapter.GetTweetsFromMultipleUserIDs(context.Background(), []int64{3, 5})

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTweetsFromMultipleUserIDs_RequestError(t *testing.T) {
	// arrange
	adapter := adapters.NewTweetServiceHTTPAdapter("http://127.0.0.1:0")

	// act
	result, err := adapter.GetTweetsFromMultipleUserIDs(context.Background(), []int64{8, 9})

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTweetsFromMultipleUserIDs_InvalidIDs(t *testing.T) {
	// arrange
	adapter := adapters.NewTweetServiceHTTPAdapter("http://localhost:1234")

	// act
	result, err := adapter.GetTweetsFromMultipleUserIDs(context.Background(), []int64{0, -1, 2})

	// assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
