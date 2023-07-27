package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"github.com/islamyakin/tester-s3-filesystem/db"
	"github.com/islamyakin/tester-s3-filesystem/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

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

	err = removeFileFromS3IfNotInDatabase(fileName)
	if err != nil {
		// Jika file tidak ditemukan di S3, berikan respons dan selesaikan penanganan request
		http.Error(w, "File not found in S3", http.StatusNotFound)
		return
	}

	database := db.GetDB()
	var file models.File
	if err := database.Where("file_name = ?", fileName).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "File not found in database", http.StatusNotFound)
			return
		}
		// Panggil fungsi removeFileFromS3IfNotInDatabase dengan fileName sebagai argumen
		err := removeFileFromS3IfNotInDatabase(fileName)
		http.Error(w, "Failed to get file from database", http.StatusInternalServerError)
		log.Println(err)
		return

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

	err = removeFileFromS3IfNotInDatabase(fileName)

	fmt.Fprintf(w, "File berhasil dihapus dari S3 dan Database")
}
func removeFileFromS3IfNotInDatabase(fileName string) error {
	database := db.GetDB()

	var count int64
	if err := database.Model(&models.File{}).Where("file_name = ?", fileName).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		// File tidak ditemukan di database, hapus dari S3
		accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
		secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
		endpoint := os.Getenv("AWS_ENDPOINT")
		bucket := os.Getenv("AWS_BUCKET")
		region := os.Getenv("AWS_REGION")

		// Configure session AWS
		sess, err := session.NewSession(&aws.Config{
			Region:           aws.String(region),
			Endpoint:         aws.String(endpoint),
			S3ForcePathStyle: aws.Bool(false),
			Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		})
		if err != nil {
			return fmt.Errorf("failed to create session: %v", err)
		}

		// Initialize service S3
		svc := s3.New(sess)

		// Configure delete
		deleteInput := &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fileName),
		}

		// Delete object
		_, err = svc.DeleteObject(deleteInput)
		if err != nil {
			return fmt.Errorf("failed to delete file from S3: %v", err)
		}
	}

	return nil
}
