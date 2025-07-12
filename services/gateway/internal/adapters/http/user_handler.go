package http

import (
	"strconv"

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
	if followerID == "" || followeeID == "" {
		WriteGenericError(c, 400)
		return
	}
	followerIDInt, err := strconv.ParseInt(followerID, 10, 64)
	if err != nil {
		WriteGenericError(c, 400)
		return
	}
	followeeIDInt, err := strconv.ParseInt(followeeID, 10, 64)
	if err != nil {
		WriteGenericError(c, 400)
		return
	}
	if err := h.FollowUseCase.Execute(c, followerIDInt, followeeIDInt); err != nil {
		WriteGenericError(c, 500)
		return
	}
	c.JSON(200, GenericResponse{Message: "Ahora sigues al usuario"})
}

// UnfollowUser handles unfollow requests.
// c: Gin context (HTTP request context).
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")
	if followerID == "" || followeeID == "" {
		WriteGenericError(c, 400)
		return
	}
	followerIDInt, err := strconv.ParseInt(followerID, 10, 64)
	if err != nil {
		WriteGenericError(c, 400)
		return
	}
	followeeIDInt, err := strconv.ParseInt(followeeID, 10, 64)
	if err != nil {
		WriteGenericError(c, 400)
		return
	}
	if err := h.UnfollowUseCase.Execute(c, followerIDInt, followeeIDInt); err != nil {
		WriteGenericError(c, 500)
		return
	}
	c.JSON(200, GenericResponse{Message: "Has dejado de seguir al usuario"})
}
