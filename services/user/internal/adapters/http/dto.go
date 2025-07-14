package http

import "github.com/Drivello/Twittah/services/user/internal/domain"

// FollowersResponseDTO is the DTO for returning followers via HTTP
type GetFollowersResponseDTO struct {
	Followers []*domain.User `json:"followers"`
}

// FollowingResponseDTO is the DTO for returning following via HTTP
type GetFollowingResponseDTO struct {
	Following []*domain.User `json:"following"`
}
