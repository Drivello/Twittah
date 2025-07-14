package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the user with AuthService using the X-User-Id header
func AuthMiddleware(authServiceURL string) gin.HandlerFunc {
	client := &http.Client{Timeout: 2 * time.Second}

	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-Id")
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "X-User-Id header required"})
			return
		}

		url := fmt.Sprintf("%s/validate", authServiceURL)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot create auth request"})
			return
		}
		req.Header.Set("X-User-Id", userID)

		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "auth service unavailable"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
			return
		}

		c.Next()
	}
}
