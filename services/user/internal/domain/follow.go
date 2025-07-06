// Follow domain entity for Twittah
package domain

// Follow represents a follow relationship between users.
type Follow struct {
	FollowerID string // The user who follows
	FolloweeID string // The user being followed
}
