package http

import (
	"net/http"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"go.uber.org/zap"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"github.com/gin-gonic/gin"
)

// UserQueryHandler maneja las rutas de consulta de followers

type UserQueryHandler struct {
	GetFollowersUseCase ports.GetFollowersUseCasePort
}

func NewUserQueryHandler(getFollowersUseCase ports.GetFollowersUseCasePort) *UserQueryHandler {
	return &UserQueryHandler{GetFollowersUseCase: getFollowersUseCase}
}

func (h *UserQueryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/followers/:user_id", h.GetFollowers)
}

func (h *UserQueryHandler) GetFollowers(c *gin.Context) {
	userID := c.Param("user_id")
	common.Logger().Debug("[Gateway] HTTP GetFollowers handler called", zap.String("user_id", userID))
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"followers": nil})
		return
	}
	followers, err := h.GetFollowersUseCase.Execute(c, userID)
	if err != nil {
		common.Logger().Debug("[Gateway] GetFollowersUseCase error", zap.String("user_id", userID), zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"followers": nil})
		return
	}
	common.Logger().Debug("[Gateway] GetFollowers handler success", zap.String("user_id", userID), zap.Any("followers", followers))
	c.JSON(http.StatusOK, gin.H{"followers": followers})
}
