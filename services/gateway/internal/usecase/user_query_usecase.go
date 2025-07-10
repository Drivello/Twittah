package usecase

import (
	"context"

	"github.com/Drivello/Twittah/services/gateway/internal/ports"
)

// UserQueryUseCase implementa la lógica de consulta de followers
// y cumple con el puerto hexagonal UserQueryPort

type UserQueryUseCase struct {
	Adapter ports.UserQueryPort
}

var _ ports.UserQueryPort = (*UserQueryUseCase)(nil)

func NewUserQueryUseCase(adapter ports.UserQueryPort) *UserQueryUseCase {
	return &UserQueryUseCase{Adapter: adapter}
}

func (uc *UserQueryUseCase) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	return uc.Adapter.GetFollowers(ctx, userID)
}
