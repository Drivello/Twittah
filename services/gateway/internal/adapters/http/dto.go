package http

// TweetRequestDTO representa el request para publicar un tweet via HTTP
// Usado solo por el handler de tweets
//
type TweetRequestDTO struct {
	AuthorID string `json:"author_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// TweetPublishResponseDTO representa la respuesta del endpoint de publicación de tweet
// Usado solo por el handler de tweets
//
type TweetPublishResponseDTO struct {
	Message string `json:"message"`
}

// RegisterUserRequestDTO is the DTO for user registration
// (adapter-level, not domain)
type RegisterUserRequestDTO struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterUserResponseDTO is the DTO for user registration response
type RegisterUserResponseDTO struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

// ValidateAuthResponseDTO is the DTO for authentication validation response
type ValidateAuthResponseDTO struct {
	Message string `json:"message"`
}

// DeactivateUserResponseDTO is the DTO for user deactivation response
type DeactivateUserResponseDTO struct {
	Message string `json:"message"`
}
