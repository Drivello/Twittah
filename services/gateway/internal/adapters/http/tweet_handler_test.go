package http_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	httpg "github.com/Drivello/Twittah/services/gateway/internal/adapters/http"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostTweet_ShouldBindJSON_Error(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockCreateUseCase := &mockCreateTweetUseCase{ExecuteFunc: func(ctx context.Context, authorID int64, content string) error { return nil }}
	mockDeleteUseCase := &mockDeleteTweetUseCase{}
	handler := httpg.NewTweetHandler(mockCreateUseCase, mockDeleteUseCase)
	router := gin.Default()
	router.POST("/", handler.PostTweet)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{invalid-json}"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid tweet request")
}

func TestPostTweet_Validate_Error(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockCreateUseCase := &mockCreateTweetUseCase{ExecuteFunc: func(ctx context.Context, authorID int64, content string) error { return nil }}
	mockDeleteUseCase := &mockDeleteTweetUseCase{}
	handler := httpg.NewTweetHandler(mockCreateUseCase, mockDeleteUseCase)
	router := gin.Default()
	router.POST("/", handler.PostTweet)
	reqBody := `{"author_id":0,"content":"Hello world"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid tweet request")
}

func TestPostTweet_UseCase_Error(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockCreateUseCase := &mockCreateTweetUseCase{ExecuteFunc: func(ctx context.Context, authorID int64, content string) error {
		return errors.New("something failed")
	}}
	mockDeleteUseCase := &mockDeleteTweetUseCase{}
	handler := httpg.NewTweetHandler(mockCreateUseCase, mockDeleteUseCase)
	router := gin.Default()
	router.POST("/", handler.PostTweet)
	reqBody := `{"author_id":1,"content":"Hello world"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 500, rec.Code)
	assert.Contains(t, rec.Body.String(), "No se pudo publicar el tweet")
}

func TestPostTweet_Success(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockCreateUseCase := &mockCreateTweetUseCase{ExecuteFunc: func(ctx context.Context, authorID int64, content string) error { return nil }}
	mockDeleteUseCase := &mockDeleteTweetUseCase{}
	handler := httpg.NewTweetHandler(mockCreateUseCase, mockDeleteUseCase)
	router := gin.Default()
	router.POST("/", handler.PostTweet)
	reqBody := `{"author_id":1,"content":"Hello world"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 202, rec.Code)
	assert.Contains(t, rec.Body.String(), "Tweet enviado para publicación")
}

func TestDeleteTweet_InvalidTweetID_Error(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockCreateUseCase := &mockCreateTweetUseCase{}
	mockDeleteUseCase := &mockDeleteTweetUseCase{}
	handler := httpg.NewTweetHandler(mockCreateUseCase, mockDeleteUseCase)
	router := gin.Default()
	router.DELETE("/:tweet_id", handler.DeleteTweet)
	req := httptest.NewRequest(http.MethodDelete, "/abc", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid tweet ID")
}

func TestDeleteTweet_UseCase_Error(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockCreateUseCase := &mockCreateTweetUseCase{}
	mockDeleteUseCase := &mockDeleteTweetUseCase{ExecuteFunc: func(ctx context.Context, tweetID int64) error {
		return errors.New("something failed")
	}}
	handler := httpg.NewTweetHandler(mockCreateUseCase, mockDeleteUseCase)
	router := gin.Default()
	router.DELETE("/:tweet_id", handler.DeleteTweet)
	req := httptest.NewRequest(http.MethodDelete, "/123", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 500, rec.Code)
	assert.Contains(t, rec.Body.String(), "No se pudo eliminar el tweet")
}

func TestDeleteTweet_Success(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockCreateUseCase := &mockCreateTweetUseCase{}
	mockDeleteUseCase := &mockDeleteTweetUseCase{ExecuteFunc: func(ctx context.Context, tweetID int64) error { return nil }}
	handler := httpg.NewTweetHandler(mockCreateUseCase, mockDeleteUseCase)
	router := gin.Default()
	router.DELETE("/:tweet_id", handler.DeleteTweet)
	req := httptest.NewRequest(http.MethodDelete, "/123", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 202, rec.Code)
	assert.Contains(t, rec.Body.String(), "Tweet eliminado")
}
