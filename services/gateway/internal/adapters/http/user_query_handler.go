package http

import (
	"net/http"

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
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"followers": nil})
		return
	}
	followers, err := h.GetFollowersUseCase.Execute(c, userID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"followers": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"followers": followers})
}
