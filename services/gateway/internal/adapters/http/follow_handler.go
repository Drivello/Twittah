package http

import (
	"io"
	"net/http"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	FollowProducer *kafka.FollowEventProducer
}

func NewFollowHandler(producer *kafka.FollowEventProducer) *FollowHandler {
	return &FollowHandler{FollowProducer: producer}
}

func (h *FollowHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/follow/:target_user_id", h.FollowUser)
	rg.DELETE("/unfollow/:target_user_id", h.UnfollowUser)
	rg.GET("/followers/:user_id", h.GetFollowers)
}

func (h *FollowHandler) FollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")
	if followerID == "" || followeeID == "" {
		c.JSON(400, gin.H{"error": "Missing user IDs"})
		return
	}
	err := h.FollowProducer.PublishFollow(followerID, followeeID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to publish follow event"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ahora sigues al usuario"})
}

func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-Id")
	followeeID := c.Param("target_user_id")
	if followerID == "" || followeeID == "" {
		c.JSON(400, gin.H{"error": "Missing user IDs"})
		return
	}
	err := h.FollowProducer.PublishUnfollow(followerID, followeeID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to publish unfollow event"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Has dejado de seguir al usuario"})
}

func (h *FollowHandler) GetFollowers(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user_id"})
		return
	}
	userServiceURL := "http://user:8082/followers/" + userID
	resp, err := http.Get(userServiceURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "No se pudo contactar UserService"})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "UserService error", "status": resp.StatusCode})
		return
	}
	c.Status(http.StatusOK)
	io.Copy(c.Writer, resp.Body)
}
