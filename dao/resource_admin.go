package dao

import (
	"reserveflow-v1/commons"
	"reserveflow-v1/models"
)

func CreateResourceAdmin(resourceAdmin *models.ResourceAdmin) error {
	return commons.DB.Create(resourceAdmin).Error
}

func IsResourceAdmin(userID uint, resourceID uint) (bool, error) {
	var count int64

	err := commons.DB.Model(&models.ResourceAdmin{}).Where("user_id = ? AND resource_id = ?", userID, resourceID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetResourceAdminsByResourceID(resourceID uint) ([]models.ResourceAdmin, error) {
	var admins []models.ResourceAdmin
	err := commons.DB.Preload("User").Preload("Resource").Where("resource_id = ?", resourceID).Find(&admins).Error
	if err != nil {
		return nil, err
	}
	return admins, nil
}

func DeleteResourceAdmin(resourceID uint, userID uint) error {
	return commons.DB.Where("resource_id = ? AND user_id = ? ", resourceID, userID).Delete(&models.ResourceAdmin{}).Error
}
func CountActiveResourceAdmins(userID uint) (int64, error) {
	var count int64
	err := commons.DB.Model(&models.ResourceAdmin{}).
		Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
