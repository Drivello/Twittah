package http

import (
	"net/http"

	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/ports"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication-related endpoints.
type AuthHandler struct {
	UserUC ports.RegisterUserUseCasesPort
}

// NewAuthHandler creates a new AuthHandler instance con inyección de usecase.
func NewAuthHandler(userUC ports.RegisterUserUseCasesPort) *AuthHandler {
	return &AuthHandler{UserUC: userUC}
}

// RegisterRoutes registers authentication-related routes with the given Gin engine.
// r: Gin engine to register routes on.
func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/validate", h.Validate)
}

// Validate is a mock endpoint that always responds OK.
// It should be replaced with real validation logic.
// c: Gin context for the HTTP request.
func (h *AuthHandler) Validate(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")

	userIDInt, err := common.ValidatePositiveIntString(userID)
	if err != nil || userIDInt <= 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "X-User-Id header required"})
		return
	}

	resp := ValidateResponseDTO{Status: "OK"}
	c.JSON(200, resp)
}
