package ports

import "github.com/Drivello/Twittah/services/tweet/internal/domain"

type UserRepository interface {
	Save(user *domain.User) error
	FindByID(userID int64) (*domain.User, error)
}
