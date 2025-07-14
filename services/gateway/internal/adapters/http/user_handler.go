package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/metrics"
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
	start := time.Now()
	metrics.ActiveRequests.Inc()
	defer func() {
		metrics.ActiveRequests.Dec()
		metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, "users/follow").Observe(time.Since(start).Seconds())
	}()
	status := http.StatusAccepted
	defer func() {
		metrics.HTTPRequestTotal.WithLabelValues(c.Request.Method, fmt.Sprintf("%d", status)).Inc()
	}()

	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")

	followerIDInt, err := common.ValidatePositiveIntString(followerID)
	if err != nil {
		status = http.StatusBadRequest
		c.JSON(status, gin.H{"error": "Invalid follower ID"})
		return
	}

	followeeIDInt, err := common.ValidatePositiveIntString(followeeID)
	if err != nil {
		status = http.StatusBadRequest
		c.JSON(status, gin.H{"error": "Invalid followee ID"})
		return
	}

	if err := h.FollowUseCase.Execute(c, followerIDInt, followeeIDInt); err != nil {
		status = http.StatusInternalServerError
		c.JSON(status, gin.H{"error": "Failed to follow user"})
		return
	}

	status = http.StatusAccepted
	c.JSON(status, gin.H{"message": "Ahora sigues al usuario"})
}

// UnfollowUser handles unfollow requests.
// c: Gin context (HTTP request context).
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	start := time.Now()
	metrics.ActiveRequests.Inc()
	defer func() {
		metrics.ActiveRequests.Dec()
		metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, "users/unfollow").Observe(time.Since(start).Seconds())
	}()
	status := http.StatusAccepted
	defer func() {
		metrics.HTTPRequestTotal.WithLabelValues(c.Request.Method, fmt.Sprintf("%d", status)).Inc()
	}()

	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")

	followerIDInt, err := common.ValidatePositiveIntString(followerID)
	if err != nil {
		status = http.StatusBadRequest
		c.JSON(status, gin.H{"error": "Invalid follower ID"})
		return
	}

	followeeIDInt, err := common.ValidatePositiveIntString(followeeID)
	if err != nil {
		status = http.StatusBadRequest
		c.JSON(status, gin.H{"error": "Invalid followee ID"})
		return
	}

	if err := h.UnfollowUseCase.Execute(c, followerIDInt, followeeIDInt); err != nil {
		status = http.StatusInternalServerError
		//TODO: change http depending on error
		c.JSON(status, gin.H{"error": "Failed to unfollow user"})
		return
	}

	status = http.StatusAccepted
	c.JSON(status, gin.H{"message": "Has dejado de seguir al usuario"})
}
