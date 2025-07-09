package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
)

// TweetHandler handles tweet-related HTTP endpoints.
// TODO: Refactor TweetHandler to use a usecase for tweet logic, following hexagonal architecture.
type TweetHandler struct {
	FollowProducer *kafka.FollowEventProducer
}

// NewTweetHandler creates a new TweetHandler.
// producer: Kafka producer for follow events.
// Returns a pointer to TweetHandler.
func NewTweetHandler(producer *kafka.FollowEventProducer) *TweetHandler {
	return &TweetHandler{FollowProducer: producer}
}

// RegisterRoutes registers tweet endpoints on the given Gin router group.
// rg: Gin router group to register routes on.
func (h *TweetHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.PostTweet)
	rg.GET("/timeline", h.GetTimeline)
}

// PostTweet handles tweet publishing requests (mock).
// c: Gin context.
func (h *TweetHandler) PostTweet(c *gin.Context) {
	// TODO: Enviar evento a Kafka y validar con AuthService
	c.JSON(http.StatusCreated, gin.H{"message": "Tweet publicado (mock)"})
}

// GetTimeline handles timeline retrieval requests (mock).
// c: Gin context.
func (h *TweetHandler) GetTimeline(c *gin.Context) {
	// TODO: Obtener timeline real
	c.JSON(http.StatusOK, gin.H{"timeline": []string{}})
}
