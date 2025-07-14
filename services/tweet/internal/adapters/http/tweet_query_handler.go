package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Drivello/Twittah/services/tweet/internal/adapters/dto"
	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/ports"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TweetQueryHandler struct {
	GetTimelineUC                  ports.GetTimelinePort
	GetTweetsFromMultipleUserIDsUC ports.GetTweetsFromMultipleUserIDsPort
	GetUserTweetsUC                ports.GetUserTweetsPort
}

func NewTweetQueryHandler(getTimelineUC ports.GetTimelinePort, getTweetsFromMultipleUserIDsUC ports.GetTweetsFromMultipleUserIDsPort, getUserTweetsUC ports.GetUserTweetsPort) *TweetQueryHandler {
	return &TweetQueryHandler{GetTimelineUC: getTimelineUC, GetTweetsFromMultipleUserIDsUC: getTweetsFromMultipleUserIDsUC, GetUserTweetsUC: getUserTweetsUC}
}

func (h *TweetQueryHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/timeline/:user_id", h.GetTimeline)
	group.GET("/multiple", h.GetTweetsFromMultipleUserIDs)
	group.GET("/user/:user_id", h.GetUserTweets)
}

func (h *TweetQueryHandler) GetTimeline(c *gin.Context) {
	userIDParam := c.Param("user_id")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query param required"})
		return
	}
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id: " + userIDParam})
		return
	}

	tweets, err := h.GetTimelineUC.Execute(c.Request.Context(), userID)
	if err != nil {
		common.Logger().Error("failed to get tweets from ids", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	common.Logger().Debug("[TweetQueryHandler] GetTimeline success", zap.Int64("user_id", userID), zap.Any("tweets", tweets))
	tweetsDTO := make([]dto.TweetDTO, len(tweets))
	for i, t := range tweets {
		tweetsDTO[i] = dto.TweetDTO{ID: t.ID, Author: t.AuthorID, Content: t.Content, CreatedAt: time.Unix(t.CreatedAt, 0)}
	}

	resp := dto.GetTimelineResponse{
		Timeline: tweetsDTO,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TweetQueryHandler) GetTweetsFromMultipleUserIDs(c *gin.Context) {
	userIDsParam := c.Query("user_ids")
	if userIDsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_ids query param required"})
		return
	}
	idStrs := strings.Split(userIDsParam, ",")
	var userIDs []int64
	for _, s := range idStrs {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id: " + s})
			return
		}
		userIDs = append(userIDs, id)
	}

	tweets, err := h.GetTweetsFromMultipleUserIDsUC.Execute(c.Request.Context(), userIDs)
	if err != nil {
		common.Logger().Error("failed to get tweets from ids", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	common.Logger().Debug("[TweetQueryHandler] GetTweetsFromMultipleUserIDs success", zap.Any("user_ids", userIDs), zap.Any("tweets", tweets))
	tweetsDTO := make([]dto.TweetDTO, len(tweets))
	for i, t := range tweets {
		tweetsDTO[i] = dto.TweetDTO{ID: t.ID, Author: t.AuthorID, Content: t.Content, CreatedAt: time.Unix(t.CreatedAt, 0)}
	}

	resp := dto.GetTweetsFromMultipleUserIDsResponse{
		Tweets: tweetsDTO,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TweetQueryHandler) GetUserTweets(c *gin.Context) {
	userIDParam := c.Param("user_id")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query param required"})
		return
	}
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id: " + userIDParam})
		return
	}

	tweets, err := h.GetUserTweetsUC.Execute(c.Request.Context(), userID)
	if err != nil {
		common.Logger().Error("failed to get user tweets", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	common.Logger().Debug("[TweetQueryHandler] GetUserTweets success", zap.Int64("user_id", userID), zap.Any("tweets", tweets))
	tweetsDTO := make([]dto.TweetDTO, len(tweets))
	for i, t := range tweets {
		tweetsDTO[i] = dto.TweetDTO{ID: t.ID, Author: t.AuthorID, Content: t.Content, CreatedAt: time.Unix(t.CreatedAt, 0)}
	}

	resp := dto.GetTweetsFromUserIDResponse{
		Tweets: tweetsDTO,
	}

	c.JSON(http.StatusOK, resp)
}
