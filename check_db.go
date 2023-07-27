package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Ambil konfigurasi koneksi dari environment variables
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")

	// Buat string koneksi
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	fmt.Println("MYSQL_USER:", os.Getenv("MYSQL_USER"))
	fmt.Println("MYSQL_PASSWORD:", os.Getenv("MYSQL_PASSWORD"))
	fmt.Println("MYSQL_HOST:", os.Getenv("MYSQL_HOST"))
	fmt.Println("MYSQL_PORT:", os.Getenv("MYSQL_PORT"))
	fmt.Println("MYSQL_DATABASE:", os.Getenv("MYSQL_DATABASE"))

	fmt.Println("Data Source Name:", dataSourceName)

	// Membuat koneksi ke database
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Testing koneksi ke database
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to the MySQL database successfully.")
}
