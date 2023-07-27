package main

import (
	"fmt"
	"github.com/islamyakin/tester-s3-filesystem/db"
	"github.com/islamyakin/tester-s3-filesystem/server"
	"gorm.io/gorm"
	"log"
)

func intiliazeDatabase() (*gorm.DB, error) {
	database, err := db.InitDB()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	if err := db.RunMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %v", err)
	}

	return database, nil
}

func main() {
	database, err := intiliazeDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if sqlDB, err := database.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	server.StartServer()
}
