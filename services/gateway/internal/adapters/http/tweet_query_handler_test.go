package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	httpg "github.com/Drivello/Twittah/services/gateway/internal/adapters/http"
	"github.com/Drivello/Twittah/services/gateway/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetTimeline_InvalidUserID_Returns400(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockTimelineUC := &mockGetTimelineUseCase{}
	mockFromIDsUC := &mockGetTweetsFromIDsUseCase{}
	mockFromUserUC := &mockGetTweetsFromUserIDUseCase{}
	handler := httpg.NewTweetQueryHandler(mockTimelineUC, mockFromIDsUC, mockFromUserUC)
	router := gin.Default()
	router.GET("/timeline/:user_id", handler.GetTimeline)
	req := httptest.NewRequest(http.MethodGet, "/timeline/abc", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 400, rec.Code)
	assert.JSONEq(t, `{"timeline": null}`, rec.Body.String())
}

func TestGetTimeline_UseCase_Error_Returns502(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockTimelineUC := &mockGetTimelineUseCase{
		ExecuteFunc: func(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
			return nil, errors.New("something failed")
		},
	}
	mockFromIDsUC := &mockGetTweetsFromIDsUseCase{}
	mockFromUserUC := &mockGetTweetsFromUserIDUseCase{}
	handler := httpg.NewTweetQueryHandler(mockTimelineUC, mockFromIDsUC, mockFromUserUC)
	router := gin.Default()
	router.GET("/timeline/:user_id", handler.GetTimeline)
	req := httptest.NewRequest(http.MethodGet, "/timeline/123", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 502, rec.Code)
	assert.JSONEq(t, `{"timeline": null}`, rec.Body.String())
}

func TestGetTimeline_Success(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockTimelineUC := &mockGetTimelineUseCase{
		ExecuteFunc: func(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
			return []*domain.Tweet{{ID: 1, Author: 2, Content: "ok"}}, nil
		},
	}
	mockFromIDsUC := &mockGetTweetsFromIDsUseCase{}
	mockFromUserUC := &mockGetTweetsFromUserIDUseCase{}
	handler := httpg.NewTweetQueryHandler(mockTimelineUC, mockFromIDsUC, mockFromUserUC)
	router := gin.Default()
	router.GET("/timeline/:user_id", handler.GetTimeline)
	req := httptest.NewRequest(http.MethodGet, "/timeline/123", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), `"timeline"`)
	assert.Contains(t, rec.Body.String(), `"id":1`)
}

func TestGetTweetsFromUserID_InvalidUserID_Returns400(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockTimelineUC := &mockGetTimelineUseCase{}
	mockFromIDsUC := &mockGetTweetsFromIDsUseCase{}
	mockFromUserUC := &mockGetTweetsFromUserIDUseCase{}
	handler := httpg.NewTweetQueryHandler(mockTimelineUC, mockFromIDsUC, mockFromUserUC)
	router := gin.Default()
	router.GET("/:user_id", handler.GetTweetsFromUserID)
	req := httptest.NewRequest(http.MethodGet, "/abc", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error":"invalid user_id: abc"`)
}

func TestGetTweetsFromUserID_UseCase_Error_Returns502(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockTimelineUC := &mockGetTimelineUseCase{}
	mockFromIDsUC := &mockGetTweetsFromIDsUseCase{}
	mockFromUserUC := &mockGetTweetsFromUserIDUseCase{
		ExecuteFunc: func(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
			return nil, errors.New("something failed")
		},
	}
	handler := httpg.NewTweetQueryHandler(mockTimelineUC, mockFromIDsUC, mockFromUserUC)
	router := gin.Default()
	router.GET("/:user_id", handler.GetTweetsFromUserID)
	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 502, rec.Code)
	assert.JSONEq(t, `{"tweets": null}`, rec.Body.String())
}

func TestGetTweetsFromUserID_Success(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockTimelineUC := &mockGetTimelineUseCase{}
	mockFromIDsUC := &mockGetTweetsFromIDsUseCase{}
	mockFromUserUC := &mockGetTweetsFromUserIDUseCase{
		ExecuteFunc: func(ctx context.Context, userID int64) ([]*domain.Tweet, error) {
			return []*domain.Tweet{{ID: 1, Author: 2, Content: "ok"}}, nil
		},
	}
	handler := httpg.NewTweetQueryHandler(mockTimelineUC, mockFromIDsUC, mockFromUserUC)
	router := gin.Default()
	router.GET("/:user_id", handler.GetTweetsFromUserID)
	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), `"tweets"`)
	assert.Contains(t, rec.Body.String(), `"id":1`)
}

func TestGetTweetsFromIDs_NoIDs_Returns400(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockTimelineUC := &mockGetTimelineUseCase{}
	mockFromIDsUC := &mockGetTweetsFromIDsUseCase{}
	mockFromUserUC := &mockGetTweetsFromUserIDUseCase{}
	handler := httpg.NewTweetQueryHandler(mockTimelineUC, mockFromIDsUC, mockFromUserUC)
	router := gin.Default()
	router.GET("/", handler.GetTweetsFromIDs)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error":"invalid user_ids: "`)
}

func TestGetTweetsFromIDs_UseCase_Error_Returns502(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockTimelineUC := &mockGetTimelineUseCase{}
	mockFromIDsUC := &mockGetTweetsFromIDsUseCase{
		ExecuteFunc: func(ctx context.Context, ids []int64) ([]*domain.Tweet, error) {
			return nil, errors.New("something failed")
		},
	}
	mockFromUserUC := &mockGetTweetsFromUserIDUseCase{}
	handler := httpg.NewTweetQueryHandler(mockTimelineUC, mockFromIDsUC, mockFromUserUC)
	router := gin.Default()
	router.GET("/", handler.GetTweetsFromIDs)
	req := httptest.NewRequest(http.MethodGet, "/?user_ids=1,2", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 502, rec.Code)
	assert.JSONEq(t, `{"tweets": null}`, rec.Body.String())
}

func TestGetTweetsFromIDs_Success(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockTimelineUC := &mockGetTimelineUseCase{}
	mockFromIDsUC := &mockGetTweetsFromIDsUseCase{
		ExecuteFunc: func(ctx context.Context, ids []int64) ([]*domain.Tweet, error) {
			return []*domain.Tweet{{ID: 1, Author: 2, Content: "ok"}}, nil
		},
	}
	mockFromUserUC := &mockGetTweetsFromUserIDUseCase{}
	handler := httpg.NewTweetQueryHandler(mockTimelineUC, mockFromIDsUC, mockFromUserUC)
	router := gin.Default()
	router.GET("/", handler.GetTweetsFromIDs)
	req := httptest.NewRequest(http.MethodGet, "/?user_ids=1,2", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), `"tweets"`)
	assert.Contains(t, rec.Body.String(), `"id":1`)
}
