package api

import (
	"net/http"
	"strconv"
	"time"

	"reserveflow-v1/dao"
	"reserveflow-v1/models"

	"github.com/gin-gonic/gin"
)

type HoldReservationRequest struct {
	ResourceID uint   `json:"resource_id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

func HoldReservation(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}

	var req HoldReservationRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_START_TIME",
				"message": "start_time must be RFC3339 format",
			},
		})
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_END_TIME",
				"message": "end_time must be RFC3339 format",
			},
		})
		return
	}

	if !endTime.After(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_TIME_RANGE",
				"message": "end_time must be after start_time",
			},
		})
		return
	}

	resource, err := dao.GetResourceByID(req.ResourceID)
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

	if !resource.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESOURCE_INACTIVE",
				"message": "resource is inactive",
			},
		})
		return
	}

	withinWorkingHours, err := isWithinResourceWorkingHours(resource.ID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "WORKING_HOURS_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	if !withinWorkingHours {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "OUTSIDE_WORKING_HOURS",
				"message": "reservation time is outside resource working hours",
			},
		})
		return
	}

	hasConflict, err := dao.HasActiveReservationConflict(resource.ID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "CONFLICT_CHECK_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	if hasConflict {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_CONFLICT",
				"message": "resource is already reserved or held for this time range",
			},
		})
		return
	}

	reservation := models.Reservation{
		UserID:     userID,
		ResourceID: resource.ID,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     models.ReservationStatusHeld,
		ExpiresAt:  time.Now().Add(10 * time.Minute),
	}

	if err := dao.CreateReservation(&reservation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_HOLD_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	createdReservation, err := dao.GetReservationByID(reservation.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_DETAIL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    createdReservation,
	})
}

func GetMyReservations(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}

	reservations, err := dao.GetReservationsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATIONS_LIST_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    reservations,
	})
}
func ConfirmReservation(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	id64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_RESERVATION_ID",
				"message": "reservation id is invalid",
			},
		})
		return
	}

	reservation, err := dao.GetReservationByID(uint(id64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_NOT_FOUND",
				"message": "reservation not found",
			},
		})
		return
	}

	if reservation.UserID != userID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_YOUR_RESERVATION",
				"message": "you can only confirm your own reservation",
			},
		})
		return
	}

	if reservation.Status != models.ReservationStatusHeld {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_NOT_HELD",
				"message": "only held reservations can be confirmed",
			},
		})
		return
	}

	if time.Now().After(reservation.ExpiresAt) {
		reservation.Status = models.ReservationStatusExpired
		_ = dao.UpdateReservation(reservation)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_EXPIRED",
				"message": "reservation hold has expired",
			},
		})
		return
	}

	hasConflict, err := dao.HasActiveReservationConflictExceptID(
		reservation.ID,
		reservation.ResourceID,
		reservation.StartTime,
		reservation.EndTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "CONFLICT_CHECK_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	if hasConflict {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_CONFLICT",
				"message": "resource is already reserved or held for this time range",
			},
		})
		return
	}

	reservation.Status = models.ReservationStatusConfirmed

	if err := dao.UpdateReservation(reservation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_CONFIRM_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	confirmedReservation, err := dao.GetReservationByID(reservation.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_DETAIL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    confirmedReservation,
	})
}

func CancelReservation(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	id64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_RESERVATION_ID",
				"message": "reservation id is invalid",
			},
		})
		return
	}

	reservation, err := dao.GetReservationByID(uint(id64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_NOT_FOUND",
				"message": "reservation not found",
			},
		})
		return
	}

	if reservation.UserID != userID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_YOUR_RESERVATION",
				"message": "you can only cancel your own reservation",
			},
		})
		return
	}

	if reservation.Status == models.ReservationStatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "ALREADY_CANCELLED",
				"message": "reservation already cancelled",
			},
		})
		return
	}

	if reservation.Status == models.ReservationStatusExpired {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_EXPIRED",
				"message": "expired reservation cannot be cancelled",
			},
		})
		return
	}

	if reservation.Status == models.ReservationStatusHeld && time.Now().After(reservation.ExpiresAt) {
		reservation.Status = models.ReservationStatusExpired
		_ = dao.UpdateReservation(reservation)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_EXPIRED",
				"message": "reservation hold has expired",
			},
		})
		return
	}

	reservation.Status = models.ReservationStatusCancelled

	if err := dao.UpdateReservation(reservation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_CANCEL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	cancelledReservation, err := dao.GetReservationByID(reservation.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "RESERVATION_DETAIL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cancelledReservation,
	})
}
func getCurrentUserID(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "USER_ID_NOT_FOUND",
				"message": "user id not found",
			},
		})
		return 0, false
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
		return 0, false
	}

	return userID, true
}

func isWithinResourceWorkingHours(resourceID uint, startTime time.Time, endTime time.Time) (bool, error) {
	dayOfWeek := getDayOfWeek(startTime)

	workingHour, err := dao.GetWorkingHourByResourceIDAndDay(resourceID, dayOfWeek)
	if err != nil {
		return false, nil
	}

	if workingHour.IsClosed {
		return false, nil
	}

	openMinutes, err := clockToMinutes(workingHour.OpenTime)
	if err != nil {
		return false, err
	}

	closeMinutes, err := clockToMinutes(workingHour.CloseTime)
	if err != nil {
		return false, err
	}

	startMinutes := startTime.Hour()*60 + startTime.Minute()
	endMinutes := endTime.Hour()*60 + endTime.Minute()

	if startMinutes < openMinutes {
		return false, nil
	}

	if endMinutes > closeMinutes {
		return false, nil
	}

	return true, nil
}

func getDayOfWeek(t time.Time) string {
	switch t.Weekday() {
	case time.Monday:
		return models.DayMonday
	case time.Tuesday:
		return models.DayTuesday
	case time.Wednesday:
		return models.DayWednesday
	case time.Thursday:
		return models.DayThursday
	case time.Friday:
		return models.DayFriday
	case time.Saturday:
		return models.DaySaturday
	case time.Sunday:
		return models.DaySunday
	default:
		return ""
	}
}

func clockToMinutes(value string) (int, error) {
	parsedTime, err := time.Parse("15:04", value)
	if err != nil {
		return 0, err
	}

	return parsedTime.Hour()*60 + parsedTime.Minute(), nil
}
