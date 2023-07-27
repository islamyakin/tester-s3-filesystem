package service

import (
	"fmt"
	"os"

	"github.com/islamyakin/tester-s3-filesystem/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DbUserAuth *gorm.DB

func main() {
	failed := godotenv.Load()
	if failed != nil {

	}

	dbUser := os.Getenv("AUTH_MYSQL_USER")
	dbPass := os.Getenv("AUTH_MYSQL_PASSWORD")
	dbHost := os.Getenv("AUTH_MYSQL_HOST")
	dbPort := os.Getenv("AUTH_MYSQL_PORT")
	dbName := os.Getenv("AUTH_MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DbUserAuth, err = gorm.Open(mysql.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		panic(err)
	}

	err_migrate := DbUserAuth.AutoMigrate(&models.User{})
	if err != nil {
		panic(err_migrate)
	}
	createAdminIfNotExists()
}

func createAdminIfNotExists() {
	// Cek apakah pengguna admin sudah ada
	var admin models.User
	result := DbUserAuth.Where("username = ?", "admin").First(&admin)
	if result.Error == nil {
		return
	}

	// Jika pengguna admin belum ada, buat pengguna admin
	admin = models.User{
		Username: "admin",
		Password: HashPassword("admin123"), // Contoh, sebaiknya gunakan hashing untuk password
		Role:     "admin",
	}
	DbUserAuth.Create(&admin)
}
func HashPassword(password string) string {
	// Implementasikan fungsi hashing password sesuai kebutuhan aplikasi Anda
	// Di sini kami hanya mengembalikan password tanpa hash sebagai contoh
	return password
}
