package http

import (
	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user HTTP endpoints.
type UserHandler struct {
	FollowUseCase   ports.FollowUserUseCasePort
	UnfollowUseCase ports.UnfollowUserUseCasePort
}

// NewUserHandler creates a new UserHandler.
// useCase: Use case for user operations.
// Returns a pointer to UserHandler.
func NewUserHandler(followUC ports.FollowUserUseCasePort, unfollowUC ports.UnfollowUserUseCasePort) *UserHandler {
	return &UserHandler{FollowUseCase: followUC, UnfollowUseCase: unfollowUC}
}

// RegisterRoutes registers follow/unfollow endpoints on the given Gin router group.
// rg: Gin router group to register routes on.
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/follow/:target_user_id", h.FollowUser)
	rg.DELETE("/unfollow/:target_user_id", h.UnfollowUser)
}

// FollowUser handles follow requests.
// c: Gin context (HTTP request context).
func (h *UserHandler) FollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")

	followerIDInt, err := common.ValidatePositiveIntString(followerID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid follower ID"})
		return
	}

	followeeIDInt, err := common.ValidatePositiveIntString(followeeID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid followee ID"})
		return
	}

	if err := h.FollowUseCase.Execute(c, followerIDInt, followeeIDInt); err != nil {
		//TODO: change http depending on error
		c.JSON(500, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(200, gin.H{"message": "Ahora sigues al usuario"})
}

// UnfollowUser handles unfollow requests.
// c: Gin context (HTTP request context).
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")

	followerIDInt, err := common.ValidatePositiveIntString(followerID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid follower ID"})
		return
	}

	followeeIDInt, err := common.ValidatePositiveIntString(followeeID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid followee ID"})
		return
	}

	if err := h.UnfollowUseCase.Execute(c, followerIDInt, followeeIDInt); err != nil {
		//TODO: change http depending on error
		c.JSON(500, gin.H{"error": "Failed to unfollow user"})
		return
	}

	c.JSON(200, gin.H{"message": "Has dejado de seguir al usuario"})
}
