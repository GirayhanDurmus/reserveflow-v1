package middleware

import (
	"net/http"

	"reserveflow-v1/models"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "ROLE_NOT_FOUND",
					"message": "role not found",
				},
			})
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_ROLE",
					"message": "role is invalid",
				},
			})
			c.Abort()
			return
		}

		if role != models.UserRoleAdmin {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "ADMIN_REQUIRED",
					"message": "admin permission required",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
