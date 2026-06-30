package dao

import (
	"reserveflow-v1/commons"
	"reserveflow-v1/models"
)

func GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	err := commons.DB.Preload("Permission").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func CreatePermissions(perm *models.Permission) error {
	return commons.DB.Create(perm).Error
}

func GetPermissionsByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := commons.DB.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}
