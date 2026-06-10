package dao

import (
	"reserveflow-v1/commons"
	"reserveflow-v1/models"
)

func CreateWorkingHour(workingHour *models.WorkingHour) error {
	return commons.DB.Create(workingHour).Error
}

func GetWorkingHoursByResourceID(resourceID uint) ([]models.WorkingHour, error) {
	var workingHours []models.WorkingHour

	err := commons.DB.
		Where("resource_id = ?", resourceID).
		Order("id asc").
		Find(&workingHours).Error

	if err != nil {
		return nil, err
	}

	return workingHours, nil
}

func DeleteWorkingHoursByResourceID(resourceID uint) error {
	return commons.DB.
		Where("resource_id = ?", resourceID).
		Delete(&models.WorkingHour{}).Error
}
