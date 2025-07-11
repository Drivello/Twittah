package domain

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidPayload = errors.New("invalid payload")

	ErrUserAlreadyExists = errors.New("user already exists")
)

func IsValidEmail(email string) bool {
	// Regex sencilla para validar email
	var re = regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// User represents an Auth user entity.
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
