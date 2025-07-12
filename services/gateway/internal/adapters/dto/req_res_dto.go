package dto

// FollowRequestDTO is the DTO for follow/unfollow requests
type FollowRequest struct {
	FollowerID string `json:"follower_id" binding:"required"`
	FolloweeID string `json:"followee_id" binding:"required"`
}

// FollowResponse is the  for follow/unfollow responses
type FollowResponse struct {
	Message string `json:"message"`
}

// GetFollowersResponse is the DTO for followers list responses
type GetFollowersResponse struct {
	Followers []FollowerUserDTO `json:"followers"`
}

type FollowerUserDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// RegisterUserRequest is the  for user registration

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterUserResponse is the  for user registration response
type RegisterUserResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

// ValidateAuthResponse is the  for authentication validation response
type ValidateAuthResponse struct {
	Message string `json:"message"`
}

// DeactivateUserResponse is the  for user deactivation response
type DeactivateUserResponse struct {
	Message string `json:"message"`
}

// TweetRequest representa el request para publicar un tweet via HTTP
// Usado solo por el handler de tweets
type TweetCreateRequest struct {
	AuthorID int64  `json:"author_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// TweetCreateResponse representa la respuesta del endpoint de publicación de tweet
// Usado solo por el handler de tweets
type TweetCreateResponse struct {
	Message string `json:"message"`
}
