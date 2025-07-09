package http

import (
	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/gin-gonic/gin"
)

// TweetHandler handles tweet-related HTTP endpoints.
// TODO: Refactor TweetHandler to use a usecase for tweet logic, following hexagonal architecture.
type TweetHandler struct {
	Usecase *usecase.TweetUsecase
}

// NewTweetHandler creates a new TweetHandler.
// producer: Kafka producer for follow events.
// Returns a pointer to TweetHandler.
func NewTweetHandler(tweetUsecase *usecase.TweetUsecase) *TweetHandler {
	return &TweetHandler{Usecase: tweetUsecase}
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
	var req TweetRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, TweetPublishResponseDTO{Message: "Invalid tweet request"})
		return
	}
	input := usecase.TweetInput{
		AuthorID: req.AuthorID,
		Content:  req.Content,
	}
	err := h.Usecase.PublishTweet(c.Request.Context(), input)
	if err != nil {
		c.JSON(500, TweetPublishResponseDTO{Message: "No se pudo publicar el tweet"})
		return
	}
	c.JSON(202, TweetPublishResponseDTO{Message: "Tweet enviado para publicación"})
}

// GetTimeline handles timeline retrieval requests (mock).
// c: Gin context.
func (h *TweetHandler) GetTimeline(c *gin.Context) {
	// TODO: Obtener timeline real
	type TimelineResponseDTO struct {
		Timeline []string `json:"timeline"`
	}
	resp := TimelineResponseDTO{Timeline: []string{}}
	c.JSON(200, resp)
}
