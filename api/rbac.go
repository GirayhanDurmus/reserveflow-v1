package api

import (
	"net/http"
	"reserveflow-v1/dao"
	"reserveflow-v1/middleware"
	"reserveflow-v1/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreatePermissionRequest struct {
	Method   string `json:"method" binding:"required"`
	Endpoint string `json:"endpoint" binding:"required"`
}

type AssignPermissionRequest struct {
	PermissionID uint `json:"permission_id"`
}

func GetAllRoles(c *gin.Context) {
	roles, err := dao.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    roles,
	})
}

func CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}

	perm := models.Permission{
		Method:   req.Method,
		Endpoint: req.Endpoint,
	}
	if err := dao.CreatePermissions(&perm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    perm,
	})
}

func AssignPermission(c *gin.Context) {
	roleIDParam := c.Param("id")
	roleID64, err := strconv.ParseUint(roleIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid roleID",
		})
		return
	}
	var req AssignPermissionRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}

	perm, err := dao.GetPermissionsByID(req.PermissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "permission not found",
		})
		return
	}
	if err := dao.AssignPermissionToRole(uint(roleID64), perm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"Mmessage": "İzin basşarıyla atandı",
	})
}

func AddBackURLs(r *gin.RouterGroup) {
	rbac := r.Group("/admin/roles")
	rbac.Use(middleware.AuthRequired(), middleware.RequirePermission())
	rbac.GET("", GetAllRoles)
	rbac.POST("/permissions", CreatePermission)
	rbac.POST("/:id/permissions", AssignPermission)

}
