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

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	// TODO: Implementar registro real
	c.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado (mock)"})
}

func (h *AuthHandler) DeactivateUser(c *gin.Context) {
	// TODO: Implementar desactivación real
	c.JSON(http.StatusOK, gin.H{"message": "Usuario desactivado (mock)"})
}

func (h *AuthHandler) ValidateAuth(c *gin.Context) {
	// TODO: Implementar validacion real
	c.JSON(http.StatusOK, gin.H{"message": "Autenticación exitosa (mock)"})
}
