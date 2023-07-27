package models

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"type:unique"`
	Password string
	Role     string
}
