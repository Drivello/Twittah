package dto

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// ---------- Constantes / Variables ----------

const MaxTweetContentLength = 280

// ---------- STRUCTS ----------

type TweetDTO struct {
	ID        int64     `json:"id"`
	Author    int64     `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (tdto *TweetDTO) Validate() error {
	if tdto.ID <= 0 {
		return errors.New("id must be a positive integer")
	}
	if tdto.Author <= 0 {
		return errors.New("author must be a positive integer")
	}
	content := strings.TrimSpace(tdto.Content)
	if content == "" {
		return errors.New("content is required")
	}
	if len(content) > MaxTweetContentLength {
		return errors.New("content must be at most 280 characters")
	}
	if tdto.CreatedAt.IsZero() {
		return errors.New("created_at is required")
	}
	return nil
}

// ---------- REQUESTS ----------

type TweetCreateRequest struct {
	AuthorID int64  `json:"author_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

func (tcr *TweetCreateRequest) Validate() error {
	if tcr.AuthorID <= 0 {
		return errors.New("author_id must be a positive integer")
	}
	content := strings.TrimSpace(tcr.Content)
	if content == "" {
		return errors.New("content is required")
	}
	if len(content) > MaxTweetContentLength {
		return errors.New("content must be at most 280 characters")
	}
	return nil
}

type GetTweetsFromUserIDRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

func (r *GetTweetsFromUserIDRequest) Validate() error {
	if r.UserID <= 0 {
		return errors.New("user_id must be a positive integer")
	}
	return nil
}

type GetTweetsFromMultipleUserIDsResponse struct {
	Tweets []TweetDTO `json:"tweets" binding:"required"`
}

func (r *GetTweetsFromMultipleUserIDsResponse) Validate() error {
	if len(r.Tweets) == 0 {
		return errors.New("tweets is required and must not be empty")
	}
	for i, t := range r.Tweets {
		if err := t.Validate(); err != nil {
			return errors.New("invalid tweet at index " + strconv.Itoa(i) + ": " + err.Error())
		}
	}
	return nil
}

type TweetDeleteRequest struct {
	TweetID int64 `json:"tweet_id" binding:"required"`
}

func (r *TweetDeleteRequest) Validate() error {
	if r.TweetID <= 0 {
		return errors.New("tweet_id must be a positive integer")
	}
	return nil
}

// ---------- RESPONSES ----------

type TweetCreateResponse struct {
	Message string `json:"message"`
}

type TweetDeleteResponse struct {
	Message string `json:"message"`
}

type GetTimelineResponse struct {
	Timeline []TweetDTO `json:"timeline"`
}

func (resp *GetTimelineResponse) Validate() error {
	for i, t := range resp.Timeline {
		if err := t.Validate(); err != nil {
			return errors.New("invalid tweet in timeline at index " + strconv.Itoa(i) + ": " + err.Error())
		}
	}
	return nil
}

type GetTweetsFromUserIDResponse struct {
	Tweets []TweetDTO `json:"tweets"`
}

func (resp *GetTweetsFromUserIDResponse) Validate() error {
	for i, t := range resp.Tweets {
		if err := t.Validate(); err != nil {
			return errors.New("invalid tweet at index " + strconv.Itoa(i) + ": " + err.Error())
		}
	}
	return nil
}

type GetTweetsFromIDsResponse struct {
	Tweets []TweetDTO `json:"tweets"`
}

func (resp *GetTweetsFromIDsResponse) Validate() error {
	for i, t := range resp.Tweets {
		if err := t.Validate(); err != nil {
			return errors.New("invalid tweet at index " + strconv.Itoa(i) + ": " + err.Error())
		}
	}
	return nil
}
