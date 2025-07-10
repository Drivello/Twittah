package http

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
