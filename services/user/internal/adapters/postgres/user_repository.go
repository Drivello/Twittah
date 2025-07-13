package postgres

import (
	"context"
	"fmt"

	"github.com/Drivello/Twittah/services/user/ent"
	userEnt "github.com/Drivello/Twittah/services/user/ent/user"
	"github.com/Drivello/Twittah/services/user/internal/common"
	"github.com/Drivello/Twittah/services/user/internal/domain"
	"go.uber.org/zap"
)

type EntUserRepository struct {
	client *ent.Client
}

// InsertUser inserts a user if not exists (id, username). Ignores duplicate key errors.
func (r *EntUserRepository) InsertUser(ctx context.Context, user *domain.User) error {
	common.Logger().Debug("[UserRepository] InsertUser", zap.Int64("user_id", user.ID), zap.String("username", user.Username))
	// Verify if user exists
	exists, err := r.client.User.
		Query().
		Where(userEnt.IDEQ(user.ID)).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %w", err)
	}

	if exists {
		// Ignore
		return nil
	}

	// Create new user
	_, err = r.client.User.
		Create().
		SetID(user.ID).
		SetUsername(user.Username).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func NewEntUserRepository(client *ent.Client) *EntUserRepository {
	return &EntUserRepository{client: client}
}

func (r *EntUserRepository) FollowUser(ctx context.Context, followerID, followeeID int64) error {
	common.Logger().Debug("[UserRepository] FollowUser", zap.Int64("follower_id", followerID), zap.Int64("followee_id", followeeID))
	if followerID == followeeID {
		return fmt.Errorf("cannot follow yourself")
	}

	_, err := r.client.User.
		UpdateOneID(followeeID).
		AddFollowerIDs(followerID).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}
	return nil
}

func (r *EntUserRepository) UnfollowUser(ctx context.Context, followerID, followeeID int64) error {
	common.Logger().Debug("[UserRepository] UnfollowUser", zap.Int64("follower_id", followerID), zap.Int64("followee_id", followeeID))
	if followerID == followeeID {
		return fmt.Errorf("cannot unfollow yourself")
	}
	_, err := r.client.User.
		UpdateOneID(followeeID).
		RemoveFollowerIDs(followerID).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}
	return nil
}

func (r *EntUserRepository) GetFollowers(ctx context.Context, userID int64) ([]*domain.User, error) {
	common.Logger().Debug("[UserRepository] GetFollowers", zap.Int64("user_id", userID))
	response, err := r.client.User.
		Query().
		Where(userEnt.IDEQ(userID)).
		QueryFollowers().
		All(ctx)
	if err != nil {
		common.Logger().Debug("[UserRepository] GetFollowers error", zap.Int64("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	followers := make([]*domain.User, len(response))
	for i, f := range response {
		followers[i] = &domain.User{
			ID:       f.ID,
			Username: f.Username,
		}
	}
	common.Logger().Debug("[UserRepository] GetFollowers success", zap.Int64("user_id", userID), zap.Any("followers", followers))
	return followers, nil
}

func (r *EntUserRepository) GetFollowing(ctx context.Context, userID int64) ([]*domain.User, error) {
	common.Logger().Debug("[UserRepository] GetFollowing", zap.Int64("user_id", userID))
	response, err := r.client.User.
		Query().
		Where(userEnt.IDEQ(userID)).
		QueryFollowing().
		All(ctx)
	if err != nil {
		common.Logger().Debug("[UserRepository] GetFollowing error", zap.Int64("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	following := make([]*domain.User, len(response))
	for i, f := range response {
		following[i] = &domain.User{
			ID:       f.ID,
			Username: f.Username,
		}
	}
	common.Logger().Debug("[UserRepository] GetFollowing success", zap.Int64("user_id", userID), zap.Any("following", following))
	return following, nil
}

func (r *EntUserRepository) Exists(ctx context.Context, userID int64) (bool, error) {
	exists, err := r.client.User.Query().Where(userEnt.IDEQ(userID)).Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return exists, nil
}
