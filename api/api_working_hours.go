package api

import (
	"net/http"
	"strconv"

	"reserveflow-v1/dao"
	"reserveflow-v1/models"

	"github.com/gin-gonic/gin"
)

type WorkingHourItemRequest struct {
	DayOfWeek string `json:"day_of_week"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
	IsClosed  bool   `json:"is_closed"`
}

type SetWorkingHoursRequest struct {
	WorkingHours []WorkingHourItemRequest `json:"working_hours"`
}

func SetWorkingHours(c *gin.Context) {
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

	var req SetWorkingHoursRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}

	if err := dao.DeleteWorkingHoursByResourceID(resource.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "WORKING_HOURS_DELETE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	var createdWorkingHours []models.WorkingHour

	for _, item := range req.WorkingHours {
		workingHour := models.WorkingHour{
			ResourceID: resource.ID,
			DayOfWeek:  item.DayOfWeek,
			OpenTime:   item.OpenTime,
			CloseTime:  item.CloseTime,
			IsClosed:   item.IsClosed,
		}

		if err := dao.CreateWorkingHour(&workingHour); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "WORKING_HOUR_CREATE_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		createdWorkingHours = append(createdWorkingHours, workingHour)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    createdWorkingHours,
	})
}

func GetWorkingHours(c *gin.Context) {
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

	workingHours, err := dao.GetWorkingHoursByResourceID(resource.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "WORKING_HOURS_LIST_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    workingHours,
	})
}
