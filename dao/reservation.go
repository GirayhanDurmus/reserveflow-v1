package dao

import (
	"time"

	"reserveflow-v1/commons"
	"reserveflow-v1/models"
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
<<<<<<< HEAD
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
=======
>>>>>>> 9d5ac127c0fd3b5689c141f3c90aa952448ce523
