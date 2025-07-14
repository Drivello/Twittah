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

func TestGetFollowers_InvalidUserID_Returns400(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockUC := &mockGetFollowersUseCase{}
	handler := httpg.NewUserQueryHandler(mockUC)
	router := gin.Default()
	router.GET("/followers/:user_id", handler.GetFollowers)
	req := httptest.NewRequest(http.MethodGet, "/followers/NaN", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 400, rec.Code)
	assert.JSONEq(t, `{"followers": null}`, rec.Body.String())
}

func TestGetFollowers_UseCase_Error_Returns502(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockUC := &mockGetFollowersUseCase{
		ExecuteFunc: func(ctx context.Context, userID int64) ([]*domain.User, error) {
			return nil, errors.New("fail")
		},
	}
	handler := httpg.NewUserQueryHandler(mockUC)
	router := gin.Default()
	router.GET("/followers/:user_id", handler.GetFollowers)
	req := httptest.NewRequest(http.MethodGet, "/followers/123", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 502, rec.Code)
	assert.JSONEq(t, `{"followers": null}`, rec.Body.String())
}

func TestGetFollowers_Success(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockUC := &mockGetFollowersUseCase{
		ExecuteFunc: func(ctx context.Context, userID int64) ([]*domain.User, error) {
			return []*domain.User{
				{ID: 1, Username: "alice"},
				{ID: 2, Username: "bob"},
			}, nil
		},
	}
	handler := httpg.NewUserQueryHandler(mockUC)
	router := gin.Default()
	router.GET("/followers/:user_id", handler.GetFollowers)
	req := httptest.NewRequest(http.MethodGet, "/followers/321", nil)
	rec := httptest.NewRecorder()
	// act
	router.ServeHTTP(rec, req)
	// assert
	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), `"followers"`)
	assert.Contains(t, rec.Body.String(), `"alice"`)
	assert.Contains(t, rec.Body.String(), `"bob"`)
}
