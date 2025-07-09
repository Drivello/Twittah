package http

import (
	"net/http"
	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP endpoints.
type AuthHandler struct {
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
	rg.POST("/register", h.RegisterUser)
	rg.POST("/validate", h.ValidateAuth)
	rg.POST("/deactivate", h.DeactivateUser)
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser handles user registration requests.
// c: Gin context.
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.JSON(400, gin.H{"error": "Missing required fields"})
		return
	}
	if err := h.UserUseCase.RegisterUser(req.Username, req.Email, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el usuario"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado", "username": req.Username})
}

// DeactivateUser handles user deactivation requests (mock).
// c: Gin context.
func (h *AuthHandler) DeactivateUser(c *gin.Context) {
	// TODO: Implementar desactivación real
	c.JSON(http.StatusOK, gin.H{"message": "Usuario desactivado (mock)"})
}

// ValidateAuth handles user authentication validation requests (mock).
// c: Gin context.
func (h *AuthHandler) ValidateAuth(c *gin.Context) {
	// TODO: Implementar validacion real
	c.JSON(http.StatusOK, gin.H{"message": "Autenticación exitosa (mock)"})
}
