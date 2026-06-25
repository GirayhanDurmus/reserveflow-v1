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

	err := commons.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := commons.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}
func UpdateUserRole(userID uint, role string) error {
	return commons.DB.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}
