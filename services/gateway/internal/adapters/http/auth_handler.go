package http

import (
	"net/http"

	"github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserProducer *kafka.UserEventProducer
}

func NewAuthHandler(producer *kafka.UserEventProducer) *AuthHandler {
	return &AuthHandler{UserProducer: producer}
}

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
	if err := h.UserProducer.PublishUserCreateRequest(req.Username, req.Email, req.Password); err != nil {
		c.JSON(500, gin.H{"error": "Failed to publish user creation request"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado", "username": req.Username})
}

func (h *AuthHandler) DeactivateUser(c *gin.Context) {
	// TODO: Implementar desactivación real
	c.JSON(http.StatusOK, gin.H{"message": "Usuario desactivado (mock)"})
}

func (h *AuthHandler) ValidateAuth(c *gin.Context) {
	// TODO: Implementar validacion real
	c.JSON(http.StatusOK, gin.H{"message": "Autenticación exitosa (mock)"})
}
