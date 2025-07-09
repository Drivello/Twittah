package http

// FollowersResponseDTO is the DTO for returning followers via HTTP
type FollowersResponseDTO struct {
	Followers []string `json:"followers"`
}

// FollowingResponseDTO is the DTO for returning following via HTTP
type FollowingResponseDTO struct {
	Following []string `json:"following"`
}
