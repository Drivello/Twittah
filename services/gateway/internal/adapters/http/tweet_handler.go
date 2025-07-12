package http

import (
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/dto"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"github.com/gin-gonic/gin"
)

// TweetHandler handles tweet-related HTTP endpoints.
// TODO: Refactor TweetHandler to use a usecase for tweet logic, following hexagonal architecture.
type TweetHandler struct {
	createTweetUseCase ports.CreateTweetUseCasePort
}

// NewTweetHandler creates a new TweetHandler.
// producer: Kafka producer for follow events.
// Returns a pointer to TweetHandler.
func NewTweetHandler(createTweetUseCase ports.CreateTweetUseCasePort) *TweetHandler {
	return &TweetHandler{createTweetUseCase: createTweetUseCase}
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
	var req dto.TweetCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.TweetCreateResponse{Message: "Invalid tweet request"})
		return
	}

	err := h.createTweetUseCase.Execute(c.Request.Context(), req.AuthorID, req.Content)
	if err != nil {
		c.JSON(500, dto.TweetCreateResponse{Message: "No se pudo publicar el tweet"})
		return
	}
	c.JSON(202, dto.TweetCreateResponse{Message: "Tweet enviado para publicación"})
}

// GetTimeline handles timeline retrieval requests (mock).
// c: Gin context.
func (h *TweetHandler) GetTimeline(c *gin.Context) {
	// TODO: Obtener timeline real
	//resp := dto.GetTimelineResponse{Timeline: []string{}}
	c.JSON(200, nil)
}
