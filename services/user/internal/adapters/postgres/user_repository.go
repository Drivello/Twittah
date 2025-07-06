package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Drivello/Twittah/services/user/ent"
	"github.com/Drivello/Twittah/services/user/ent/user"
)

// EntUserRepository implements ports.UserRepository using Ent ORM and Postgres.
type EntUserRepository struct {
	client *ent.Client
}

func NewEntUserRepository(client *ent.Client) *EntUserRepository {
	return &EntUserRepository{client: client}
}

func (r *EntUserRepository) FollowUser(ctx context.Context, followerID, followeeID string) error {
	if followerID == followeeID {
		return fmt.Errorf("cannot follow yourself")
	}
	followerInt, err := strconv.Atoi(followerID)
	if err != nil {
		return fmt.Errorf("invalid follower id: %w", err)
	}
	followeeInt, err := strconv.Atoi(followeeID)
	if err != nil {
		return fmt.Errorf("invalid followee id: %w", err)
	}
	_, err = r.client.User.
		UpdateOneID(followerInt).
		AddFollowingIDs(followeeInt).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}
	return nil
}

func (r *EntUserRepository) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	if followerID == followeeID {
		return fmt.Errorf("cannot unfollow yourself")
	}
	followerInt, err := strconv.Atoi(followerID)
	if err != nil {
		return fmt.Errorf("invalid follower id: %w", err)
	}
	followeeInt, err := strconv.Atoi(followeeID)
	if err != nil {
		return fmt.Errorf("invalid followee id: %w", err)
	}
	_, err = r.client.User.
		UpdateOneID(followerInt).
		RemoveFollowingIDs(followeeInt).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}
	return nil
}

func (r *EntUserRepository) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	followers, err := r.client.User.
		Query().
		Where(user.IDEQ(uid)).
		QueryFollowers().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	ids := make([]string, len(followers))
	for i, f := range followers {
		ids[i] = fmt.Sprintf("%d", f.ID)
	}
	return ids, nil
}

func (r *EntUserRepository) GetFollowing(ctx context.Context, userID string) ([]string, error) {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	following, err := r.client.User.
		Query().
		Where(user.IDEQ(uid)).
		QueryFollowing().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	ids := make([]string, len(following))
	for i, f := range following {
		ids[i] = fmt.Sprintf("%d", f.ID)
	}
	return ids, nil
}

func (r *EntUserRepository) Exists(ctx context.Context, userID string) (bool, error) {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user id: %w", err)
	}
	exists, err := r.client.User.Query().Where(user.IDEQ(uid)).Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return exists, nil
}
