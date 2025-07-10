package http

import (
	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user HTTP endpoints.
type UserHandler struct {
	UserUseCase *usecase.UserUseCase
}

// NewUserHandler creates a new UserHandler.
// useCase: Use case for user operations.
// Returns a pointer to UserHandler.
func NewUserHandler(useCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{UserUseCase: useCase}
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
	if followerID == "" || followeeID == "" {
		c.JSON(400, FollowResponseDTO{Message: "Missing user IDs"})
		return
	}
	if err := h.UserUseCase.FollowUser(c, followerID, followeeID); err != nil {
		c.JSON(500, FollowResponseDTO{Message: "No se pudo seguir al usuario"})
		return
	}
	c.JSON(200, FollowResponseDTO{Message: "Ahora sigues al usuario"})
}

// UnfollowUser handles unfollow requests.
// c: Gin context (HTTP request context).
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")
	if followerID == "" || followeeID == "" {
		c.JSON(400, FollowResponseDTO{Message: "Missing user IDs"})
		return
	}
	if err := h.UserUseCase.UnfollowUser(c, followerID, followeeID); err != nil {
		c.JSON(500, FollowResponseDTO{Message: "No se pudo dejar de seguir al usuario"})
		return
	}
	c.JSON(200, FollowResponseDTO{Message: "Has dejado de seguir al usuario"})
}
