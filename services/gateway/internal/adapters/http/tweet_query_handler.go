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
	"go.uber.org/zap"
)

type TweetQueryHandler struct {
	GetTimelineUC         ports.GetTimelinePort
	GetTweetsFromIDsUC    ports.GetTweetsFromIDsPort
	GetTweetsFromUserIDUC ports.GetTweetsFromUserIDPort
}

func NewTweetQueryHandler(getTimelineUC ports.GetTimelinePort, getTweetsFromIDsUC ports.GetTweetsFromIDsPort, getTweetsFromUserIDUC ports.GetTweetsFromUserIDPort) *TweetQueryHandler {
	return &TweetQueryHandler{GetTimelineUC: getTimelineUC, GetTweetsFromIDsUC: getTweetsFromIDsUC, GetTweetsFromUserIDUC: getTweetsFromUserIDUC}
}

func (h *TweetQueryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/", h.GetTweetsFromIDs)
	rg.GET("/:user_id", h.GetTweetsFromUserID)
	rg.GET("/timeline/:user_id", h.GetTimeline)
}

func (h *TweetQueryHandler) GetTimeline(c *gin.Context) {
	start := time.Now()
	metrics.ActiveRequests.Inc()
	defer func() {
		metrics.ActiveRequests.Dec()
		metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, "tweets/timeline").Observe(time.Since(start).Seconds())
	}()
	status := http.StatusOK
	defer func() {
		metrics.HTTPRequestTotal.WithLabelValues(c.Request.Method, fmt.Sprintf("%d", status)).Inc()
	}()

	userID := c.Param("user_id")
	userIDInt, err := common.ValidatePositiveIntString(userID)
	if err != nil {
		status = http.StatusBadRequest
		c.JSON(status, gin.H{"timeline": nil})
		return
	}
	timeline, err := h.GetTimelineUC.Execute(c, userIDInt)
	if err != nil {
		status = http.StatusBadGateway
		common.Logger().Debug("[Gateway] GetTimelineUseCase error", zap.String("user_id", userID), zap.Error(err))
		c.JSON(status, gin.H{"timeline": nil})
		return
	}
	tweetDtos := make([]dto.TweetDTO, len(timeline))
	for i, t := range timeline {
		tweetDtos[i] = dto.TweetDTO{ID: t.ID, Author: t.Author, Content: t.Content, CreatedAt: t.CreatedAt}
	}
	common.Logger().Debug("[Gateway] GetTimeline handler success", zap.String("user_id", userID), zap.Any("timeline", tweetDtos))
	c.JSON(http.StatusOK, dto.GetTimelineResponse{Timeline: tweetDtos})
}

func (h *TweetQueryHandler) GetTweetsFromIDs(c *gin.Context) {
	start := time.Now()
	metrics.ActiveRequests.Inc()
	defer func() {
		metrics.ActiveRequests.Dec()
		metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, "tweets/from_ids").Observe(time.Since(start).Seconds())
	}()
	status := http.StatusOK
	defer func() {
		metrics.HTTPRequestTotal.WithLabelValues(c.Request.Method, fmt.Sprintf("%d", status)).Inc()
	}()
	userIDsParam := c.Query("user_ids")

	userIDs, err := common.ValidateIDListFromString(userIDsParam)
	if err != nil {
		status = http.StatusBadRequest
		c.JSON(status, gin.H{"error": "invalid user_ids: " + userIDsParam})
		return
	}

	tweets, err := h.GetTweetsFromIDsUC.Execute(c, userIDs)
	if err != nil {
		common.Logger().Debug("[Gateway] GetTweetsFromIDsUseCase error", zap.String("user_ids", userIDsParam), zap.Error(err))
		status = http.StatusBadGateway
		c.JSON(status, gin.H{"tweets": nil})
		return
	}

	dtos := make([]dto.TweetDTO, len(tweets))
	for i, t := range tweets {
		dtos[i] = dto.TweetDTO{ID: t.ID, Author: t.Author, Content: t.Content, CreatedAt: t.CreatedAt}
	}

	common.Logger().Debug("[Gateway] GetTweetsFromIDs handler success", zap.String("user_ids", userIDsParam), zap.Any("tweets", dtos))
	status = http.StatusOK
	c.JSON(status, dto.GetTweetsFromIDsResponse{Tweets: dtos})
}

func (h *TweetQueryHandler) GetTweetsFromUserID(c *gin.Context) {
	start := time.Now()
	metrics.ActiveRequests.Inc()
	defer func() {
		metrics.ActiveRequests.Dec()
		metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, "tweets/:user_id").Observe(time.Since(start).Seconds())
	}()
	status := http.StatusOK
	defer func() {
		metrics.HTTPRequestTotal.WithLabelValues(c.Request.Method, fmt.Sprintf("%d", status)).Inc()
	}()

	userID := c.Param("user_id")
	userIDInt, err := common.ValidatePositiveIntString(userID)
	if err != nil {
		status = http.StatusBadRequest
		c.JSON(status, gin.H{"error": "invalid user_id: " + userID})
		return
	}
	tweets, err := h.GetTweetsFromUserIDUC.Execute(c, userIDInt)
	if err != nil {
		common.Logger().Debug("[Gateway] GetTweetsFromIDsUseCase error", zap.String("user_id", userID), zap.Error(err))
		status = http.StatusBadGateway
		c.JSON(status, gin.H{"tweets": nil})
		return
	}
	dtos := make([]dto.TweetDTO, len(tweets))
	for i, t := range tweets {
		dtos[i] = dto.TweetDTO{ID: t.ID, Author: t.Author, Content: t.Content, CreatedAt: t.CreatedAt}
	}
	common.Logger().Debug("[Gateway] GetTweetsFromUserID handler success", zap.String("user_id", userID), zap.Any("tweets", dtos))
	status = http.StatusOK
	c.JSON(status, dto.GetTweetsFromIDsResponse{Tweets: dtos})
}
