package usecase

import "github.com/Drivello/Twittah/services/gateway/internal/ports"

// UserUseCase provides the application logic for user-related operations.
type UserUseCase struct {
	Producer ports.UserEventProducerPort
}

// NewUserUseCase creates a new UserUseCase.
func NewUserUseCase(producer ports.UserEventProducerPort) *UserUseCase {
	return &UserUseCase{Producer: producer}
}

// RegisterUser handles user registration logic.
func (uc *UserUseCase) RegisterUser(username, email, password string) error {
	return uc.Producer.PublishUserCreateRequest(username, email, password)
}
