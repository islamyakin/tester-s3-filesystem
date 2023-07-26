package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
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

	// Konfigurasi sesi AWS
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

	// Inisialisasi layanan S3
	svc := s3.New(sess)

	// Konfigurasi upload
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(handler.Filename),
		Body:   file,
		ACL:    aws.String("public-read"),
	}

	// Melakukan upload file
	_, err = svc.PutObject(uploadInput)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/%s/%s", endpoint, bucket, handler.Filename)
	fmt.Fprintf(w, "File berhasil diunggah ke S3. URL: %s", url)
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

	vars := mux.Vars(r)
	fileName := vars["filename"]

	// Konfigurasi sesi AWS
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

	// Inisialisasi layanan S3
	svc := s3.New(sess)

	// Konfigurasi delete
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}

	// Melakukan delete objek
	_, err = svc.DeleteObject(deleteInput)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File berhasil dihapus dari S3: %s", fileName)
}
