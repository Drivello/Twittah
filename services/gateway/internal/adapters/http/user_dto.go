package http

// FollowRequestDTO is the DTO for follow/unfollow requests
type FollowRequestDTO struct {
	FollowerID string `json:"follower_id" binding:"required"`
	FolloweeID string `json:"followee_id" binding:"required"`
}

// FollowResponseDTO is the DTO for follow/unfollow responses
type FollowResponseDTO struct {
	Message string `json:"message"`
}

// FollowersResponseDTO is the DTO for followers list responses
type FollowersResponseDTO struct {
	Followers []string `json:"followers"`
}
