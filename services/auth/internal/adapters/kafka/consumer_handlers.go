package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Drivello/Twittah/services/auth/internal/domain"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
)

// Handler for user creation events (modularizado)
func HandleUserCreate(ctx context.Context, authUC ports.AuthUseCasesPort, payload map[string]interface{}) error {
	var data struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	b, _ := json.Marshal(payload)
	if err := json.Unmarshal(b, &data); err != nil {
		return fmt.Errorf("invalid payload structure: %w", err)
	}

	if data.Username == "" || data.Email == "" || data.Password == "" {
		return domain.ErrInvalidPayload
	}

	user := domain.User{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	}
	_, err := authUC.RegisterUser(ctx, user)
	return err
}
