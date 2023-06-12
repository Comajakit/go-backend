package db

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/your-username/your-project-name/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	config.InitConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.name"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db

	// AutoMigrate will create the User table if it doesn't exist
	db.AutoMigrate(&models.User{})
}
