package db

import (
	"github.com/islamyakin/tester-s3-filesystem/models"
)

func RunMigrations() error {
	// Get the GORM database instance
	database, err := InitDB()
	if err != nil {
		return err
	}

	// Auto-migrate the File model to create the 'files' table and columns
	err = database.AutoMigrate(&models.File{})
	if err != nil {
		return err
	}
	return nil
}
