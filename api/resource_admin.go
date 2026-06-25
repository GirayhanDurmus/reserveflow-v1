package api

import (
	"net/http"
	"reserveflow-v1/dao"
	"reserveflow-v1/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssignResourceAdminRequest struct {
	UserID uint `json:"user_id"`
}

func AssignResourceAdmin(c *gin.Context) {
	resourceID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
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

	var req AssignResourceAdminRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}

	if req.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "user_id is required",
			},
		})
		return
	}

	_, err = dao.GetResourceByID(uint(resourceID64))
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

	user, err := dao.GetUserByID(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "USER_NOT_FOUND",
				"message": "user not found",
			},
		})
		return
	}

	resourceAdmin := models.ResourceAdmin{
		UserID:     req.UserID,
		ResourceID: uint(resourceID64),
	}

	if err := dao.CreateResourceAdmin(&resourceAdmin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESOURCE_ADMIN_CREATE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Sadece "user" rolündekileri resource-admin'e yükselt.
	// "manager" ve "admin" rollerine dokunma.
	if user.Role == models.UserRoleUser {
		if err := dao.UpdateUserRole(user.ID, models.UserRoleResourceAdmin); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "USER_ROLE_UPDATE_ERROR",
					"message": err.Error(),
				},
			})
			return
		}
	}

	// Kayıt oluşturuldu; User ve Resource alanlarını dolu döndürmek için reload yap.
	fullRecord, err := dao.GetResourceAdminByID(resourceAdmin.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": resourceAdmin})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    fullRecord,
	})
}

func ListResourceAdmins(c *gin.Context) {
	resourceID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
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

	admins, err := dao.GetResourceAdminsByResourceID(uint(resourceID64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESOURCE_ADMINS_LIST_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    admins,
	})
}

func RemoveResourceAdmin(c *gin.Context) {
	resourceID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
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

	userID64, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "user id is invalid",
			},
		})
		return
	}

	if err := dao.DeleteResourceAdmin(uint(resourceID64), uint(userID64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESOURCE_ADMIN_DELETE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Silme sonrası başka aktif ataması kalmadıysa rolü "user"a döndür.
	count, err := dao.CountActiveResourceAdmins(uint(userID64))
	if err == nil && count == 0 {
		_ = dao.UpdateUserRole(uint(userID64), models.UserRoleUser)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message": "resource admin removed",
		},
	})
}
