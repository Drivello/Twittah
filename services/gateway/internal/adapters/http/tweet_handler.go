package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/dto"
	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/metrics"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"github.com/gin-gonic/gin"
)

// TweetHandler handles tweet-related HTTP endpoints.

type TweetHandler struct {
	createTweetUseCase ports.CreateTweetUseCasePort
	deleteTweetUseCase ports.DeleteTweetUseCasePort
}

// NewTweetHandler creates a new TweetHandler.
// producer: Kafka producer for follow events.
// Returns a pointer to TweetHandler.
func NewTweetHandler(createTweetUseCase ports.CreateTweetUseCasePort, deleteTweetUseCase ports.DeleteTweetUseCasePort) *TweetHandler {
	return &TweetHandler{createTweetUseCase: createTweetUseCase, deleteTweetUseCase: deleteTweetUseCase}
}

// RegisterRoutes registers tweet endpoints on the given Gin router group.
// rg: Gin router group to register routes on.
func (h *TweetHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.PostTweet)
	rg.DELETE("/:tweet_id", h.DeleteTweet)

}

func (h *TweetHandler) PostTweet(c *gin.Context) {
	start := time.Now()
	metrics.ActiveRequests.Inc()
	defer func() {
		metrics.ActiveRequests.Dec()
		metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, "/tweets").Observe(time.Since(start).Seconds())
	}()
	status := http.StatusAccepted
	defer func() {
		metrics.HTTPRequestTotal.WithLabelValues(c.Request.Method, fmt.Sprintf("%d", status)).Inc()
	}()

	var req dto.TweetCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status = http.StatusBadRequest
		c.JSON(status, dto.TweetCreateResponse{Message: "Invalid tweet request"})
		return
	}

	if err := req.Validate(); err != nil {
		status = http.StatusBadRequest
		c.JSON(status, dto.TweetCreateResponse{Message: "Invalid tweet request"})
		return
	}

	err := h.createTweetUseCase.Execute(c.Request.Context(), req.AuthorID, req.Content)
	if err != nil {
		status = http.StatusInternalServerError
		c.JSON(status, dto.TweetCreateResponse{Message: "No se pudo publicar el tweet"})
		return
	}
	status = http.StatusAccepted
	c.JSON(status, dto.TweetCreateResponse{Message: "Tweet enviado para publicación"})
}

func (h *TweetHandler) DeleteTweet(c *gin.Context) {
	start := time.Now()
	metrics.ActiveRequests.Inc()
	defer func() {
		metrics.ActiveRequests.Dec()
		metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, "/tweets/:tweet_id").Observe(time.Since(start).Seconds())
	}()
	status := http.StatusAccepted
	defer func() {
		metrics.HTTPRequestTotal.WithLabelValues(c.Request.Method, fmt.Sprintf("%d", status)).Inc()
	}()

	tweetID := c.Param("tweet_id")

	if tweetID == "" {
		status = http.StatusBadRequest
		c.JSON(status, dto.TweetDeleteResponse{Message: "Invalid tweet ID"})
		return
	}

	tweetIDInt, err := common.ValidatePositiveIntString(tweetID)

	if err != nil {
		status = http.StatusBadRequest
		c.JSON(status, dto.TweetDeleteResponse{Message: "Invalid tweet ID"})
		return
	}

	err = h.deleteTweetUseCase.Execute(c.Request.Context(), tweetIDInt)
	if err != nil {
		status = http.StatusInternalServerError
		c.JSON(status, dto.TweetDeleteResponse{Message: "No se pudo eliminar el tweet"})
		return
	}
	status = http.StatusAccepted
	c.JSON(status, dto.TweetDeleteResponse{Message: "Tweet eliminado"})
}
