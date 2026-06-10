package dao

import (
	"reserveflow-v1/commons"
	"reserveflow-v1/models"
)

func CreateResource(resource *models.Resource) error {
	return commons.DB.Create(resource).Error
}

func GetAllResources() ([]models.Resource, error) {
	var resources []models.Resource

	err := commons.DB.Find(&resources).Error
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func GetResourceByID(id uint) (*models.Resource, error) {
	var resource models.Resource

	err := commons.DB.First(&resource, id).Error
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

func UpdateResource(resource *models.Resource) error {
	return commons.DB.Save(resource).Error
}

func DeleteResource(id uint) error {
	return commons.DB.Delete(&models.Resource{}, id).Error
}
