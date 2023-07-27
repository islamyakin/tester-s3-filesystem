package models

type File struct {
	ID       uint   `gorm:"primary_key"`
	FileName string `gorm:"type:varchar(255);not null"`
	FileType string `gorm:"type:varchar(100)"`
	S3URL    string `gorm:"type:varchar(255)"`
	Bucket   string `gorm:"type:varchar(255)"`
}

func (File) TableName() string {
	return "files"
}
