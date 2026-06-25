package main

import (
	"fmt"

	"reserveflow-v1/api"
	"reserveflow-v1/commons"
	"reserveflow-v1/middleware"
	"reserveflow-v1/models"

	"github.com/gin-gonic/gin"
)

func main() {
	commons.LoadConfig()
	commons.ConnectPostgres()

	if err := commons.DB.AutoMigrate(&models.User{}, &models.Resource{}, &models.WorkingHour{}, &models.Reservation{}, &models.ResourceAdmin{}); err != nil {
		panic(err)
	}
	if err := commons.DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_resource_admins_active ON resources_admin(user_id, resource_id) WHERE deleted_at IS NULL`).Error; err != nil {
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

	port := commons.AppConfig.AppPort

	fmt.Println("Server running on port:", port)

	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
