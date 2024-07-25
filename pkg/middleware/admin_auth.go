package middleware

import (
	"net/http"

	"github.com/blanc42/ecms/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth-token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization cookie is required"})
			c.Abort()
			return
		}

		adminID, err := utils.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Add admin ID to the context
		c.Set("admin_id", adminID)

		c.Next()
	}
}
