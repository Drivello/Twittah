package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	httpg "github.com/Drivello/Twittah/services/gateway/internal/adapters/http"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFollowUser_InvalidFollowerID_Returns400(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockFollow := &mockFollowUserUseCase{}
	mockUnfollow := &mockUnfollowUserUseCase{}
	handler := httpg.NewUserHandler(mockFollow, mockUnfollow)
	router := gin.Default()
	router.POST("/follow/:target_user_id", handler.FollowUser)
	req := httptest.NewRequest(http.MethodPost, "/follow/10", nil)
	req.Header.Set("X-User-Id", "not_a_number")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid follower ID")
}

func TestFollowUser_InvalidFolloweeID_Returns400(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockFollow := &mockFollowUserUseCase{}
	mockUnfollow := &mockUnfollowUserUseCase{}
	handler := httpg.NewUserHandler(mockFollow, mockUnfollow)
	router := gin.Default()
	router.POST("/follow/:target_user_id", handler.FollowUser)
	req := httptest.NewRequest(http.MethodPost, "/follow/abc", nil)
	req.Header.Set("X-User-Id", "123")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid followee ID")
}

func TestFollowUser_UseCase_Error_Returns500(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockFollow := &mockFollowUserUseCase{
		ExecuteFunc: func(ctx context.Context, followerID, followeeID int64) error {
			return errors.New("fail follow")
		},
	}
	mockUnfollow := &mockUnfollowUserUseCase{}
	handler := httpg.NewUserHandler(mockFollow, mockUnfollow)
	router := gin.Default()
	router.POST("/follow/:target_user_id", handler.FollowUser)
	req := httptest.NewRequest(http.MethodPost, "/follow/22", nil)
	req.Header.Set("X-User-Id", "11")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 500, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to follow user")
}

func TestFollowUser_Success(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockFollow := &mockFollowUserUseCase{
		ExecuteFunc: func(ctx context.Context, followerID, followeeID int64) error {
			return nil
		},
	}
	mockUnfollow := &mockUnfollowUserUseCase{}
	handler := httpg.NewUserHandler(mockFollow, mockUnfollow)
	router := gin.Default()
	router.POST("/follow/:target_user_id", handler.FollowUser)
	req := httptest.NewRequest(http.MethodPost, "/follow/22", nil)
	req.Header.Set("X-User-Id", "11")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), "Ahora sigues al usuario")
}

func TestUnfollowUser_InvalidFollowerID_Returns400(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockFollow := &mockFollowUserUseCase{}
	mockUnfollow := &mockUnfollowUserUseCase{}
	handler := httpg.NewUserHandler(mockFollow, mockUnfollow)
	router := gin.Default()
	router.DELETE("/unfollow/:target_user_id", handler.UnfollowUser)
	req := httptest.NewRequest(http.MethodDelete, "/unfollow/10", nil)
	req.Header.Set("X-User-Id", "NaN")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid follower ID")
}

func TestUnfollowUser_InvalidFolloweeID_Returns400(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockFollow := &mockFollowUserUseCase{}
	mockUnfollow := &mockUnfollowUserUseCase{}
	handler := httpg.NewUserHandler(mockFollow, mockUnfollow)
	router := gin.Default()
	router.DELETE("/unfollow/:target_user_id", handler.UnfollowUser)
	req := httptest.NewRequest(http.MethodDelete, "/unfollow/abc", nil)
	req.Header.Set("X-User-Id", "123")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid followee ID")
}

func TestUnfollowUser_UseCase_Error_Returns500(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockFollow := &mockFollowUserUseCase{}
	mockUnfollow := &mockUnfollowUserUseCase{
		ExecuteFunc: func(ctx context.Context, followerID, followeeID int64) error {
			return errors.New("fail unfollow")
		},
	}
	handler := httpg.NewUserHandler(mockFollow, mockUnfollow)
	router := gin.Default()
	router.DELETE("/unfollow/:target_user_id", handler.UnfollowUser)
	req := httptest.NewRequest(http.MethodDelete, "/unfollow/22", nil)
	req.Header.Set("X-User-Id", "11")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 500, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to unfollow user")
}

func TestUnfollowUser_Success(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockFollow := &mockFollowUserUseCase{}
	mockUnfollow := &mockUnfollowUserUseCase{
		ExecuteFunc: func(ctx context.Context, followerID, followeeID int64) error {
			return nil
		},
	}
	handler := httpg.NewUserHandler(mockFollow, mockUnfollow)
	router := gin.Default()
	router.DELETE("/unfollow/:target_user_id", handler.UnfollowUser)
	req := httptest.NewRequest(http.MethodDelete, "/unfollow/22", nil)
	req.Header.Set("X-User-Id", "11")
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), "Has dejado de seguir al usuario")
}
