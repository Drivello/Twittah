package postgres

import (
	"context"
	"time"

	"github.com/Drivello/Twittah/services/auth/ent"
	"github.com/Drivello/Twittah/services/auth/internal/domain"
	"go.uber.org/zap"
)

type PostgresUserRepository struct {
	Client *ent.Client
}

func NewPostgresUserRepository(client *ent.Client) *PostgresUserRepository {
	return &PostgresUserRepository{Client: client}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *domain.User) (int64, error) {
	zap.L().Info("[Repo] Persisting user to database", zap.String("username", user.Username))
	createdUser, err := r.Client.User.
		Create().
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetPassword(user.Password).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		zap.L().Error("[Repo] Failed to persist user", zap.Error(err))
		return 0, err
	}
	zap.L().Info("[Repo] User persisted to database", zap.Int64("id", createdUser.ID))
	return createdUser.ID, nil
}
