package http

import (
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/validate", h.Validate)
}

// Validate is a mock endpoint that always responds OK
func (h *AuthHandler) Validate(c *gin.Context) {
	//TODO: Implement validation logic	c.JSON(200, gin.H{"status": "OK"})
}
