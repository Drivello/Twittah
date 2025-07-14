// Package domain contains core business entities for the User Service.
package domain

// User represents a user in the Twittah platform domain.

type User struct {
	ID       int64  // Unique identifier for the user
	Username string // Username for the user
}
