package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/dto"
	"github.com/Drivello/Twittah/services/gateway/internal/metrics"
	"github.com/Drivello/Twittah/services/gateway/internal/ports"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP endpoints.
type AuthHandler struct {
	registerUserUseCase ports.RegisterUserUseCasePort
}

// NewAuthHandler creates a new AuthHandler.
// producer: Kafka producer for user events.
// Returns a pointer to AuthHandler.
func NewAuthHandler(registerUserUseCase ports.RegisterUserUseCasePort) *AuthHandler {
	return &AuthHandler{registerUserUseCase: registerUserUseCase}
}

// RegisterRoutes registers auth endpoints on the given Gin router group.
// rg: Gin router group to register routes on.
func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {

	rg.POST("/register", h.RegisterUser)
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	start := time.Now()
	metrics.ActiveRequests.Inc()
	defer func() {
		metrics.ActiveRequests.Dec()
		metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, "users/register").Observe(time.Since(start).Seconds())
	}()
	status := http.StatusCreated
	defer func() {
		metrics.HTTPRequestTotal.WithLabelValues(c.Request.Method, fmt.Sprintf("%d", status)).Inc()
	}()

	var reqDTO dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		status = http.StatusBadRequest
		c.JSON(status, dto.RegisterUserResponse{Message: "Invalid request"})
		return
	}

	if err := reqDTO.Validate(); err != nil {
		status = http.StatusBadRequest
		c.JSON(status, dto.RegisterUserResponse{Message: "Invalid request"})
		return
	}

	if err := h.registerUserUseCase.Execute(c, reqDTO.Username, reqDTO.Email, reqDTO.Password); err != nil {
		status = http.StatusInternalServerError
		//TODO: change http depending on error
		c.JSON(status, dto.RegisterUserResponse{Message: "No se pudo registrar el usuario"})
		return
	}

	resp := dto.RegisterUserResponse{Message: "Usuario registrado", Username: reqDTO.Username}
	status = http.StatusCreated
	c.JSON(status, resp)
}
