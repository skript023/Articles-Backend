package database

import (
	"ArticleBackend/config"
	"ArticleBackend/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Env("DB_USERNAME"), config.Env("DB_PASSWORD"), config.Env("DB_HOST"), config.Env("DB_PORT"), config.Env("DB_DATABASE"))
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Connection with database, failed.")
		os.Exit(2)
	}

	log.Println("Connected to database, success.")
	DB.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Migrations Running")
	DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Category{}, &models.Comment{}, &models.Contact{}, &models.Role{})
}
