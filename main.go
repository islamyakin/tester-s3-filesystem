package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/joho/godotenv"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env file", err)
		return
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("AWS_ENDPOINT")
	bucket := os.Getenv("AWS_BUCKET")
	region := os.Getenv("AWS_REGION")

	// Konfigurasi sesi AWS
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(false),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		fmt.Println("Failed to create session:", err)
		return
	}

	// mulai sesi s3
	svc := s3.New(sess)

	// nama file di lokal
	localFilePath := "testing1.txt"

	// nama untuk di s3
	s3ObjectKey := "testing1.txt"

	// cek file di lokal
	file, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	// Konfig upload
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3ObjectKey),
		Body:   file,
		ACL:    aws.String("public-read"),
	}

	// Upload
	_, err = svc.PutObject(uploadInput)
	if err != nil {
		fmt.Println("Failed to upload file:", err)
		return
	}
	url := fmt.Sprintf("%s/%s/%s", endpoint, bucket, s3ObjectKey)
	fmt.Println("URL Public object: ", url)
}
