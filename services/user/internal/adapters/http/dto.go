package http

// FollowersResponseDTO is the DTO for returning followers via HTTP
type FollowersResponseDTO struct {
	Followers []int64 `json:"followers"`
}

// FollowingResponseDTO is the DTO for returning following via HTTP
type FollowingResponseDTO struct {
	Following []int64 `json:"following"`
}
