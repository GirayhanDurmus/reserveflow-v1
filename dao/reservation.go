package dao

import (
	"errors"
	"time"

	"reserveflow-v1/commons"
	"reserveflow-v1/models"

	"gorm.io/gorm"
)

func CreateReservation(reservation *models.Reservation) error {
	return commons.DB.Create(reservation).Error
}

func GetReservationsByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation

	err := commons.DB.
		Preload("Resource").
		Where("user_id = ?", userID).
		Order("start_time desc").
		Find(&reservations).Error

	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func GetReservationByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation

	err := commons.DB.
		Preload("User").
		Preload("Resource").
		First(&reservation, id).Error

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func UpdateReservation(reservation *models.Reservation) error {
	return commons.DB.Save(reservation).Error
}

func HasActiveReservationConflict(resourceID uint, startTime time.Time, endTime time.Time) (bool, error) {
	var count int64

	now := time.Now()

	err := commons.DB.Model(&models.Reservation{}).
		Where("resource_id = ?", resourceID).
		Where("start_time < ? AND end_time > ?", endTime, startTime).
		Where(
			"status = ? OR (status = ? AND expires_at > ?)",
			models.ReservationStatusConfirmed,
			models.ReservationStatusHeld,
			now,
		).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetWorkingHourByResourceIDAndDay(resourceID uint, dayOfWeek string) (*models.WorkingHour, error) {
	var workingHour models.WorkingHour

	err := commons.DB.
		Where("resource_id = ?", resourceID).
		Where("day_of_week = ?", dayOfWeek).
		First(&workingHour).Error

	if err != nil {
		return nil, err
	}

	return &workingHour, nil
}

func HasActiveReservationConflictExceptID(reservationID uint, resourceID uint, startTime time.Time, endTime time.Time) (bool, error) {
	var count int64

	now := time.Now()

	err := commons.DB.Model(&models.Reservation{}).
		Where("id <> ?", reservationID).
		Where("resource_id = ?", resourceID).
		Where("start_time < ? AND end_time > ?", endTime, startTime).
		Where(
			"status = ? OR (status = ? AND expires_at > ?)",
			models.ReservationStatusConfirmed,
			models.ReservationStatusHeld,
			now,
		).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetExpiredReservationIDs(currentTime time.Time) ([]uint, error) {
	var expiredIDs []uint
	err := commons.DB.
		Model(&models.Reservation{}).
		Where("status = ? AND expires_at < ?", models.ReservationStatusHeld, currentTime).
		Pluck("id", &expiredIDs).Error

	return expiredIDs, err
}

func MarkReservationAsExpired(resID uint) error {
	return commons.DB.
		Model(&models.Reservation{}).
		Where("id = ? AND status = ?", resID, models.ReservationStatusHeld).
		Update("status", models.ReservationStatusExpired).Error
}

func HoldReservationWithTx(reservation *models.Reservation) error {
	return commons.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		now := time.Now()

		err := tx.Model(&models.Reservation{}).
			Where("resource_id = ?", reservation.ResourceID).
			Where("start_time < ? AND end_time > ?", reservation.EndTime, reservation.StartTime).
			Where(
				"status = ? OR (status = ? AND expires_at > ?)",
				models.ReservationStatusConfirmed,
				models.ReservationStatusHeld,
				now,
			).
			Count(&count).Error

		if err != nil {
			return err
		}

		if count > 0 {
			return errors.New("RESERVATION_CONFLICT")
		}

		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		return nil
	})
}
