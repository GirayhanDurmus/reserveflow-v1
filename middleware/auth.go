package middleware

import (
	"net/http"
	"reserveflow-v1/dao"
	"strings"

	"reserveflow-v1/commons"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "TOKEN_REQUIRED",
					"message": "token is required",
				},
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(commons.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "token is invalid",
				},
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN_CLAIMS",
					"message": "token claims are invalid",
				},
			})
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
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

		c.Set("user_id", uint(userIDFloat))
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

func RequirePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED_CODE",
					"message": "User Not Found",
				},
			})
			c.Abort()
			return
		}
		user, err := dao.GetUserByID(userIDValue.(uint))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "USER_NOT_FOUND",
					"message": "User Not Found",
				},
			})
			c.Abort()
			return
		}

		role, err := dao.GetWithPermission(user.RoleID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED_CODE",
					"message": "role setup is invalid",
				},
			})
			c.Abort()
			return
		}

		requsetMethod := c.Request.Method
		requestEndpoint := c.FullPath()

		hashpermission := false

		for _, perm := range role.Permission {
			if strings.ToUpper(perm.Method) == strings.ToUpper(requsetMethod) && perm.Endpoint == requestEndpoint {
				hashpermission = true
				break
			}
		}
		if !hashpermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "You have no permission to perform this action",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
