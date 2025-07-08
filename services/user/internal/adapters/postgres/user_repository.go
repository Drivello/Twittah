package postgres

import (
	"context"
	"fmt"

	"github.com/Drivello/Twittah/services/user/ent"
	"github.com/Drivello/Twittah/services/user/ent/user"
)

type EntUserRepository struct {
	client *ent.Client
}

// InsertUser inserts a user if not exists (id, username). Ignores duplicate key errors.
func (r *EntUserRepository) InsertUser(ctx context.Context, id, username string) error {
	// Verify if user exists
	exists, err := r.client.User.
		Query().
		Where(user.IDEQ(id)).
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
		SetID(id).
		SetUsername(username).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func NewEntUserRepository(client *ent.Client) *EntUserRepository {
	return &EntUserRepository{client: client}
}

func (r *EntUserRepository) FollowUser(ctx context.Context, followerID, followeeID string) error {
	if followerID == followeeID {
		return fmt.Errorf("cannot follow yourself")
	}

	_, err := r.client.User.
		UpdateOneID(followerID).
		AddFollowingIDs(followeeID).
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
	_, err := r.client.User.
		UpdateOneID(followerID).
		RemoveFollowingIDs(followeeID).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}
	return nil
}

func (r *EntUserRepository) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	followers, err := r.client.User.
		Query().
		Where(user.IDEQ(userID)).
		QueryFollowers().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	ids := make([]string, len(followers))
	for i, f := range followers {
		ids[i] = f.ID
	}
	return ids, nil
}

func (r *EntUserRepository) GetFollowing(ctx context.Context, userID string) ([]string, error) {
	following, err := r.client.User.
		Query().
		Where(user.IDEQ(userID)).
		QueryFollowing().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	ids := make([]string, len(following))
	for i, f := range following {
		ids[i] = f.ID
	}
	return ids, nil
}

func (r *EntUserRepository) Exists(ctx context.Context, userID string) (bool, error) {
	exists, err := r.client.User.Query().Where(user.IDEQ(userID)).Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return exists, nil
}
