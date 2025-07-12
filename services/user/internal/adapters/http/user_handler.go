package http

import (
	"strconv"

	"github.com/Drivello/Twittah/services/user/internal/common"
	"go.uber.org/zap"
	"github.com/Drivello/Twittah/services/user/internal/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	GetFollowersUseCase *usecase.GetFollowersUseCase
	GetFollowingUseCase *usecase.GetFollowingUseCase
}

func NewUserHandler(getFollowersUC *usecase.GetFollowersUseCase, getFollowingUC *usecase.GetFollowingUseCase) *UserHandler {
	return &UserHandler{
		GetFollowersUseCase: getFollowersUC,
		GetFollowingUseCase: getFollowingUC,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/followers/:user_id", h.GetFollowers)
	r.GET("/following/:user_id", h.GetFollowing)
}

func (h *UserHandler) GetFollowers(c *gin.Context) {
	userIDStr := c.Param("user_id")
	common.Logger().Debug("[UserService] HTTP GetFollowers handler called", zap.String("user_id", userIDStr))
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		common.Logger().Debug("[UserService] Invalid user_id", zap.String("user_id", userIDStr), zap.Error(err))
		c.JSON(400, gin.H{"error": "invalid user_id"})
		return
	}
	followers, err := h.GetFollowersUseCase.Execute(c.Request.Context(), userID)
	if err != nil {
		common.Logger().Debug("[UserService] GetFollowersUseCase error", zap.Int64("user_id", userID), zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	common.Logger().Debug("[UserService] GetFollowers handler success", zap.Int64("user_id", userID), zap.Any("followers", followers))
	resp := FollowersResponseDTO{Followers: followers}
	c.JSON(200, resp)
}

func (h *UserHandler) GetFollowing(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user_id"})
		return
	}
	following, err := h.GetFollowingUseCase.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resp := FollowingResponseDTO{Following: following}
	c.JSON(200, resp)
}
