package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Drivello/Twittah/services/user/internal/usecase"
)

type UserHandler struct {
	FollowUseCase   *usecase.FollowUserUseCase
	UnfollowUseCase *usecase.UnfollowUserUseCase
}

func NewUserHandler(followUC *usecase.FollowUserUseCase, unfollowUC *usecase.UnfollowUserUseCase) *UserHandler {
	return &UserHandler{
		FollowUseCase:   followUC,
		UnfollowUseCase: unfollowUC,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/follow/:target_user_id", h.FollowUser)
	r.DELETE("/follow/:target_user_id", h.UnfollowUser)
}

func (h *UserHandler) FollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")
	if err := h.FollowUseCase.Execute(c.Request.Context(), followerID, followeeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})
}

func (h *UserHandler) UnfollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")
	if err := h.UnfollowUseCase.Execute(c.Request.Context(), followerID, followeeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}
