package server

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/gorilla/mux"
	"github.com/islamyakin/tester-s3-filesystem/db"
	"github.com/islamyakin/tester-s3-filesystem/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
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
func HandleS3Delete(w http.ResponseWriter, r *http.Request) {
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

	database := db.GetDB()
	vars := mux.Vars(r)
	fileName := vars["filename"]

	var file models.File
	if err := database.Where("file_name = ?", fileName).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "File not found in database", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get file from database", http.StatusInternalServerError)
		log.Println(err)
		return
	}

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

	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}

	_, err = svc.HeadObject(headInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			http.Error(w, "File not found in S3", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get file in S3", http.StatusInternalServerError)
	}

	// Configure delete
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}

	// Delete object
	_, err = svc.DeleteObject(deleteInput)
	if err != nil {
		http.Error(w, "Failed to delete file from s3", http.StatusInternalServerError)
		return
	}

	deleteQuery := "DELETE FROM files WHERE file_name = ?"
	if err := database.Exec(deleteQuery, fileName).Error; err != nil {
		http.Error(w, "Failed to delete data from database", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	fmt.Fprintf(w, "File berhasil dihapus dari S3 dan Database")
}
