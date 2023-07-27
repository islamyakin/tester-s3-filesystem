package service

import (
	"fmt"
	"github.com/islamyakin/tester-s3-filesystem/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DbUserAuth *gorm.DB

func InitDBAuth() (*gorm.DB, error) {
	failed := godotenv.Load()
	if failed != nil {
		fmt.Errorf("failed to load .env file: %w", failed)
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

	createAdminIfNotExists()
	return DbUserAuth, nil
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

func RunMigrationsUser() error {
	// Jalankan migrasi untuk tabel pengguna (user) di sini
	// Contoh implementasi:
	err := DbUserAuth.AutoMigrate(&models.User{})
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	// Tambahkan migrasi untuk tabel-tabel lain yang diperlukan di sini

	return nil
}
func HashPassword(password string) string {
	// Implementasikan fungsi hashing password sesuai kebutuhan aplikasi Anda
	// Di sini kami hanya mengembalikan password tanpa hash sebagai contoh
	return password
}
