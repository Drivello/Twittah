package http_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	httpg "github.com/Drivello/Twittah/services/gateway/internal/adapters/http"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser_ShouldBindJSON_Error(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockUseCase := &mockRegisterUserUseCase{
		ExecuteFunc: func(ctx context.Context, username, email, password string) error { return nil },
	}
	handler := httpg.NewAuthHandler(mockUseCase)
	router := gin.Default()
	router.POST("/register", handler.RegisterUser)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("{invalid-json}"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// act
	router.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid request")
}

func TestRegisterUser_Validate_Error(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockUseCase := &mockRegisterUserUseCase{
		ExecuteFunc: func(ctx context.Context, username, email, password string) error { return nil },
	}
	handler := httpg.NewAuthHandler(mockUseCase)
	router := gin.Default()
	router.POST("/register", handler.RegisterUser)

	// Provoca error de validación: username vacío
	reqBody := `{"username":"","email":"test@example.com","password":"Valid$Pass1"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// act
	router.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid request")
}

func TestRegisterUser_UseCase_Error(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockUseCase := &mockRegisterUserUseCase{
		ExecuteFunc: func(ctx context.Context, username, email, password string) error {
			return errors.New("something failed")
		},
	}
	handler := httpg.NewAuthHandler(mockUseCase)
	router := gin.Default()
	router.POST("/register", handler.RegisterUser)

	reqBody := `{"username":"ValidUser","email":"test@example.com","password":"Valid$Pass1"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// act
	router.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, 500, rec.Code)
	assert.Contains(t, rec.Body.String(), "No se pudo registrar el usuario")
}

func TestRegisterUser_Success(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	mockUseCase := &mockRegisterUserUseCase{
		ExecuteFunc: func(ctx context.Context, username, email, password string) error {
			return nil
		},
	}
	handler := httpg.NewAuthHandler(mockUseCase)
	router := gin.Default()
	router.POST("/register", handler.RegisterUser)

	reqBody := `{"username":"ValidUser","email":"test@example.com","password":"Valid$Pass1"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// act
	router.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, 201, rec.Code)
	assert.Contains(t, rec.Body.String(), "Usuario registrado")
}
