package http

import (
	"net/http"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/dto"
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
	// Map domain.User to DTO
	dtos := make([]dto.FollowerUserDTO, len(followers))
	for i, f := range followers {
		dtos[i] = dto.FollowerUserDTO{ID: f.ID, Username: f.Username}
	}
	common.Logger().Debug("[Gateway] GetFollowers handler success", zap.String("user_id", userID), zap.Any("followers", dtos))
	c.JSON(http.StatusOK, dto.GetFollowersResponse{Followers: dtos})
}
