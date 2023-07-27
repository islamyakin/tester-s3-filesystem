package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	FileName string
	FileType string
	S3URL    string
}
