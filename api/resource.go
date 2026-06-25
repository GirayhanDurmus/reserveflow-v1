package api

import (
	"net/http"
	"reserveflow-v1/middleware"
	"strconv"

	"reserveflow-v1/dao"
	"reserveflow-v1/models"

	"github.com/gin-gonic/gin"
)

type CreateResourceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
}

type UpdateResourceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
	IsActive    *bool  `json:"is_active"`
}

func CreateResource(c *gin.Context) {
	var req CreateResourceRequest

	if err := c.BindJSON(&req); err != nil {
		return
	}

	resource := models.Resource{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		IsActive:    true,
	}

	if err := dao.CreateResource(&resource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESOURCE_CREATE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resource,
	})
}

func GetAllResources(c *gin.Context) {
	resources, err := dao.GetAllResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESOURCES_LIST_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resources,
	})
}

func GetResourceByID(c *gin.Context) {
	id := c.Param("id")

	id64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_RESOURCE_ID",
				"message": "resource id is invalid",
			},
		})
		return
	}

	resource, err := dao.GetResourceByID(uint(id64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESOURCE_NOT_FOUND",
				"message": "resource not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resource,
	})
}

func UpdateResource(c *gin.Context) {
	id := c.Param("id")

	id64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_RESOURCE_ID",
				"message": "resource id is invalid",
			},
		})
		return
	}

	resource, err := dao.GetResourceByID(uint(id64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESOURCE_NOT_FOUND",
				"message": "resource not found",
			},
		})
		return
	}

	var req UpdateResourceRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}

	if req.Name != "" {
		resource.Name = req.Name
	}

	if req.Description != "" {
		resource.Description = req.Description
	}

	if req.Capacity > 0 {
		resource.Capacity = req.Capacity
	}

	if req.IsActive != nil {
		resource.IsActive = *req.IsActive
	}

	if err := dao.UpdateResource(resource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESOURCE_UPDATE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resource,
	})
}

func DeleteResource(c *gin.Context) {
	id := c.Param("id")

	id64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_RESOURCE_ID",
				"message": "resource id is invalid",
			},
		})
		return
	}

	if err := dao.DeleteResource(uint(id64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESOURCE_DELETE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message": "resource deleted",
		},
	})
}
func AddResourceURLs(r *gin.RouterGroup) {
	resources := r.Group("/resources")

	resources.GET("", GetAllResources)
	resources.GET("/:id", GetResourceByID)

	adminResources := r.Group("/admin/resources")
	adminResources.Use(middleware.AuthRequired())

	adminResources.POST("", middleware.RequireAdmin(), CreateResource)
	adminResources.PATCH("/:id", middleware.RequireResourceAdmin(), UpdateResource)
	adminResources.DELETE("/:id", middleware.RequireAdmin(), DeleteResource)

	adminResources.GET("/:id/admins", middleware.RequireResourceAdmin(), ListResourceAdmins)
	adminResources.POST("/:id/admins", middleware.RequireResourceAdmin(), AssignResourceAdmin)
	adminResources.DELETE("/:id/admins/:user_id", middleware.RequireResourceAdmin(), RemoveResourceAdmin)
}
