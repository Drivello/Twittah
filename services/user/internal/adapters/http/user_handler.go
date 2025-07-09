package http

import (
	"context"

	"github.com/Drivello/Twittah/services/user/internal/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	FollowUseCase   *usecase.FollowUserUseCase
	UnfollowUseCase *usecase.UnfollowUserUseCase
}

func (h *UserHandler) repo() interface {
	GetFollowers(ctx context.Context, userID string) ([]string, error)
	GetFollowing(ctx context.Context, userID string) ([]string, error)
} {
	if h.FollowUseCase != nil && h.FollowUseCase.Repo != nil {
		return h.FollowUseCase.Repo
	}
	if h.UnfollowUseCase != nil && h.UnfollowUseCase.Repo != nil {
		return h.UnfollowUseCase.Repo
	}
	return nil
}

func NewUserHandler(followUC *usecase.FollowUserUseCase, unfollowUC *usecase.UnfollowUserUseCase) *UserHandler {
	return &UserHandler{
		FollowUseCase:   followUC,
		UnfollowUseCase: unfollowUC,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/followers/:user_id", h.GetFollowers)
	r.GET("/following/:user_id", h.GetFollowing)
}

func (h *UserHandler) GetFollowers(c *gin.Context) {
	repo := h.repo()
	if repo == nil {
		c.JSON(500, gin.H{"error": "repository not available"})
		return
	}
	userID := c.Param("user_id")
	followers, err := repo.GetFollowers(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resp := FollowersResponseDTO{Followers: followers}
	c.JSON(200, resp)
}

func (h *UserHandler) GetFollowing(c *gin.Context) {
	repo := h.repo()
	if repo == nil {
		c.JSON(500, gin.H{"error": "repository not available"})
		return
	}
	userID := c.Param("user_id")
	following, err := repo.GetFollowing(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resp := FollowingResponseDTO{Following: following}
	c.JSON(200, resp)
}
