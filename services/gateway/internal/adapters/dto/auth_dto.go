package dto

import (
	"errors"
	"strings"

	"github.com/Drivello/Twittah/services/gateway/internal/common"
)

// ---------- REQUESTS ----------
type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r *RegisterUserRequest) Validate() error {
	username := strings.TrimSpace(r.Username)
	if username == "" {
		return errors.New("username is required")
	}
	if len(username) > 50 {
		return errors.New("username must be at most 50 characters")
	}
	if !common.UsernameRegex.MatchString(username) {
		return errors.New("username must contain only letters (a-z, A-Z)")
	}

	email := strings.TrimSpace(r.Email)
	if email == "" {
		return errors.New("email is required")
	}
	if len(email) > 100 {
		return errors.New("email must be at most 100 characters")
	}
	if !common.EmailRegex.MatchString(email) {
		return errors.New("email is not valid")
	}

	password := strings.TrimSpace(r.Password)
	if password == "" {
		return errors.New("password is required")
	}
	if !common.IsValidPassword(password) {
		return errors.New("password must be at least 8 characters, contain at least one uppercase letter and one special character")
	}

	return nil
}

// ---------- RESPONSES ----------

type RegisterUserResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

type ValidateAuthResponse struct {
	Message string `json:"message"`
}

type DeactivateUserResponse struct {
	Message string `json:"message"`
}
