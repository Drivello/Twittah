package http

import (
	"encoding/json"
	"net/http"

	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/gin-gonic/gin"
)

// FollowHandler handles follow/unfollow HTTP endpoints.
type FollowHandler struct {
	FollowUseCase *usecase.FollowUseCase
}

// NewFollowHandler creates a new FollowHandler.
// producer: Kafka producer for follow events.
// Returns a pointer to FollowHandler.
func NewFollowHandler(useCase *usecase.FollowUseCase) *FollowHandler {
	return &FollowHandler{FollowUseCase: useCase}
}

// RegisterRoutes registers follow endpoints on the given Gin router group.
// rg: Gin router group to register routes on.
func (h *FollowHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/follow/:target_user_id", h.FollowUser)
	rg.DELETE("/unfollow/:target_user_id", h.UnfollowUser)
	rg.GET("/followers/:user_id", h.GetFollowers)
}

// FollowUser handles follow requests.
// c: Gin context.
func (h *FollowHandler) FollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")
	if followerID == "" || followeeID == "" {
		c.JSON(400, FollowResponseDTO{Message: "Missing user IDs"})
		return
	}
	if err := h.FollowUseCase.FollowUser(followerID, followeeID); err != nil {
		c.JSON(500, FollowResponseDTO{Message: "No se pudo seguir al usuario"})
		return
	}
	c.JSON(200, FollowResponseDTO{Message: "Ahora sigues al usuario"})
}

// UnfollowUser handles unfollow requests.
// c: Gin context.
func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")
	if followerID == "" || followeeID == "" {
		c.JSON(400, FollowResponseDTO{Message: "Missing user IDs"})
		return
	}
	if err := h.FollowUseCase.UnfollowUser(followerID, followeeID); err != nil {
		c.JSON(500, FollowResponseDTO{Message: "No se pudo dejar de seguir al usuario"})
		return
	}
	c.JSON(200, FollowResponseDTO{Message: "Has dejado de seguir al usuario"})
}

// GetFollowers handles requests to get a user's followers.
// c: Gin context.
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(400, FollowersResponseDTO{Followers: nil})
		return
	}
	userServiceURL := "http://user:8082/followers/" + userID
	resp, err := http.Get(userServiceURL)
	if err != nil || resp.StatusCode != 200 {
		c.JSON(502, FollowersResponseDTO{Followers: nil})
		return
	}
	defer resp.Body.Close()
	var followersDTO FollowersResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&followersDTO); err != nil {
		c.JSON(502, FollowersResponseDTO{Followers: nil})
		return
	}
	c.JSON(200, followersDTO)
}
