// Package http provides the REST API handlers for TweetService.
package http

import (
	"strconv"
	"time"

	"github.com/Drivello/Twittah/services/tweet/internal/domain"
	"github.com/Drivello/Twittah/services/tweet/internal/usecase"
	"github.com/gin-gonic/gin"
)

// TweetHandler handles HTTP requests for tweets and timelines.
type TweetHandler struct {
	PublishUC  *usecase.PublishTweet
	TimelineUC *usecase.GetTimeline
}

// NewTweetHandler creates a new TweetHandler.
func NewTweetHandler(publishUC *usecase.PublishTweet, timelineUC *usecase.GetTimeline) *TweetHandler {
	return &TweetHandler{PublishUC: publishUC, TimelineUC: timelineUC}
}

// RegisterRoutes registers tweet endpoints on the Gin router.
func (h *TweetHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/tweets", h.PostTweet)
	r.GET("/timeline", h.GetTimeline)
}

// PostTweet handles POST /tweets.
func (h *TweetHandler) PostTweet(c *gin.Context) {
	var req PostTweetRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	tweet := &domain.Tweet{
		ID:        "",
		AuthorID:  req.AuthorID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}
	if err := tweet.Validate(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := h.PublishUC.Execute(c.Request.Context(), tweet, req.Followers)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resp := TweetResponseDTO{
		ID:        tweet.ID,
		AuthorID:  tweet.AuthorID,
		Content:   tweet.Content,
		CreatedAt: tweet.CreatedAt.Format(time.RFC3339),
	}
	c.JSON(201, gin.H{"message": "tweet published", "tweet": resp})
}

// GetTimeline handles GET /timeline.
func (h *TweetHandler) GetTimeline(c *gin.Context) {
	userID := c.Query("user_id")
	limit := 50
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	if userID == "" {
		c.JSON(400, gin.H{"error": "user_id required"})
		return
	}
	tweets, err := h.TimelineUC.Execute(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resp := make([]TweetResponseDTO, 0, len(tweets))
	for _, t := range tweets {
		resp = append(resp, TweetResponseDTO{
			ID:        t.ID,
			AuthorID:  t.AuthorID,
			Content:   t.Content,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
		})
	}
	c.JSON(200, gin.H{"timeline": resp})
}
