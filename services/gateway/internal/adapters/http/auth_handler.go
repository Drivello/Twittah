package http

import "strings"

import (
	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP endpoints.
type AuthHandler struct {
	Kafka KafkaTopicLister // Interfaz que exponga Topics() ([]string, error)

	UserUseCase *usecase.UserUseCase
}

// NewAuthHandler creates a new AuthHandler.
// producer: Kafka producer for user events.
// Returns a pointer to AuthHandler.
func NewAuthHandler(useCase *usecase.UserUseCase) *AuthHandler {
	return &AuthHandler{UserUseCase: useCase}
}

// RegisterRoutes registers auth endpoints on the given Gin router group.
// rg: Gin router group to register routes on.
func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/health", h.Health)

	rg.POST("/register", h.RegisterUser)
	rg.POST("/validate", h.ValidateAuth)
	rg.POST("/deactivate", h.DeactivateUser)
}



// RegisterUser handles user registration requests.
// c: Gin context.
// Health handles GET /health for deep readiness check.
func (h *AuthHandler) Health(c *gin.Context) {
	status := map[string]string{}
	ok := true
	requiredTopics := []string{"tweets.published", "follows.created", "follows.deleted", "timelines.updated"}
	if h.Kafka != nil {
		topics, err := h.Kafka.Topics()
		topicsMap := map[string]bool{}
		for _, t := range topics {
			topicsMap[t] = true
		}
		if err != nil {
			status["kafka"] = err.Error()
			ok = false
		} else {
			missing := []string{}
			for _, t := range requiredTopics {
				if !topicsMap[t] {
					missing = append(missing, t)
				}
			}
			if len(missing) > 0 {
				status["kafka"] = "missing topics: " + strings.Join(missing, ", ")
				ok = false
			} else {
				status["kafka"] = "ok"
			}
		}
	} else {
		status["kafka"] = "not configured"
		ok = false
	}

	if ok {
		c.JSON(200, status)
	} else {
		c.JSON(500, status)
	}
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
	if err := h.UserUseCase.RegisterUser(reqDTO.Username, reqDTO.Email, reqDTO.Password); err != nil {
		c.JSON(500, RegisterUserResponseDTO{Message: "No se pudo registrar el usuario"})
		return
	}
	resp := RegisterUserResponseDTO{Message: "Usuario registrado", Username: reqDTO.Username}
	c.JSON(201, resp)
}

// DeactivateUser handles user deactivation requests (mock).
// c: Gin context.
func (h *AuthHandler) DeactivateUser(c *gin.Context) {
	// TODO: Implementar desactivación real
	resp := DeactivateUserResponseDTO{Message: "Usuario desactivado (mock)"}
	c.JSON(200, resp)
}

// ValidateAuth handles user authentication validation requests (mock).
// c: Gin context.
func (h *AuthHandler) ValidateAuth(c *gin.Context) {
	// TODO: Implementar validación real
	resp := ValidateAuthResponseDTO{Message: "Autenticación exitosa (mock)"}
	c.JSON(200, resp)
}
