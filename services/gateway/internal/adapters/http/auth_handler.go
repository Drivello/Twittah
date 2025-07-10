package http

import (
	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP endpoints.
type AuthHandler struct {
	Kafka KafkaTopicLister // Interfaz que exponga Topics() ([]string, error)

	UserUseCase *usecase.AuthUseCase
}

// NewAuthHandler creates a new AuthHandler.
// producer: Kafka producer for user events.
// Returns a pointer to AuthHandler.
func NewAuthHandler(useCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{UserUseCase: useCase}
}

// RegisterRoutes registers auth endpoints on the given Gin router group.
// rg: Gin router group to register routes on.
func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {

	rg.POST("/register", h.RegisterUser)
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var reqDTO RegisterUserRequestDTO
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		c.JSON(400, RegisterUserResponseDTO{Message: "Invalid request"})
		return
	}
	if reqDTO.Username == "" || reqDTO.Email == "" || reqDTO.Password == "" {
		c.JSON(400, RegisterUserResponseDTO{Message: "Missing required fields"})
		return
	}
	if err := h.UserUseCase.RegisterUser(c, reqDTO.Username, reqDTO.Email, reqDTO.Password); err != nil {
		c.JSON(500, RegisterUserResponseDTO{Message: "No se pudo registrar el usuario"})
		return
	}
	resp := RegisterUserResponseDTO{Message: "Usuario registrado", Username: reqDTO.Username}
	c.JSON(201, resp)
}
