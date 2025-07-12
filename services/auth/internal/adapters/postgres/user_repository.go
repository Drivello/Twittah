package postgres

import (
	"context"
	"time"

	"github.com/Drivello/Twittah/services/auth/ent"
	userent "github.com/Drivello/Twittah/services/auth/ent/user"
	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type PostgresUserRepository struct {
	Client *ent.Client
}

func NewPostgresUserRepository(client *ent.Client) *PostgresUserRepository {
	return &PostgresUserRepository{Client: client}
}

func (r *PostgresUserRepository) CreateOrGetUser(ctx context.Context, user domain.User) (int64, bool, error) {

	common.Logger().Debug("[UserRepo] Starting user creation", zap.String("username", user.Username), zap.String("email", user.Email))

	if user.Username == "" || user.Email == "" || !domain.IsValidEmail(user.Email) || user.Password == "" {
		return 0, false, domain.ErrInvalidPayload
	}

	// Hash de la contraseña
	hashed, err := hashPassword(user.Password)
	if err != nil {
		common.Logger().Error("[AuthUseCase] Failed to hash password", zap.Error(err))
		return 0, false, err
	}
	user.Password = hashed

	createdUser, err := r.Client.User.
		Create().
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetPassword(user.Password).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)

	if err == nil {
		common.Logger().Info("[UserRepo] User created successfully", zap.Int64("id", createdUser.ID))
		return createdUser.ID, true, nil
	}

	if ent.IsConstraintError(err) {
		common.Logger().Warn("[UserRepo] User already exists, fetching existing user", zap.String("email", user.Email))
		existingUser, fetchErr := r.Client.User.
			Query().
			Where(userent.EmailEQ(user.Email)).
			Only(ctx)
		if fetchErr != nil {
			common.Logger().Error("[UserRepo] Failed to fetch existing user", zap.Error(fetchErr))
			return 0, false, fetchErr
		}
		return existingUser.ID, false, nil
	}

	common.Logger().Error("[UserRepo] Failed to create user", zap.Error(err))
	return 0, false, err
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
