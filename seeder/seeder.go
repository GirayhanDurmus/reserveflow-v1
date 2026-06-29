package seeder

import (
	"fmt"
	"reserveflow-v1/commons"
	"reserveflow-v1/models"

	"gorm.io/gorm"
)

func SeedRolesPermisiions() {

	roles := []models.Role{
		{Name: models.UserRoleUser, Description: "Normal User"},
		{Name: models.UserRoleResourceAdmin, Description: "System Manager"},
		{Name: models.UserRoleAdmin, Description: "Normal Admin"},
		{Name: models.UserRoleResourceAdmin, Description: "Resource Manager"},
		{Name: models.UserRoleManager, Description: "Manager"},
	}
	for i, role := range roles {
		var existingRole models.Role
		err := commons.DB.Where("name = ?", role.Name).First(&existingRole).Error
		if err == gorm.ErrRecordNotFound {
			commons.DB.Create(&roles[i])
			fmt.Printf("Role '%s' created\n", role.Name)
		}
	}
	adminPermissions := []models.Permission{
		{Method: "GET", Endpoint: "/admin/roles"},
		{Method: "POST", Endpoint: "/admin/roles/permision"},
		{Method: "POST", Endpoint: "/admin/roles/:id/permision"},
	}
	for i, perm := range adminPermissions {
		var existingPerm models.Permission
		err := commons.DB.Where("method = ? AND endpoint = ?", perm.Method, perm.Endpoint).First(&existingPerm).Error
		if err == gorm.ErrRecordNotFound {
			commons.DB.Create(&adminPermissions[i])
			fmt.Printf("Permission '%s %s' created\n", perm.Method, perm.Endpoint)
		} else {
			adminPermissions[i] = existingPerm
		}
	}
	var AdminRole models.Role
	if err := commons.DB.Where("name = ?", models.UserRoleAdmin).First(&AdminRole).Error; err == nil {
		commons.DB.Model(&AdminRole).Association("Permission").Append(&adminPermissions)
		fmt.Println("Admin permissions assigned to Admin role")
	}

}
