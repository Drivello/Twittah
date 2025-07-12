package postgres

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/ent"
	"github.com/Drivello/Twittah/services/tweet/internal/domain"
)

type UserRepository struct {
	Client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{Client: client}
}

func (r *UserRepository) Save(u *domain.User) error {
	_, err := r.Client.User.
		Create().
		SetID(u.ID).
		SetUsername(u.Username).
		Save(context.Background())
	return err
}

func (r *UserRepository) FindByID(userID int64) (*domain.User, error) {
	e, err := r.Client.User.Get(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:       e.ID,
		Username: e.Username,
	}, nil
}
