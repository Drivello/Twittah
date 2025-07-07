package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
)

type TweetHandler struct {
	FollowProducer *kafka.FollowEventProducer
}

func NewTweetHandler(producer *kafka.FollowEventProducer) *TweetHandler {
	return &TweetHandler{FollowProducer: producer}
}

func (h *TweetHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.PostTweet)
	rg.GET("/timeline", h.GetTimeline)
}

func (h *TweetHandler) PostTweet(c *gin.Context) {
	// TODO: Enviar evento a Kafka y validar con AuthService
	c.JSON(http.StatusCreated, gin.H{"message": "Tweet publicado (mock)"})
}

func (h *TweetHandler) GetTimeline(c *gin.Context) {
	// TODO: Obtener timeline real
	c.JSON(http.StatusOK, gin.H{"timeline": []string{}})
}
