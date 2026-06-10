package commons

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {
	dsn := AppConfig.PostgresDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed the connect database: %v\n", err)
	}
	DB = db
	fmt.Println("Succesfully connected to database")
}
