package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var database *gorm.DB

func InitDB() (*gorm.DB, error) {

	gagal := godotenv.Load()
	if gagal != nil {
		log.Fatal("Error loading .env file", gagal)
	}

	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	database, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return database, nil
}

func GetDB() *gorm.DB {
	return database
}

/*
func (database *gorm.DB) Close() {
	database, err := database.DB()
	if err != nil {
		log.Fatal("Error getting database connection:", err)
	}
	database.Close()
}

/*
func CloseDB() {
	if database != nil {
		sqlDB, err := database.DB()
		if err != nil {
			log.Fatal("Error closing the database: ", err)
			return
		}
		sqlDB.Close()
	}
}
*/
