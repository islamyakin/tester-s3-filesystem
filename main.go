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
	if len(os.Args) < 4 {
		fmt.Println("usage:go run main.go <name file in local> <name file in s3>")
	}
	nameFileLocal := os.Args[1]
	nameFileS3 := os.Args[2]

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

	err = uploadFileToS3(accessKeyID, secretAccessKey, endpoint, bucket, region, nameFileLocal, nameFileS3)
	if err != nil {
		fmt.Println("Failed upload to s3,", err)
		return
	}
}

func uploadFileToS3(accessKeyID, secretAccessKey, endpoint, bucket, region, nameFileLocal, nameFileS3 string) error {

	// Konfigurasi sesi AWS
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(false),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	// mulai sesi s3
	svc := s3.New(sess)

	file, err := os.Open(nameFileLocal)
	if err != nil {
		return err
	}
	defer file.Close()

	// Konfig upload
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(nameFileS3),
		Body:   file,
		ACL:    aws.String("public-read"),
	}

	// Upload
	_, err = svc.PutObject(uploadInput)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s/%s", endpoint, bucket, nameFileS3)
	fmt.Println("URL Public object: ", url)

	return nil
}
