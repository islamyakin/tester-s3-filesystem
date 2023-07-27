package models

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"unique"`
	Password string
	Role     string
}
