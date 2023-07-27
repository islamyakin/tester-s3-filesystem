package server

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/islamyakin/tester-s3-filesystem/db"
	"github.com/islamyakin/tester-s3-filesystem/models"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func HandleS3Upload(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("AWS_ENDPOINT")
	bucket := os.Getenv("AWS_BUCKET")
	region := os.Getenv("AWS_REGION")

	authHeader := r.Header.Get("Authorization")
	if authHeader != "your_auth_token" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB maksimum ukuran file
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Configure session AWS
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(false),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Initialize service S3
	svc := s3.New(sess)

	// Configure upload
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(handler.Filename),
		Body:   file,
		ACL:    aws.String("public-read"),
	}

	// Upload file to S3
	_, err = svc.PutObject(uploadInput)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/%s/%s", endpoint, bucket, handler.Filename)

	database := db.GetDB()

	var count int64
	database.Model(&models.File{}).Where("file_name = ?", handler.Filename).Count(&count)
	if count > 0 {
		http.Error(w, "File sudah ada di database", http.StatusBadRequest)
		return
	}

	insertQuery := "INSERT INTO files (file_name, file_type, s3_url, bucket) VALUES (?, ?, ?, ?)"
	if err := database.Exec(insertQuery, handler.Filename, handler.Header.Get("Content-Type"), url, bucket).Error; err != nil {
		http.Error(w, "Failed to insert data to database", http.StatusInternalServerError)
		log.Println("Failed to insert data to database:", err)
		return
	}
	fmt.Fprintf(w, "File berhasil diupload ke S3: %s", url)
}
