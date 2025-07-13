package dto

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
)

// --------- STRUCTS ---------

type FollowerUserDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (fud *FollowerUserDTO) Validate() error {
	if fud.ID <= 0 {
		return errors.New("id must be a positive integer")
	}
	username := strings.TrimSpace(fud.Username)
	if username == "" {
		return errors.New("username is required")
	}
	if len(username) > 50 {
		return errors.New("username must be at most 50 characters long")
	}
	if !common.UsernameRegex.MatchString(username) {
		return errors.New("username must contain only letters (a-z, A-Z), no spaces or symbols")
	}
	return nil
}

// --------- REQUESTS ---------

type FollowRequest struct {
	FollowerID string `json:"follower_id" binding:"required"`
	FolloweeID string `json:"followee_id" binding:"required"`
}

func (fr *FollowRequest) Validate() error {
	if strings.TrimSpace(fr.FollowerID) == "" {
		return errors.New("follower_id is required")
	}
	if strings.TrimSpace(fr.FolloweeID) == "" {
		return errors.New("followee_id is required")
	}
	followerID, err := strconv.ParseInt(fr.FollowerID, 10, 64)
	if err != nil || followerID <= 0 {
		return errors.New("follower_id must be a positive integer")
	}
	followeeID, err := strconv.ParseInt(fr.FolloweeID, 10, 64)
	if err != nil || followeeID <= 0 {
		return errors.New("followee_id must be a positive integer")
	}
	if followerID == followeeID {
		return errors.New("follower_id and followee_id cannot be the same")
	}
	return nil
}

// --------- RESPONSES ---------

type GetFollowersResponse struct {
	Followers []FollowerUserDTO `json:"followers"`
}

func (gfr *GetFollowersResponse) Validate() error {
	for i, follower := range gfr.Followers {
		if err := follower.Validate(); err != nil {
			return fmt.Errorf("invalid follower at index %d: %w", i, err)
		}
	}
	return nil
}

type FollowResponse struct {
	Message string `json:"message"`
}

func (fr *FollowResponse) Validate() error {
	if strings.TrimSpace(fr.Message) == "" {
		return errors.New("message is required")
	}
	return nil
}
