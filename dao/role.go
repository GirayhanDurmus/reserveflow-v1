package dao

import (
	"reserveflow-v1/commons"
	"reserveflow-v1/models"
)

func GetWithPermission(roleID uint) (*models.Role, error) {
	var role models.Role
	err := commons.DB.Preload("Permissions").Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := commons.DB.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func AssignPermissionToRole(roleID uint, perm *models.Permission) error {
	role := models.Role{}
	role.ID = roleID
	return commons.DB.Model(&role).Association("Permission").Append(perm)

}
