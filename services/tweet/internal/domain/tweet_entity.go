// Package domain defines the core Tweet entity for Twittah.
package domain

import (
	"errors"
	"strings"
	"time"
)

// Tweet represents a user tweet in the domain model.
type Tweet struct {
	ID        string    // UUID
	AuthorID  string    // UUID of the tweet's author
	Content   string    // Tweet content (max 280 chars)
	CreatedAt time.Time // Creation timestamp
}

// Validate checks the tweet content for domain rules.
func (t *Tweet) Validate() error {
	if strings.TrimSpace(t.Content) == "" {
		return errors.New("tweet content cannot be empty")
	}
	if len([]rune(t.Content)) > 280 {
		return errors.New("tweet content exceeds 280 characters")
	}

	return nil
}
