package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware checks if the user is an admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
			c.Abort()
			return
		}

		if role.(string) != "Admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access, admin role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
