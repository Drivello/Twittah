package http

import (
	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication-related endpoints.
type AuthHandler struct {
}

// NewAuthHandler creates a new AuthHandler instance.
// Returns a pointer to AuthHandler.
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
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
	resp := ValidateResponseDTO{Status: "OK"}
	c.JSON(200, resp)
}
