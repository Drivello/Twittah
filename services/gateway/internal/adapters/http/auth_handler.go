package http

import (
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/dto"
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
	var reqDTO dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		WriteGenericError(c, 400)
		return
	}
	if reqDTO.Username == "" || reqDTO.Email == "" || reqDTO.Password == "" {
		WriteGenericError(c, 400)
		return
	}
	if err := h.registerUserUseCase.Execute(c, reqDTO.Username, reqDTO.Email, reqDTO.Password); err != nil {
		WriteGenericError(c, 500)
		return
	}
	resp := dto.RegisterUserResponse{Message: "Usuario registrado", Username: reqDTO.Username}
	c.JSON(201, resp)
}
