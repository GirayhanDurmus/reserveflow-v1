package main

import (
	"fmt"
	"reserveflow-v1/seeder"

	"reserveflow-v1/api"
	"reserveflow-v1/commons"
	"reserveflow-v1/middleware"
	"reserveflow-v1/models"

	"github.com/gin-gonic/gin"
)

func main() {
	commons.LoadConfig()
	commons.ConnectPostgres()

	if err := commons.DB.AutoMigrate(&models.Role{}, &models.Permission{}, &models.User{}, &models.Resource{}, &models.WorkingHour{}, &models.Reservation{}, &models.ResourceAdmin{}); err != nil {
		panic(err)
	}
	seeder.SeedRolesPermisiions()
	if err := commons.DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_resource_admins_active ON resource_admins(user_id, resource_id) WHERE deleted_at IS NULL`).Error; err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/health", api.Health)

	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired(), middleware.RequireAdmin())
	admin.GET("/ping", api.AdminPing)

	baseGroup := r.Group("/")
	api.AddAuthURLs(baseGroup)
	api.AddResourceURLs(baseGroup)
	api.AddWorkingHoursURLs(baseGroup)
	api.AddReservationURLs(baseGroup)

	api.AddBackURLs(baseGroup)

	port := commons.AppConfig.AppPort

	fmt.Println("Server running on port:", port)

	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
