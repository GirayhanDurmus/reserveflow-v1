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

	if err := commons.DB.AutoMigrate(&models.User{}, &models.Resource{}, &models.WorkingHour{}, &models.Reservation{}, &models.Reservation{}); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/health", api.Health)
	r.GET("/resources", api.GetAllResources)
	r.GET("/resources/:id", api.GetResourceByID)

	r.POST("/auth/register", api.Register)
	r.POST("/auth/login", api.Login)

	auth := r.Group("/auth")
	auth.Use(middleware.AuthRequired())
	auth.GET("/me", api.Me)

	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired())
	admin.Use(middleware.RequireAdmin())
	admin.GET("/ping", api.AdminPing)
	admin.POST("/resources", api.CreateResource)
	admin.PATCH("/resources/:id", api.UpdateResource)
	admin.DELETE("/resources/:id", api.DeleteResource)

	admin.POST("/resources/:id/working-hours", api.SetWorkingHours)
	admin.GET("/resources/:id/working-hours", api.GetWorkingHours)

	reservation := r.Group("/reservations")
	reservation.Use(middleware.AuthRequired())
	reservation.POST("/hold", api.HoldReservation)
	reservation.POST("/:id/confirm", api.ConfirmReservation)
	reservation.GET("/my", api.GetMyReservations)

	port := commons.AppConfig.AppPort

	fmt.Println("Server running on port:", port)

	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
