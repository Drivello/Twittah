package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type GatewayHandler struct{}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{}
}

func (h *GatewayHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/tweets", h.PostTweet)
	r.POST("/follow/:target_user_id", h.FollowUser)
	r.DELETE("/unfollow/:target_user_id", h.UnfollowUser)
	r.GET("/timeline", h.GetTimeline)
}

func (h *GatewayHandler) PostTweet(c *gin.Context) {
	// TODO: Enviar evento a Kafka y validar con AuthService
	c.JSON(http.StatusCreated, gin.H{
		"message": "Tweet publicado (mock)",
	})
}

func (h *GatewayHandler) FollowUser(c *gin.Context) {
	// TODO: Enviar evento follow a Kafka y validar con AuthService/UserService
	c.JSON(http.StatusOK, gin.H{
		"message": "Ahora sigues al usuario (mock)",
	})
}

func (h *GatewayHandler) UnfollowUser(c *gin.Context) {
	// TODO: Enviar evento unfollow a Kafka y validar con AuthService/UserService
	c.JSON(http.StatusOK, gin.H{
		"message": "Has dejado de seguir al usuario (mock)",
	})
}

func (h *GatewayHandler) GetTimeline(c *gin.Context) {
	// TODO: Consultar timeline real al TweetService (por ahora mock)
	c.JSON(http.StatusOK, gin.H{
		"user_id": "123",
		"timeline": []gin.H{
			{
				"tweet_id": "789",
				"author_id": "456",
				"content": "Hola mundo desde nuestra API!",
				"created_at": "2025-07-07T14:32:00Z",
			},
		},
	})
}
