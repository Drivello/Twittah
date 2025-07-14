package usecase_test

import (
	"context"
	"testing"

	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser_Success(t *testing.T) {
	// arrange
	uc := usecase.NewRegisterUserUseCase(&mockAuthProducer{})
	// act
	err := uc.Execute(context.Background(), "user", "user@mail.com", "secret")
	// assert
	assert.NoError(t, err)
}

func TestRegisterUser_ProducerError_ReturnsError(t *testing.T) {
	// arrange
	uc := usecase.NewRegisterUserUseCase(&errorAuthProducer{})
	// act
	err := uc.Execute(context.Background(), "user", "user@mail.com", "secret")
	// assert
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}
