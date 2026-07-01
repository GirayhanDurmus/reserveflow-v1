package dao

import (
	"reserveflow-v1/commons"
	"reserveflow-v1/models"
)

func CreateUser(user *models.User) error {
	return commons.DB.Create(user).Error
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := commons.DB.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := commons.DB.Preload("Role").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}
func UpdateUserRole(userID uint, roleName string) error {
	var role models.Role
	if err := commons.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}
	return commons.DB.Model(&models.User{}).Where("id = ?", userID).Update("role_id", role.ID).Error
}
