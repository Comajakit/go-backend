package db

import (
	"fmt"
	"go-backend/config"
	"go-backend/database/models"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	config.InitConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	db.AutoMigrate(
		&models.User{},
		&models.Trip{},
		&models.TripCategory{},
		&models.TripDetail{},
	)
}

func GetHashedPassword(username string) (string, error) {
	var user models.User

	// Find the user with the given username
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User not found
			fmt.Println("user not found")
			return "", nil
		}
		// An error occurred during the query
		return "", err
	}

	return user.Password, nil
}
