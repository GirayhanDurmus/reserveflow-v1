package middleware

import (
	"net/http"
	"reserveflow-v1/dao"
	"strconv"

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
					"code":    "SUPER_ADMIN_REQUIRED",
					"message": "super admin permission required",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireResourceAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "ROLE_NOTFOUND",
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
					"code":    "INVAL_ID_CODE",
					"message": "role is invalid",
				},
			})
			c.Abort()
			return
		}
		if role == models.UserRoleAdmin {
			c.Next()
			return
		}
		if role != models.UserRoleResourceAdmin && role != models.UserRoleManager {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INSUFFICIENT_ROLE",
					"message": "resource admin permission required",
				},
			})
			c.Abort()
			return
		}
		userIDValue, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "USER_ID_NOTFOUND",
					"message": "user not found",
				},
			})
			c.Abort()
			return
		}
		userID, ok := userIDValue.(uint)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_USER_ID",
					"message": "user id is invalid",
				},
			})
			c.Abort()
			return
		}
		resourceIDParam := c.Param("id")

		resourceID64, err := strconv.ParseUint(resourceIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_RESOURCE_ID",
					"message": "resource id is invalid",
				},
			})
			c.Abort()
			return
		}
		isResourceAdmin, err := dao.IsResourceAdmin(userID, uint(resourceID64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RESOURCE_ADMIN_CHECK_ERROR",
					"message": err.Error(),
				},
			})
			c.Abort()
			return
		}

		if !isResourceAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED_ACCESS",
					"message": "You are not authorized to access this resource",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}

}
